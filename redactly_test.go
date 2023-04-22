package redactly

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserName  string `json:"userName" redactly:"[REDACTED]"`
	Email     string `json:"email" redactly:"[REDACTED]"`
}

type Billing struct {
	Street  string `json:"street,omitempty" redactly:"[REDACTED]"`
	City    string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
	Zip     string `json:"zip,omitempty"`
}

type Order struct {
	ID        string  `json:"id"`
	Quantity  int     `json:"amount"`
	User      User    `json:"user"`
	Billing   Billing `json:"billing"`
	PinNumber string  `json:"pin_number" redactly:"[REDACTED]"`
}

type UnitTestSuite struct {
	suite.Suite
	order *Order
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &UnitTestSuite{})
}

func (uts *UnitTestSuite) SetupSuite() {
	order := &Order{
		ID:       "a5cd5a4a-b7c1-4439-9ceb-7b4e6e6fbe7d",
		Quantity: 10,
		User: User{
			FirstName: "Jimmy",
			LastName:  "Sticks",
			UserName:  "JimmySticks1234",
			Email:     "jimmysticks@fatmail.com",
		},
		Billing: Billing{
			Street:  "22 Sputnik Park Road",
			City:    "Paignton",
			Country: "England",
			Zip:     "TQ3 3DT",
		},
		PinNumber: "1234",
	}
	uts.order = order
}

func (uts *UnitTestSuite) TestSomething() {

	expectedResult := &Order{
		ID:       "a5cd5a4a-b7c1-4439-9ceb-7b4e6e6fbe7d",
		Quantity: 10,
		User: User{
			FirstName: "Jimmy",
			LastName:  "Sticks",
			UserName:  "[REDACTED]",
			Email:     "[REDACTED]",
		},
		Billing: Billing{
			Street:  "[REDACTED]",
			City:    "Paignton",
			Country: "England",
			Zip:     "TQ3 3DT",
		},
		PinNumber: "[REDACTED]",
	}

	payloadBytes, err := json.Marshal(uts.order)
	if err != nil {
		uts.FailNowf("unable to marshal order", err.Error())
	}

	result, ok := Unmarshal(payloadBytes, uts.order).(*Order)
	if !ok {
		uts.FailNowf("failed to cast response to Order", err.Error())
	}

	uts.Equal(expectedResult, result)
}
