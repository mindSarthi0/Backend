package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"myproject/models"
	"myproject/response"
	"net/http"
)

// TODO A single test should have one Id
func postSubmit(c *gin.Context) {

	var submit response.Submit

	if err := c.ShouldBindJSON(&submit); err != nil {
		// If there is an error, respond with 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Received person data: %+v\n", submit.Answers)
	// Find the user
	isExistingUsers := []models.User{}

	fmt.Printf("Received person data: %+v\n", isExistingUsers)
	mgm.Coll(&models.User{}).SimpleFind(&isExistingUsers, bson.M{"email": submit.Email})

	if isExistingUsers == nil || len(isExistingUsers) == 0 {
		user := models.NewUser(submit.Name, submit.Email, submit.Gender, submit.Age)

		err := mgm.Coll(user).Create(user)

		if err != nil {
			c.IndentedJSON(http.StatusOK, err)
			log.Fatalf("Failed to create a new err: %v", err) // Fatal will log and stop the program
		}
		isExistingUsers = append(isExistingUsers, *user)

	}

	userId := isExistingUsers[0].ID

	// TODO Generate payment link from razorPay

	// Create a test entry
	newTest := models.NewTest("BIG_5", userId, "PENDING", "https://google.com")

	newTestError := mgm.Coll(&models.Test{}).Create(newTest)

	if newTestError != nil {
		c.IndentedJSON(http.StatusOK, newTestError)
		log.Fatalf("Failed to create a new test: %v", newTestError) // Fatal will log and stop the program
	}

	// Store Score
	coll := mgm.Coll(&models.Score{})

	scoreDocs := []models.Score{}

	for _, item := range submit.Answers {
		questionId, err := primitive.ObjectIDFromHex(item.Id)
		if err != nil {
			log.Fatalf("Failed to convert string to ObjectID: %v", err)
		}
		score := models.NewScore(userId, questionId, item.Answer, newTest.ID)
		scoreDocs = append(scoreDocs, *score)
	}

	var docs []interface{}
	for _, score := range scoreDocs {
		docs = append(docs, score)
	}

	_, err := coll.InsertMany(c, docs)

	if err != nil {
		log.Fatalf("Failed to insert multiple documents: %v", err)
	}

	// Create a user

	log.Printf("Summited sucessfully")
	c.IndentedJSON(http.StatusOK, err)
}

func init() {
	// Setup the mgm default config
	err := mgm.SetDefaultConfig(nil, "cognify", options.Client().ApplyURI("mongodb+srv://cognify:dEQGVwIY24QzdUu6@cluster0.cjyqt.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	// Error handling
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err) // Fatal will log and stop the program
	}

	fmt.Println("Successfully connected to MongoDB!")
}

func postQuestions(c *gin.Context) {
	var question []response.Question

	if err := c.ShouldBindJSON(&question); err != nil {
		// If there is an error, respond with 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coll := mgm.Coll(&models.Question{})

	questions := []models.Question{}

	for _, item := range question {

		questionToSave := models.NewQuestion(item.TestName, item.Question)
		questions = append(questions, *questionToSave)
	}

	var docs []interface{}
	for _, question := range questions {
		docs = append(docs, question)
	}

	_, err := coll.InsertMany(context.TODO(), docs)

	if err != nil {
		log.Printf("Failed to insert multiple question documents: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert documents"})
		return
	}

	log.Printf("Created sucessfully")
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Questions created successfully"})

}

func getQuestions(c *gin.Context) {
	// Get the MongoDB collection for questions
	coll := mgm.Coll(&models.Question{})

	// Create a slice to hold the retrieved questions
	var questions []models.Question

	// Query all documents in the collection
	err := coll.SimpleFind(&questions, bson.D{})
	if err != nil {
		log.Printf("Error retrieving questions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve questions"})
		return
	}

	// Check if no documents are found
	if len(questions) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No questions found"})
		return
	}

	// Return the found questions
	c.IndentedJSON(http.StatusOK, questions)
}

func main() {
	router := gin.Default()

	router.POST("/questions", postQuestions)
	router.GET("/questions", getQuestions) // For retrieving questions
	router.POST("/submit", postSubmit)

	// Start the server on localhost:8080
	router.Run("localhost:8080")
}
