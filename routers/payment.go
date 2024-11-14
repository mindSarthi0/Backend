package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"myproject/API"
	"myproject/controller"
	"myproject/models"
	"net/http"
	"strings"
	// "os"
	// "time"
)

func ExtractTestId(referenceId string) (string, error) {
	// Split the input string by the underscore delimiter
	parts := strings.Split(referenceId, "_")
	if len(parts) != 2 {
		return "", fmt.Errorf("input string is not in the expected format")
	}

	// Return the second part of the split string
	return parts[1], nil
}

func HandlePaymentCallback(c *gin.Context) {
	// print all the query param, handle error as well

	// Get all the query parameters
	queryParams := c.Request.URL.Query()

	// Print all the query parameters
	for key, values := range queryParams {
		for _, value := range values {
			fmt.Printf("Key: %s, Value: %s\n", key, value)
		}
	}

	referenceId := queryParams["razorpay_payment_link_reference_id"][0]
	paymentLintStatus := queryParams["razorpay_payment_link_status"][0]

	params := map[string]interface{}{
		"payment_link_id":           queryParams["razorpay_payment_link_id"][0],
		"razorpay_payment_id":       queryParams["razorpay_payment_id"][0],
		"payment_link_reference_id": referenceId,
		"payment_link_status":       paymentLintStatus,
	}

	signature := queryParams["razorpay_signature"][0]
	if API.VerifyPaymentLink(params, signature) {

		// Payment verification successfull
		// Mark payment status of test to successfull

		testIdString, err := ExtractTestId(referenceId)

		if err != nil {
			log.Println(":: Error : " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to parse reference Id",
			})
			return
		}

		log.Println(testIdString)
		testId, err := primitive.ObjectIDFromHex(testIdString)

		if err != nil {
			log.Println(":: Error : " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to convert the given test in hex",
			})
			return
		}

		var test *models.Test
		var user models.User

		test, err = models.UpdateTestPaymentStatus(testId, paymentLintStatus)

		if err != nil {
			log.Println(":: Error : " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch test with id " + testIdString,
			})
			return
		}

		log.Println(test.TestGiver)
		user = models.FetchUserUsingId(test.UserId)
		// Now generate report

		if paymentLintStatus == "paid" {
			go controller.GenerateNewReport(c, *test, user)
		}

		// Fetch user
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
		return
	}
	// Handle any specific errors or respond with a status
	// For demonstration, we'll just check if a specific required parameter is present
	// requiredParam := "transaction_id"
	// if _, ok := queryParams[requiredParam]; !ok {
	// 	// If the required parameter is missing, respond with an error
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": fmt.Sprintf("Missing required query parameter: %s", requiredParam),
	// 	})
	// 	return
	// }

	// If everything is fine, respond with a success message
	c.JSON(http.StatusOK, gin.H{
		"status": "failed",
	})
}
