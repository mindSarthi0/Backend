package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Question model with fields for MongoDB
type Test struct {
	// DefaultModel includes the MongoDB ID (_id), createdAt, and updatedAt fields.
	mgm.DefaultModel `bson:",inline"`

	ID                primitive.ObjectID `json:"_id" bson:"_id"`
	TestGiver         string             `json:"testGiven" bson:"testGiven"`
	TestGiverAge      int                `json:"testGiverAge" bson:"testGiverAge"`
	TestGiverGender   string             `json:"testGiverGender" bson:"testGiverGender"`
	TestName          string             `json:"testName" bson:"testName"` // Name of the test
	UserId            primitive.ObjectID `json:"userId" bson:"userId"`     // The actual question text
	PaymentStatus     string             `json:"paymentStatus" bson:"paymentStatus"`
	PaymentLink       string             `json:"paymentLink" bson:"paymentLink"`
	ExternalPaymentId string             `json:"externalPaymentId" bson:"externalPaymentId"`
	ReportSent        string             `json:"reportSent" bson:"reportSent"`
}

// NewQuestion creates a new instance of the Question model
func NewTest(
	testId primitive.ObjectID,
	testGiver string,
	testGiverAge int,
	testGiverGender string,
	testName string,
	userId primitive.ObjectID,
	paymentStatus string,
	paymentLink string,
	externalPaymentId string,
	reportSent string,
) *Test {
	return &Test{
		ID:                testId,
		TestGiverAge:      testGiverAge,
		TestGiver:         testGiver,
		TestGiverGender:   testGiverGender,
		TestName:          testName,
		UserId:            userId,
		PaymentStatus:     paymentStatus,
		PaymentLink:       paymentLink,
		ExternalPaymentId: externalPaymentId,
		ReportSent:        reportSent,
	}
}

func FetchTestById(testId primitive.ObjectID) (*Test, error) {
	var test Test

	err := mgm.Coll(&Test{}).First(bson.M{"_id": testId}, &test)
	if err != nil {
		// If no document is found, return a custom error
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("test with ID %s not found", testId.Hex())
		}
		// Handle other potential errors
		return nil, err
	}

	return &test, nil
}

func UpdateTestPaymentStatus(testId primitive.ObjectID, status string) (*Test, error) {
	var test Test

	// Define the update document
	update := bson.M{
		"$set": bson.M{
			"paymentStatus": status,
		},
	}

	// Find one document and update it, returning the updated document
	err := mgm.Coll(&Test{}).FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": testId},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&test)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no document found with the given ID")
		}
		return nil, err
	}

	return &test, nil
}

func UpdateTestReportSent(testId primitive.ObjectID, reportSent string) (*Test, error) {
	var test Test

	// Define the update document
	update := bson.M{
		"$set": bson.M{
			"reportSent": reportSent,
		},
	}

	// Find one document and update it, returning the updated document
	err := mgm.Coll(&Test{}).FindOneAndUpdate(
		context.TODO(),
		bson.M{"_id": testId},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&test)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no document found with the given ID")
		}
		return nil, err
	}

	return &test, nil
}
