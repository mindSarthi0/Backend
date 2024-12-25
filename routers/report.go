package routers

import (
	controller "myproject/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleBig5Report(c *gin.Context) {
	// Retrieve the dynamic testId from the route parameter
	testId := c.Param("testId")

	reports, err := controller.GetReportsByTestId(testId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reports", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reports)
}
