package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Question model with fields for MongoDB
type Test struct {
	// DefaultModel includes the MongoDB ID (_id), createdAt, and updatedAt fields.
	mgm.DefaultModel `bson:",inline"`

	TestName    string             `json:"testName" bson:"testName"` // Name of the test
	UserId      primitive.ObjectID `json:"userId" bson:"userId"`     // The actual question text
	Paid        string             `json:"paid" bson:"paid"`
	PaymentLink string             `json:"paymentLink" bson:"paymentLink"`
}

// NewQuestion creates a new instance of the Question model
func NewTest(testName string, userId primitive.ObjectID, paid string, paymentLink string) *Test {
	return &Test{
		TestName:    testName,
		UserId:      userId,
		Paid:        paid,
		PaymentLink: paymentLink,
	}
}
