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
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.0-pro-latest:generateContent"

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
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent"

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

// Function to generate content using Google Cloud's Vertex AI (client library-based approach)
// func GenerateContentFromText(prompt string) error {
// 	// Define the model name you want to use from Vertex AI
// 	modelName := "gemini-1.5-flash-001" // This is the model to use for content generation

// 	// Set up the context and initialize the Vertex AI client using the API key
// 	ctx := context.Background()

// 	// Retrieve the API key from the environment variable 'API_KEY'
// 	// apiKey := os.Getenv("API_KEY")
// 	// if apiKey == "" {
// 	// 	return fmt.Errorf("API key is not set. Please ensure the environment variable 'API_KEY' is set")
// 	// }

// 	// Create a new client for Vertex AI Generative models using the API key
// 	projectID := "cognify-438322" // Replace with your actual project ID
// 	location := "asia-south1"     // Replace with your actual location

// 	client, err := genai.NewClient(ctx, projectID, location, option.WithCredentialsFile("cognify-438322-90e04392ce12.json"))
// 	// client, err := genai.NewClient(ctx, projectID, location, option.WithCredentialsFile("cognify-438322-90e04392ce12.json"))
// 	if err != nil {
// 		return fmt.Errorf("error creating Vertex AI client: %w", err)
// 	}

// 	// Create the prompt text for the model to generate content
// 	gemini := client.GenerativeModel(modelName)
// 	generatePrompt := genai.Text(prompt)

// 	// Generate content using the defined model and prompt
// 	resp, err := gemini.GenerateContent(ctx, generatePrompt)
// 	if err != nil {
// 		return fmt.Errorf("error generating content from Vertex AI: %w", err)
// 	}

// 	// Print the generated content as formatted JSON
// 	rb, err := json.MarshalIndent(resp, "", "  ")
// 	if err != nil {
// 		return fmt.Errorf("json.MarshalIndent: %w", err)
// 	}
// 	fmt.Println(string(rb))

// 	return nil
// }

// Function to create a Big 5 personality prompt (can be used as input to the API)
func CreatePrompt(D1, N1, N2, N3, N4, N5, N6, D2, E1, E2, E3, E4, E5, E6, D3, O1, O2, O3, O4, O5, O6, D4, A1, A2, A3, A4, A5, A6, D5, C1, C2, C3, C4, C5, C6 string) string {
	// Build the Big 5 Personality prompt with all domain and subdomain scores
	prompt := "Instructions - \n" +
		"1) Write a personalized 6 page (5 dedicated to each domain and 6th titled unlocking your potential) report for my client with this BIG5 test score.\n" +
		"2) Do nwot refer directly to any subdomain score. \n" +
		"3) Keep the Tone Professional, Encouraging, Warm, empathetic, positive & solution focused\n" +
		"4) Use Second-person pronoun\n" +
		"5) Keept the Structure as follows: for 1 to 5 page dedicated to each domain - \n" +
		"Introduction (Should be 80 to 100 words) : Explain the trait and its impact on the client's experiences\n" +
		"Career & Academia (Should be 30 to 40 words) : Impact on clinet's professional & student life\n" +
		"Relationship (Should be 30 to 40 words) : Impact on Client's Personal Relationships\n" +
		"Strength & Weakness (30 to 40 words) : Highlight the client's strengths and areas for growth, focusing on positivity and potential.\n\n" +
		"For 6th page: write a 200 word summary giving insight to help the client in thier self development\n\n" +
		"6) Print the output in JSON for each page and their sub parts\n\n" +
		"Most Important Instruction : Adhere to ord limit mentioned in 5th point of Instructions \n\n\n\n" +

		"Domain: Neuroticism: " + D1 + "\n" +
		"Subdomains-\n" +
		"  Anxiety: " + N1 + "\n" +
		"  Anger: " + N2 + "\n" +
		"  Depression: " + N3 + "\n" +
		"  Self-consciousness: " + N4 + "\n" +
		"  Immoderation: " + N5 + "\n" +
		"  Vulnerability: " + N6 + "\n\n" +

		"Domain: Extraversion: " + D2 + "\n" +
		"Subdomains-\n" +
		"  Friendliness: " + E1 + "\n" +
		"  Gregariousness: " + E2 + "\n" +
		"  Assertiveness: " + E3 + "\n" +
		"  Activity Level: " + E4 + "\n" +
		"  Excitement Seeking: " + E5 + "\n" +
		"  Cheerfulness: " + E6 + "\n\n" +

		"Domain: Openness: " + D3 + "\n" +
		"Subdomains-\n" +
		"  Imagination: " + O1 + "\n" +
		"  Artistic Interests: " + O2 + "\n" +
		"  Emotionality: " + O3 + "\n" +
		"  Adventurousness: " + O4 + "\n" +
		"  Intellect: " + O5 + "\n" +
		"  Liberalism: " + O6 + "\n\n" +

		"Domain: Agreeableness: " + D4 + "\n" +
		"Subdomains-\n" +
		"  Trust: " + A1 + "\n" +
		"  Morality: " + A2 + "\n" +
		"  Altruism: " + A3 + "\n" +
		"  Cooperation: " + A4 + "\n" +
		"  Modesty: " + A5 + "\n" +
		"  Sympathy: " + A6 + "\n\n" +

		"Domain: Conscientiousness: " + D5 + "\n" +
		"Subdomains-\n" +
		"  Self Efficacy: " + C1 + "\n" +
		"  Orderliness: " + C2 + "\n" +
		"  Dutifulness: " + C3 + "\n" +
		"  Achievement Striving: " + C4 + "\n" +
		"  Self Discipline: " + C5 + "\n" +
		"  Cautiousness: " + C6 + "\n\n" +

		"Note: For Domain (0 to 60) if score is <=20, it is low, <=30 is below average, <40 is average, <50 is above average, <=60 is high.\n" +
		"For Subdomain (0 to 10) if score is <=3, it is low, <=4 is below average, <=6 is average, <=8 is above average, <=10 is high.\n\n"

	return prompt
}

