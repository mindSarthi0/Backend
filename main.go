package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"myproject/models"
	"myproject/response"
	"net/http"
	"sort"
)

// ["neuroticism", "n1", "Anxiety", "1", "2"]

type Subdomain struct {
	Name  string
	Score int
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

func getScoresWithQuestions(testId primitive.ObjectID) ([]ScoreQuestion, error) {
	// Initialize empty slice to store the final merged data
	var mergedData []ScoreQuestion

	// Step 1: Find all scores that match the given testId
	var scores []models.Score
	err := mgm.Coll(&models.Score{}).SimpleFind(&scores, bson.M{"testId": testId})
	if err != nil {
		return nil, fmt.Errorf("failed to find scores for testId %s: %v", testId.Hex(), err)
	}

	// Step 2: For each score entry, pull the corresponding question data
	for _, score := range scores {
		var question models.Question

		// Find the question that matches the questionId in the score
		err := mgm.Coll(&models.Question{}).FindByID(score.QuestionId.Hex(), &question)

		if err != nil {
			log.Printf("No question found for questionId: %s", score.QuestionId.Hex())
			continue
		} else if err != nil {
			return nil, fmt.Errorf("failed to find question for questionId %s: %v", score.QuestionId.Hex(), err)
		}

		// Step 3: Merge the score and question data
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

	// Step 4: Sort the merged data based on the "No" field in ascending order
	sort.Slice(mergedData, func(i, j int) bool {
		return mergedData[i].No < mergedData[j].No
	})

	return mergedData, nil
}

// func calculateSubdomainScore(subdomain string, score1 string, flow1 string, score2 string, flow2 stirng) (string, int) {

// 	// calculation
// 	subdomainScore := 3
// 	return subdomain, subdomainScore
// }

// func calculateProccessedScore(scoreQuestion []ScoreQuestion) {

// 	rule := map[string][][]string{
// 		"neuroticism": {
// 			{"n1", "Anxiety", "1", "N", "2", "R"},
// 			{"n2", "ABCD", "4", "N", "6", "R"},
// 		},
// 		"domain2": {
// 			{"n1", "Anxiety", "1", "N", "2", "R"},
// 			{"n2", "ABCD", "4", "N", "6", "R"},
// 		},
// 	}

// 	fmt.Println(scoreQuestion)

// for domain, values := range rule {
// 	fmt.Println("Domain:", domain)

// 	var pSubdomain []Subdomain

// 	for _, value := range values {

// 		fmt.Println("Subdomain:", value)
// 		subdomain := value[1]
// 		no1 := value[2]
// 		flow1 := value[3]

// 		// Score form db array
// 		score1 = scoreQuestion[no1+1]

// 		no2 := data[4]

// 		flow2 = datas[5]

// 		score2 = scores[no2+1]

// 		_, subdomainScore := calculateSubdomainScore(subdomain, score1, flow1, score2, flow2)

// 		pSubdomain = append(pSubdomain, Subdomain{subdomain, subdomainScore})

// 		fmt.Println(pSubdomain)
// 	}
// }

// TODO A single test should have one Id
func postSubmit(c *gin.Context) {

	var submit response.Submit

	if err := c.ShouldBindJSON(&submit); err != nil {
		// If there is an error, respond with 400 Bad Request
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Received person data: %+v\n", submit.Answers)
	// Find the user
	isExistingUsers := []models.User{}

	fmt.Printf("Received person data: %+v\n", isExistingUsers)
	mgm.Coll(&models.User{}).SimpleFind(&isExistingUsers, bson.M{"email": submit.Email})

	if isExistingUsers == nil || len(isExistingUsers) == 0 {
		user := models.NewUser(submit.Name, submit.Email, submit.Gender, submit.Age)

		err := mgm.Coll(user).Create(user)

		if err != nil {
			c.IndentedJSON(http.StatusOK, err)
			log.Fatalf("Failed to create a new err: %v", err) // Fatal will log and stop the program
		}
		isExistingUsers = append(isExistingUsers, *user)

	}

	userId := isExistingUsers[0].ID

	// TODO Generate payment link from razorPay

	// Create a test entry
	newTest := models.NewTest("BIG_5", userId, "PENDING", "https://google.com")

	newTestError := mgm.Coll(&models.Test{}).Create(newTest)

	if newTestError != nil {
		c.IndentedJSON(http.StatusOK, newTestError)
		log.Fatalf("Failed to create a new test: %v", newTestError) // Fatal will log and stop the program
	}

	// Store Score
	coll := mgm.Coll(&models.Score{})

	scoreDocs := []models.Score{}

	for _, item := range submit.Answers {
		questionId, err := primitive.ObjectIDFromHex(item.Id)
		if err != nil {
			log.Fatalf("Failed to convert string to ObjectID: %v", err)
		}
		score := models.NewScore(userId, questionId, item.Answer, newTest.ID)
		scoreDocs = append(scoreDocs, *score)
	}

	var docs []interface{}

	for _, score := range scoreDocs {
		docs = append(docs, score)
	}

	_, err := coll.InsertMany(c, docs)

	if err != nil {
		log.Fatalf("Failed to insert multiple documents: %v", err)
	}

	for _, item := range submit.Answers {
		questionId, err := primitive.ObjectIDFromHex(item.Id)
		if err != nil {
			log.Fatalf("Failed to convert string to ObjectID: %v", err)
		}
		score := models.NewScore(userId, questionId, item.Answer, newTest.ID)
		scoreDocs = append(scoreDocs, *score)
	}

	// calculateProccessedScore(scoresAndQuestion)

	log.Printf("Summited sucessfully")
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

	scoresAndQuestion, err := getScoresWithQuestions(primitive.ObjectID(changedTestId))

	if err != nil {
		log.Fatalf("Failed to get questions and score: %v", err)
	}

	log.Printf("Scores and Questions: %v", scoresAndQuestion)

	// Return the found questions
	c.IndentedJSON(http.StatusOK, scoresAndQuestion)

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

func getQuestions(c *gin.Context) {
	// Get the MongoDB collection for questions
	coll := mgm.Coll(&models.Question{})

	// Create a slice to hold the retrieved questions
	var questions []models.Question

	// Query all documents in the collection
	err := coll.SimpleFind(&questions, bson.D{})
	if err != nil {
		log.Printf("Error retrieving questions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve questions"})
		return
	}

	// Check if no documents are found
	if len(questions) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No questions found"})
		return
	}

	// Return the found questions
	c.IndentedJSON(http.StatusOK, questions)
}

func main() {
	router := gin.Default()

	router.POST("/questions", postQuestions)
	router.GET("/report", getReport)
	router.GET("/questions", getQuestions) // For retrieving questions
	router.POST("/submit", postSubmit)

	// Start the server on localhost:8080
	router.Run("localhost:8080")
}
