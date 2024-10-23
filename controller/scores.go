package controller

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"myproject/lib"
	"myproject/models"
	"sort"
	"strconv"
)

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
func FetchScoresWithQuestions(testId primitive.ObjectID) ([]ScoreQuestion, error) {
	var mergedData []ScoreQuestion
	var scores []models.Score

	// Fetch all scores matching the testId
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

func CalculateProcessedScore(scoreQuestions []ScoreQuestion) []Domain {
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
