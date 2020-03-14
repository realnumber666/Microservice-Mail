package main

import "github.com/gin-gonic/gin"

func InitHttpServerRouter() *gin.Engine {
	server := gin.Default()

	// Get one's mails
	server.GET("/mails/:username", handleGetMailList)

	// Get single mail
	server.GET("/mail", handleGetMail)

	// Delete mail
	server.DELETE("/mail", handleDeleteMail)

	// Add mail in incoming pool
	server.POST("/mail", handleAddMail)

	return server
}
