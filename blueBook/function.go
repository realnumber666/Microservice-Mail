package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func handleAddAddress(c *gin.Context) {
	var addressBody AddressBody
	err := c.BindJSON(&addressBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "POST DATA ERROR"})
		return
	}

	addressMap[addressBody.MailAddress] = addressBody.NetworkAddress

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Suc to Add address.",
	})
	return
}

func handleGetAddress(c *gin.Context) {
	mailAddress := c.Param("mailAddress")
	networkAddress, ok := addressMap[mailAddress]

	if ok {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": networkAddress,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "No network address with mail address " + mailAddress + ".",
		})
	}

	return
}
