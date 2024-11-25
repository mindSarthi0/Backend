package API

import (
	"log"
	"os"
	"sync"

	razorpay "github.com/razorpay/razorpay-go"
	utils "github.com/razorpay/razorpay-go/utils"
)

// Declare a package-level variable for the client
var razorpayClient *razorpay.Client

// Declare a sync.Once to ensure client is initialized only once
var once sync.Once

// getClient returns the singleton Razorpay client
func GetClient() *razorpay.Client {

	razorpayApiKey := os.Getenv("RAZORPAY_API_KEY")
	razorpayApiSecret := os.Getenv("RAZORPAY_API_SECRET")

	once.Do(func() {
		// Initialize the Razorpay client only once
		razorpayClient = razorpay.NewClient(razorpayApiKey, razorpayApiSecret)
	})
	// Return the initialized client (or the already initialized client)
	return razorpayClient
}

// Function to create the data object based on input parameters
func CreatePaymentLinkData(upiLink bool, amount int, currency string, acceptPartial bool, minPartialAmount int, expireBy int64, referenceID string, description string, customerName string, customerContact string, customerEmail string, notifySMS bool, notifyEmail bool, reminderEnable bool, policyName string, callbackURL string, callbackMethod string) (map[string]interface{}, error) {

	client := GetClient()
	// Create the data map based on input parameters
	data := map[string]interface{}{
		"upi_link":                 upiLink,
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

	body, err := client.PaymentLink.Create(data, nil)

	if err != nil && err.Error() != "" {
		log.Printf("::Razorpay Error : %v", err.Error())
	}

	return body, err
}

// {
// 	"data": {
// 	  "accept_partial": false,
// 	  "amount": 21000,
// 	  "amount_paid": 0,
// 	  "callback_method": "get",
// 	  "callback_url": "https://cognify-api-gateway-976411241646.asia-south1.run.app/paymentStatus",
// 	  "cancelled_at": 0,
// 	  "created_at": 1731574341,
// 	  "currency": "INR",
// 	  "customer": {
// 		"contact": "1234567890",
// 		"email": "john.doe@example.com",
// 		"name": "John Doe"
// 	  },
// 	  "description": "Payment for services",
// 	  "expire_by": 1732179140,
// 	  "expired_at": 0,
// 	  "first_min_partial_amount": 0,
// 	  "id": "plink_PL7qXGWIlOHnlL", --- Payment Link Id
// 	  "notes": {
// 		"policy_name": "Standard Policy"
// 	  },
// 	  "notify": {
// 		"email": true,
// 		"sms": true,
// 		"whatsapp": false
// 	  },
// 	  "payments": null,
// 	  "reference_id": "ref126",
// 	  "reminder_enable": true,
// 	  "reminders": [],
// 	  "short_url": "https://rzp.io/rzp/IlGei0sw",
// 	  "status": "created",
// 	  "updated_at": 1731574341,
// 	  "upi_link": false,
// 	  "user_id": "",
// 	  "whatsapp_link": false
// 	},
// 	"message": "Payment Link created successfully"
//   }

// https://cognify-api-gateway-976411241646.asia-south1.run.app/paymentStatus?razorpay_payment_id=pay_PLBEF4VzmRD992&razorpay_payment_link_id=plink_PL7qXGWIlOHnlL&razorpay_payment_link_reference_id=ref126&razorpay_payment_link_status=paid&razorpay_signature=7d2588410ae3faa83a39fb8dda0a784bcd1204ebfad038c49d90cd772303fd27

// TODO Check payment status based on reference id

func GetPaymentLink(paymentLinkId string) (map[string]interface{}, error) {

	client := GetClient()

	body, err := client.PaymentLink.Fetch(paymentLinkId, nil, nil)

	if err != nil && err.Error() != "" {
		log.Printf("::Razorpay Error : %v", err.Error())
	}

	return body, err

}

func VerifyPaymentLink(params map[string]interface{}, signature string) bool {

	// params := map[string]interface{} {
	// 	"payment_link_id": "plink_IH3cNucfVEgV68",
	// 	"razorpay_payment_id": "pay_IH3d0ara9bSsjQ",
	// 	"payment_link_reference_id": "TSsd1989",
	// 	"payment_link_status": "paid",
	// }

	razorpayApiSecret := os.Getenv("RAZORPAY_API_SECRET")

	isVerified := utils.VerifyPaymentLinkSignature(params, signature, razorpayApiSecret)

	return isVerified
}
