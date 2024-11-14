package response

// Define the struct for questions
type Answers struct {
	Id     string `json:"id"` // Assuming ID is a string
	Answer string `json:"answer"`
}

// Define the main struct
type Submit struct {
	Email   string    `json:"email"`
	Name    string    `json:"name"`
	Age     int       `json:"age"`    // Assuming age is an integer
	Gender  string    `json:"gender"` // Assuming gender is a string
	Answers []Answers `json:"answers"`
	PMode   string    `json:"pMode"`
}
