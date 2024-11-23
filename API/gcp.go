package API

import (
	"bytes"
	//"context"
	"encoding/json"
	"fmt"
	//"github.com/google/generative-ai-go/genai"
	// "github.com/joho/godotenv"
	//"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"os"
	//"strings"
)

type GeminiPromptRequest struct {
	ID       string         `json:"id"`       // Adjusted field name capitalization
	Response OpenAIResponse `json:"response"` // Updated to OpenAI-specific response
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`    // "assistant" or "user"
			Content string `json:"content"` // Actual text response
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// type GeminiPromptRequest struct {
// 	Id       string
// 	Response ContentResponse
// }

// type ContentResponse struct {
// 	Candidates []struct {
// 		Content struct {
// 			Parts []struct {
// 				Text string `json:"text"`
// 			} `json:"parts"`
// 			Role string `json:"role"`
// 		} `json:"content"`
// 		FinishReason  string `json:"finishReason"`
// 		Index         int    `json:"index"`
// 		SafetyRatings []struct {
// 			Category    string `json:"category"`
// 			Probability string `json:"probability"`
// 		} `json:"safetyRatings"`
// 	} `json:"candidates"`
// 	UsageMetadata struct {
// 		PromptTokenCount     int `json:"promptTokenCount"`
// 		CandidatesTokenCount int `json:"candidatesTokenCount"`
// 		TotalTokenCount      int `json:"totalTokenCount"`
// 	} `json:"usageMetadata"`
// }

// // Function to parse markdown code and extract JSON
// func ParseMarkdownCode(markdown string) (string, error) {
// 	// Define the struct to store the parsed data
// 	var report string

// 	// Extract the JSON part (between the "```json" block)
// 	start := strings.Index(markdown, "```json")
// 	end := strings.LastIndex(markdown, "```")

// 	if start == -1 || end == -1 {
// 		return report, fmt.Errorf("invalid markdown format")
// 	}

// 	// // Extract the JSON part and trim spaces
// 	jsonPart := markdown[start+len("```json") : end]
// 	// jsonPart := strings.TrimSpace(markdown)

// 	// Remove all newlines
// 	jsonPart = strings.ReplaceAll(jsonPart, "\n", "")

// 	// Unmarshal the JSON into the struct
// 	err := json.Unmarshal([]byte(jsonPart), &report)
// 	if err != nil {
// 		return report, err
// 	}

// 	return report, nil
// }

// Function to make the POST request to Google API (REST-based approach)
// func GenerateContentFromTextGCP(prompt string) (string, error) {
// 	// Define the URL for the Google API endpoint (generative language model)
// 	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-pro-latest:generateContent"

// 	// Prepare the request payload in JSON format with the prompt content
// 	requestBody, err := json.Marshal(map[string]interface{}{
// 		"contents": []map[string]interface{}{
// 			{
// 				"parts": []map[string]string{
// 					{"text": prompt},
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		return "", fmt.Errorf("failed to marshal request body: %v", err)
// 	}

// 	// Retrieve the API key from the environment variable 'API_KEY'
// 	apiKey := os.Getenv("API_KEY")
// 	if apiKey == "" {
// 		return "", fmt.Errorf("API key is not set. Please ensure the environment variable 'API_KEY' is set")
// 	}

// 	// Create a new HTTP POST request with the request body and API key
// 	req, err := http.NewRequest("POST", url+"?key="+apiKey, bytes.NewBuffer(requestBody))
// 	if err != nil {
// 		return "", fmt.Errorf("failed to create request: %v", err)
// 	}

// 	// Set appropriate headers for the request
// 	req.Header.Set("Content-Type", "application/json")

// 	// Send the HTTP request using an HTTP client
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to send request: %v", err)
// 	}
// 	defer resp.Body.Close() // Ensure the response body is closed

// 	// Read and parse the response body
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to read response body: %v", err)
// 	}

// 	// Check if the response status is not 200 (OK), and handle the error accordingly
// 	if resp.StatusCode != http.StatusOK {
// 		return "", fmt.Errorf("non-200 response status: %d, body: %s", resp.StatusCode, string(body))
// 	}

// 	// Parse the JSON response into the struct
// 	var contentResponse ContentResponse
// 	err = json.Unmarshal(body, &contentResponse)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to unmarshal response body: %v", err)
// 	}

// 	// Return the response body as a string
// 	return string(body), nil
// }

func GenerateContentFromTextGCP(prompt string) (string, error) {
	// Define the URL for the OpenAI API endpoint
	url := "https://api.openai.com/v1/chat/completions"

	// Prepare the request payload in JSON format with the prompt content
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-4o", // Specify the OpenAI model
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Retrieve the API key from the environment variable 'OPENAI_API_KEY'
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("API key is not set. Please ensure the environment variable 'OPENAI_API_KEY' is set")
	}

	// Create a new HTTP POST request with the request body and API key
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set appropriate headers for the request
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

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
	var contentResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	err = json.Unmarshal(body, &contentResponse)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	// Extract the response content from the first choice
	if len(contentResponse.Choices) == 0 {
		return "", fmt.Errorf("no choices found in the response")
	}
	return contentResponse.Choices[0].Message.Content, nil
}