func CreatePromptNeuroticism(D, S1, S2, S3, S4, S5, S6, string) string {
	prompt := "Create a personalised BIG5 Personlaity Assessment Report for the Domain:" + D + "\n" +
		"Keept the Structure as follows:\n" +
		"Introduction (Should be 80 to 100 words) : Explain the trait and its impact on the client's experiences\n" +
		"Career & Academia (Should be 30 to 40 words) : Impact on clinet's professional & student life\n" +
		"Relationship (Should be 30 to 40 words) : Impact on Client's Personal Relationships\n" +
		"Strength & Weakness (30 to 40 words) : Highlight the client's strengths and areas for growth, focusing on positivity and potential.\n\n" +
		"Domain: Neuroticism: " + D + "\n" +
		"Subdomains-\n" +
		"  Anxiety: " + S1 + "\n" +
		"  Anger: " + S2 + "\n" +
		"  Depression: " + S3 + "\n" +
		"  Self-consciousness: " + S4 + "\n" +
		"  Immoderation: " + S5 + "\n" +
		"  Vulnerability: " + S6 + "\n\n" +

		"Note: For Domain (0 to 60) if score is <=20, it is low, <=30 is below average, <40 is average, <50 is above average, <=60 is high.\n" +
		"For Subdomain (0 to 10) if score is <=3, it is low, <=4 is below average, <=6 is average, <=8 is above average, <=10 is high.\n\n"

	return prompt

}

func CreatePromptExtraversion(D, S1, S2, S3, S4, S5, S6, string) string {
	prompt := "Create a personalised BIG5 Personlaity Assessment Report for the Domain:" + D + "\n" +
		"Keept the Structure as follows:\n" +
		"Introduction (Should be 80 to 100 words) : Explain the trait and its impact on the client's experiences\n" +
		"Career & Academia (Should be 30 to 40 words) : Impact on clinet's professional & student life\n" +
		"Relationship (Should be 30 to 40 words) : Impact on Client's Personal Relationships\n" +
		"Strength & Weakness (30 to 40 words) : Highlight the client's strengths and areas for growth, focusing on positivity and potential.\n\n" +
		"Domain: Neuroticism: " + D + "\n" +
		"Subdomains-\n" +
		"  Anxiety: " + S1 + "\n" +
		"  Anger: " + S2 + "\n" +
		"  Depression: " + S3 + "\n" +
		"  Self-consciousness: " + S4 + "\n" +
		"  Immoderation: " + S5 + "\n" +
		"  Vulnerability: " + S6 + "\n\n" +

		"Note: For Domain (0 to 60) if score is <=20, it is low, <=30 is below average, <40 is average, <50 is above average, <=60 is high.\n" +
		"For Subdomain (0 to 10) if score is <=3, it is low, <=4 is below average, <=6 is average, <=8 is above average, <=10 is high.\n\n"

	return prompt
}

