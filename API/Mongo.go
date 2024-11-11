package API

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Subdomain struct {
	Name      string
	Score     int
	Intensity string
}

type Domain struct {
	Name      string
	Score     int
	Subdomain []Subdomain
	UserId    primitive.ObjectID `json:"userId" bson:"userId"`
	TestId    primitive.ObjectID `json:"testId" bson:"testId"`
	Intensity string
}

func fetchPersonalityData(client *mongo.Client, userID string) ([]Domain, error) {
	collection := client.Database("personality_db").Collection("assessments")
	filter := bson.M{"user_id": userID}

	var result []Domain
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("error fetching personality data: %v", err)
	}

	return result, nil
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://cognify:dEQGVwIY24QzdUu6@cluster0.cjyqt.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	userID := "12345"
	data, err := fetchPersonalityData(client, userID)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}

	// Process the fetched data for report generation
	fmt.Printf("Personality Data: %+v\n", data)
}