// func GenerateContentFromTextGOAPIJSON(promt string) (*genai.GenerateContentResponse, error) {
// 	ctx := context.Background()

// 	apiKey := os.Getenv("API_KEY")
// 	if apiKey == "" {
// 		return nil, fmt.Errorf("API key is not set. Please ensure the environment variable 'API_KEY' is set")
// 	}

// 	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Close()

// 	model := client.GenerativeModel("gemini-1.5-pro-latest")
// 	// Ask the model to respond with JSON.
// 	model.ResponseMIMEType = "application/json"
// 	// Specify the schema.
// 	model.ResponseSchema = &genai.Schema{
// 		Type:  genai.TypeArray,
// 		Items: &genai.Schema{Type: genai.TypeString},
// 	}
// 	resp, err := model.GenerateContent(ctx, genai.Text(promt))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for _, part := range resp.Candidates[0].Content.Parts {
// 		if txt, ok := part.(genai.Text); ok {
// 			var recipes []string
// 			if err := json.Unmarshal([]byte(txt), &recipes); err != nil {
// 				log.Fatal(err)
// 			}
// 			fmt.Println(recipes)
// 		}
// 	}

// 	return resp, err
// }

// // Function to make the POST request to Google API (REST-based approach)
// func GenerateContentFromTextGCPJSON(prompt string) (*ContentResponse, error) {

// 	// Define the URL for the Google API endpoint (generative language model)
// 	url := ""

// 	// Prepare the request payload in JSON format with the prompt content
// 	requestBody, err := json.Marshal(map[string]interface{}{
// 		"contents": []map[string]interface{}{
// 			{
// 				"parts": []map[string]string{
// 					{"text": prompt},
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to marshal request body: %v", err)
// 	}

// 	// Retrieve the API key from the environment variable 'API_KEY'
// 	apiKey := os.Getenv("API_KEY")
// 	if apiKey == "" {
// 		return nil, fmt.Errorf("API key is not set. Please ensure the environment variable 'API_KEY' is set")
// 	}

// 	// Create a new HTTP POST request with the request body and API key
// 	req, err := http.NewRequest("POST", url+"?key="+apiKey, bytes.NewBuffer(requestBody))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create request: %v", err)
// 	}

// 	// Set appropriate headers for the request
// 	req.Header.Set("Content-Type", "application/json")

// 	// Send the HTTP request using an HTTP client
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to send request: %v", err)
// 	}
// 	defer resp.Body.Close() // Ensure the response body is closed

// 	// Read and parse the response body
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read response body: %v", err)
// 	}

// 	// Check if the response status is not 200 (OK), and handle the error accordingly
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("non-200 response status: %d, body: %s", resp.StatusCode, string(body))
// 	}

// 	// Parse the JSON response into the struct
// 	var contentResponse ContentResponse
// 	err = json.Unmarshal(body, &contentResponse)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
// 	}

