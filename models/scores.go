package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Question model with fields for MongoDB
type Score struct {
	// DefaultModel includes the MongoDB ID (_id), createdAt, and updatedAt fields.
	mgm.DefaultModel `bson:",inline"`

	UserId     primitive.ObjectID `json:"userId" bson:"userId"` // Name of the test
	TestId     primitive.ObjectID `json:"testId" bson:"testId"`
	QuestionId primitive.ObjectID `json:"questionId" bson:"questionId"` // The actual question text
	RawScore   string             `json:"rawScore" bson:"rawScore"`
}

func NewScore(userId primitive.ObjectID, questionId primitive.ObjectID, rawScore string, testId primitive.ObjectID) *Score {
	return &Score{
		UserId:     userId,
		QuestionId: questionId,
		RawScore:   rawScore,
		TestId:     testId,
	}
}
