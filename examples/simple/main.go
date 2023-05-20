package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/squeakycheese75/redact"
)

type Guest struct {
	Name    string `json:"name"`
	Species string `json:"species" redact:"[REDACTED]"`
}

type Party struct {
	GuestList []Guest `json:"guest_list"`
	Date      string  `json:"date"`
	Venue     string  `json:"venue" redact:"[REDACTED]"`
}

func main() {
	log.SetFormatter(&log.JSONFormatter{PrettyPrint: true})

	barry := Guest{Name: "barry", Species: "seal"}
	patricia := Guest{Name: "patricia", Species: "penguin"}
	morris := Guest{Name: "morris", Species: "mackeral"}

	myParty := &Party{
		GuestList: []Guest{barry, patricia, morris},
		Date:      "Tuesday",
		Venue:     "Barrys Place",
	}

	log.WithFields(log.Fields{
		"animal": "walrus",
		"party":  myParty,
	}).Info("A walrus is having a party")

	log.WithFields(log.Fields{
		"animal": "walrus",
		"party":  redact.Redact(myParty),
	}).Info("A walrus is having a [REDACTED] party")
}
