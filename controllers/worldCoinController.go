package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Guesstrain/ethglobal/models"
	"github.com/gin-gonic/gin"
)

func VerifyHandler(c *gin.Context) {
	var payload models.RequestPayload

	// Bind the JSON payload from the request body to the struct
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Forward the request to the external API
	externalAPIURL := "https://example.com/api/v2/verify/{app_id}" // Replace with your actual API URL

	// Convert the payload struct to JSON for forwarding
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payload"})
		return
	}

	// Create a new POST request to the external API
	req, err := http.NewRequest("POST", externalAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create external request"})
		return
	}

	// Set the content type to application/json
	req.Header.Set("Content-Type", "application/json")

	// Send the request to the external API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request to external API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request to external API"})
		return
	}
	defer resp.Body.Close()

	// Read the response from the external API
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response from external API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read external response"})
		return
	}

	// Forward the external API's response back to the client
	c.String(resp.StatusCode, string(body))
}
