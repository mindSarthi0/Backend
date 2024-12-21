package models

import (
	mgm "github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FinalReport struct {
	// DefaultModel includes the MongoDB ID (_id), createdAt, and updatedAt fields.
	mgm.DefaultModel `bson:",inline"`

	UserId           primitive.ObjectID     `json:"userId" bson:"userId"`
	TestId           primitive.ObjectID     `json:"testId" bson:"testId"`
	GeneratedContent map[string]interface{} `json:"generatedContent" bson:"generatedContent"`
}

func NewFinalReport(userId primitive.ObjectID, testId primitive.ObjectID, generatedContent map[string]interface{}) *FinalReport {
	return &FinalReport{
		UserId:           userId,
		TestId:           testId,
		GeneratedContent: generatedContent,
	}
}
