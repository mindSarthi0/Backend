package models

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subdomain struct {
	Name      string `json:"name" bson:"name"`
	Score     int    `json:"score" bson:"score"`
	Intensity string `json:"intensity" bson:"intensity"`
}

type Report struct {
	// DefaultModel includes the MongoDB ID (_id), createdAt, and updatedAt fields.
	mgm.DefaultModel `bson:",inline"`

	Name      string             `json:"name" bson:"name"`
	Score     int                `json:"score" bson:"score"`
	Subdomain []Subdomain        `json:"subdomain" bson:"subdomain"`
	UserId    primitive.ObjectID `json:"userId" bson:"userId"`
	TestId    primitive.ObjectID `json:"testId" bson:"testId"`
}

func NewReport(name string, score int, subdomain []Subdomain, userId primitive.ObjectID, testId primitive.ObjectID) *Report {
	return &Report{
		Name:      name,
		Score:     score,
		Subdomain: subdomain,
		TestId:    testId,
		UserId:    userId,
	}
}

func NewSubdomain(name string, score int, intensity string) *Subdomain {
	return &Subdomain{
		Name:      name,
		Score:     score,
		Intensity: intensity,
	}
}
