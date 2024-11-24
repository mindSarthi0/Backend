package controller

import (
	"fmt"
	"log"
	"myproject/lib"
	"myproject/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"

	// "strconv"
	"myproject/API"
	"myproject/constants"
	"os"
	"time"
)

type MyError struct {
	Code    int
	Message string
}

// var OutputPageMap = []string{"result", "relationsh`ip", "career_academic", "strength_weakness"}

func GenerateNewReport(c *gin.Context, test models.Test, user models.User) *MyError {
	startTime := time.Now()

	// Fetch Scores and Questions
	scoresAndQuestions, err := FetchScoresWithQuestions(test.ID)
	if err != nil {
		return &MyError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch scores and questions",
		}
	}
	fmt.Println("Time taken by Fetch Scores with Questions API calls:", time.Since(startTime))

	// Process Scores
	processedScores := CalculateProcessedScore(scoresAndQuestions)

	newDbReports := map[string]models.Report{}
	prompts := map[string]string{}

	// Generate Reports for Each Domain
	for _, value := range processedScores {
		newSubdomainReports := []models.Subdomain{}
		_domainName := value.Name
		_domainIntensity := value.Intensity

		// Create Subdomain Reports
		for _, subdomain := range value.Subdomain {
			newDbSubdomain := models.NewSubdomain(subdomain.Name, subdomain.Score, subdomain.Intensity)
			newSubdomainReports = append(newSubdomainReports, *newDbSubdomain)
		}

		// Create Report for Domain
		newDbReport := models.NewReport(value.Name, value.Score, newSubdomainReports, value.UserId, value.TestId, _domainIntensity, "")
		newDbReports[_domainName] = *newDbReport
	}

	// Generate Prompts for AI Model
	for page := range constants.BIG_5_Report {
		prompt := API.CreatePrompt(constants.BIG_5_Report[page], processedScores)
		prompts[constants.BIG_5_Report[page]] = prompt
	}

	// Concurrent API Calls for AI Responses
	startTime = time.Now()
	channel := make(chan API.PromptRequest)
	for page, prompt := range prompts {
		go API.WorkerOpenAIGPT(page, prompt, channel)
	}

	results := map[string]string{}
	for range prompts {
		result := <-channel
		results[result.Id] = result.Response.Choices[0].Message.Content
	}
	fmt.Println("Time taken by GCP Worker to generate response from gemini:", time.Since(startTime))

	// Generate PDF Content
	pdfGenerationContent := map[string]string{}
	for page, _ := range prompts {
		generatedResponseString := results[page]
		generatedContent, err := lib.ParseMarkdownCode(generatedResponseString)
		if err != nil {
			// Respond with an error message if content generation failed
		}
		fmt.Println("page:", page, generatedContent, generatedResponseString)
		pdfGenerationContent[page] = generatedResponseString
	}

	log.Println("Generating PDF")
	startTime = time.Now()
	reportPdfFilename := "report_" + test.ID.Hex()
	log.Println("Tester Name: " + test.TestGiver)

	// Generate PDF
	errInPdfGeneration := API.GenerateBigFivePDF(pdfGenerationContent, test.TestGiver, reportPdfFilename)
	if errInPdfGeneration != nil {
		return &MyError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate pdf",
		}
	}
	fmt.Println("Time taken to generate PDF:", time.Since(startTime))

	// Send Report via Email
	startTime = time.Now()
	log.Println("Sending Report via Email to user")
	API.SendBIG5Report(user.Email, test.TestGiver, "./"+reportPdfFilename+".pdf")
	fmt.Println("Time taken to send email:", time.Since(startTime))

	// Save Report to Database
	startTime = time.Now()
	var docs []interface{}
	for _, q := range newDbReports {
		docs = append(docs, q)
	}
	_, err = mgm.Coll(&models.Report{}).InsertMany(c, docs)
	fmt.Println("Time taken to save Report in db:", time.Since(startTime))
	if err != nil {
		return &MyError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to insert questions",
		}
	}

	return nil
}

func GeneratePaymentLink(
	amount int,
	description string,
	name string,
	email string,
	referenceID string,
) (map[string]interface{}, error) {

	backendAapiDomain := os.Getenv("BACKEND_API_DOMAIN")
	callbackPath := os.Getenv("CALLBACK_PATH")

	currency := "INR"
	acceptPartial := false
	minPartialAmount := 0
	expireBy := time.Now().AddDate(0, 0, 7).Unix() // Expire in 7 days
	customerName := name
	customerContact := ""
	customerEmail := email
	notifySMS := true
	notifyEmail := true
	reminderEnable := true
	policyName := "Standard Policy"
	callbackURL := backendAapiDomain + callbackPath
	callbackMethod := "get"
	upiLink := false

	return API.CreatePaymentLinkData(upiLink, amount, currency, acceptPartial, minPartialAmount, expireBy, referenceID, description, customerName, customerContact, customerEmail, notifySMS, notifyEmail, reminderEnable, policyName, callbackURL, callbackMethod)

}