// 	// Return the response body as a string
// 	return &contentResponse, nil
// }

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
	neuroticismIntensity := neuroticismDomain.Intensity

	n1 := neuroticismDomain.Subdomain[0].Score
	n2 := neuroticismDomain.Subdomain[1].Score
	n3 := neuroticismDomain.Subdomain[2].Score
	n4 := neuroticismDomain.Subdomain[3].Score
	n5 := neuroticismDomain.Subdomain[4].Score
	n6 := neuroticismDomain.Subdomain[5].Score

	n1I := neuroticismDomain.Subdomain[0].Intensity
	n2I := neuroticismDomain.Subdomain[1].Intensity
	n3I := neuroticismDomain.Subdomain[2].Intensity
	n4I := neuroticismDomain.Subdomain[3].Intensity
	n5I := neuroticismDomain.Subdomain[4].Intensity
	n6I := neuroticismDomain.Subdomain[5].Intensity

	extraversionDomain := score[1]

	extraversionScore := extraversionDomain.Score
	extraversionIntensity := extraversionDomain.Intensity

	e1 := extraversionDomain.Subdomain[0].Score
	e2 := extraversionDomain.Subdomain[1].Score
	e3 := extraversionDomain.Subdomain[2].Score
	e4 := extraversionDomain.Subdomain[3].Score
	e5 := extraversionDomain.Subdomain[4].Score
	e6 := extraversionDomain.Subdomain[5].Score

	e1I := extraversionDomain.Subdomain[0].Intensity
	e2I := extraversionDomain.Subdomain[1].Intensity
	e3I := extraversionDomain.Subdomain[2].Intensity
	e4I := extraversionDomain.Subdomain[3].Intensity
	e5I := extraversionDomain.Subdomain[4].Intensity
	e6I := extraversionDomain.Subdomain[5].Intensity

	opennessDomain := score[2]

	opennessScore := opennessDomain.Score
	opennessIntensity := opennessDomain.Intensity

	o1 := opennessDomain.Subdomain[0].Score
	o2 := opennessDomain.Subdomain[1].Score
	o3 := opennessDomain.Subdomain[2].Score
	o4 := opennessDomain.Subdomain[3].Score
	o5 := opennessDomain.Subdomain[4].Score
	o6 := opennessDomain.Subdomain[5].Score

	o1I := opennessDomain.Subdomain[0].Intensity
	o2I := opennessDomain.Subdomain[1].Intensity
	o3I := opennessDomain.Subdomain[2].Intensity
	o4I := opennessDomain.Subdomain[3].Intensity
	o5I := opennessDomain.Subdomain[4].Intensity
	o6I := opennessDomain.Subdomain[5].Intensity

	agreeablenessDomain := score[3]

	agreeablenessScore := agreeablenessDomain.Score
	agreeablenessIntensity := agreeablenessDomain.Intensity

	a1 := agreeablenessDomain.Subdomain[0].Score
	a2 := agreeablenessDomain.Subdomain[1].Score
	a3 := agreeablenessDomain.Subdomain[2].Score
	a4 := agreeablenessDomain.Subdomain[3].Score
	a5 := agreeablenessDomain.Subdomain[4].Score
	a6 := agreeablenessDomain.Subdomain[5].Score

	a1I := agreeablenessDomain.Subdomain[0].Intensity
	a2I := agreeablenessDomain.Subdomain[1].Intensity
	a3I := agreeablenessDomain.Subdomain[2].Intensity
	a4I := agreeablenessDomain.Subdomain[3].Intensity
	a5I := agreeablenessDomain.Subdomain[4].Intensity
	a6I := agreeablenessDomain.Subdomain[5].Intensity

	conscientiousnessDomain := score[4]

	conscientiousnessScore := conscientiousnessDomain.Score
	conscientiousnessIntensity := conscientiousnessDomain.Intensity

	c1 := conscientiousnessDomain.Subdomain[0].Score
	c2 := conscientiousnessDomain.Subdomain[1].Score
	c3 := conscientiousnessDomain.Subdomain[2].Score
	c4 := conscientiousnessDomain.Subdomain[3].Score
	c5 := conscientiousnessDomain.Subdomain[4].Score
	c6 := conscientiousnessDomain.Subdomain[5].Score

	c1I := conscientiousnessDomain.Subdomain[0].Intensity
	c2I := conscientiousnessDomain.Subdomain[1].Intensity
	c3I := conscientiousnessDomain.Subdomain[2].Intensity
	c4I := conscientiousnessDomain.Subdomain[3].Intensity
	c5I := conscientiousnessDomain.Subdomain[4].Intensity
	c6I := conscientiousnessDomain.Subdomain[5].Intensity

	prompt := fmt.Sprintf("Using the Big 5 Assessment score given below, create Summary for each domain in around 300-400 words in total"+

		"Domain: Neuroticism Score: %s/60 (%s)\n"+
		"  Subdomains of Neuroticism-\n"+
		"    Anxiety Score : %s\n"+
		"    Anger Score: %s\n"+
		"    Depression Score: %s\n"+
		"    Self-consciousness Score: %s\n"+
		"    Immoderation Score: %s\n"+
		"    Vulnerability Score: %s\n\n"+

		"Domain: Extraversion Score: %s/60 (%s)\n"+
		"  Subdomains of Extraversion-\n"+
		"    Friendliness Score: %s\n"+
		"    Gregariousness Score: %s\n"+
		"    Assertiveness Score: %s\n"+
		"    Activity Level Score: %s\n"+
		"    Excitement Seeking Score: %s\n"+
		"    Cheerfulness Score: %s\n\n"+

		"Domain: Openness: %s/60 (%s)\n"+
		"  Subdomains of Openness-\n"+
		"    Imagination Score: %s\n"+
		"    Artistic Interests Score: %s\n"+
		"    Emotionality Score: %s\n"+
		"    Adventurousness Score: %s\n"+
		"    Intellect Score: %s\n"+
		"    Liberalism Score: %s\n\n"+

		"Domain: Agreeableness: %s/60 (%s)\n"+
		"  Subdomains of Agreeableness-\n"+
		"    Trust Score: %s\n"+
		"    Morality Score: %s\n"+
		"    Altruism Score: %s\n"+
		"    Cooperation Score: %s\n"+
		"    Modesty Score: %s\n"+
		"    Sympathy Score: %s\n\n"+

		"Domain: Conscientiousness: %s/60 (%s)\n"+
		"  Subdomains of Conscientiousness-\n"+
		"    Self Efficacy Score: %s\n"+
		"    Orderliness Score: %s\n"+
		"    Dutifulness Score: %s\n"+
		"    Achievement Striving Score: %s\n"+
		"    Self Discipline Score: %s\n"+
		"    Cautiousness Score: %s\n",

		neuroticismScore, neuroticismIntensity, n1I, n2I, n3I, n4I, n5I, n6I,
		extraversionScore, extraversionIntensity, e1I, e2I, e3I, e4I, e5I, e6I,
		opennessScore, opennessIntensity, o1I, o2I, o3I, o4I, o5I, o6I,
		agreeablenessScore, agreeablenessIntensity, a1I, a2I, a3I, a4I, a5I, a6I,
		conscientiousnessScore, conscientiousnessIntensity, c1I, c2I, c3I, c4I, c5I, c6I)

	return prompt
}

