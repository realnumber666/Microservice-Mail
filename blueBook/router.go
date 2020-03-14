package main

import "github.com/gin-gonic/gin"

func InitHttpServerRouter() *gin.Engine  {
	server := gin.Default()

	// Add network address
	server.POST("/address", handleAddAddress)

	// Get network address
	server.GET("/address/:mailAddress", handleGetAddress)


	return server
}