package redactly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
)

func redactStructFields(obj interface{}) {
	val := reflect.ValueOf(obj).Elem()
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tagVal, tagExists := field.Tag.Lookup("redactly")

		switch field.Type.Kind() {
		case reflect.Pointer:
			if tagExists {
				if field.Type.String() == "*string" {
					if !val.Field(i).IsNil() {
						val.Field(i).Set(reflect.ValueOf(&tagVal))
					}
				}
			}
		case reflect.String:
			if tagExists {
				val.Field(i).SetString(tagVal)
			}
		case reflect.Struct:
			redactStructFields(val.Field(i).Addr().Interface())
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.Struct {
				slice := val.Field(i)
				for j := 0; j < slice.Len(); j++ {
					elem := slice.Index(j).Addr().Interface()
					redactStructFields(elem)
				}
			}
		}
	}
}

func redactFormFields(obj interface{}) {
	val := reflect.ValueOf(obj).Elem()
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tagVal, tagExists := field.Tag.Lookup("redactly")

		switch field.Type.Kind() {
		case reflect.Pointer:
			if tagExists {
				if field.Type.String() == "*string" {
					if !val.Field(i).IsNil() {
						val.Field(i).Set(reflect.ValueOf(&tagVal))
					}
				}
			}
		case reflect.String:
			if tagExists {
				val.Field(i).SetString(tagVal)
			}
		}
	}
}

func redactHeaderField(k string, redactList *[]string) bool {
	for _, redactionKey := range *redactList {
		if k == redactionKey {
			return true
		}
	}
	return false
}

func DecodeForm[T any](target *T, querystring string) interface{} {
	form, err := url.ParseQuery(querystring)
	if err != nil {
		return err.Error()
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(target, form); err != nil {
		return err.Error()
	}

	redactFormFields(target)
	return target
}

func Unmarshal(data []byte, v any) interface{} {
	var bMap map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		if err := json.Unmarshal(data, &bMap); err != nil {
			return fmt.Sprint(bMap)
		}
	}

	redactStructFields(v)

	return v
}

func Redact(v any) interface{} {
	payloadBytes, err := json.Marshal(v)
	if err != nil {
		return v
	}

	var bMap map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &v); err != nil {
		if err := json.Unmarshal(payloadBytes, &bMap); err != nil {
			return fmt.Sprint(bMap)
		}
	}

	redactStructFields(v)

	return v
}

func RedactHttpHeader(header *http.Header, supplimentaryRedactionList *[]string) *http.Header {
	redactList := []string{"Authorization"}

	if supplimentaryRedactionList != nil {
		redactList = append(redactList, *supplimentaryRedactionList...)
	}

	redactedRequest := header.Clone()
	for k := range *header {
		if redactHeaderField(k, &redactList) {
			redactedRequest.Set(k, "REDACTED")
		}
	}
	return &redactedRequest
}

type FieldLogger interface {
	logrus.FieldLogger
}
