package API

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/vertexai/apiv1beta1"
	genpb "cloud.google.com/go/vertexai/apiv1beta1/genprotopb"
	"google.golang.org/api/option"
)

var (
	client    *vertexai.TextGenerationClient
	modelName string
	ctx       context.Context
)

// init function to initialize the client and set up context
func init() {
	projectID := "cognify-438322"
	location := "us-central1"
	modelName = fmt.Sprintf("projects/%s/locations/%s/publishers/google/models/gemini-1.5-flash-001", projectID, location)

	// Initialize context
	ctx = context.Background()

	// Initialize the Vertex AI client
	var err error
	client, err = vertexai.NewTextGenerationClient(ctx, option.WithEndpoint(fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)))
	if err != nil {
		log.Fatalf("Error creating Vertex AI client: %v", err)
	}
}

func generateContentFromText(w io.Writer) error {
	// Prepare the request with the desired prompt
	prompt := "Create a report "
	req := &genpb.GenerateTextRequest{
		Model:     modelName,
		InputText: prompt,
	}

	// Call the model to generate content
	resp, err := client.GenerateText(ctx, req)
	if err != nil {
		return fmt.Errorf("error generating content: %w", err)
	}

	// Check if candidates are returned
	if len(resp.GetCandidates()) == 0 {
		return fmt.Errorf("no candidates received in response")
	}

	// Get the output from the first candidate
	outputText := resp.GetCandidates()[0].GetOutput()

	// Log the generated output for debugging purposes
	log.Printf("Generated text: %s", outputText)

	// Write the response directly into JSON format
	err = json.NewEncoder(w).Encode(struct {
		GeneratedText string `json:"generated_text"`
	}{
		GeneratedText: outputText,
	})
	if err != nil {
		return fmt.Errorf("error encoding response to JSON: %w", err)
	}

	return nil
}