func CreatePromptCareerAcademic(score []Domain) string {

	neuroticismDomain := score[0]

	neuroticismScore := neuroticismDomain.Score
	neuroticismIntensity := neuroticismDomain.Intensity

	n1 := neuroticismDomain.Subdomain[0].Score
	n2 := neuroticismDomain.Subdomain[1].Score
	n3 := neuroticismDomain.Subdomain[2].Score
	n4 := neuroticismDomain.Subdomain[3].Score
	n5 := neuroticismDomain.Subdomain[4].Score
	n6 := neuroticismDomain.Subdomain[5].Score

	n1I := neuroticismDomain.Subdomain[0].Intensity
	n2I := neuroticismDomain.Subdomain[1].Intensity
	n3I := neuroticismDomain.Subdomain[2].Intensity
	n4I := neuroticismDomain.Subdomain[3].Intensity
	n5I := neuroticismDomain.Subdomain[4].Intensity
	n6I := neuroticismDomain.Subdomain[5].Intensity

	extraversionDomain := score[1]

	extraversionScore := extraversionDomain.Score
	extraversionIntensity := extraversionDomain.Intensity

	e1 := extraversionDomain.Subdomain[0].Score
	e2 := extraversionDomain.Subdomain[1].Score
	e3 := extraversionDomain.Subdomain[2].Score
	e4 := extraversionDomain.Subdomain[3].Score
	e5 := extraversionDomain.Subdomain[4].Score
	e6 := extraversionDomain.Subdomain[5].Score

	e1I := extraversionDomain.Subdomain[0].Intensity
	e2I := extraversionDomain.Subdomain[1].Intensity
	e3I := extraversionDomain.Subdomain[2].Intensity
	e4I := extraversionDomain.Subdomain[3].Intensity
	e5I := extraversionDomain.Subdomain[4].Intensity
	e6I := extraversionDomain.Subdomain[5].Intensity

	opennessDomain := score[2]

	opennessScore := opennessDomain.Score
	opennessIntensity := opennessDomain.Intensity

	o1 := opennessDomain.Subdomain[0].Score
	o2 := opennessDomain.Subdomain[1].Score
	o3 := opennessDomain.Subdomain[2].Score
	o4 := opennessDomain.Subdomain[3].Score
	o5 := opennessDomain.Subdomain[4].Score
	o6 := opennessDomain.Subdomain[5].Score

	o1I := opennessDomain.Subdomain[0].Intensity
	o2I := opennessDomain.Subdomain[1].Intensity
	o3I := opennessDomain.Subdomain[2].Intensity
	o4I := opennessDomain.Subdomain[3].Intensity
	o5I := opennessDomain.Subdomain[4].Intensity
	o6I := opennessDomain.Subdomain[5].Intensity

	agreeablenessDomain := score[3]

	agreeablenessScore := agreeablenessDomain.Score
	agreeablenessIntensity := agreeablenessDomain.Intensity

	a1 := agreeablenessDomain.Subdomain[0].Score
	a2 := agreeablenessDomain.Subdomain[1].Score
	a3 := agreeablenessDomain.Subdomain[2].Score
	a4 := agreeablenessDomain.Subdomain[3].Score
	a5 := agreeablenessDomain.Subdomain[4].Score
	a6 := agreeablenessDomain.Subdomain[5].Score

	a1I := agreeablenessDomain.Subdomain[0].Intensity
	a2I := agreeablenessDomain.Subdomain[1].Intensity
	a3I := agreeablenessDomain.Subdomain[2].Intensity
	a4I := agreeablenessDomain.Subdomain[3].Intensity
	a5I := agreeablenessDomain.Subdomain[4].Intensity
	a6I := agreeablenessDomain.Subdomain[5].Intensity

	conscientiousnessDomain := score[4]

	conscientiousnessScore := conscientiousnessDomain.Score
	conscientiousnessIntensity := conscientiousnessDomain.Intensity

	c1 := conscientiousnessDomain.Subdomain[0].Score
	c2 := conscientiousnessDomain.Subdomain[1].Score
	c3 := conscientiousnessDomain.Subdomain[2].Score
	c4 := conscientiousnessDomain.Subdomain[3].Score
	c5 := conscientiousnessDomain.Subdomain[4].Score
	c6 := conscientiousnessDomain.Subdomain[5].Score

	c1I := conscientiousnessDomain.Subdomain[0].Intensity
	c2I := conscientiousnessDomain.Subdomain[1].Intensity
	c3I := conscientiousnessDomain.Subdomain[2].Intensity
	c4I := conscientiousnessDomain.Subdomain[3].Intensity
	c5I := conscientiousnessDomain.Subdomain[4].Intensity
	c6I := conscientiousnessDomain.Subdomain[5].Intensity

	prompt := fmt.Sprintf("Using the Big 5 Assessment score given below, create Career & Academia Page under 200 words for the Report\n\n"+

		"Domain: Neuroticism Score: %s/60 (%s)\n"+
		"  Subdomains of Neuroticism-\n"+
		"    Anxiety Score : %s\n"+
		"    Anger Score: %s\n"+
		"    Depression Score: %s\n"+
		"    Self-consciousness Score: %s\n"+
		"    Immoderation Score: %s\n"+
		"    Vulnerability Score: %s\n\n"+

		"Domain: Extraversion Score: %s/60 (%s)\n"+
		"  Subdomains of Extraversion-\n"+
		"    Friendliness Score: %s\n"+
		"    Gregariousness Score: %s\n"+
		"    Assertiveness Score: %s\n"+
		"    Activity Level Score: %s\n"+
		"    Excitement Seeking Score: %s\n"+
		"    Cheerfulness Score: %s\n\n"+

		"Domain: Openness: %s/60 (%s)\n"+
		"  Subdomains of Openness-\n"+
		"    Imagination Score: %s\n"+
		"    Artistic Interests Score: %s\n"+
		"    Emotionality Score: %s\n"+
		"    Adventurousness Score: %s\n"+
		"    Intellect Score: %s\n"+
		"    Liberalism Score: %s\n\n"+

		"Domain: Agreeableness: %s/60 (%s)\n"+
		"  Subdomains of Agreeableness-\n"+
		"    Trust Score: %s\n"+
		"    Morality Score: %s\n"+
		"    Altruism Score: %s\n"+
		"    Cooperation Score: %s\n"+
		"    Modesty Score: %s\n"+
		"    Sympathy Score: %s\n\n"+

		"Domain: Conscientiousness: %s/60 (%s)\n"+
		"  Subdomains of Conscientiousness-\n"+
		"    Self Efficacy Score: %s\n"+
		"    Orderliness Score: %s\n"+
		"    Dutifulness Score: %s\n"+
		"    Achievement Striving Score: %s\n"+
		"    Self Discipline Score: %s\n"+
		"    Cautiousness Score: %s\n",

		neuroticismScore, neuroticismIntensity, n1I, n2I, n3I, n4I, n5I, n6I,
		extraversionScore, extraversionIntensity, e1I, e2I, e3I, e4I, e5I, e6I,
		opennessScore, opennessIntensity, o1I, o2I, o3I, o4I, o5I, o6I,
		agreeablenessScore, agreeablenessIntensity, a1I, a2I, a3I, a4I, a5I, a6I,
		conscientiousnessScore, conscientiousnessIntensity, c1I, c2I, c3I, c4I, c5I, c6I)

	return prompt
}

