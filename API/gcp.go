package API

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"os"
)

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

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
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

func CreatePrompt(d string, ds string, s1 string, s2 string, s3 string, s4 string, s5 string, s6 string) string {
	switch d {
	case "neuroticism":
		return CreatePromptNeuroticism(d, ds, s1, s2, s3, s4, s5, s6)
	case "extraversion":
		return CreatePromptExtraversion(d, ds, s1, s2, s3, s4, s5, s6)
	case "openness":
		return CreatePromptOpenness(d, ds, s1, s2, s3, s4, s5, s6)
	case "agreeableness":
		return CreatePromptAgreeableness(d, ds, s1, s2, s3, s4, s5, s6)
	case "conscientiousness":
		return CreatePromptConscientiousness(d, ds, s1, s2, s3, s4, s5, s6)
	}
	return ""
}

func CreatePromptNeuroticism(d string, ds string, s1 string, s2 string, s3 string, s4 string, s5 string, s6 string) string {
	prompt :=

		"Domain: Neuroticism:" + ds + "/60" + "\n" +
			"Subdomains-\n" +
			"  Anxiety: " + s1 + "/10" + "\n" +
			"  Anger: " + s2 + "/10" + "\n" +
			"  Depression: " + s3 + "/10" + "\n" +
			"  Self-consciousness: " + s4 + "/10" + "\n" +
			"  Immoderation: " + s5 + "/10" + "\n" +
			"  Vulnerability: " + s6 + "/10" + "\n\n" +

			"Create a personalised BIG5 Personlaity Assessment Report in JSON format for the Domain:" + d + ", while taking insight from subdomain score\n" +
			"Keept the Structure as follows:\n" +
			"Introduction in 100 words : Explain the trait and its impact on the client's experiences\n" +
			"Career & Academia in 40 words : Impact on clinet's professional & student life\n" +
			"Relationship in 40 words : Impact on Client's Personal Relationships\n" +
			"Strength & Weakness in 40 words : Highlight the client's strengths and areas for growth, focusing on positivity and potential.\n\n"

	return prompt
}

func CreatePromptExtraversion(d string, ds string, s1 string, s2 string, s3 string, s4 string, s5 string, s6 string) string {
	prompt :=

		"Domain: Extraversion: " + ds + "\n" +
			"Subdomains-\n" +
			"  Friendliness: " + s1 + "\n" +
			"  Gregariousness: " + s2 + "\n" +
			"  Assertiveness: " + s3 + "\n" +
			"  Activity Level: " + s4 + "\n" +
			"  Excitement Seeking: " + s5 + "\n" +
			"  Cheerfulness: " + s6 + "\n\n" +

			"Create a personalised BIG5 Personlaity Assessment Report in JSON format for the Domain:" + d + ", while taking insight from subdomain score\n" +
			"Keept the Structure as follows:\n" +
			"Introduction in 100 words : Explain the trait and its impact on the client's experiences\n" +
			"Career & Academia in 40 words : Impact on clinet's professional & student life\n" +
			"Relationship in 40 words : Impact on Client's Personal Relationships\n" +
			"Strength & Weakness in 40 words : Highlight the client's strengths and areas for growth, focusing on positivity and potential.\n\n"

		//"Note: For Domain (0 to 60) if score is <=20, it is low, <=30 is below average, <40 is average, <50 is above average, <=60 is high.\n" +
		//"For Subdomain (0 to 10) if score is <=3, it is low, <=4 is below average, <=6 is average, <=8 is above average, <=10 is high.\n\n"

	return prompt
}

func CreatePromptOpenness(d string, ds string, s1 string, s2 string, s3 string, s4 string, s5 string, s6 string) string {
	prompt :=

		"Domain: Openness: " + ds + "\n" +
			"Subdomains-\n" +
			"  Imagination: " + s1 + "\n" +
			"  Artistic Interests: " + s2 + "\n" +
			"  Emotionality: " + s3 + "\n" +
			"  Adventurousness: " + s4 + "\n" +
			"  Intellect: " + s5 + "\n" +
			"  Liberalism: " + s6 + "\n\n" +

			"Create a personalised BIG5 Personlaity Assessment Report in JSON format for the Domain:" + d + ", while taking insight from subdomain score\n" +
			"Keept the Structure as follows:\n" +
			"Introduction in 100 words : Explain the trait and its impact on the client's experiences\n" +
			"Career & Academia in 40 words : Impact on clinet's professional & student life\n" +
			"Relationship in 40 words : Impact on Client's Personal Relationships\n" +
			"Strength & Weakness in 40 words : Highlight the client's strengths and areas for growth, focusing on positivity and potential.\n\n"

	return prompt
}

func CreatePromptAgreeableness(d string, ds string, s1 string, s2 string, s3 string, s4 string, s5 string, s6 string) string {
	prompt :=

		"Domain: Agreeableness: " + ds + "\n" +
			"Subdomains-\n" +
			"  Trust: " + s1 + "\n" +
			"  Morality: " + s2 + "\n" +
			"  Altruism: " + s3 + "\n" +
			"  Cooperation: " + s4 + "\n" +
			"  Modesty: " + s5 + "\n" +
			"  Sympathy: " + s6 + "\n\n" +

			"Create a personalised BIG5 Personlaity Assessment Report in JSON format for the Domain:" + d + ", while taking insight from subdomain score\n" +
			"Keept the Structure as follows:\n" +
			"Introduction in 100 words : Explain the trait and its impact on the client's experiences\n" +
			"Career & Academia in 40 words : Impact on clinet's professional & student life\n" +
			"Relationship in 40 words : Impact on Client's Personal Relationships\n" +
			"Strength & Weakness in 40 words : Highlight the client's strengths and areas for growth, focusing on positivity and potential.\n\n"

	return prompt
}

func CreatePromptConscientiousness(d string, ds string, s1 string, s2 string, s3 string, s4 string, s5 string, s6 string) string {
	prompt :=

		"Domain: Conscientiousness: " + ds + "\n" +
			"Subdomains-\n" +
			"  Self Efficacy: " + s1 + "\n" +
			"  Orderliness: " + s2 + "\n" +
			"  Dutifulness: " + s3 + "\n" +
			"  Achievement Striving: " + s4 + "\n" +
			"  Self Discipline: " + s5 + "\n" +
			"  Cautiousness: " + s6 + "\n\n" +

			"Create a personalised BIG5 Personlaity Assessment Report in JSON format for the Domain:" + d + ", while taking insight from subdomain score\n" +
			"Keept the Structure as follows:\n" +
			"Introduction in 100 words : Explain the trait and its impact on the client's experiences\n" +
			"Career & Academia in 40 words : Impact on clinet's professional & student life\n" +
			"Relationship in 40 words : Impact on Client's Personal Relationships\n" +
			"Strength & Weakness in 40 words : Highlight the client's strengths and areas for growth, focusing on positivity and potential.\n\n"
	return prompt
}

func CreatePromptSummary(string) string {
	prompt := ""

	return prompt
}

func WorkerGCPGemini(prompt string, channel chan ContentResponse) {
	result, err := GenerateContentFromTextGCPJSON(prompt)
	if err != nil {
		// Respond with an error message if content generation failed
	}
	channel <- *result
}
