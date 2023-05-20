# Redact

## Why

We developed Redact as a faster and more resource-efficient alternative to [Redactrus](https://github.com/whuang8/redactrus). While Redactrus is a great project that offers considerable flexibility, we encountered challenges with its use of regular expressions. Specifically, when dealing with our log volume, Redactrus consumed significant computational resources. In contrast, Redact utilizes Go's struct tagging feature, enabling the tagging of specific fields for value redaction.

## Usage

To use Redact, with logrus

```go
package main

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/squeakycheese75/redact"
)

type Guest struct {
	Name    string `json:"name"`
	Species string `json:"species" redact:"[REDACTED]"`
}

type Party struct {
	GuestList []*Guest `json:"guest_list"`
	Date      string   `json:"date"`
	Venue     string   `json:"venue" redact:"[REDACTED]"`
}

func main() {
	log.SetFormatter(&log.JSONFormatter{PrettyPrint: true})

	barry := &Guest{Name: "barry", Species: "seal"}
	patricia := &Guest{Name: "patricia", Species: "penguin"}
	morris := &Guest{Name: "morris", Species: "mackeral"}

	myParty := &Party{
		GuestList: []*Guest{barry, patricia, morris},
		Date:      "Tuesday",
		Venue:     "Barrys Place",
	}

	payloadBytes, err := json.Marshal(myParty)
	if err != nil {
		panic("error marshalling Party")
	}

    log.WithFields(log.Fields{
		"animal": "walrus",
		"party":  redact.Redact(myParty),
	}).Info("A walrus is having a [REDACTED] party")
}
```

results...

```bash
{
  "animal": "walrus",
  "level": "info",
  "msg": "A walrus is having a party",
  "party": {
    "guest_list": [
      {
        "name": "barry",
        "species": "seal"
      },
      {
        "name": "patricia",
        "species": "penguin"
      },
      {
        "name": "morris",
        "species": "mackeral"
      }
    ],
    "date": "Tuesday",
    "venue": "[REDACTED]"
  },
  "time": "2023-04-22T12:49:46+02:00"
}

```
