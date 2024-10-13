package models

import (
	"github.com/kamva/mgm/v3"
)

// Question model with fields for MongoDB
type User struct {
	// DefaultModel includes the MongoDB ID (_id), createdAt, and updatedAt fields.
	mgm.DefaultModel `bson:",inline"`

	Name   string `json:"name" bson:"name"`   // Name of the test
	Email  string `json:"email" bson:"email"` // The actual question text
	Gender string `json:"gender" bson:"gender"`
	Age    int    `json:"age" bson:"age"`
}

// NewQuestion creates a new instance of the Question model
func NewUser(name string, email string, gender string, age int) *User {
	return &User{
		Name:   name,
		Email:  email,
		Gender: gender,
		Age:    age,
	}
}
