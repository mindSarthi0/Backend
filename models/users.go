package models

import (
	"context"
	"fmt"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Question model with fields for MongoDB
type User struct {
	// DefaultModel includes the MongoDB ID (_id), createdAt, and updatedAt fields.
	mgm.DefaultModel `bson:",inline"`

	Name             string `json:"name" bson:"name"`   // Name of the test
	Email            string `json:"email" bson:"email"` // The actual question text
	Gender           string `json:"gender" bson:"gender"`
	Age              int    `json:"age" bson:"age"`
	Password         string `json:"password" bson:"password"`
	Role             string `json:"role" bson:"role"`
	Status           string `json:"status" bson:"status"`
	OnBoardingStatus string `json:"onboarding_status" bson:"onboarding_status"`
}

// NewQuestion creates a new instance of the Question model
func NewUser(name string, email string, gender string, age int, password string, role string, status string, onboardingStatus string) *User {
	return &User{
		Name:             name,
		Email:            email,
		Gender:           gender,
		Age:              age,
		Password:         password,
		Role:             role,
		Status:           status,
		OnBoardingStatus: onboardingStatus,
	}
}

func UpdateUserOnBoardingStatus(userId primitive.ObjectID, onBoardingStatus string) (*User, error) {
	var user User

	// Define the update document
	update := bson.M{
		"$set": bson.M{
			"OnBoardingStatus": onBoardingStatus,
		},
	}

	// Find one document and update it, returning the updated document
	err := mgm.Coll(&User{}).FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": userId},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no document found with the given ID")
		}
		return nil, err
	}

	return &user, nil
}

func FetchUserUsingId(id primitive.ObjectID) User {
	collRef := mgm.Coll(&User{})

	var user User

	collRef.FindByID(id.Hex(), &user)

	return user

}
