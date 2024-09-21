package main

import (
	"github.com/Guesstrain/ethglobal/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	router.POST("/addwallet", controllers.InsertWallet)
	router.GET("/wallet/:address", controllers.SelectWallet)
	router.GET("/wallets/time-period", controllers.QueryWalletsByTimePeriod)
	router.DELETE("/wallet/:address", controllers.DeleteWallet)
	router.POST("/addPrize", controllers.InsertPrize)
	router.PUT("/prizes/:prize_name", controllers.UpdatePrize)
	router.GET("/distribute-prizes", controllers.DistributePrize)
	router.POST("/Verify", controllers.VerifyHandler)

	router.Run(":8080")
}
