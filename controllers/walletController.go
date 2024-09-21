package controllers

import (
	"net/http"
	"time"

	"github.com/Guesstrain/ethglobal/database"
	"github.com/Guesstrain/ethglobal/models"
	"github.com/Guesstrain/ethglobal/services"
	"github.com/gin-gonic/gin"
)

var dsn = "root:password@tcp(127.0.0.1:3306)/eth_global?parseTime=true"
var databaseService database.DatabaseService = database.NewGormDatabaseService(dsn)
var walletService services.WalletService = services.NewWalletService(databaseService)

func InsertWallet(c *gin.Context) {
	var wallet models.Wallet

	// Bind JSON to wallet struct
	if err := c.ShouldBindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet.Timestamp = time.Now()

	err := walletService.InsertWallet(wallet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Wallet created successfully",
		"wallet":  wallet,
	})
}

func SelectWallet(c *gin.Context) {
	address := c.Param("address") // Extract the wallet address from the URL parameters

	// Call the wallet service to get the wallet by address
	wallet, err := walletService.SelectWalletByAddress(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve wallet: " + err.Error()})
		return
	}

	if wallet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wallet": wallet,
	})
}

func DeleteWallet(c *gin.Context) {
	address := c.Param("address") // Extract the wallet address from the URL parameters

	// Call the wallet service to delete the wallet by address
	err := walletService.DeleteWalletByAddress(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete wallet: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Wallet deleted successfully",
	})
}

func QueryWalletsByTimePeriod(c *gin.Context) {
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format"})
		return
	}
	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format"})
		return
	}

	wallets, err := databaseService.QueryWalletsByTimePeriod(startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query wallets: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallets": wallets})
}
