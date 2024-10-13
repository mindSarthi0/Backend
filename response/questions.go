package response

// Define the struct for questions
type Question struct {
	TestName string `json:"testName"`
	Question string `json:"question"`
	No       int    `json:"no"`
}
