package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"myproject/models"
)

// Define the album structure
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// Seed some initial data (stored in memory)
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func init() {
	// Setup the mgm default config
	err := mgm.SetDefaultConfig(nil, "mgm_lab", options.Client().ApplyURI("mongodb://root:12345@localhost:27017"))
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Bind the JSON data in the request to the newAlbum object.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the list of albums.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID matches the id parameter sent by the client.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Look through the albums slice for an album with a matching ID.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func getQuestions(c *gin.Context) {
	question := models.NewQuestion("Math Test", "What is 2 + 2?")
	c.IndentedJSON(http.StatusOK, albums)
}
func main() {
	router := gin.Default()

	// Define the endpoints and handlers
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.GET("/albums/:id", getAlbumByID)

	router.GET("/questions", getQuestions)
	// router.POST("/questions", postQuestions)

	// Start the server on localhost:8080
	router.Run("localhost:8080")
}
