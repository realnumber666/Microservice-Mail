package main

import "github.com/gin-gonic/gin"

func InitHttpServerRouter() *gin.Engine {
	server := gin.Default()

	// Add mail in incoming pool
	server.POST("/mails", handleAddMails)

	return server
}