func CreatePromptRelationship(score []Domain) string {

	neuroticismDomain := score[0]

	neuroticismScore := neuroticismDomain.Score
	neuroticismIntensity := neuroticismDomain.Intensity

	n1 := neuroticismDomain.Subdomain[0].Score
	n2 := neuroticismDomain.Subdomain[1].Score
	n3 := neuroticismDomain.Subdomain[2].Score
	n4 := neuroticismDomain.Subdomain[3].Score
	n5 := neuroticismDomain.Subdomain[4].Score
	n6 := neuroticismDomain.Subdomain[5].Score

	n1I := neuroticismDomain.Subdomain[0].Intensity
	n2I := neuroticismDomain.Subdomain[1].Intensity
	n3I := neuroticismDomain.Subdomain[2].Intensity
	n4I := neuroticismDomain.Subdomain[3].Intensity
	n5I := neuroticismDomain.Subdomain[4].Intensity
	n6I := neuroticismDomain.Subdomain[5].Intensity

	extraversionDomain := score[1]

	extraversionScore := extraversionDomain.Score
	extraversionIntensity := extraversionDomain.Intensity

	e1 := extraversionDomain.Subdomain[0].Score
	e2 := extraversionDomain.Subdomain[1].Score
	e3 := extraversionDomain.Subdomain[2].Score
	e4 := extraversionDomain.Subdomain[3].Score
	e5 := extraversionDomain.Subdomain[4].Score
	e6 := extraversionDomain.Subdomain[5].Score

	e1I := extraversionDomain.Subdomain[0].Intensity
	e2I := extraversionDomain.Subdomain[1].Intensity
	e3I := extraversionDomain.Subdomain[2].Intensity
	e4I := extraversionDomain.Subdomain[3].Intensity
	e5I := extraversionDomain.Subdomain[4].Intensity
	e6I := extraversionDomain.Subdomain[5].Intensity

	opennessDomain := score[2]

	opennessScore := opennessDomain.Score
	opennessIntensity := opennessDomain.Intensity

	o1 := opennessDomain.Subdomain[0].Score
	o2 := opennessDomain.Subdomain[1].Score
	o3 := opennessDomain.Subdomain[2].Score
	o4 := opennessDomain.Subdomain[3].Score
	o5 := opennessDomain.Subdomain[4].Score
	o6 := opennessDomain.Subdomain[5].Score

	o1I := opennessDomain.Subdomain[0].Intensity
	o2I := opennessDomain.Subdomain[1].Intensity
	o3I := opennessDomain.Subdomain[2].Intensity
	o4I := opennessDomain.Subdomain[3].Intensity
	o5I := opennessDomain.Subdomain[4].Intensity
	o6I := opennessDomain.Subdomain[5].Intensity

	agreeablenessDomain := score[3]

	agreeablenessScore := agreeablenessDomain.Score
	agreeablenessIntensity := agreeablenessDomain.Intensity

	a1 := agreeablenessDomain.Subdomain[0].Score
	a2 := agreeablenessDomain.Subdomain[1].Score
	a3 := agreeablenessDomain.Subdomain[2].Score
	a4 := agreeablenessDomain.Subdomain[3].Score
	a5 := agreeablenessDomain.Subdomain[4].Score
	a6 := agreeablenessDomain.Subdomain[5].Score

	a1I := agreeablenessDomain.Subdomain[0].Intensity
	a2I := agreeablenessDomain.Subdomain[1].Intensity
	a3I := agreeablenessDomain.Subdomain[2].Intensity
	a4I := agreeablenessDomain.Subdomain[3].Intensity
	a5I := agreeablenessDomain.Subdomain[4].Intensity
	a6I := agreeablenessDomain.Subdomain[5].Intensity

	conscientiousnessDomain := score[4]

	conscientiousnessScore := conscientiousnessDomain.Score
	conscientiousnessIntensity := conscientiousnessDomain.Intensity

	c1 := conscientiousnessDomain.Subdomain[0].Score
	c2 := conscientiousnessDomain.Subdomain[1].Score
	c3 := conscientiousnessDomain.Subdomain[2].Score
	c4 := conscientiousnessDomain.Subdomain[3].Score
	c5 := conscientiousnessDomain.Subdomain[4].Score
	c6 := conscientiousnessDomain.Subdomain[5].Score

	c1I := conscientiousnessDomain.Subdomain[0].Intensity
	c2I := conscientiousnessDomain.Subdomain[1].Intensity
	c3I := conscientiousnessDomain.Subdomain[2].Intensity
	c4I := conscientiousnessDomain.Subdomain[3].Intensity
	c5I := conscientiousnessDomain.Subdomain[4].Intensity
	c6I := conscientiousnessDomain.Subdomain[5].Intensity

	prompt := fmt.Sprintf("Using the Big 5 Assessment score given below, create Relationship page under 200 words for the Report\n\n"+

		"Domain: Neuroticism Score: %s/60 (%s)\n"+
		"  Subdomains of Neuroticism-\n"+
		"    Anxiety Score : %s\n"+
		"    Anger Score: %s\n"+
		"    Depression Score: %s\n"+
		"    Self-consciousness Score: %s\n"+
		"    Immoderation Score: %s\n"+
		"    Vulnerability Score: %s\n\n"+

		"Domain: Extraversion Score: %s/60 (%s)\n"+
		"  Subdomains of Extraversion-\n"+
		"    Friendliness Score: %s\n"+
		"    Gregariousness Score: %s\n"+
		"    Assertiveness Score: %s\n"+
		"    Activity Level Score: %s\n"+
		"    Excitement Seeking Score: %s\n"+
		"    Cheerfulness Score: %s\n\n"+

		"Domain: Openness: %s/60 (%s)\n"+
		"  Subdomains of Openness-\n"+
		"    Imagination Score: %s\n"+
		"    Artistic Interests Score: %s\n"+
		"    Emotionality Score: %s\n"+
		"    Adventurousness Score: %s\n"+
		"    Intellect Score: %s\n"+
		"    Liberalism Score: %s\n\n"+

		"Domain: Agreeableness: %s/60 (%s)\n"+
		"  Subdomains of Agreeableness-\n"+
		"    Trust Score: %s\n"+
		"    Morality Score: %s\n"+
		"    Altruism Score: %s\n"+
		"    Cooperation Score: %s\n"+
		"    Modesty Score: %s\n"+
		"    Sympathy Score: %s\n\n"+

		"Domain: Conscientiousness: %s/60 (%s)\n"+
		"  Subdomains of Conscientiousness-\n"+
		"    Self Efficacy Score: %s\n"+
		"    Orderliness Score: %s\n"+
		"    Dutifulness Score: %s\n"+
		"    Achievement Striving Score: %s\n"+
		"    Self Discipline Score: %s\n"+
		"    Cautiousness Score: %s\n",

		neuroticismScore, neuroticismIntensity, n1I, n2I, n3I, n4I, n5I, n6I,
		extraversionScore, extraversionIntensity, e1I, e2I, e3I, e4I, e5I, e6I,
		opennessScore, opennessIntensity, o1I, o2I, o3I, o4I, o5I, o6I,
		agreeablenessScore, agreeablenessIntensity, a1I, a2I, a3I, a4I, a5I, a6I,
		conscientiousnessScore, conscientiousnessIntensity, c1I, c2I, c3I, c4I, c5I, c6I)

	return prompt
}

