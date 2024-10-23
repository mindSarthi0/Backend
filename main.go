package main

import (
	"context"
	"fmt"
	"log"
	"myproject/API"
	"myproject/controller"
	"myproject/lib"
	"myproject/models"
	"myproject/response"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"os/signal"
	"time"
	//"google.golang.org/genproto/googleapis/actions/sdk/v2/interactionmodel/prompt"
)

// Handle test submissions
func handleTestSubmission(c *gin.Context) {
	var submission response.Submit

	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission format"})
		return
	}

	// Find or create the user
	var existingUsers []models.User
	mgm.Coll(&models.User{}).SimpleFind(&existingUsers, bson.M{"email": submission.Email})

	if len(existingUsers) == 0 {
		newUser := models.NewUser(submission.Name, submission.Email, submission.Gender, submission.Age)
		if err := mgm.Coll(newUser).Create(newUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		existingUsers = append(existingUsers, *newUser)
	}

	userId := existingUsers[0].ID

	// Create a new test entry
	newTest := models.NewTest("BIG_5", userId, "PENDING", "https://google.com")
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
		scoreDocs = append(scoreDocs, *models.NewScore(userId, questionId, answer.Answer, newTest.ID))
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

	c.JSON(http.StatusOK, gin.H{"message": "Submission successful"})
}

// Generate report
func handleReportGeneration(c *gin.Context) {
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

	println("::::::::::::::count::::::::::::", count, testId.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get count of already generated reports"})
		return
	}

	if count != 0 {
		c.IndentedJSON(http.StatusAlreadyReported, gin.H{"message": "report is already generated for test id"})
		return
	}

	scoresAndQuestions, err := controller.FetchScoresWithQuestions(testId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scores and questions"})
		return
	}

	processedScores := controller.CalculateProcessedScore(scoresAndQuestions)

	newDbReports := map[string]models.Report{}

	prompts := map[string]string{}

	for domain, value := range processedScores {

		fmt.Println("Domain and Value:", domain, value)
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

		log.Println("Formated JSON", formatedJson.Introduction)
		pdfGenerationContent[domain] = formatedJson

		combinedDBReports[domain] = *combinedDBReport
	}

	// TODO save to db, MAKE SURE YOU ARE CHECKING THE DOMAIN NAME CORRECTLY
	// TODO generate content from from GCP
	// reportDbColumn := mgm.Coll(&models.Report{})

	log.Println("Generating PDF")

	reportPdfFilename := "report_" + reportRequest.TestId
	API.GenerateBigFivePDF(pdfGenerationContent, reportPdfFilename)
	// TODO handle error

	log.Println("Sending Report via Email to user")
	API.SendBIG5Report(user.Email, "./"+reportPdfFilename+".pdf")
	// TODO handle error

	log.Println("Submmited sucessfully", combinedDBReports)

	var docs []interface{}
	for _, q := range combinedDBReports {
		docs = append(docs, q) // Add each question as an interface{}
	}

	responseDb, err := mgm.Coll(&models.Report{}).InsertMany(c, docs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert questions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prompt generated successfully", "prompt": prompts, "Gemini Response": responseDb})
}

// Fetch all questions
func fetchAllQuestions(c *gin.Context) {
	var questions []models.Question
	err := mgm.Coll(&models.Question{}).SimpleFind(&questions, bson.D{})
	if err != nil {
		fmt.Println("Failed to retrieve questions:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve questions"})
		return
	}

	if len(questions) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No questions found"})
		return
	}

	c.JSON(http.StatusOK, questions)
}

