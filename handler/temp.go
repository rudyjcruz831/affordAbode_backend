package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type TemporaryLocalListing struct {
	Address   string  `json:"address"`
	Price     float64 `json:"price"`
	Bedrooms  float64 `json:"beds"`
	Bathrooms float64 `json:"baths"`
}

// TemporaryLocalListingHandler handles reading from local JSON file
func (h *Handler) TemporaryLocalListingHandler(c *gin.Context) {
	file, err := os.Open("/Users/rudy/Desktop/Rudy-macbook/CMPE_272_Spring2025/cmpe_272_project_spring2025-/affordAbode_backend/handler/message.json") // Update with correct file path
	if err != nil {
		log.Printf("Error opening JSON file: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not open JSON file: " + err.Error(),
		})
		return
	}
	defer file.Close()

	var tempListing []TemporaryLocalListing
	if err := json.NewDecoder(file).Decode(&tempListing); err != nil {
		log.Printf("Error decoding JSON file: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not decode JSON file: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Listings loaded successfully",
		"listings": tempListing,
	})
}