func CreatePromptStrengthWeakness(score []Domain) string {

	neuroticismDomain := score[0]

	neuroticismScore := neuroticismDomain.Score
	neuroticismIntensity := neuroticismDomain.Intensity

	n1 := neuroticismDomain.Subdomain[0].Score
	n2 := neuroticismDomain.Subdomain[1].Score
	n3 := neuroticismDomain.Subdomain[2].Score
	n4 := neuroticismDomain.Subdomain[3].Score
	n5 := neuroticismDomain.Subdomain[4].Score
	n6 := neuroticismDomain.Subdomain[5].Score

	n1I := neuroticismDomain.Subdomain[0].Intensity
	n2I := neuroticismDomain.Subdomain[1].Intensity
	n3I := neuroticismDomain.Subdomain[2].Intensity
	n4I := neuroticismDomain.Subdomain[3].Intensity
	n5I := neuroticismDomain.Subdomain[4].Intensity
	n6I := neuroticismDomain.Subdomain[5].Intensity

	extraversionDomain := score[1]

	extraversionScore := extraversionDomain.Score
	extraversionIntensity := extraversionDomain.Intensity

	e1 := extraversionDomain.Subdomain[0].Score
	e2 := extraversionDomain.Subdomain[1].Score
	e3 := extraversionDomain.Subdomain[2].Score
	e4 := extraversionDomain.Subdomain[3].Score
	e5 := extraversionDomain.Subdomain[4].Score
	e6 := extraversionDomain.Subdomain[5].Score

	e1I := extraversionDomain.Subdomain[0].Intensity
	e2I := extraversionDomain.Subdomain[1].Intensity
	e3I := extraversionDomain.Subdomain[2].Intensity
	e4I := extraversionDomain.Subdomain[3].Intensity
	e5I := extraversionDomain.Subdomain[4].Intensity
	e6I := extraversionDomain.Subdomain[5].Intensity

	opennessDomain := score[2]

	opennessScore := opennessDomain.Score
	opennessIntensity := opennessDomain.Intensity

	o1 := opennessDomain.Subdomain[0].Score
	o2 := opennessDomain.Subdomain[1].Score
	o3 := opennessDomain.Subdomain[2].Score
	o4 := opennessDomain.Subdomain[3].Score
	o5 := opennessDomain.Subdomain[4].Score
	o6 := opennessDomain.Subdomain[5].Score

	o1I := opennessDomain.Subdomain[0].Intensity
	o2I := opennessDomain.Subdomain[1].Intensity
	o3I := opennessDomain.Subdomain[2].Intensity
	o4I := opennessDomain.Subdomain[3].Intensity
	o5I := opennessDomain.Subdomain[4].Intensity
	o6I := opennessDomain.Subdomain[5].Intensity

	agreeablenessDomain := score[3]

	agreeablenessScore := agreeablenessDomain.Score
	agreeablenessIntensity := agreeablenessDomain.Intensity

	a1 := agreeablenessDomain.Subdomain[0].Score
	a2 := agreeablenessDomain.Subdomain[1].Score
	a3 := agreeablenessDomain.Subdomain[2].Score
	a4 := agreeablenessDomain.Subdomain[3].Score
	a5 := agreeablenessDomain.Subdomain[4].Score
	a6 := agreeablenessDomain.Subdomain[5].Score

	a1I := agreeablenessDomain.Subdomain[0].Intensity
	a2I := agreeablenessDomain.Subdomain[1].Intensity
	a3I := agreeablenessDomain.Subdomain[2].Intensity
	a4I := agreeablenessDomain.Subdomain[3].Intensity
	a5I := agreeablenessDomain.Subdomain[4].Intensity
	a6I := agreeablenessDomain.Subdomain[5].Intensity

	conscientiousnessDomain := score[4]

	conscientiousnessScore := conscientiousnessDomain.Score
	conscientiousnessIntensity := conscientiousnessDomain.Intensity

	c1 := conscientiousnessDomain.Subdomain[0].Score
	c2 := conscientiousnessDomain.Subdomain[1].Score
	c3 := conscientiousnessDomain.Subdomain[2].Score
	c4 := conscientiousnessDomain.Subdomain[3].Score
	c5 := conscientiousnessDomain.Subdomain[4].Score
	c6 := conscientiousnessDomain.Subdomain[5].Score

	c1I := conscientiousnessDomain.Subdomain[0].Intensity
	c2I := conscientiousnessDomain.Subdomain[1].Intensity
	c3I := conscientiousnessDomain.Subdomain[2].Intensity
	c4I := conscientiousnessDomain.Subdomain[3].Intensity
	c5I := conscientiousnessDomain.Subdomain[4].Intensity
	c6I := conscientiousnessDomain.Subdomain[5].Intensity

	prompt := fmt.Sprintf("Using the Big 5 Assessment score given below, create Strength & Weakness page under 200 words for the Report\n\n"+

		"Domain: Neuroticism Score: %s/60 (%s)\n"+
		"  Subdomains of Neuroticism-\n"+
		"    Anxiety Score : %s\n"+
		"    Anger Score: %s\n"+
		"    Depression Score: %s\n"+
		"    Self-consciousness Score: %s\n"+
		"    Immoderation Score: %s\n"+
		"    Vulnerability Score: %s\n\n"+

		"Domain: Extraversion Score: %s/60 (%s)\n"+
		"  Subdomains of Extraversion-\n"+
		"    Friendliness Score: %s\n"+
		"    Gregariousness Score: %s\n"+
		"    Assertiveness Score: %s\n"+
		"    Activity Level Score: %s\n"+
		"    Excitement Seeking Score: %s\n"+
		"    Cheerfulness Score: %s\n\n"+

		"Domain: Openness: %s/60 (%s)\n"+
		"  Subdomains of Openness-\n"+
		"    Imagination Score: %s\n"+
		"    Artistic Interests Score: %s\n"+
		"    Emotionality Score: %s\n"+
		"    Adventurousness Score: %s\n"+
		"    Intellect Score: %s\n"+
		"    Liberalism Score: %s\n\n"+

		"Domain: Agreeableness: %s/60 (%s)\n"+
		"  Subdomains of Agreeableness-\n"+
		"    Trust Score: %s\n"+
		"    Morality Score: %s\n"+
		"    Altruism Score: %s\n"+
		"    Cooperation Score: %s\n"+
		"    Modesty Score: %s\n"+
		"    Sympathy Score: %s\n\n"+

		"Domain: Conscientiousness: %s/60 (%s)\n"+
		"  Subdomains of Conscientiousness-\n"+
		"    Self Efficacy Score: %s\n"+
		"    Orderliness Score: %s\n"+
		"    Dutifulness Score: %s\n"+
		"    Achievement Striving Score: %s\n"+
		"    Self Discipline Score: %s\n"+
		"    Cautiousness Score: %s\n",

		neuroticismScore, neuroticismIntensity, n1I, n2I, n3I, n4I, n5I, n6I,
		extraversionScore, extraversionIntensity, e1I, e2I, e3I, e4I, e5I, e6I,
		opennessScore, opennessIntensity, o1I, o2I, o3I, o4I, o5I, o6I,
		agreeablenessScore, agreeablenessIntensity, a1I, a2I, a3I, a4I, a5I, a6I,
		conscientiousnessScore, conscientiousnessIntensity, c1I, c2I, c3I, c4I, c5I, c6I)

	return prompt
}