func CreatePromptOpenness(D, S1, S2, S3, S4, S5, S6, string) string {
	prompt := "Create a personalised BIG5 Personlaity Assessment Report for the Domain:" + D + "\n" +
		"Keept the Structure as follows:\n" +
		"Introduction (Should be 80 to 100 words) : Explain the trait and its impact on the client's experiences\n" +
		"Career & Academia (Should be 30 to 40 words) : Impact on clinet's professional & student life\n" +
		"Relationship (Should be 30 to 40 words) : Impact on Client's Personal Relationships\n" +
		"Strength & Weakness (30 to 40 words) : Highlight the client's strengths and areas for growth, focusing on positivity and potential.\n\n" +
		"Domain: Neuroticism: " + D + "\n" +
		"Subdomains-\n" +
		"  Anxiety: " + S1 + "\n" +
		"  Anger: " + S2 + "\n" +
		"  Depression: " + S3 + "\n" +
		"  Self-consciousness: " + S4 + "\n" +
		"  Immoderation: " + S5 + "\n" +
		"  Vulnerability: " + S6 + "\n\n" +

		"Note: For Domain (0 to 60) if score is <=20, it is low, <=30 is below average, <40 is average, <50 is above average, <=60 is high.\n" +
		"For Subdomain (0 to 10) if score is <=3, it is low, <=4 is below average, <=6 is average, <=8 is above average, <=10 is high.\n\n"

	return prompt
}

func CreatePromptAgreeableness(D, S1, S2, S3, S4, S5, S6, string) string {
	prompt := "Create a personalised BIG5 Personlaity Assessment Report for the Domain:" + D + "\n" +
		"Keept the Structure as follows:\n" +
		"Introduction (Should be 80 to 100 words) : Explain the trait and its impact on the client's experiences\n" +
		"Career & Academia (Should be 30 to 40 words) : Impact on clinet's professional & student life\n" +
		"Relationship (Should be 30 to 40 words) : Impact on Client's Personal Relationships\n" +
		"Strength & Weakness (30 to 40 words) : Highlight the client's strengths and areas for growth, focusing on positivity and potential.\n\n" +
		"Domain: Neuroticism: " + D + "\n" +
		"Subdomains-\n" +
		"  Anxiety: " + S1 + "\n" +
		"  Anger: " + S2 + "\n" +
		"  Depression: " + S3 + "\n" +
		"  Self-consciousness: " + S4 + "\n" +
		"  Immoderation: " + S5 + "\n" +
		"  Vulnerability: " + S6 + "\n\n" +

		"Note: For Domain (0 to 60) if score is <=20, it is low, <=30 is below average, <40 is average, <50 is above average, <=60 is high.\n" +
		"For Subdomain (0 to 10) if score is <=3, it is low, <=4 is below average, <=6 is average, <=8 is above average, <=10 is high.\n\n"

	return prompt
}

func CreatePromptConscientiousness(D, S1, S2, S3, S4, S5, S6, string) string {
	prompt := "Create a personalised BIG5 Personlaity Assessment Report for the Domain:" + D + "\n" +
		"Keept the Structure as follows:\n" +
		"Introduction (Should be 80 to 100 words) : Explain the trait and its impact on the client's experiences\n" +
		"Career & Academia (Should be 30 to 40 words) : Impact on clinet's professional & student life\n" +
		"Relationship (Should be 30 to 40 words) : Impact on Client's Personal Relationships\n" +
		"Strength & Weakness (30 to 40 words) : Highlight the client's strengths and areas for growth, focusing on positivity and potential.\n\n" +
		"Domain: Neuroticism: " + D + "\n" +
		"Subdomains-\n" +
		"  Anxiety: " + S1 + "\n" +
		"  Anger: " + S2 + "\n" +
		"  Depression: " + S3 + "\n" +
		"  Self-consciousness: " + S4 + "\n" +
		"  Immoderation: " + S5 + "\n" +
		"  Vulnerability: " + S6 + "\n\n" +

		"Note: For Domain (0 to 60) if score is <=20, it is low, <=30 is below average, <40 is average, <50 is above average, <=60 is high.\n" +
		"For Subdomain (0 to 10) if score is <=3, it is low, <=4 is below average, <=6 is average, <=8 is above average, <=10 is high.\n\n" +
		"Print the output in JSON for each page and their sub parts"

	return prompt
}

func CreatePromptSummary(string) string {
	prompt := ""

	return prompt
}
