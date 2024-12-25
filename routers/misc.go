package routers

import (
	"fmt"
	"myproject/models"
	"myproject/response"
	"net/http"

	"myproject/controller"

	"context"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Handle test submissions
func HandleSubmission(c *gin.Context) {
	var submission response.Submit

	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission format"})
		return
	}

	println("::: PMODE :::" + submission.PMode)

	// Find or create the user
	var existingUsers []models.User
	mgm.Coll(&models.User{}).SimpleFind(&existingUsers, bson.M{"email": submission.Email})

	if len(existingUsers) == 0 {
		newUser := models.NewUser(submission.Name, submission.Email, submission.Gender, submission.Age, "", "USER", "ACTIVE", "PENDING")
		if err := mgm.Coll(newUser).Create(newUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		existingUsers = append(existingUsers, *newUser)
	}

	user := existingUsers[0]

	// Create a new test entry

	var testPaymentStatus string = "PENDING"
	var testPaymentLink string = ""
	var paymentLinkId string = ""
	var testId primitive.ObjectID = primitive.NewObjectID()

	if submission.PMode != "" && submission.PMode == "pass" {
		// Just Generate New Report
		testPaymentStatus = "BYPASS_PAYMENT"
	} else {
		// Go through payment mode
		referenceId := "big5_" + testId.Hex()
		amount := os.Getenv("BIG_5_REPORT_PRICE")

		amountInt, err := strconv.Atoi(amount)
		if err != nil {
			fmt.Println(":: ERROR : " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generated payment link"})
			return
		}
		data, err := controller.GeneratePaymentLink(amountInt, "For BIG 5 report generator", submission.Name, user.Email, referenceId)

		if err != nil {
			fmt.Println(":: ERROR : " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generated payment link"})
			return
		}

		// Type assertion to get the string value from the map
		shortURL, ok := data["short_url"].(string)
		if !ok {
			fmt.Println(":: ERROR : short_url is not a string")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generated payment link"})
			return
		}

		id, ok := data["id"].(string)
		if !ok {
			fmt.Println(":: ERROR : id is not a string")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generated payment link"})
			return
		}

		testPaymentLink = shortURL
		paymentLinkId = id
	}

	newTest := models.NewTest(testId, submission.Name, submission.Age, submission.Gender, "BIG_5", user.ID, testPaymentStatus, testPaymentLink, paymentLinkId, "PENDING")

	if err := mgm.Coll(&models.Test{}).Create(newTest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create test"})
		return
	}

	// Store scores
	var scoreDocs []models.Score
	for _, answer := range submission.Answers {
		questionId, err := primitive.ObjectIDFromHex(answer.Id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
			return
		}
		scoreDocs = append(scoreDocs, *models.NewScore(user.ID, questionId, answer.Answer, newTest.ID))
	}

	var docs []interface{}
	for _, q := range scoreDocs {
		docs = append(docs, q) // Add each question as an interface{}
	}

	_, err := mgm.Coll(&models.Score{}).InsertMany(c, docs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store scores"})
		return
	}

	if submission.PMode != "" && submission.PMode == "pass" {
		// Just Generate New Report
		go controller.GenerateNewReport(c, *newTest, user)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Submission successful", "paymentLink": testPaymentLink})
}

// Generate report
func HandleReportGeneration(c *gin.Context) {
	var reportRequest response.Report

	if err := c.ShouldBindJSON(&reportRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report request"})
		return
	}

	testId, err := primitive.ObjectIDFromHex(reportRequest.TestId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test ID"})
		return
	}

	test := models.Test{}
	// Get test
	mgm.Coll(&models.Test{}).FindByID(testId, &test)

	if test == (models.Test{}) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get test for test id"})
		return
	}

	// Get user
	// TODO remove this once the middleware is implemented
	user := models.FetchUserUsingId(test.UserId)

	if user == (models.User{}) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user for test id"})
		return
	}

	// Checking if report already generated

	filter := bson.D{{Key: "testId", Value: testId}}

	count, err := mgm.Coll(&models.Report{}).CountDocuments(context.TODO(), filter)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	if count != 0 {
		c.IndentedJSON(http.StatusAlreadyReported, gin.H{"message": "report is already generated for test id"})
		return
	}

	errFromRequest := controller.GenerateNewReport(c, test, user)

	if errFromRequest != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errFromRequest.Message})
		return
	}

	_, err = models.UpdateTestReportSent(test.ID, "DONE")

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report generated successfully"})
}
