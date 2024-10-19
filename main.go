package main

import (
	"context"
	"fmt"
	"log"
	"myproject/API"
	"myproject/lib"
	"myproject/models"
	"myproject/response"
	"net/http"
	"sort"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"google.golang.org/genproto/googleapis/actions/sdk/v2/interactionmodel/prompt"
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

// ["neuroticism", "n1", "Anxiety", "1", "2"]

type Domain struct {
	Name      string
	Score     int
	Subdomain []Subdomain
	UserId    primitive.ObjectID `json:"userId" bson:"userId"`
	TestId    primitive.ObjectID `json:"testId" bson:"testId"`
	Intensity string
}

type Subdomain struct {
	Name      string
	Score     int
	Intensity string
}

type ScoreQuestion struct {
	UserId     primitive.ObjectID `json:"userId" bson:"userId"`
	TestId     primitive.ObjectID `json:"testId" bson:"testId"`
	QuestionId primitive.ObjectID `json:"questionId" bson:"questionId"`
	RawScore   string             `json:"rawScore" bson:"rawScore"`
	TestName   string             `json:"testName" bson:"testName"`
	Question   string             `json:"question" bson:"question"`
	No         int                `json:"no" bson:"no"`
}

// Fetch scores and corresponding questions based on testId
func fetchScoresWithQuestions(testId primitive.ObjectID) ([]ScoreQuestion, error) {
	var mergedData []ScoreQuestion
	var scores []models.Score

	// Fetch all scores matching the testId
	fmt.Println("Failed to get Score", testId)
	err := mgm.Coll(&models.Score{}).SimpleFind(&scores, bson.M{"testId": testId})
	if err != nil {
		fmt.Println("Failed to get Score", err)
		return nil, fmt.Errorf("failed to find scores for testId %s: %v", testId.Hex(), err)
	}

	// Merge score and question data
	for _, score := range scores {
		var question models.Question
		err := mgm.Coll(&models.Question{}).FindByID(score.QuestionId.Hex(), &question)
		if err != nil {
			log.Printf("No question found for questionId: %s", score.QuestionId.Hex())
			continue
		}

		mergedData = append(mergedData, ScoreQuestion{
			UserId:     score.UserId,
			TestId:     score.TestId,
			QuestionId: score.QuestionId,
			RawScore:   score.RawStore,
			TestName:   question.TestName,
			Question:   question.Question,
			No:         question.No,
		})
	}

	// Sort by question number
	sort.Slice(mergedData, func(i, j int) bool {
		return mergedData[i].No < mergedData[j].No
	})

	return mergedData, nil
}

// Calculate subdomain score (placeholder logic)
func calculateSubdomainScore(subdomain, score1, flow1, score2, flow2 string) (string, int, string) {
	subdomainScore := 3 // Placeholder score
	intensity := "low"  // Placeholder intensity
	return subdomain, subdomainScore, intensity
}

// Process and calculate score for domains and subdomains
func calculateProcessedScore(scoreQuestions []ScoreQuestion) []Domain {
	rules := map[string][][]string{
		"neuroticism": {
			{"n1", "Anxiety", "1", "N", "2", "R"},
			{"n2", "Anger", "3", "N", "4", "R"},
		},
		"domain_2": {
			{"n1", "Anxiety", "1", "N", "2", "R"},
			{"n2", "Anger", "3", "N", "4", "R"},
		},
	}

	var domains []Domain
	for domainName, subdomains := range rules {
		var domainScore int
		var processedSubdomains []Subdomain
		var testId, userId primitive.ObjectID

		for _, rule := range subdomains {
			subdomainName := rule[1]
			no1, flow1 := rule[2], rule[3]
			cNo1, err1 := lib.ConvertToInt(no1)
			if err1 != nil {
				log.Printf("Error converting question number: %v", err1)
				continue
			}

			score1 := scoreQuestions[cNo1-1]
			testId = score1.TestId
			userId = score1.UserId

			no2, flow2 := rule[4], rule[5]
			cNo2, err2 := lib.ConvertToInt(no2)
			if err2 != nil {
				log.Printf("Error converting question number: %v", err2)
				continue
			}
			score2 := scoreQuestions[cNo2-1]

			_, subdomainScore, intensity := calculateSubdomainScore(subdomainName, score1.RawScore, flow1, score2.RawScore, flow2)
			processedSubdomains = append(processedSubdomains, Subdomain{subdomainName, subdomainScore, intensity})
			domainScore += subdomainScore
		}

		domainIntensity := calculateDomainIntensity(domainName, domainScore)
		domains = append(domains, Domain{domainName, domainScore, processedSubdomains, testId, userId, domainIntensity})
	}

	return domains
}

// Placeholder for domain intensity calculation
func calculateDomainIntensity(domain string, score int) string {
	return "High" // Placeholder logic
}

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

	scoresAndQuestions, err := fetchScoresWithQuestions(testId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scores and questions"})
		return
	}

	processedScores := calculateProcessedScore(scoresAndQuestions)

	reportAlreadyGenerated := []models.Test{}
	mgm.Coll(&models.Test{}).SimpleFind(&reportAlreadyGenerated, bson.M{"testId": testId})

	if reportAlreadyGenerated != nil && len(reportAlreadyGenerated) != 0 {
		// Report already exit

		// Fetch the existing report using testId

		// Send to gcp to create report

		c.IndentedJSON(http.StatusOK, processedScores)
		return
	}

	newDbReports := []models.Report{}

	for domain, value := range processedScores {

		fmt.Println("Domain and Value:", domain, value)
		newSubdomainReports := []models.Subdomain{}

		for _, value := range value.Subdomain {
			newDbSubdomain := models.NewSubdomain(value.Name, value.Score, value.Intensity)

			newSubdomainReports = append(newSubdomainReports, *newDbSubdomain)
		}
		newDbResport := models.NewReport(value.Name, value.Score, newSubdomainReports, value.UserId, value.TestId)

		newDbReports = append(newDbReports, *newDbResport)
	}

	reportDbColumn := mgm.Coll(&models.Report{})
	log.Printf("Summited sucessfully", reportDbColumn)
	c.IndentedJSON(http.StatusOK, err)
}

