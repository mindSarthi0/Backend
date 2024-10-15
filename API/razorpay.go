package API

import (
	razorpay "github.com/razorpay/razorpay-go"
	"sync"
)

// Declare a package-level variable for the client
var razorpayClient *razorpay.Client

// Declare a sync.Once to ensure client is initialized only once
var once sync.Once

// getClient returns the singleton Razorpay client
func GetClient() *razorpay.Client {
	once.Do(func() {
		// Initialize the Razorpay client only once
		razorpayClient = razorpay.NewClient("<YOUR_API_KEY>", "<YOUR_API_SECRET>")
	})
	// Return the initialized client (or the already initialized client)
	return razorpayClient
}

// Function to create the data object based on input parameters
func CreatePaymentLinkData(amount int, currency string, acceptPartial bool, minPartialAmount int, expireBy int64, referenceID string, description string, customerName string, customerContact string, customerEmail string, notifySMS bool, notifyEmail bool, reminderEnable bool, policyName string, callbackURL string, callbackMethod string) map[string]interface{} {
	// Create the data map based on input parameters
	data := map[string]interface{}{
		"amount":                   amount,
		"currency":                 currency,
		"accept_partial":           acceptPartial,
		"first_min_partial_amount": minPartialAmount,
		"expire_by":                expireBy, // Can be set dynamically or passed in
		"reference_id":             referenceID,
		"description":              description,
		"customer": map[string]interface{}{
			"name":    customerName,
			"contact": customerContact,
			"email":   customerEmail,
		},
		"notify": map[string]interface{}{
			"sms":   notifySMS,
			"email": notifyEmail,
		},
		"reminder_enable": reminderEnable,
		"notes": map[string]interface{}{
			"policy_name": policyName,
		},
		"callback_url":    callbackURL,
		"callback_method": callbackMethod,
	}

	return data
}
