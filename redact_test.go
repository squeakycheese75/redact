package redact

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UnitTestSuite struct {
	suite.Suite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &UnitTestSuite{})
}

func (uts *UnitTestSuite) TestWeReplaceString() {
	type SmallOrder struct {
		ID        string `json:"id"`
		Quantity  int    `json:"amount"`
		Name      string `json:"name"`
		Address   string `json:"address" redact:"[REDACTED]"`
		PinNumber string `json:"pin_number" redact:"[REDACTED]"`
	}

	order := &SmallOrder{
		ID:        "a5cd5a4a-b7c1-4439-9ceb-7b4e6e6fbe7d",
		Quantity:  10,
		Name:      "Mr Bob",
		Address:   "Some Address",
		PinNumber: "1234",
	}

	expectedResult := &SmallOrder{
		ID:        "a5cd5a4a-b7c1-4439-9ceb-7b4e6e6fbe7d",
		Quantity:  10,
		Name:      "Mr Bob",
		Address:   "[REDACTED]",
		PinNumber: "[REDACTED]",
	}

	result, ok := Redact(order).(*SmallOrder)
	if !ok {
		uts.Fail("failed to cast response to Order")
	}

	uts.Equal(expectedResult, result)
}

func (uts *UnitTestSuite) TestWeReplaceNestedString() {
	type SmallOrder struct {
		ID       string `json:"id"`
		Quantity int    `json:"amount"`
		User     struct {
			Name      string
			Address   string `json:"address" redact:"[REDACTED]"`
			PinNumber string `json:"pin_number" redact:"[REDACTED]"`
		} `json:"user"`
	}

	order := &SmallOrder{
		ID:       "a5cd5a4a-b7c1-4439-9ceb-7b4e6e6fbe7d",
		Quantity: 10,
		User: struct {
			Name      string
			Address   string "json:\"address\" redact:\"[REDACTED]\""
			PinNumber string `json:"pin_number" redact:"[REDACTED]"`
		}{
			Name:      "Mr Bob",
			Address:   "Some Address",
			PinNumber: "1234",
		},
	}

	expectedResult := &SmallOrder{
		ID:       "a5cd5a4a-b7c1-4439-9ceb-7b4e6e6fbe7d",
		Quantity: 10,
		User: struct {
			Name      string
			Address   string "json:\"address\" redact:\"[REDACTED]\""
			PinNumber string `json:"pin_number" redact:"[REDACTED]"`
		}{
			Name:      "Mr Bob",
			Address:   "[REDACTED]",
			PinNumber: "[REDACTED]",
		},
	}

	result, ok := Redact(order).(*SmallOrder)
	if !ok {
		uts.Fail("failed to cast response to Order")
	}

	uts.Equal(expectedResult, result)
}

func (uts *UnitTestSuite) TestWeReplaceStringPointer() {
	type TestOrder struct {
		ID        string  `json:"id"`
		Quantity  int     `json:"amount"`
		Name      string  `json:"name"`
		Address   *string `json:"address" redact:"[REDACTED]"`
		PinNumber string  `json:"pin_number" redact:"[REDACTED]"`
	}

	address := "Some address"
	order := &TestOrder{
		ID:        "a5cd5a4a-b7c1-4439-9ceb-7b4e6e6fbe7d",
		Quantity:  10,
		Name:      "Mr Bob",
		Address:   &address,
		PinNumber: "1234",
	}

	expectedAddress := "[REDACTED]"
	expectedResult := &TestOrder{
		ID:        "a5cd5a4a-b7c1-4439-9ceb-7b4e6e6fbe7d",
		Quantity:  10,
		Name:      "Mr Bob",
		Address:   &expectedAddress,
		PinNumber: "[REDACTED]",
	}

	result, ok := Redact(order).(*TestOrder)
	if !ok {
		uts.Fail("failed to cast response to Order")
	}

	uts.Equal(expectedResult, result)
}
