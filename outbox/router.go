package main

import "github.com/gin-gonic/gin"

func InitHttpServerRouter() *gin.Engine {
	server := gin.Default()

	// Add mail in pending pool
	server.POST("/mail", handleAddMail)

	// Return all pending mail and clear mails
	server.GET("/pendingMail", handleReturnAndClearMails)

	// Return all pending mail
	server.GET("/userPendingMail/:username", handleReturnMails)

	return server
}
