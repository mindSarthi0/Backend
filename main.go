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
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Handle test submissions
func handleSubmission(c *gin.Context) {
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

	user := existingUsers[0]

	// Create a new test entry
	newTest := models.NewTest("BIG_5", user.ID, "PENDING", "https://google.com")
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

	// TODO add logger
	go controller.GenerateNewReport(c, newTest.ID.Hex(), user)

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

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	if count != 0 {
		c.IndentedJSON(http.StatusAlreadyReported, gin.H{"message": "report is already generated for test id"})
		return
	}

	errFromRequest := controller.GenerateNewReport(c, test.ID.String(), user)

	if errFromRequest != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get count of already generated reports"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report generated successfully"})
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

	API.GenerateBigFivePDF(testContent, "User name", "report")
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
	router.POST("/submit", handleSubmission)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback to port 8080 if not set
	}

	srv := &http.Server{
		Addr:    ":" + port,
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
