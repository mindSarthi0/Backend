package API

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	// "github.com/joho/godotenv"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type GeminiPromptRequest struct {
	Id       string
	Response ContentResponse
}

type ContentResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason  string `json:"finishReason"`
		Index         int    `json:"index"`
		SafetyRatings []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
}

// Function to parse markdown code and extract JSON
func ParseMarkdownCode(markdown string) (string, error) {
	// Define the struct to store the parsed data
	var report string

	// Extract the JSON part (between the "```json" block)
	start := strings.Index(markdown, "```json")
	end := strings.LastIndex(markdown, "```")

	if start == -1 || end == -1 {
		return report, fmt.Errorf("invalid markdown format")
	}

	// // Extract the JSON part and trim spaces
	jsonPart := markdown[start+len("```json") : end]
	// jsonPart := strings.TrimSpace(markdown)

	// Remove all newlines
	jsonPart = strings.ReplaceAll(jsonPart, "\n", "")

	// Unmarshal the JSON into the struct
	err := json.Unmarshal([]byte(jsonPart), &report)
	if err != nil {
		return report, err
	}

	return report, nil
}

// Function to make the POST request to Google API (REST-based approach)
func GenerateContentFromTextGCP(prompt string) (string, error) {
	// Define the URL for the Google API endpoint (generative language model)
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-pro-latest:generateContent"

	// Prepare the request payload in JSON format with the prompt content
	requestBody, err := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Retrieve the API key from the environment variable 'API_KEY'
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("API key is not set. Please ensure the environment variable 'API_KEY' is set")
	}

	// Create a new HTTP POST request with the request body and API key
	req, err := http.NewRequest("POST", url+"?key="+apiKey, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set appropriate headers for the request
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request using an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Check if the response status is not 200 (OK), and handle the error accordingly
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 response status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response into the struct
	var contentResponse ContentResponse
	err = json.Unmarshal(body, &contentResponse)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	// Return the response body as a string
	return string(body), nil
}

func GenerateContentFromTextGOAPIJSON(promt string) (*genai.GenerateContentResponse, error) {
	ctx := context.Background()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key is not set. Please ensure the environment variable 'API_KEY' is set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro-latest")
	// Ask the model to respond with JSON.
	model.ResponseMIMEType = "application/json"
	// Specify the schema.
	model.ResponseSchema = &genai.Schema{
		Type:  genai.TypeArray,
		Items: &genai.Schema{Type: genai.TypeString},
	}
	resp, err := model.GenerateContent(ctx, genai.Text(promt))
	if err != nil {
		log.Fatal(err)
	}
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			var recipes []string
			if err := json.Unmarshal([]byte(txt), &recipes); err != nil {
				log.Fatal(err)
			}
			fmt.Println(recipes)
		}
	}

	return resp, err
}

