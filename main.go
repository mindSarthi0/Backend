package main

import (
	"fmt"
	"log"
	"myproject/API"
	"myproject/lib"
	"myproject/models"
	"myproject/response"
	"net/http"
	"sort"
	"strconv"

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

	if flow1 == "R" {
		score1Int, err := strconv.Atoi(score1)
		if err != nil {
			log.Printf("Error converting score1 to int: %v", err)
			return subdomain, 0, "Error"
		}
		score1Int = 6 - score1Int
		score1 = strconv.Itoa(score1Int)
	}

	// Adjust score2 based on flow2
	if flow2 == "R" {
		score2Int, err := strconv.Atoi(score2)
		if err != nil {
			log.Printf("Error converting score2 to int: %v", err)
			return subdomain, 0, "Error"
		}
		score2Int = 6 - score2Int
		score2 = strconv.Itoa(score2Int)
	}

	// Calculate the average subdomain score
	score1Int, _ := strconv.Atoi(score1)
	score2Int, _ := strconv.Atoi(score2)
	subdomainScore := score1Int + score2Int

	// Determine the intensity based on subdomain score
	var intensity string

	if subdomainScore > 8 {
		intensity = "High"
	} else if subdomainScore > 6 {
		intensity = "Above Average"
	} else if subdomainScore > 4 {
		intensity = "Average"
	} else if subdomainScore > 3 {
		intensity = "Below Average"
	} else {
		intensity = "Low"
	}

	return subdomain, subdomainScore, intensity
}

// Process and calculate score for domains and subdomains
func calculateProcessedScore(scoreQuestions []ScoreQuestion) []Domain {
	rules := map[string][][]string{
		"neuroticism": {
			{"n1", "Anxiety", "1", "N", "2", "N"},
			{"n2", "Anger", "3", "N", "4", "N"},
			{"n3", "Depression", "5", "N", "6", "N"},
			{"n4", "Self-consciousness", "7", "N", "8", "N"},
			{"n5", "Immoderation", "9", "R", "10", "R"},
			{"n6", "Vulnerability", "11", "R", "12", "R"},
		},
		"extraversion": {
			{"e1", "Friendliness", "13", "N", "14", "N"},
			{"e2", "Gregariousness", "15", "N", "16", "R"},
			{"e3", "Assertiveness", "17", "N", "18", "N"},
			{"e4", "Activity Level", "19", "N", "20", "N"},
			{"e5", "Excitement Seeking", "21", "N", "22", "N"},
			{"e6", "Cheerfulness", "23", "N", "24", "N"},
		},
		"openness": {
			{"o1", "Imagination", "25", "N", "26", "N"},
			{"o2", "Artistic Interests", "27", "N", "28", "R"},
			{"o3", "Emotionality", "29", "N", "30", "R"},
			{"o4", "Adventurousness", "31", "R", "32", "R"},
			{"o5", "Intellect", "33", "R", "34", "R"},
			{"o6", "Liberalism", "35", "N", "36", "R"},
		},
		"agreeableness": {
			{"a1", "Trust", "37", "N", "38", "N"},
			{"a2", "Morality", "39", "R", "40", "R"},
			{"a3", "Altruism", "41", "N", "42", "N"},
			{"a4", "Cooperation", "43", "R", "44", "R"},
			{"a5", "Modesty", "45", "R", "46", "R"},
			{"a6", "Sympathy", "47", "N", "48", "N"},
		},
		"conscientiousness": {
			{"c1", "Self Efficacy", "49", "N", "50", "N"},
			{"c2", "Orderliness", "51", "N", "52", "R"},
			{"c3", "Dutifulness", "53", "N", "54", "R"},
			{"c4", "Achievement Striving", "55", "N", "56", "N"},
			{"c5", "Self Discipline", "57", "N", "58", "R"},
			{"c6", "Cautiousness", "59", "R", "60", "R"},
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

		domainIntensity := calculateDomainIntensity(domainScore)
		domains = append(domains, Domain{domainName, domainScore, processedSubdomains, testId, userId, domainIntensity})
	}

	return domains
}

// Placeholder for domain intensity calculation
func calculateDomainIntensity(domainscore int) string {
	var domainIntensity string
	if domainscore >= 50 {
		domainIntensity = "High"
	} else if domainscore >= 40 {
		domainIntensity = "Above Average"
	} else if domainscore >= 30 {
		domainIntensity = "Average"
	} else if domainscore >= 20 {
		domainIntensity = "Below Average"
	} else if domainscore >= 10 {
		domainIntensity = "Low"
	}
	return domainIntensity
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

	if len(reportAlreadyGenerated) != 0 {
		// Report already exit
		// Fetch the existing report using testId
		// Send to gcp to create report
		c.IndentedJSON(http.StatusOK, processedScores)
		return
	}

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

	for domain, value := range newDbReports {

		// TODO add failure case
		// TODO added generated result to db
		// results[result] = result.Candidates[0].Content.Parts[0].Text
		combinedDBReport := models.NewReport(value.Name, value.Score, value.Subdomain, value.UserId, value.TestId, value.Intensity, results[domain])
		combinedDBReports[domain] = *combinedDBReport
	}

	// TODO save to db, MAKE SURE YOU ARE CHECKING THE DOMAIN NAME CORRECTLY
	// TODO generate content from from GCP
	// reportDbColumn := mgm.Coll(&models.Report{})

	log.Println("Summited sucessfully", combinedDBReports)

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

func init() {
	// Setup the mgm default config
	err := mgm.SetDefaultConfig(nil, "cognify", options.Client().ApplyURI("mongodb+srv://cognify:dEQGVwIY24QzdUu6@cluster0.cjyqt.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	// Error handling
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err) // Fatal will log and stop the program
	}

	fmt.Println("Successfully connected to MongoDB!")
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
	// router.GET("/testprompt", getPrompt)
	// Start the server on localhost:8080
	router.Run("localhost:8080")
}

func testMail(c *gin.Context) {
	API.Mail()
}

func creatingPdf(c *gin.Context) {
	lib.CreatePdfWithBg()
}

// func getPrompt(c *gin.Context) {

var values = map[string][]string{"neuroticism": {"7", "4", "6", "5", "8", "6", "4"}, "extraversion": {"3", "5", "2", "6", "6", "2", "3"}, "openness": {"7", "4", "6", "5", "8", "6", "4"}, "agreeableness": {"7", "4", "6", "5", "8", "6", "4"}, "conscientiousness": {"7", "4", "6", "5", "8", "6", "4"}}

// 	prompts := []string{}

// 	for domain, value := range values {
// 		promt := API.CreatePrompt(domain, value[0], value[1], value[2], value[3], value[4], value[5], value[6])
// 		prompts = append(prompts, promt)
// 	}

// 	channel := make(chan API.ContentResponse)

// 	for _, prompt := range prompts {
// 		go API.WorkerGCPGemini(prompt, channel)
// 	}

// 	results := []string{}

// 	for range prompts {
// 		result := <-channel // Read the result from the channel
// 		// TODO add failure case
// 		results = append(results, result.Candidates[0].Content.Parts[0].Text)
// 	}

// 	// Respond with a success message
// 	c.JSON(http.StatusOK, gin.H{"message": "Prompt generated successfully", "prompt": prompts, "Gemini Response": results})
// }
