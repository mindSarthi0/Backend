package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"log"
	"myproject/models"
	"net/http"
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

// var OutputPageMap = []string{"result", "relationship", "career_academic", "strength_weakness"}

func GenerateNewReport(c *gin.Context, test models.Test, user models.User) *MyError {

	startTime := time.Millisecond
	scoresAndQuestions, err := FetchScoresWithQuestions(test.ID)

	if err != nil {
		return &MyError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch scores and questions",
		}
	}

	fmt.Println("Time taken by Fetch Scores with Questions API calls", time.Millisecond-startTime)
	processedScores := CalculateProcessedScore(scoresAndQuestions)

	newDbReports := map[string]models.Report{}

	prompts := map[string]string{}

	for _, value := range processedScores {

		// fmt.Println("Domain and Value:", domain, value)
		newSubdomainReports := []models.Subdomain{}

		_domainName := value.Name
		// _domainScore := value.Score
		_domainIntensity := value.Intensity

		for _, value := range value.Subdomain {
			newDbSubdomain := models.NewSubdomain(value.Name, value.Score, value.Intensity)
			newSubdomainReports = append(newSubdomainReports, *newDbSubdomain)
		}

		newDbReport := models.NewReport(value.Name, value.Score, newSubdomainReports, value.UserId, value.TestId, _domainIntensity, "")
		newDbReports[_domainName] = *newDbReport
	}

	for page := range constants.BIG_5_Report {

		prompt := API.CreatePrompt(constants.BIG_5_Report[page], processedScores)

		prompts[constants.BIG_5_Report[page]] = prompt

	}

	/**

	{
		"result": promptResult,
		"relationship": promptRelationship,
		"career_academic" : promptCareerAcademic,
		"strength_weakness" : promptStrengthWeakness
	}

	**/

	startTime = time.Millisecond
	// Now prompt generation starts
	channel := make(chan API.GeminiPromptRequest)

	for page, prompt := range prompts {
		go API.WorkerGCPGemini(page, prompt, channel)
	}

	results := map[string]string{}

	for range prompts {
		result := <-channel // Read the result from the channel
		// TODO add failure case
		// TODO added generated result to db
		// results[result] = result.Candidates[0].Content.Parts[0].Text
		results[result.Id] = result.Response.Candidates[0].Content.Parts[0].Text
	}

	fmt.Println("Time taken by GCP Worker to generate response from gemini", time.Millisecond-startTime)

	combinedDBReports := map[string]models.Report{}

	pdfGenerationContent := map[string]string{}

	for page, _ := range prompts {
		generatedResponseString := results[page]
		generatedContent, err := API.ParseMarkdownCode(generatedResponseString)

		if err != nil {
			// Respond with an error message if content generation failed
		}

		println("page:", page, generatedContent, generatedResponseString)
		pdfGenerationContent[page] = generatedResponseString
	}

	// TODO save to db, MAKE SURE YOU ARE CHECKING THE DOMAIN NAME CORRECTLY
	// TODO generate content from from GCP
	// reportDbColumn := mgm.Coll(&models.Report{})

	log.Println("Generating PDF")

	startTime = time.Second
	reportPdfFilename := "report_" + test.ID.Hex()
	errInPdfGeneration := API.GenerateBigFivePDF(pdfGenerationContent, test.TestGiver, reportPdfFilename)

	if errInPdfGeneration != nil {
		// Add logger
		return &MyError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate pdf",
		}
	}

	fmt.Println("Time taken to generate PDF", time.Second-startTime)

	startTime = time.Second
	log.Println("Sending Report via Email to user")
	API.SendBIG5Report(user.Email, test.TestGiver, "./"+reportPdfFilename+".pdf")

	fmt.Println("Time taken to send email", time.Second-startTime)

	startTime = time.Second
	var docs []interface{}

	for _, q := range combinedDBReports {
		docs = append(docs, q) // Add each question as an interface{}
	}

	_, err = mgm.Coll(&models.Report{}).InsertMany(c, docs)

	fmt.Println("Time taken to save Report in db", time.Second-startTime)

	if err != nil {
		// Add logger
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