// Function to make the POST request to Google API (REST-based approach)
func GenerateContentFromTextGCPJSON(prompt string) (*ContentResponse, error) {

	// Define the URL for the Google API endpoint (generative language model)
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-pro-latest:generateContent"

	// Prepare the request payload in JSON format with the prompt content
	requestBody, err := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Retrieve the API key from the environment variable 'API_KEY'
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key is not set. Please ensure the environment variable 'API_KEY' is set")
	}

	// Create a new HTTP POST request with the request body and API key
	req, err := http.NewRequest("POST", url+"?key="+apiKey, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set appropriate headers for the request
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request using an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check if the response status is not 200 (OK), and handle the error accordingly
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response into the struct
	var contentResponse ContentResponse
	err = json.Unmarshal(body, &contentResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	// Return the response body as a string
	return &contentResponse, nil
}

// func CreatePrompt(d string, ds string, s1 string, s2 string, s3 string, s4 string, s5 string, s6 string) string {
// 	switch d {
// 	case "neuroticism":
// 		return CreatePromptNeuroticism(d, ds, s1, s2, s3, s4, s5, s6)
// 	case "extraversion":
// 		return CreatePromptExtraversion(d, ds, s1, s2, s3, s4, s5, s6)
// 	case "openness":
// 		return CreatePromptOpenness(d, ds, s1, s2, s3, s4, s5, s6)
// 	case "agreeableness":
// 		return CreatePromptAgreeableness(d, ds, s1, s2, s3, s4, s5, s6)
// 	case "conscientiousness":
// 		return CreatePromptConscientiousness(d, ds, s1, s2, s3, s4, s5, s6)
// 	}
// 	return ""
// }

func CreatePrompt(page string, score []Domain) string {
	switch page {
	case "result":
		return CreatePromptResult(score)
	case "career_academic":
		return CreatePromptCareerAcademic(score)
	case "relationship":
		return CreatePromptRelationship(score)
	case "strength_weakness":
		return CreatePromptStrengthWeakness(score)
	}

	return ""
}

func CreatePromptResult(score []Domain) string {

	neuroticismDomain := score[0]

	neuroticismScore := neuroticismDomain.Score

	n1 := neuroticismDomain.Subdomain[0].Score
	n2 := neuroticismDomain.Subdomain[1].Score
	n3 := neuroticismDomain.Subdomain[2].Score
	n4 := neuroticismDomain.Subdomain[3].Score
	n5 := neuroticismDomain.Subdomain[4].Score
	n6 := neuroticismDomain.Subdomain[5].Score

	extraversionDomain := score[1]

	extraversionScore := extraversionDomain.Score

	e1 := extraversionDomain.Subdomain[0].Score
	e2 := extraversionDomain.Subdomain[1].Score
	e3 := extraversionDomain.Subdomain[2].Score
	e4 := extraversionDomain.Subdomain[3].Score
	e5 := extraversionDomain.Subdomain[4].Score
	e6 := extraversionDomain.Subdomain[5].Score

	opennessDomain := score[2]

	opennessScore := opennessDomain.Score

	o1 := opennessDomain.Subdomain[0].Score
	o2 := opennessDomain.Subdomain[1].Score
	o3 := opennessDomain.Subdomain[2].Score
	o4 := opennessDomain.Subdomain[3].Score
	o5 := opennessDomain.Subdomain[4].Score
	o6 := opennessDomain.Subdomain[5].Score

	agreeablenessDomain := score[3]

	agreeablenessScore := agreeablenessDomain.Score

	a1 := agreeablenessDomain.Subdomain[0].Score
	a2 := agreeablenessDomain.Subdomain[1].Score
	a3 := agreeablenessDomain.Subdomain[2].Score
	a4 := agreeablenessDomain.Subdomain[3].Score
	a5 := agreeablenessDomain.Subdomain[4].Score
	a6 := agreeablenessDomain.Subdomain[5].Score

	conscientiousnessDomain := score[4]

	conscientiousnessScore := conscientiousnessDomain.Score

	c1 := conscientiousnessDomain.Subdomain[0].Score
	c2 := conscientiousnessDomain.Subdomain[1].Score
	c3 := conscientiousnessDomain.Subdomain[2].Score
	c4 := conscientiousnessDomain.Subdomain[3].Score
	c5 := conscientiousnessDomain.Subdomain[4].Score
	c6 := conscientiousnessDomain.Subdomain[5].Score

	prompt := fmt.Sprintf("Using the Big 5 score given below, create Sumarry for each domain for the Report"+

		"Domain Score to Intensity map { [0 - 30] : 'Low', [31 - 39]: 'Average',[40 - 60]: 'high'}"+
		"Subdomain Score to Intensity map { [0 - 4] : 'Low', [5 - 7]: 'Average',[8 - 10]: 'high'}"+

		"Personality Assessment:\n\n"+
		"Domain: Neuroticism Score: %s\n"+
		"  Subdomains-\n"+
		"    Anxiety Score: %s\n"+
		"    Anger Score: %s\n"+
		"    Depression Score: %s\n"+
		"    Self-consciousness Score: %s\n"+
		"    Immoderation Score: %s\n"+
		"    Vulnerability Score: %s\n\n"+

		"Domain: Extraversion Score: %s\n"+
		"  Subdomains-\n"+
		"    Friendliness Score: %s\n"+
		"    Gregariousness Score: %s\n"+
		"    Assertiveness Score: %s\n"+
		"    Activity Level Score: %s\n"+
		"    Excitement Seeking Score: %s\n"+
		"    Cheerfulness Score: %s\n\n"+

		"Domain: Openness: %s\n"+
		"  Subdomains-\n"+
		"    Imagination Score: %s\n"+
		"    Artistic Interests Score: %s\n"+
		"    Emotionality Score: %s\n"+
		"    Adventurousness Score: %s\n"+
		"    Intellect Score: %s\n"+
		"    Liberalism Score: %s\n\n"+

		"Domain: Agreeableness: %s\n"+
		"  Subdomains-\n"+
		"    Trust Score: %s\n"+
		"    Morality Score: %s\n"+
		"    Altruism Score: %s\n"+
		"    Cooperation Score: %s\n"+
		"    Modesty Score: %s\n"+
		"    Sympathy Score: %s\n\n"+

		"Domain: Conscientiousness: %s\n"+
		"  Subdomains-\n"+
		"    Self Efficacy Score: %s\n"+
		"    Orderliness Score: %s\n"+
		"    Dutifulness Score: %s\n"+
		"    Achievement Striving Score: %s\n"+
		"    Self Discipline Score: %s\n"+
		"    Cautiousness Score: %s\n",

		neuroticismScore, n1, n2, n3, n4, n5, n6,
		extraversionScore, e1, e2, e3, e4, e5, e6,
		opennessScore, o1, o2, o3, o4, o5, o6,
		agreeablenessScore, a1, a2, a3, a4, a5, a6,
		conscientiousnessScore, c1, c2, c3, c4, c5, c6)

	return prompt
}

func CreatePromptCareerAcademic(score []Domain) string {

	neuroticismDomain := score[0]

	neuroticismScore := neuroticismDomain.Score

	n1 := neuroticismDomain.Subdomain[0].Score
	n2 := neuroticismDomain.Subdomain[1].Score
	n3 := neuroticismDomain.Subdomain[2].Score
	n4 := neuroticismDomain.Subdomain[3].Score
	n5 := neuroticismDomain.Subdomain[4].Score
	n6 := neuroticismDomain.Subdomain[5].Score

	extraversionDomain := score[1]

	extraversionScore := extraversionDomain.Score

	e1 := extraversionDomain.Subdomain[0].Score
	e2 := extraversionDomain.Subdomain[1].Score
	e3 := extraversionDomain.Subdomain[2].Score
	e4 := extraversionDomain.Subdomain[3].Score
	e5 := extraversionDomain.Subdomain[4].Score
	e6 := extraversionDomain.Subdomain[5].Score

	opennessDomain := score[2]

	opennessScore := opennessDomain.Score

	o1 := opennessDomain.Subdomain[0].Score
	o2 := opennessDomain.Subdomain[1].Score
	o3 := opennessDomain.Subdomain[2].Score
	o4 := opennessDomain.Subdomain[3].Score
	o5 := opennessDomain.Subdomain[4].Score
	o6 := opennessDomain.Subdomain[5].Score

	agreeablenessDomain := score[3]

	agreeablenessScore := agreeablenessDomain.Score

	a1 := agreeablenessDomain.Subdomain[0].Score
	a2 := agreeablenessDomain.Subdomain[1].Score
	a3 := agreeablenessDomain.Subdomain[2].Score
	a4 := agreeablenessDomain.Subdomain[3].Score
	a5 := agreeablenessDomain.Subdomain[4].Score
	a6 := agreeablenessDomain.Subdomain[5].Score

	conscientiousnessDomain := score[4]

	conscientiousnessScore := conscientiousnessDomain.Score

	c1 := conscientiousnessDomain.Subdomain[0].Score
	c2 := conscientiousnessDomain.Subdomain[1].Score
	c3 := conscientiousnessDomain.Subdomain[2].Score
	c4 := conscientiousnessDomain.Subdomain[3].Score
	c5 := conscientiousnessDomain.Subdomain[4].Score
	c6 := conscientiousnessDomain.Subdomain[5].Score

	prompt := fmt.Sprintf("Using the Big 5 score given below, create Career & Academia Page for the Report\n\n"+
		"Domain Score to Intensity map { [0 - 30] : 'Low', [31 - 39]: 'Average',[40 - 60]: 'high'}"+
		"Subdomain Score to Intensity map { [0 - 4] : 'Low', [5 - 7]: 'Average',[8 - 10]: 'high'}"+

		"Personality Assessment:\n\n"+
		"Domain: Neuroticism: %s\n"+
		"  Subdomains-\n"+
		"    Anxiety: %s\n"+
		"    Anger: %s\n"+
		"    Depression: %s\n"+
		"    Self-consciousness: %s\n"+
		"    Immoderation: %s\n"+
		"    Vulnerability: %s\n\n"+

		"Domain: Extraversion: %s\n"+
		"  Subdomains-\n"+
		"    Friendliness: %s\n"+
		"    Gregariousness: %s\n"+
		"    Assertiveness: %s\n"+
		"    Activity Level: %s\n"+
		"    Excitement Seeking: %s\n"+
		"    Cheerfulness: %s\n\n"+

		"Domain: Openness: %s\n"+
		"  Subdomains-\n"+
		"    Imagination: %s\n"+
		"    Artistic Interests: %s\n"+
		"    Emotionality: %s\n"+
		"    Adventurousness: %s\n"+
		"    Intellect: %s\n"+
		"    Liberalism: %s\n\n"+

		"Domain: Agreeableness: %s\n"+
		"  Subdomains-\n"+
		"    Trust: %s\n"+
		"    Morality: %s\n"+
		"    Altruism: %s\n"+
		"    Cooperation: %s\n"+
		"    Modesty: %s\n"+
		"    Sympathy: %s\n\n"+

		"Domain: Conscientiousness: %s\n"+
		"  Subdomains-\n"+
		"    Self Efficacy: %s\n"+
		"    Orderliness: %s\n"+
		"    Dutifulness: %s\n"+
		"    Achievement Striving: %s\n"+
		"    Self Discipline: %s\n"+
		"    Cautiousness: %s\n",

		neuroticismScore, n1, n2, n3, n4, n5, n6,
		extraversionScore, e1, e2, e3, e4, e5, e6,
		opennessScore, o1, o2, o3, o4, o5, o6,
		agreeablenessScore, a1, a2, a3, a4, a5, a6,
		conscientiousnessScore, c1, c2, c3, c4, c5, c6)

	return prompt
}

func CreatePromptRelationship(score []Domain) string {

	neuroticismDomain := score[0]

	neuroticismScore := neuroticismDomain.Score

	n1 := neuroticismDomain.Subdomain[0].Score
	n2 := neuroticismDomain.Subdomain[1].Score
	n3 := neuroticismDomain.Subdomain[2].Score
	n4 := neuroticismDomain.Subdomain[3].Score
	n5 := neuroticismDomain.Subdomain[4].Score
	n6 := neuroticismDomain.Subdomain[5].Score

	extraversionDomain := score[1]

	extraversionScore := extraversionDomain.Score

	e1 := extraversionDomain.Subdomain[0].Score
	e2 := extraversionDomain.Subdomain[1].Score
	e3 := extraversionDomain.Subdomain[2].Score
	e4 := extraversionDomain.Subdomain[3].Score
	e5 := extraversionDomain.Subdomain[4].Score
	e6 := extraversionDomain.Subdomain[5].Score

	opennessDomain := score[2]

	opennessScore := opennessDomain.Score

	o1 := opennessDomain.Subdomain[0].Score
	o2 := opennessDomain.Subdomain[1].Score
	o3 := opennessDomain.Subdomain[2].Score
	o4 := opennessDomain.Subdomain[3].Score
	o5 := opennessDomain.Subdomain[4].Score
	o6 := opennessDomain.Subdomain[5].Score

	agreeablenessDomain := score[3]

	agreeablenessScore := agreeablenessDomain.Score

	a1 := agreeablenessDomain.Subdomain[0].Score
	a2 := agreeablenessDomain.Subdomain[1].Score
	a3 := agreeablenessDomain.Subdomain[2].Score
	a4 := agreeablenessDomain.Subdomain[3].Score
	a5 := agreeablenessDomain.Subdomain[4].Score
	a6 := agreeablenessDomain.Subdomain[5].Score

	conscientiousnessDomain := score[4]

	conscientiousnessScore := conscientiousnessDomain.Score

	c1 := conscientiousnessDomain.Subdomain[0].Score
	c2 := conscientiousnessDomain.Subdomain[1].Score
	c3 := conscientiousnessDomain.Subdomain[2].Score
	c4 := conscientiousnessDomain.Subdomain[3].Score
	c5 := conscientiousnessDomain.Subdomain[4].Score
	c6 := conscientiousnessDomain.Subdomain[5].Score

	prompt := fmt.Sprintf("Using the Big 5 score given below, create Relationship page for the Report\n\n"+

		"Domain Score to Intensity map { [0 - 30] : 'Low', [31 - 39]: 'Average',[40 - 60]: 'high'}"+
		"Subdomain Score to Intensity map { [0 - 4] : 'Low', [5 - 7]: 'Average',[8 - 10]: 'high'}"+

		"Personality Assessment:\n\n"+
		"Domain: Neuroticism: %v\n"+
		"Domain: Neuroticism: %s\n"+
		"  Subdomains-\n"+
		"    Anxiety: %s\n"+
		"    Anger: %s\n"+
		"    Depression: %s\n"+
		"    Self-consciousness: %s\n"+
		"    Immoderation: %s\n"+
		"    Vulnerability: %s\n\n"+

		"Domain: Extraversion: %s\n"+
		"  Subdomains-\n"+
		"    Friendliness: %s\n"+
		"    Gregariousness: %s\n"+
		"    Assertiveness: %s\n"+
		"    Activity Level: %s\n"+
		"    Excitement Seeking: %s\n"+
		"    Cheerfulness: %s\n\n"+

		"Domain: Openness: %s\n"+
		"  Subdomains-\n"+
		"    Imagination: %s\n"+
		"    Artistic Interests: %s\n"+
		"    Emotionality: %s\n"+
		"    Adventurousness: %s\n"+
		"    Intellect: %s\n"+
		"    Liberalism: %s\n\n"+

		"Domain: Agreeableness: %s\n"+
		"  Subdomains-\n"+
		"    Trust: %s\n"+
		"    Morality: %s\n"+
		"    Altruism: %s\n"+
		"    Cooperation: %s\n"+
		"    Modesty: %s\n"+
		"    Sympathy: %s\n\n"+

		"Domain: Conscientiousness: %s\n"+
		"  Subdomains-\n"+
		"    Self Efficacy: %s\n"+
		"    Orderliness: %s\n"+
		"    Dutifulness: %s\n"+
		"    Achievement Striving: %s\n"+
		"    Self Discipline: %s\n"+
		"    Cautiousness: %s\n",

		neuroticismScore, n1, n2, n3, n4, n5, n6,
		extraversionScore, e1, e2, e3, e4, e5, e6,
		opennessScore, o1, o2, o3, o4, o5, o6,
		agreeablenessScore, a1, a2, a3, a4, a5, a6,
		conscientiousnessScore, c1, c2, c3, c4, c5, c6)

	return prompt
}

func CreatePromptStrengthWeakness(score []Domain) string {

	neuroticismDomain := score[0]

	neuroticismScore := neuroticismDomain.Score

	n1 := neuroticismDomain.Subdomain[0].Score
	n2 := neuroticismDomain.Subdomain[1].Score
	n3 := neuroticismDomain.Subdomain[2].Score
	n4 := neuroticismDomain.Subdomain[3].Score
	n5 := neuroticismDomain.Subdomain[4].Score
	n6 := neuroticismDomain.Subdomain[5].Score

	extraversionDomain := score[1]

	extraversionScore := extraversionDomain.Score

	e1 := extraversionDomain.Subdomain[0].Score
	e2 := extraversionDomain.Subdomain[1].Score
	e3 := extraversionDomain.Subdomain[2].Score
	e4 := extraversionDomain.Subdomain[3].Score
	e5 := extraversionDomain.Subdomain[4].Score
	e6 := extraversionDomain.Subdomain[5].Score

	opennessDomain := score[2]

	opennessScore := opennessDomain.Score

	o1 := opennessDomain.Subdomain[0].Score
	o2 := opennessDomain.Subdomain[1].Score
	o3 := opennessDomain.Subdomain[2].Score
	o4 := opennessDomain.Subdomain[3].Score
	o5 := opennessDomain.Subdomain[4].Score
	o6 := opennessDomain.Subdomain[5].Score

	agreeablenessDomain := score[3]

	agreeablenessScore := agreeablenessDomain.Score

	a1 := agreeablenessDomain.Subdomain[0].Score
	a2 := agreeablenessDomain.Subdomain[1].Score
	a3 := agreeablenessDomain.Subdomain[2].Score
	a4 := agreeablenessDomain.Subdomain[3].Score
	a5 := agreeablenessDomain.Subdomain[4].Score
	a6 := agreeablenessDomain.Subdomain[5].Score

	conscientiousnessDomain := score[4]

	conscientiousnessScore := conscientiousnessDomain.Score

	c1 := conscientiousnessDomain.Subdomain[0].Score
	c2 := conscientiousnessDomain.Subdomain[1].Score
	c3 := conscientiousnessDomain.Subdomain[2].Score
	c4 := conscientiousnessDomain.Subdomain[3].Score
	c5 := conscientiousnessDomain.Subdomain[4].Score
	c6 := conscientiousnessDomain.Subdomain[5].Score

	prompt := fmt.Sprintf("Using the Big 5 score given below, create Strength & Weakness for the Report\n\n"+

		"Domain Score to Intensity map { [0 - 30] : 'Low', [31 - 39]: 'Average',[40 - 60]: 'high'}"+
		"Subdomain Score to Intensity map { [0 - 4] : 'Low', [5 - 7]: 'Average',[8 - 10]: 'high'}"+

		"Personality Assessment:\n\n"+
		"Domain: Neuroticism: %s\n"+
		"  Subdomains-\n"+
		"    Anxiety: %s\n"+
		"    Anger: %s\n"+
		"    Depression: %s\n"+
		"    Self-consciousness: %s\n"+
		"    Immoderation: %s\n"+
		"    Vulnerability: %s\n\n"+

		"Domain: Extraversion: %s\n"+
		"  Subdomains-\n"+
		"    Friendliness: %s\n"+
		"    Gregariousness: %s\n"+
		"    Assertiveness: %s\n"+
		"    Activity Level: %s\n"+
		"    Excitement Seeking: %s\n"+
		"    Cheerfulness: %s\n\n"+

		"Domain: Openness: %s\n"+
		"  Subdomains-\n"+
		"    Imagination: %s\n"+
		"    Artistic Interests: %s\n"+
		"    Emotionality: %s\n"+
		"    Adventurousness: %s\n"+
		"    Intellect: %s\n"+
		"    Liberalism: %s\n\n"+

		"Domain: Agreeableness: %s\n"+
		"  Subdomains-\n"+
		"    Trust: %s\n"+
		"    Morality: %s\n"+
		"    Altruism: %s\n"+
		"    Cooperation: %s\n"+
		"    Modesty: %s\n"+
		"    Sympathy: %s\n\n"+

		"Domain: Conscientiousness: %s\n"+
		"  Subdomains-\n"+
		"    Self Efficacy: %s\n"+
		"    Orderliness: %s\n"+
		"    Dutifulness: %s\n"+
		"    Achievement Striving: %s\n"+
		"    Self Discipline: %s\n"+
		"    Cautiousness: %s\n",

		neuroticismScore, n1, n2, n3, n4, n5, n6,
		extraversionScore, e1, e2, e3, e4, e5, e6,
		opennessScore, o1, o2, o3, o4, o5, o6,
		agreeablenessScore, a1, a2, a3, a4, a5, a6,
		conscientiousnessScore, c1, c2, c3, c4, c5, c6)

	return prompt
}

func WorkerGCPGemini(id string, prompt string, channel chan GeminiPromptRequest) {
	result, err := GenerateContentFromTextGCPJSON(prompt)
	if err != nil {
		// Respond with an error message if content generation failed
		log.Println("Error on generating content from GCP", err)

	}

	resultWithId := GeminiPromptRequest{id, *result}
	channel <- resultWithId
}
