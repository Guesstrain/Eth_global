package controllers

import (
	"net/http"

	"github.com/Guesstrain/ethglobal/models"
	"github.com/Guesstrain/ethglobal/services"
	"github.com/gin-gonic/gin"
)

var prizeService services.PrizeService = services.NewPrizeService(databaseService)

func InsertPrize(c *gin.Context) {
	var prize models.PrizeList

	// Bind JSON input to the prize struct
	if err := c.ShouldBindJSON(&prize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to insert the prize
	if err := prizeService.InsertPrize(prize); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert prize: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Prize created successfully",
		"prize":   prize,
	})
}

func UpdatePrize(c *gin.Context) {
	prizeName := c.Param("prize_name")
	var updatedPrize models.PrizeList

	// Bind the JSON input to the updatedPrize struct
	if err := c.ShouldBindJSON(&updatedPrize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to update the prize
	if err := prizeService.UpdatePrize(prizeName, updatedPrize); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update prize: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Prize updated successfully",
		"prize":   updatedPrize,
	})
}

func DistributePrize(c *gin.Context) {
	// Call the service to distribute prizes
	prizes := prizeService.DistributePrize()

	if len(prizes) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to distribute prizes or no wallets found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Prizes distributed successfully",
		"prizes":  prizes,
	})
}