// Submit questions
func submitQuestions(c *gin.Context) {
	var questions []response.Question
	if err := c.ShouldBindJSON(&questions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question data"})
		return
	}

	var questionDocs []models.Question
	for _, item := range questions {
		questionDocs = append(questionDocs, *models.NewQuestion(item.TestName, item.Question, item.No))
	}

	var docs []interface{}
	for _, q := range questionDocs {
		docs = append(docs, q) // Add each question as an interface{}
	}

	_, err := mgm.Coll(&models.Question{}).InsertMany(c, docs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert questions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Questions submitted successfully"})
}

func generatepdf(c *gin.Context) {
	testContent := map[string]API.JSONOutputFormat{
		"neuroticism": {
			Introduction:     "This is intro to BIG 5",
			CareerAcademia:   "Please check with career councellor",
			Relationship:     "Good with relationship",
			StrengthWeakness: "Good with undertanding Black and White thinking",
		},
		"extraversion": {
			Introduction:     "This is intro to BIG 5",
			CareerAcademia:   "Please check with career councellor",
			Relationship:     "Good with relationship",
			StrengthWeakness: "Good with undertanding Black and White thinking",
		},
		"openness": {
			Introduction:     "This is intro to BIG 5",
			CareerAcademia:   "Please check with career councellor",
			Relationship:     "Good with relationship",
			StrengthWeakness: "Good with undertanding Black and White thinking",
		},
		"agreeableness": {
			Introduction:     "This is intro to BIG 5",
			CareerAcademia:   "Please check with career councellor",
			Relationship:     "Good with relationship",
			StrengthWeakness: "Good with undertanding Black and White thinking",
		},
		"conscientiousness": {
			Introduction:     "This is intro to BIG 5",
			CareerAcademia:   "Please check with career councellor",
			Relationship:     "Good with relationship",
			StrengthWeakness: "Good with undertanding Black and White thinking",
		},
	}

	API.GenerateBigFivePDF(testContent, "report")
}

func testMail(c *gin.Context) {
	API.Mail()
}

func creatingPdf(c *gin.Context) {
	lib.CreatePdfWithBg()
}

func getPrompt(c *gin.Context) {

	var values = map[string][]string{"neuroticism": {"7", "4", "6", "5", "8", "6", "4"}, "extraversion": {"3", "5", "2", "6", "6", "2", "3"}, "openness": {"7", "4", "6", "5", "8", "6", "4"}, "agreeableness": {"7", "4", "6", "5", "8", "6", "4"}, "conscientiousness": {"7", "4", "6", "5", "8", "6", "4"}}

	prompts := map[string]string{}

	for domain, value := range values {
		promt := API.CreatePrompt(domain, value[0], value[1], value[2], value[3], value[4], value[5], value[6])
		prompts[domain] = promt
	}

	channel := make(chan API.GeminiPromptRequest)

	for domain, prompt := range prompts {
		go API.WorkerGCPGemini(domain, prompt, channel)
	}

	results := map[string]string{}

	for range prompts {
		result := <-channel // Read the result from the channel
		// TODO add failure case
		generatedResponseString := result.Response.Candidates[0].Content.Parts[0].Text
		formatedJson, err := API.ParseMarkdownCode(generatedResponseString)

		if err != nil {
			// Respond with an error message if content generation failed
		}

		log.Println("Formated JSON", formatedJson.Introduction)
		results[result.Id] = generatedResponseString
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Prompt generated successfully", "prompt": prompts, "Gemini Response": results})
}

func init() {
	errLoadingEnv := godotenv.Load()
	if errLoadingEnv != nil {
		log.Fatalf("Error loading .env file")
	}

	// Setup the mgm default config
	err := mgm.SetDefaultConfig(nil, "cognify", options.Client().ApplyURI("mongodb+srv://cognify:dEQGVwIY24QzdUu6@cluster0.cjyqt.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	// Error handling
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err) // Fatal will log and stop the program
	}

	fmt.Println("Successfully connected to MongoDB!")
}

func main() {
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
	}))

	// Routes
	router.POST("/questions", submitQuestions)
	router.GET("/questions", fetchAllQuestions)
	router.POST("/submit", handleTestSubmission)
	router.GET("/report", handleReportGeneration)
	//Pdf test route
	router.POST("/pdf", creatingPdf)
	router.POST("/mail", testMail)
	router.POST("/generatepdf", generatepdf)

	// Start server
	router.GET("/testprompt", getPrompt)
	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	gin.SetMode(gin.ReleaseMode)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	// Timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
