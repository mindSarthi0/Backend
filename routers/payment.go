package routers

import (
	"fmt"
	"log"
	apis "myproject/apis"
	"myproject/controller"
	"myproject/models"
	"net/http"
	"os"
	"strings"

	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	webappDomain := os.Getenv("WEBAPP_DOMAIN")
	webappPaymentStatusPath := webappDomain + os.Getenv("WEBAPP_PAYMENT_STATUS_PATH")

	// Get all the query parameters
	queryParams := c.Request.URL.Query()

	// Print all the query parameters
	for key, values := range queryParams {
		for _, value := range values {
			fmt.Printf("Key: %s, Value: %s\n", key, value)
		}
	}

	// Extract required query parameters
	referenceId := queryParams.Get("razorpay_payment_link_reference_id")
	paymentLinkStatus := queryParams.Get("razorpay_payment_link_status")
	paymentLinkId := queryParams.Get("razorpay_payment_link_id")
	razorpayPaymentId := queryParams.Get("razorpay_payment_id")
	signature := queryParams.Get("razorpay_signature")

	// Check if all required parameters are present
	if referenceId == "" || paymentLinkStatus == "" || paymentLinkId == "" || razorpayPaymentId == "" || signature == "" {
		c.Redirect(http.StatusFound, webappPaymentStatusPath+"?status=failed")
		return
	}

	params := map[string]interface{}{
		"payment_link_id":           paymentLinkId,
		"razorpay_payment_id":       razorpayPaymentId,
		"payment_link_reference_id": referenceId,
		"payment_link_status":       paymentLinkStatus,
	}

	// Verify the payment link
	if apis.VerifyPaymentLink(params, signature) {
		// Payment verification successful
		// Mark payment status of test as successful

		testIdString, err := ExtractTestId(referenceId)
		if err != nil {
			log.Println(":: Error : " + err.Error())
			c.Redirect(http.StatusFound, webappPaymentStatusPath+"?status=pending&message=Your payment is being processed")
			return
		}

		log.Println(testIdString)
		testId, err := primitive.ObjectIDFromHex(testIdString)
		if err != nil {
			log.Println(":: Error : " + err.Error())
			c.Redirect(http.StatusFound, webappPaymentStatusPath+"?status=pending&message=Your payment is being processed")
			return
		}

		test, err := models.FetchTestById(testId)
		if err != nil {
			log.Println(":: Error : " + err.Error())
			c.Redirect(http.StatusFound, webappPaymentStatusPath+"?status=pending&message=Your payment is being processed")
			return
		}

		if test.PaymentStatus == "paid" {
			c.Redirect(http.StatusFound, webappPaymentStatusPath+"?status=success&message=Already marked paid")
			return
		}

		updatedTest, err := models.UpdateTestPaymentStatus(testId, paymentLinkStatus)
		if err != nil {
			log.Println(":: Error : " + err.Error())
			c.Redirect(http.StatusFound, webappPaymentStatusPath+"?status=pending&message=Your payment is being processed")
			return
		}

		log.Println(updatedTest.TestGiver)
		user := models.FetchUserUsingId(updatedTest.UserId)

		// Generate report if the payment status is "paid"
		if paymentLinkStatus == "paid" {
			go controller.GenerateNewReport(c, *test, user)
		}
		link := os.Getenv("WEBAPP_DOMAIN") + os.Getenv("REPORT_PATH") + test.ID.Hex()
		time.Sleep(2 * time.Second)
		c.Redirect(http.StatusFound, link)
		// c.Redirect(http.StatusFound, webappPaymentStatusPath+"?status=success&message=Thank you for your purchase. Your response is being analyzed by our scientific algorithm and will be sent to you within 5 minutes. We appreciate your interest in understanding yourself better!&link="+link)
		return
	}

	// In case of failed signature
	c.Redirect(http.StatusFound, webappPaymentStatusPath+"?status=failed&message=Sorry! Please try again")
}
