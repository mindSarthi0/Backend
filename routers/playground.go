package routers

import (
	"github.com/gin-gonic/gin"
	"myproject/API"
	"net/http"
	"os"
	"time"
)

func PaymentTest(c *gin.Context) {

	backendAapiDomain := os.Getenv("BACKEND_API_DOMAIN")
	callbackPath := os.Getenv("CALLBACK_PATH")

	amount := 2100
	currency := "INR"
	acceptPartial := false
	minPartialAmount := 0
	expireBy := time.Now().AddDate(0, 0, 7).Unix() // Expire in 7 days
	referenceID := "ref127"
	description := "Payment for services"
	customerName := "John Doe"
	customerContact := "1234567890"
	customerEmail := "john.doe@example.com"
	notifySMS := true
	notifyEmail := true
	reminderEnable := true
	policyName := "Standard Policy"
	callbackURL := backendAapiDomain + callbackPath
	callbackMethod := "get"
	upiLink := false

	println(callbackURL)

	data, _ := API.CreatePaymentLinkData(upiLink, amount, currency, acceptPartial, minPartialAmount, expireBy, referenceID, description, customerName, customerContact, customerEmail, notifySMS, notifyEmail, reminderEnable, policyName, callbackURL, callbackMethod)

	c.JSON(http.StatusOK, gin.H{"message": "Payment Link created successfully", "data": data})
}

type PaymentLinkFetchRequest struct {
	PaymentLinkId string `json:"paymentLinkId"`
}

func PaymentLinkFetch(c *gin.Context) {

	var paymentLinkFetchRequest PaymentLinkFetchRequest

	if err := c.ShouldBindJSON(&paymentLinkFetchRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Payment Link Fetch Request data"})
		return
	}

	data, _ := API.GetPaymentLink(paymentLinkFetchRequest.PaymentLinkId)

	c.JSON(http.StatusOK, gin.H{"message": "Payment Link fetched successfully", "data": data})
}
