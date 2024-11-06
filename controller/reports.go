package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"log"
	"myproject/API"
	"myproject/models"
	"net/http"
	"strconv"
	"time"
)

type MyError struct {
	Code    int
	Message string
}

func GenerateNewReport(c *gin.Context, test models.Test, user models.User) *MyError {

	startTime := time.Second
	scoresAndQuestions, err := FetchScoresWithQuestions(test.ID)

	if err != nil {
		return &MyError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch scores and questions",
		}
	}

	fmt.Println("Time taken by Fetch Scores with Questions API calls", time.Second-startTime)
	processedScores := CalculateProcessedScore(scoresAndQuestions)

	newDbReports := map[string]models.Report{}

	prompts := map[string]string{}

	for _, value := range processedScores {

		// fmt.Println("Domain and Value:", domain, value)
		newSubdomainReports := []models.Subdomain{}

		_domainName := value.Name
		_domainScore := value.Score
		_domainIntensity := value.Intensity

		for _, value := range value.Subdomain {
			newDbSubdomain := models.NewSubdomain(value.Name, value.Score, value.Intensity)
			newSubdomainReports = append(newSubdomainReports, *newDbSubdomain)
		}

		// TODO do error handling
		s1 := strconv.Itoa(newSubdomainReports[0].Score)
		s2 := strconv.Itoa(newSubdomainReports[1].Score)
		s3 := strconv.Itoa(newSubdomainReports[2].Score)
		s4 := strconv.Itoa(newSubdomainReports[3].Score)
		s5 := strconv.Itoa(newSubdomainReports[4].Score)
		s6 := strconv.Itoa(newSubdomainReports[5].Score)

		// TODO create a better interface here
		prompt := API.CreatePrompt(_domainName, strconv.Itoa(_domainScore), s1, s2, s3, s4, s5, s6)

		prompts[_domainName] = prompt

		newDbReport := models.NewReport(value.Name, value.Score, newSubdomainReports, value.UserId, value.TestId, _domainIntensity, "")
		newDbReports[_domainName] = *newDbReport
	}

	startTime = time.Second
	// Now prompt generation starts
	channel := make(chan API.GeminiPromptRequest)

	for domain, prompt := range prompts {
		go API.WorkerGCPGemini(domain, prompt, channel)
	}

	results := map[string]string{}

	for range newDbReports {
		result := <-channel // Read the result from the channel
		// TODO add failure case
		// TODO added generated result to db
		// results[result] = result.Candidates[0].Content.Parts[0].Text
		results[result.Id] = result.Response.Candidates[0].Content.Parts[0].Text
	}

	fmt.Println("Time taken by GCP Worker to generate response from gemini", time.Second-startTime)

	combinedDBReports := map[string]models.Report{}

	pdfGenerationContent := map[string]API.JSONOutputFormat{}

	for domain, value := range newDbReports {

		// TODO add failure case
		// TODO added generated result to db
		// results[result] = result.Candidates[0].Content.Parts[0].Text

		generatedResponseString := results[domain]

		combinedDBReport := models.NewReport(value.Name, value.Score, value.Subdomain, value.UserId, value.TestId, value.Intensity, generatedResponseString)

		formatedJson, err := API.ParseMarkdownCode(generatedResponseString)

		if err != nil {
			// Respond with an error message if content generation failed
		}

		pdfGenerationContent[domain] = formatedJson

		combinedDBReports[domain] = *combinedDBReport
	}

	// TODO save to db, MAKE SURE YOU ARE CHECKING THE DOMAIN NAME CORRECTLY
	// TODO generate content from from GCP
	// reportDbColumn := mgm.Coll(&models.Report{})

	log.Println("Generating PDF")

	startTime = time.Second
	reportPdfFilename := "report_" + test.ID.Hex()
	errInPdfGeneration := API.GenerateBigFivePDF(pdfGenerationContent, user.Name, reportPdfFilename)

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
