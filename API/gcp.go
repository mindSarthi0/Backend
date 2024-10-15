package API

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/vertexai/genai"
)

// ai integration
func CreatePrompt(D1, N1, N2, N3, N4, N5, N6, D2, E1, E2, E3, E4, E5, E6, D3, O1, O2, O3, O4, O5, O6, D4, A1, A2, A3, A4, A5, A6, D5, C1, C2, C3, C4, C5, C6 string) string {
	prompt := "Generate a Big 5 Personality Assessment report based on the following data:\n" +
		"Domain: Neuroticism: " + D1 + "\n" +
		"Subdomains:\n" +
		"  Anxiety: " + N1 + "\n" +
		"  Anger: " + N2 + "\n" +
		"  Depression: " + N3 + "\n" +
		"  Self-consciousness: " + N4 + "\n" +
		"  Immoderation: " + N5 + "\n" +
		"  Vulnerability: " + N6 + "\n\n" +

		"Domain: Extraversion: " + D2 + "\n" +
		"Subdomains:\n" +
		"  Friendliness: " + E1 + "\n" +
		"  Gregariousness: " + E2 + "\n" +
		"  Assertiveness: " + E3 + "\n" +
		"  Activity Level: " + E4 + "\n" +
		"  Excitement Seeking: " + E5 + "\n" +
		"  Cheerfulness: " + E6 + "\n\n" +

		"Domain: Openness: " + D3 + "\n" +
		"Subdomains:\n" +
		"  Imagination: " + O1 + "\n" +
		"  Artistic Interests: " + O2 + "\n" +
		"  Emotionality: " + O3 + "\n" +
		"  Adventurousness: " + O4 + "\n" +
		"  Intellect: " + O5 + "\n" +
		"  Liberalism: " + O6 + "\n\n" +

		"Domain: Agreeableness: " + D4 + "\n" +
		"Subdomains:\n" +
		"  Trust: " + A1 + "\n" +
		"  Morality: " + A2 + "\n" +
		"  Altruism: " + A3 + "\n" +
		"  Cooperation: " + A4 + "\n" +
		"  Modesty: " + A5 + "\n" +
		"  Sympathy: " + A6 + "\n\n" +

		"Domain: Conscientiousness: " + D5 + "\n" +
		"Subdomains:\n" +
		"  Self Efficacy: " + C1 + "\n" +
		"  Orderliness: " + C2 + "\n" +
		"  Dutifulness: " + C3 + "\n" +
		"  Achievement Striving: " + C4 + "\n" +
		"  Self Discipline: " + C5 + "\n" +
		"  Cautiousness: " + C6 + "\n\n" +

		"Note: Use this interpretation,\n" +
		"For Domain if score is <=20, it is low, <=30 is below average, <40 is average, <50 is above average, <=60 is high.\n" +
		"For Subdomain if score is <=3, it is low, <=4 is below average, <=6 is average, <=8 is above average, <=10 is high."

	return prompt
}

func GenerateContentFromText(pormpt string) error {
	location := "us-central1"
	modelName := "gemini-1.5-flash-001"

	ctx := context.Background()
	client, err := genai.NewClient(ctx, "cognify-438322", location)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	gemini := client.GenerativeModel(modelName)
	Genaiprompt := genai.Text(pormpt)
	resp, err := gemini.GenerateContent(ctx, Genaiprompt)
	if err != nil {
		return fmt.Errorf("error generating content: %w", err)
	}
	// See the JSON response in
	// https://pkg.go.dev/cloud.google.com/go/vertexai/genai#GenerateContentResponse.
	rb, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}
	fmt.Println(string(rb))
	return nil
}