func getReport(c *gin.Context) {
	var report response.Report

	if err := c.ShouldBindJSON(&report); err != nil {
		// If there is an error, respond with 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var testId = report.TestId

	changedTestId, err := primitive.ObjectIDFromHex(testId)

	if err != nil {
		fmt.Errorf("invalid testId: %v", err)
	}

	scoresAndQuestion, err := fetchScoresWithQuestions(primitive.ObjectID(changedTestId))

	if err != nil {
		log.Fatalf("Failed to get questions and score: %v", err)
	}

	log.Printf("Scores and Questions: %v", scoresAndQuestion)

	pDomain := calculateProcessedScore(scoresAndQuestion)
	c.IndentedJSON(http.StatusOK, pDomain)

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

func postQuestions(c *gin.Context) {
	var question []response.Question

	if err := c.ShouldBindJSON(&question); err != nil {
		// If there is an error, respond with 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coll := mgm.Coll(&models.Question{})

	questions := []models.Question{}

	for _, item := range question {
		questionToSave := models.NewQuestion(item.TestName, item.Question, item.No)
		questions = append(questions, *questionToSave)
	}

	var docs []interface{}
	for _, question := range questions {
		docs = append(docs, question)
	}

	_, err := coll.InsertMany(context.TODO(), docs)

	if err != nil {
		log.Printf("Failed to insert multiple question documents: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert documents"})
		return
	}

	log.Printf("Created sucessfully")
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Questions created successfully"})

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

func init() {
	err := mgm.SetDefaultConfig(nil, "cognify", options.Client().ApplyURI("mongodb+srv://cognify:dEQGVwIY24QzdUu6@cluster0.cjyqt.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
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

	// Start server
	router.GET("/testprompt", getPrompt)
	// Start the server on localhost:8080
	router.Run("localhost:8080")
}

func testMail(c *gin.Context) {
	API.Mail()
}

func creatingPdf(c *gin.Context) {
	lib.CreatePdf()
}

func getPrompt(c *gin.Context) {

	var values = map[string][]string{"neuroticism": []string{"7", "4", "6", "5", "8", "6", "4"}, "extraversion": []string{"7", "4", "6", "5", "8", "6", "4"}}

	prompts := []string{}

	for domain, value := range values {
		promt := API.CreatePrompt(domain, value[0], value[1], value[2], value[3], value[4], value[5], value[6])
		prompts = append(prompts, promt)
	}

	channel := make(chan API.ContentResponse)

	for _, prompt := range prompts {
		go worker(prompt, channel)
	}

	results := []string{}
	for range prompts {
		result := <-channel // Read the result from the channel
		// TODO add failure case
		results = append(results, result.Candidates[0].Content.Parts[0].Text)
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Prompt generated successfully", "prompt": prompts, "Gemini Response": results})
}

func worker(prompt string, channel chan API.ContentResponse) {
	result, err := API.GenerateContentFromTextGCPJSON(prompt)
	if err != nil {
		// Respond with an error message if content generation failed
	}
	channel <- *result
}
