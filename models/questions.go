package models

import (
	"github.com/kamva/mgm/v3"
)

// Question model with fields for MongoDB
type Question struct {
	// DefaultModel includes the MongoDB ID (_id), createdAt, and updatedAt fields.
	mgm.DefaultModel `bson:",inline"`

	TestName string `json:"testName" bson:"testName"` // Name of the test
	Question string `json:"question" bson:"question"` // The actual question text
}

// NewQuestion creates a new instance of the Question model
func NewQuestion(testName string, question string) *Question {
	return &Question{
		TestName: testName,
		Question: question,
	}
}
