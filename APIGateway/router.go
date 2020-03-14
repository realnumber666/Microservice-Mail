package main

import "github.com/gin-gonic/gin"

func InitHttpServerRouter() *gin.Engine {
	server := gin.Default()

	// Send mail
	server.POST("/mail", handleSendMail)

	// Get mail list in inbox
	server.GET("/mailList/:username", handleGetMailList)

	// Get mail list in outbox
	server.GET("/outboxMailList/:username", handleGetOutboxMailList)

	// Get single mail content
	server.GET("/mail", handleGetMail)

	// Delete single mail
	server.DELETE("/mail", handleDeleteMail)

	return server
}