// func WorkerGCPGemini(id string, prompt string, channel chan GeminiPromptRequest) {
// 	result, err := GenerateContentFromTextGCPJSON(prompt)
// 	if err != nil {
// 		// Respond with an error message if content generation failed
// 		log.Println("Error on generating content from GCP", err)

// 	}

//		resultWithId := GeminiPromptRequest{id, *result}
//		channel <- resultWithId
//	}

func WorkerOpenAIGPT(id string, prompt string, channel chan GeminiPromptRequest) {
	// Call the OpenAI API
	result, err := GenerateContentFromTextGCP(prompt)
	if err != nil {
		log.Printf("Error generating content from OpenAI for ID %s: %v", id, err)
		channel <- GeminiPromptRequest{
			ID: id,
			Response: OpenAIResponse{
				Choices: []struct {
					Message struct {
						Role    string `json:"role"`
						Content string `json:"content"`
					} `json:"message"`
					FinishReason string `json:"finish_reason"`
					Index        int    `json:"index"`
				}{
					{
						Message: struct {
							Role    string `json:"role"`
							Content string `json:"content"`
						}{
							Role:    "system",
							Content: fmt.Sprintf("Error: %v", err),
						},
						FinishReason: "error",
						Index:        0,
					},
				},
			},
		}
		return
	}

	// Construct the GeminiPromptRequest
	response := GeminiPromptRequest{
		ID: id,
		Response: OpenAIResponse{
			Choices: []struct {
				Message struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				} `json:"message"`
				FinishReason string `json:"finish_reason"`
				Index        int    `json:"index"`
			}{
				{
					Message: struct {
						Role    string `json:"role"`
						Content string `json:"content"`
					}{
						Role:    "assistant",
						Content: result,
					},
					FinishReason: "stop",
					Index:        0,
				},
			},
			Usage: struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
				TotalTokens      int `json:"total_tokens"`
			}{
				PromptTokens:     0, // Replace with actual token counts if available
				CompletionTokens: 0,
				TotalTokens:      0,
			},
		},
	}

	// Send response to the channel
	channel <- response
}
