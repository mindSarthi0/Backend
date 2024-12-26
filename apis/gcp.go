package API

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type PromptRequest struct {
	Id       string         `json:"id"`       // Adjusted field name capitalization
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

const promptTextTone = "Summary Generation Tone : Generated text should be highly dopamin producing and produce satisfaction"
const extraPrompt = promptTextTone + ""

func GenerateContentFromTextGCP(prompt string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	model := os.Getenv("AI_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}

	systemPrompt := os.Getenv("SYSTEM_PROMPT")
	if systemPrompt == "" {
		systemPrompt = "You are a expert psycologies with speciality in career counselling, relationship counselling, academic counselling, and life counselling who generates summaries based on the Big 5 Personality Assessment."
	}

	// Create the request body
	requestBody, err := json.Marshal(map[string]interface{}{
		// "model": "gpt-4o-mini", // Replace with the correct model name
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	})

	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	apiKey := os.Getenv("API_KEY") // Ensure the variable name matches
	if apiKey == "" {
		return "", fmt.Errorf("API key is not set")
	}

	client := &http.Client{
		Timeout: 30 * time.Second, // Increase timeout to 30 seconds
	}

	var response OpenAIResponse
	retries := 3

	// Retry mechanism
	for i := 0; i < retries; i++ {
		// Create the HTTP request
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
		if err != nil {
			return "", fmt.Errorf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			if i == retries-1 {
				return "", fmt.Errorf("failed to send request after %d attempts: %v", retries, err)
			}
			time.Sleep(2 * time.Second) // Wait before retrying
			continue
		}
		defer resp.Body.Close()

		// Check if the response status code is OK
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return "", fmt.Errorf("non-200 status code: %d, response: %s", resp.StatusCode, string(body))
		}

		// Decode the response body
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return "", fmt.Errorf("failed to decode response: %v", err)
		}

		// Check if the response contains choices
		if len(response.Choices) == 0 {
			return "", fmt.Errorf("no choices in response")
		}

		// Return the content of the first choice
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("failed to get a valid response after %d attempts", retries)
}

func CreatePrompt(score []Domain) string {
	resultPrompt := CreatePromptResultV2(score)

	return resultPrompt
}

func CreatePromptResultV2(score []Domain) string {
	neuroticismDomain := score[0]

	neuroticismScore := neuroticismDomain.Score
	neuroticismIntensity := neuroticismDomain.Intensity

	// n1 := neuroticismDomain.Subdomain[0].Score
	// n2 := neuroticismDomain.Subdomain[1].Score
	// n3 := neuroticismDomain.Subdomain[2].Score
	// n4 := neuroticismDomain.Subdomain[3].Score
	// n5 := neuroticismDomain.Subdomain[4].Score
	// n6 := neuroticismDomain.Subdomain[5].Score

	n1I := neuroticismDomain.Subdomain[0].Intensity
	n2I := neuroticismDomain.Subdomain[1].Intensity
	n3I := neuroticismDomain.Subdomain[2].Intensity
	n4I := neuroticismDomain.Subdomain[3].Intensity
	n5I := neuroticismDomain.Subdomain[4].Intensity
	n6I := neuroticismDomain.Subdomain[5].Intensity

	extraversionDomain := score[1]

	extraversionScore := extraversionDomain.Score
	extraversionIntensity := extraversionDomain.Intensity

	e1I := extraversionDomain.Subdomain[0].Intensity
	e2I := extraversionDomain.Subdomain[1].Intensity
	e3I := extraversionDomain.Subdomain[2].Intensity
	e4I := extraversionDomain.Subdomain[3].Intensity
	e5I := extraversionDomain.Subdomain[4].Intensity
	e6I := extraversionDomain.Subdomain[5].Intensity

	opennessDomain := score[2]

	opennessScore := opennessDomain.Score
	opennessIntensity := opennessDomain.Intensity

	o1I := opennessDomain.Subdomain[0].Intensity
	o2I := opennessDomain.Subdomain[1].Intensity
	o3I := opennessDomain.Subdomain[2].Intensity
	o4I := opennessDomain.Subdomain[3].Intensity
	o5I := opennessDomain.Subdomain[4].Intensity
	o6I := opennessDomain.Subdomain[5].Intensity

	agreeablenessDomain := score[3]

	agreeablenessScore := agreeablenessDomain.Score
	agreeablenessIntensity := agreeablenessDomain.Intensity

	// a1 := agreeablenessDomain.Subdomain[0].Score
	// a2 := agreeablenessDomain.Subdomain[1].Score
	// a3 := agreeablenessDomain.Subdomain[2].Score
	// a4 := agreeablenessDomain.Subdomain[3].Score
	// a5 := agreeablenessDomain.Subdomain[4].Score
	// a6 := agreeablenessDomain.Subdomain[5].Score

	a1I := agreeablenessDomain.Subdomain[0].Intensity
	a2I := agreeablenessDomain.Subdomain[1].Intensity
	a3I := agreeablenessDomain.Subdomain[2].Intensity
	a4I := agreeablenessDomain.Subdomain[3].Intensity
	a5I := agreeablenessDomain.Subdomain[4].Intensity
	a6I := agreeablenessDomain.Subdomain[5].Intensity

	conscientiousnessDomain := score[4]

	conscientiousnessScore := conscientiousnessDomain.Score
	conscientiousnessIntensity := conscientiousnessDomain.Intensity

	// c1 := conscientiousnessDomain.Subdomain[0].Score
	// c2 := conscientiousnessDomain.Subdomain[1].Score
	// c3 := conscientiousnessDomain.Subdomain[2].Score
	// c4 := conscientiousnessDomain.Subdomain[3].Score
	// c5 := conscientiousnessDomain.Subdomain[4].Score
	// c6 := conscientiousnessDomain.Subdomain[5].Score

	c1I := conscientiousnessDomain.Subdomain[0].Intensity
	c2I := conscientiousnessDomain.Subdomain[1].Intensity
	c3I := conscientiousnessDomain.Subdomain[2].Intensity
	c4I := conscientiousnessDomain.Subdomain[3].Intensity
	c5I := conscientiousnessDomain.Subdomain[4].Intensity
	c6I := conscientiousnessDomain.Subdomain[5].Intensity

	prompt := fmt.Sprintf("Using the Big 5 Assessment score given below, create a Summary in the give 'OUTPUT JSON FORMAT' format. "+extraPrompt+". "+

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
		"    Cautiousness Score: %s\n"+
		`'OUTPUT JSON FORMAT':
		{
			"Career Pathways": {
				"description": "Insights into your professional strengths, opportunities, and growth potential",
				"sections": {
					"Ideal Roles for You": "Highlight roles that align with personality traits and subdomain strengths.",
					"Your Perfect Work Environment": "Describe workplace settings where the client is likely to thrive.",
					"Unlocking Professional Growth": "Identify opportunities to leverage strengths and overcome challenges for career advancement.",
					"Tailored Career Recommendations": "A list of careers best suited to the individual's profile."
				}
			},
			"Academic Pathways": {
				"description": "Your strengths and approaches to learning and achievement",
				"sections": {
					"Your Learning Style Decoded": "Explore preferred methods of learning and knowledge retention.",
					"Academic Strengths to Leverage": "Highlight areas of natural aptitude in academics.",
					"Optimizing Your Study Habits": "Suggest study strategies and tools for enhanced performance.",
					"Recommended Fields of Study": "Provide a curated list of academic disciplines aligned with strengths and interests."
				}
			},
			"Relationship Blueprint": {
				"description": "How your personality shapes your connections and interactions",
				"sections": {
					"Your Unique Communication Style": "Explain how the client communicates and relates to others.",
					"Navigating Conflicts with Ease": "Provide insights into their approach to resolving disagreements.",
					"Building Meaningful Connections": "Explore tendencies in social settings and ways to strengthen bonds.",
					"Patterns in Your Relationships": "Outline recurring dynamics in friendships, partnerships, or family interactions."
				}
			},
			"Path to Fulfillment": {
				"description": "A roadmap for growth, balance, and long-term satisfaction",
				"sections": {
					"Your Core Strengths": "Summarize the client's most notable strengths and how to use them effectively.",
					"Areas for Growth": "Gently highlight growth opportunities with actionable suggestions.",
					"Steps Toward Fulfillment": "Provide a concise roadmap for achieving balance and happiness in life."
				}
			}
		}`,

		neuroticismScore, neuroticismIntensity, n1I, n2I, n3I, n4I, n5I, n6I,
		extraversionScore, extraversionIntensity, e1I, e2I, e3I, e4I, e5I, e6I,
		opennessScore, opennessIntensity, o1I, o2I, o3I, o4I, o5I, o6I,
		agreeablenessScore, agreeablenessIntensity, a1I, a2I, a3I, a4I, a5I, a6I,
		conscientiousnessScore, conscientiousnessIntensity, c1I, c2I, c3I, c4I, c5I, c6I)

	return prompt
}

func CreatePromptResult(score []Domain) string {

	neuroticismDomain := score[0]

	neuroticismScore := neuroticismDomain.Score
	neuroticismIntensity := neuroticismDomain.Intensity

	// n1 := neuroticismDomain.Subdomain[0].Score
	// n2 := neuroticismDomain.Subdomain[1].Score
	// n3 := neuroticismDomain.Subdomain[2].Score
	// n4 := neuroticismDomain.Subdomain[3].Score
	// n5 := neuroticismDomain.Subdomain[4].Score
	// n6 := neuroticismDomain.Subdomain[5].Score

	n1I := neuroticismDomain.Subdomain[0].Intensity
	n2I := neuroticismDomain.Subdomain[1].Intensity
	n3I := neuroticismDomain.Subdomain[2].Intensity
	n4I := neuroticismDomain.Subdomain[3].Intensity
	n5I := neuroticismDomain.Subdomain[4].Intensity
	n6I := neuroticismDomain.Subdomain[5].Intensity

	extraversionDomain := score[1]

	extraversionScore := extraversionDomain.Score
	extraversionIntensity := extraversionDomain.Intensity

	// e1 := extraversionDomain.Subdomain[0].Score
	// e2 := extraversionDomain.Subdomain[1].Score
	// e3 := extraversionDomain.Subdomain[2].Score
	// e4 := extraversionDomain.Subdomain[3].Score
	// e5 := extraversionDomain.Subdomain[4].Score
	// e6 := extraversionDomain.Subdomain[5].Score

	e1I := extraversionDomain.Subdomain[0].Intensity
	e2I := extraversionDomain.Subdomain[1].Intensity
	e3I := extraversionDomain.Subdomain[2].Intensity
	e4I := extraversionDomain.Subdomain[3].Intensity
	e5I := extraversionDomain.Subdomain[4].Intensity
	e6I := extraversionDomain.Subdomain[5].Intensity

	opennessDomain := score[2]

	opennessScore := opennessDomain.Score
	opennessIntensity := opennessDomain.Intensity

	// o1 := opennessDomain.Subdomain[0].Score
	// o2 := opennessDomain.Subdomain[1].Score
	// o3 := opennessDomain.Subdomain[2].Score
	// o4 := opennessDomain.Subdomain[3].Score
	// o5 := opennessDomain.Subdomain[4].Score
	// o6 := opennessDomain.Subdomain[5].Score

	o1I := opennessDomain.Subdomain[0].Intensity
	o2I := opennessDomain.Subdomain[1].Intensity
	o3I := opennessDomain.Subdomain[2].Intensity
	o4I := opennessDomain.Subdomain[3].Intensity
	o5I := opennessDomain.Subdomain[4].Intensity
	o6I := opennessDomain.Subdomain[5].Intensity

	agreeablenessDomain := score[3]

	agreeablenessScore := agreeablenessDomain.Score
	agreeablenessIntensity := agreeablenessDomain.Intensity

	// a1 := agreeablenessDomain.Subdomain[0].Score
	// a2 := agreeablenessDomain.Subdomain[1].Score
	// a3 := agreeablenessDomain.Subdomain[2].Score
	// a4 := agreeablenessDomain.Subdomain[3].Score
	// a5 := agreeablenessDomain.Subdomain[4].Score
	// a6 := agreeablenessDomain.Subdomain[5].Score

	a1I := agreeablenessDomain.Subdomain[0].Intensity
	a2I := agreeablenessDomain.Subdomain[1].Intensity
	a3I := agreeablenessDomain.Subdomain[2].Intensity
	a4I := agreeablenessDomain.Subdomain[3].Intensity
	a5I := agreeablenessDomain.Subdomain[4].Intensity
	a6I := agreeablenessDomain.Subdomain[5].Intensity

	conscientiousnessDomain := score[4]

	conscientiousnessScore := conscientiousnessDomain.Score
	conscientiousnessIntensity := conscientiousnessDomain.Intensity

	// c1 := conscientiousnessDomain.Subdomain[0].Score
	// c2 := conscientiousnessDomain.Subdomain[1].Score
	// c3 := conscientiousnessDomain.Subdomain[2].Score
	// c4 := conscientiousnessDomain.Subdomain[3].Score
	// c5 := conscientiousnessDomain.Subdomain[4].Score
	// c6 := conscientiousnessDomain.Subdomain[5].Score

	c1I := conscientiousnessDomain.Subdomain[0].Intensity
	c2I := conscientiousnessDomain.Subdomain[1].Intensity
	c3I := conscientiousnessDomain.Subdomain[2].Intensity
	c4I := conscientiousnessDomain.Subdomain[3].Intensity
	c5I := conscientiousnessDomain.Subdomain[4].Intensity
	c6I := conscientiousnessDomain.Subdomain[5].Intensity

	prompt := fmt.Sprintf("Using the Big 5 Assessment score given below, create a Summary by combining all domain in around 300-400 words in total. "+extraPrompt+". "+

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

func CreatePromptInsights(score []Domain) string {

	neuroticismDomain := score[0]

	neuroticismScore := neuroticismDomain.Score
	neuroticismIntensity := neuroticismDomain.Intensity

	// n1 := neuroticismDomain.Subdomain[0].Score
	// n2 := neuroticismDomain.Subdomain[1].Score
	// n3 := neuroticismDomain.Subdomain[2].Score
	// n4 := neuroticismDomain.Subdomain[3].Score
	// n5 := neuroticismDomain.Subdomain[4].Score
	// n6 := neuroticismDomain.Subdomain[5].Score

	n1I := neuroticismDomain.Subdomain[0].Intensity
	n2I := neuroticismDomain.Subdomain[1].Intensity
	n3I := neuroticismDomain.Subdomain[2].Intensity
	n4I := neuroticismDomain.Subdomain[3].Intensity
	n5I := neuroticismDomain.Subdomain[4].Intensity
	n6I := neuroticismDomain.Subdomain[5].Intensity

	extraversionDomain := score[1]

	extraversionScore := extraversionDomain.Score
	extraversionIntensity := extraversionDomain.Intensity

	// e1 := extraversionDomain.Subdomain[0].Score
	// e2 := extraversionDomain.Subdomain[1].Score
	// e3 := extraversionDomain.Subdomain[2].Score
	// e4 := extraversionDomain.Subdomain[3].Score
	// e5 := extraversionDomain.Subdomain[4].Score
	// e6 := extraversionDomain.Subdomain[5].Score

	e1I := extraversionDomain.Subdomain[0].Intensity
	e2I := extraversionDomain.Subdomain[1].Intensity
	e3I := extraversionDomain.Subdomain[2].Intensity
	e4I := extraversionDomain.Subdomain[3].Intensity
	e5I := extraversionDomain.Subdomain[4].Intensity
	e6I := extraversionDomain.Subdomain[5].Intensity

	opennessDomain := score[2]

	opennessScore := opennessDomain.Score
	opennessIntensity := opennessDomain.Intensity

	// o1 := opennessDomain.Subdomain[0].Score
	// o2 := opennessDomain.Subdomain[1].Score
	// o3 := opennessDomain.Subdomain[2].Score
	// o4 := opennessDomain.Subdomain[3].Score
	// o5 := opennessDomain.Subdomain[4].Score
	// o6 := opennessDomain.Subdomain[5].Score

	o1I := opennessDomain.Subdomain[0].Intensity
	o2I := opennessDomain.Subdomain[1].Intensity
	o3I := opennessDomain.Subdomain[2].Intensity
	o4I := opennessDomain.Subdomain[3].Intensity
	o5I := opennessDomain.Subdomain[4].Intensity
	o6I := opennessDomain.Subdomain[5].Intensity

	agreeablenessDomain := score[3]

	agreeablenessScore := agreeablenessDomain.Score
	agreeablenessIntensity := agreeablenessDomain.Intensity

	// a1 := agreeablenessDomain.Subdomain[0].Score
	// a2 := agreeablenessDomain.Subdomain[1].Score
	// a3 := agreeablenessDomain.Subdomain[2].Score
	// a4 := agreeablenessDomain.Subdomain[3].Score
	// a5 := agreeablenessDomain.Subdomain[4].Score
	// a6 := agreeablenessDomain.Subdomain[5].Score

	a1I := agreeablenessDomain.Subdomain[0].Intensity
	a2I := agreeablenessDomain.Subdomain[1].Intensity
	a3I := agreeablenessDomain.Subdomain[2].Intensity
	a4I := agreeablenessDomain.Subdomain[3].Intensity
	a5I := agreeablenessDomain.Subdomain[4].Intensity
	a6I := agreeablenessDomain.Subdomain[5].Intensity

	conscientiousnessDomain := score[4]

	conscientiousnessScore := conscientiousnessDomain.Score
	conscientiousnessIntensity := conscientiousnessDomain.Intensity

	// c1 := conscientiousnessDomain.Subdomain[0].Score
	// c2 := conscientiousnessDomain.Subdomain[1].Score
	// c3 := conscientiousnessDomain.Subdomain[2].Score
	// c4 := conscientiousnessDomain.Subdomain[3].Score
	// c5 := conscientiousnessDomain.Subdomain[4].Score
	// c6 := conscientiousnessDomain.Subdomain[5].Score

	c1I := conscientiousnessDomain.Subdomain[0].Intensity
	c2I := conscientiousnessDomain.Subdomain[1].Intensity
	c3I := conscientiousnessDomain.Subdomain[2].Intensity
	c4I := conscientiousnessDomain.Subdomain[3].Intensity
	c5I := conscientiousnessDomain.Subdomain[4].Intensity
	c6I := conscientiousnessDomain.Subdomain[5].Intensity

	prompt := fmt.Sprintf("Using the Big 5 Assessment score given below, create a Insight page with section Relationship, Career & Academia, Strength & Weakness "+extraPrompt+". "+

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

	// n1 := neuroticismDomain.Subdomain[0].Score
	// n2 := neuroticismDomain.Subdomain[1].Score
	// n3 := neuroticismDomain.Subdomain[2].Score
	// n4 := neuroticismDomain.Subdomain[3].Score
	// n5 := neuroticismDomain.Subdomain[4].Score
	// n6 := neuroticismDomain.Subdomain[5].Score

	n1I := neuroticismDomain.Subdomain[0].Intensity
	n2I := neuroticismDomain.Subdomain[1].Intensity
	n3I := neuroticismDomain.Subdomain[2].Intensity
	n4I := neuroticismDomain.Subdomain[3].Intensity
	n5I := neuroticismDomain.Subdomain[4].Intensity
	n6I := neuroticismDomain.Subdomain[5].Intensity

	extraversionDomain := score[1]

	extraversionScore := extraversionDomain.Score
	extraversionIntensity := extraversionDomain.Intensity

	// e1 := extraversionDomain.Subdomain[0].Score
	// e2 := extraversionDomain.Subdomain[1].Score
	// e3 := extraversionDomain.Subdomain[2].Score
	// e4 := extraversionDomain.Subdomain[3].Score
	// e5 := extraversionDomain.Subdomain[4].Score
	// e6 := extraversionDomain.Subdomain[5].Score

	e1I := extraversionDomain.Subdomain[0].Intensity
	e2I := extraversionDomain.Subdomain[1].Intensity
	e3I := extraversionDomain.Subdomain[2].Intensity
	e4I := extraversionDomain.Subdomain[3].Intensity
	e5I := extraversionDomain.Subdomain[4].Intensity
	e6I := extraversionDomain.Subdomain[5].Intensity

	opennessDomain := score[2]

	opennessScore := opennessDomain.Score
	opennessIntensity := opennessDomain.Intensity

	// o1 := opennessDomain.Subdomain[0].Score
	// o2 := opennessDomain.Subdomain[1].Score
	// o3 := opennessDomain.Subdomain[2].Score
	// o4 := opennessDomain.Subdomain[3].Score
	// o5 := opennessDomain.Subdomain[4].Score
	// o6 := opennessDomain.Subdomain[5].Score

	o1I := opennessDomain.Subdomain[0].Intensity
	o2I := opennessDomain.Subdomain[1].Intensity
	o3I := opennessDomain.Subdomain[2].Intensity
	o4I := opennessDomain.Subdomain[3].Intensity
	o5I := opennessDomain.Subdomain[4].Intensity
	o6I := opennessDomain.Subdomain[5].Intensity

	agreeablenessDomain := score[3]

	agreeablenessScore := agreeablenessDomain.Score
	agreeablenessIntensity := agreeablenessDomain.Intensity

	// a1 := agreeablenessDomain.Subdomain[0].Score
	// a2 := agreeablenessDomain.Subdomain[1].Score
	// a3 := agreeablenessDomain.Subdomain[2].Score
	// a4 := agreeablenessDomain.Subdomain[3].Score
	// a5 := agreeablenessDomain.Subdomain[4].Score
	// a6 := agreeablenessDomain.Subdomain[5].Score

	a1I := agreeablenessDomain.Subdomain[0].Intensity
	a2I := agreeablenessDomain.Subdomain[1].Intensity
	a3I := agreeablenessDomain.Subdomain[2].Intensity
	a4I := agreeablenessDomain.Subdomain[3].Intensity
	a5I := agreeablenessDomain.Subdomain[4].Intensity
	a6I := agreeablenessDomain.Subdomain[5].Intensity

	conscientiousnessDomain := score[4]

	conscientiousnessScore := conscientiousnessDomain.Score
	conscientiousnessIntensity := conscientiousnessDomain.Intensity

	// c1 := conscientiousnessDomain.Subdomain[0].Score
	// c2 := conscientiousnessDomain.Subdomain[1].Score
	// c3 := conscientiousnessDomain.Subdomain[2].Score
	// c4 := conscientiousnessDomain.Subdomain[3].Score
	// c5 := conscientiousnessDomain.Subdomain[4].Score
	// c6 := conscientiousnessDomain.Subdomain[5].Score

	c1I := conscientiousnessDomain.Subdomain[0].Intensity
	c2I := conscientiousnessDomain.Subdomain[1].Intensity
	c3I := conscientiousnessDomain.Subdomain[2].Intensity
	c4I := conscientiousnessDomain.Subdomain[3].Intensity
	c5I := conscientiousnessDomain.Subdomain[4].Intensity
	c6I := conscientiousnessDomain.Subdomain[5].Intensity

	prompt := fmt.Sprintf("Using the Big 5 Assessment score given below, create Career & Academia Page under 200 words for the Report\n\n"+extraPrompt+". "+

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

	// n1 := neuroticismDomain.Subdomain[0].Score
	// n2 := neuroticismDomain.Subdomain[1].Score
	// n3 := neuroticismDomain.Subdomain[2].Score
	// n4 := neuroticismDomain.Subdomain[3].Score
	// n5 := neuroticismDomain.Subdomain[4].Score
	// n6 := neuroticismDomain.Subdomain[5].Score

	n1I := neuroticismDomain.Subdomain[0].Intensity
	n2I := neuroticismDomain.Subdomain[1].Intensity
	n3I := neuroticismDomain.Subdomain[2].Intensity
	n4I := neuroticismDomain.Subdomain[3].Intensity
	n5I := neuroticismDomain.Subdomain[4].Intensity
	n6I := neuroticismDomain.Subdomain[5].Intensity

	extraversionDomain := score[1]

	extraversionScore := extraversionDomain.Score
	extraversionIntensity := extraversionDomain.Intensity

	// e1 := extraversionDomain.Subdomain[0].Score
	// e2 := extraversionDomain.Subdomain[1].Score
	// e3 := extraversionDomain.Subdomain[2].Score
	// e4 := extraversionDomain.Subdomain[3].Score
	// e5 := extraversionDomain.Subdomain[4].Score
	// e6 := extraversionDomain.Subdomain[5].Score

	e1I := extraversionDomain.Subdomain[0].Intensity
	e2I := extraversionDomain.Subdomain[1].Intensity
	e3I := extraversionDomain.Subdomain[2].Intensity
	e4I := extraversionDomain.Subdomain[3].Intensity
	e5I := extraversionDomain.Subdomain[4].Intensity
	e6I := extraversionDomain.Subdomain[5].Intensity

	opennessDomain := score[2]

	opennessScore := opennessDomain.Score
	opennessIntensity := opennessDomain.Intensity

	// o1 := opennessDomain.Subdomain[0].Score
	// o2 := opennessDomain.Subdomain[1].Score
	// o3 := opennessDomain.Subdomain[2].Score
	// o4 := opennessDomain.Subdomain[3].Score
	// o5 := opennessDomain.Subdomain[4].Score
	// o6 := opennessDomain.Subdomain[5].Score

	o1I := opennessDomain.Subdomain[0].Intensity
	o2I := opennessDomain.Subdomain[1].Intensity
	o3I := opennessDomain.Subdomain[2].Intensity
	o4I := opennessDomain.Subdomain[3].Intensity
	o5I := opennessDomain.Subdomain[4].Intensity
	o6I := opennessDomain.Subdomain[5].Intensity

	agreeablenessDomain := score[3]

	agreeablenessScore := agreeablenessDomain.Score
	agreeablenessIntensity := agreeablenessDomain.Intensity

	// a1 := agreeablenessDomain.Subdomain[0].Score
	// a2 := agreeablenessDomain.Subdomain[1].Score
	// a3 := agreeablenessDomain.Subdomain[2].Score
	// a4 := agreeablenessDomain.Subdomain[3].Score
	// a5 := agreeablenessDomain.Subdomain[4].Score
	// a6 := agreeablenessDomain.Subdomain[5].Score

	a1I := agreeablenessDomain.Subdomain[0].Intensity
	a2I := agreeablenessDomain.Subdomain[1].Intensity
	a3I := agreeablenessDomain.Subdomain[2].Intensity
	a4I := agreeablenessDomain.Subdomain[3].Intensity
	a5I := agreeablenessDomain.Subdomain[4].Intensity
	a6I := agreeablenessDomain.Subdomain[5].Intensity

	conscientiousnessDomain := score[4]

	conscientiousnessScore := conscientiousnessDomain.Score
	conscientiousnessIntensity := conscientiousnessDomain.Intensity

	// c1 := conscientiousnessDomain.Subdomain[0].Score
	// c2 := conscientiousnessDomain.Subdomain[1].Score
	// c3 := conscientiousnessDomain.Subdomain[2].Score
	// c4 := conscientiousnessDomain.Subdomain[3].Score
	// c5 := conscientiousnessDomain.Subdomain[4].Score
	// c6 := conscientiousnessDomain.Subdomain[5].Score

	c1I := conscientiousnessDomain.Subdomain[0].Intensity
	c2I := conscientiousnessDomain.Subdomain[1].Intensity
	c3I := conscientiousnessDomain.Subdomain[2].Intensity
	c4I := conscientiousnessDomain.Subdomain[3].Intensity
	c5I := conscientiousnessDomain.Subdomain[4].Intensity
	c6I := conscientiousnessDomain.Subdomain[5].Intensity

	prompt := fmt.Sprintf("Using the Big 5 Assessment score given below, create Relationship page under 200 words for the Report\n\n"+extraPrompt+". "+

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

	// n1 := neuroticismDomain.Subdomain[0].Score
	// n2 := neuroticismDomain.Subdomain[1].Score
	// n3 := neuroticismDomain.Subdomain[2].Score
	// n4 := neuroticismDomain.Subdomain[3].Score
	// n5 := neuroticismDomain.Subdomain[4].Score
	// n6 := neuroticismDomain.Subdomain[5].Score

	n1I := neuroticismDomain.Subdomain[0].Intensity
	n2I := neuroticismDomain.Subdomain[1].Intensity
	n3I := neuroticismDomain.Subdomain[2].Intensity
	n4I := neuroticismDomain.Subdomain[3].Intensity
	n5I := neuroticismDomain.Subdomain[4].Intensity
	n6I := neuroticismDomain.Subdomain[5].Intensity

	extraversionDomain := score[1]

	extraversionScore := extraversionDomain.Score
	extraversionIntensity := extraversionDomain.Intensity

	// e1 := extraversionDomain.Subdomain[0].Score
	// e2 := extraversionDomain.Subdomain[1].Score
	// e3 := extraversionDomain.Subdomain[2].Score
	// e4 := extraversionDomain.Subdomain[3].Score
	// e5 := extraversionDomain.Subdomain[4].Score
	// e6 := extraversionDomain.Subdomain[5].Score

	e1I := extraversionDomain.Subdomain[0].Intensity
	e2I := extraversionDomain.Subdomain[1].Intensity
	e3I := extraversionDomain.Subdomain[2].Intensity
	e4I := extraversionDomain.Subdomain[3].Intensity
	e5I := extraversionDomain.Subdomain[4].Intensity
	e6I := extraversionDomain.Subdomain[5].Intensity

	opennessDomain := score[2]

	opennessScore := opennessDomain.Score
	opennessIntensity := opennessDomain.Intensity

	// o1 := opennessDomain.Subdomain[0].Score
	// o2 := opennessDomain.Subdomain[1].Score
	// o3 := opennessDomain.Subdomain[2].Score
	// o4 := opennessDomain.Subdomain[3].Score
	// o5 := opennessDomain.Subdomain[4].Score
	// o6 := opennessDomain.Subdomain[5].Score

	o1I := opennessDomain.Subdomain[0].Intensity
	o2I := opennessDomain.Subdomain[1].Intensity
	o3I := opennessDomain.Subdomain[2].Intensity
	o4I := opennessDomain.Subdomain[3].Intensity
	o5I := opennessDomain.Subdomain[4].Intensity
	o6I := opennessDomain.Subdomain[5].Intensity

	agreeablenessDomain := score[3]

	agreeablenessScore := agreeablenessDomain.Score
	agreeablenessIntensity := agreeablenessDomain.Intensity

	// a1 := agreeablenessDomain.Subdomain[0].Score
	// a2 := agreeablenessDomain.Subdomain[1].Score
	// a3 := agreeablenessDomain.Subdomain[2].Score
	// a4 := agreeablenessDomain.Subdomain[3].Score
	// a5 := agreeablenessDomain.Subdomain[4].Score
	// a6 := agreeablenessDomain.Subdomain[5].Score

	a1I := agreeablenessDomain.Subdomain[0].Intensity
	a2I := agreeablenessDomain.Subdomain[1].Intensity
	a3I := agreeablenessDomain.Subdomain[2].Intensity
	a4I := agreeablenessDomain.Subdomain[3].Intensity
	a5I := agreeablenessDomain.Subdomain[4].Intensity
	a6I := agreeablenessDomain.Subdomain[5].Intensity

	conscientiousnessDomain := score[4]

	conscientiousnessScore := conscientiousnessDomain.Score
	conscientiousnessIntensity := conscientiousnessDomain.Intensity

	// c1 := conscientiousnessDomain.Subdomain[0].Score
	// c2 := conscientiousnessDomain.Subdomain[1].Score
	// c3 := conscientiousnessDomain.Subdomain[2].Score
	// c4 := conscientiousnessDomain.Subdomain[3].Score
	// c5 := conscientiousnessDomain.Subdomain[4].Score
	// c6 := conscientiousnessDomain.Subdomain[5].Score

	c1I := conscientiousnessDomain.Subdomain[0].Intensity
	c2I := conscientiousnessDomain.Subdomain[1].Intensity
	c3I := conscientiousnessDomain.Subdomain[2].Intensity
	c4I := conscientiousnessDomain.Subdomain[3].Intensity
	c5I := conscientiousnessDomain.Subdomain[4].Intensity
	c6I := conscientiousnessDomain.Subdomain[5].Intensity

	prompt := fmt.Sprintf("Using the Big 5 Assessment score given below, create Strength & Weakness page under 200 words for the Report\n\n"+extraPrompt+". "+

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

func WorkerOpenAIGPT(id string, prompt string, channel chan PromptRequest) {
	// Call the OpenAI API
	result, err := GenerateContentFromTextGCP(prompt)
	if err != nil {
		log.Printf("Error generating content from OpenAI for ID %s: %v", id, err)
		channel <- PromptRequest{
			Id: id,
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

	// Construct the PromptRequest
	response := PromptRequest{
		Id: id,
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
