package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func handleGetMailList(c *gin.Context) {
	username := c.Param("username")

	incomingPool.RLock()
	mails := incomingPool.m[username].m
	incomingPool.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": mails,
	})
	return
}

func handleGetMail(c *gin.Context) {
	username := c.Query("username")
	id := c.Query("id")

	incomingPool.RLock()
	mails := incomingPool.m[username]
	incomingPool.RUnlock()

	mails.RLock()
	if mails.m == nil {
		mails.Unlock()
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Fail to get mail",
		})
	} else {
		mails.RUnlock()

		mails.Lock()
		mail, ok := mails.m[id]
		mails.Unlock()

		if ok == false {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "Fail to get mail",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": mail,
			})
		}
	}
	return
}

func handleDeleteMail(c *gin.Context) {
	username := c.Query("username")
	id := c.Query("id")

	incomingPool.RLock()
	mails := incomingPool.m[username]
	incomingPool.RUnlock()

	mails.RLock()
	if mails.m == nil {
		mails.RUnlock()
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "No such mail",
		})
	} else {
		mails.RUnlock()

		mails.Lock()
		delete(mails.m, id)
		mails.Unlock()

		incomingPool.Lock()
		incomingPool.m[username] = mails
		incomingPool.Unlock()

		log.Printf("%+v", incomingPool.m)
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "Suc to delete mail",
		})
	}

	return
}

func handleAddMail(c *gin.Context) {
	var mailBody MailBody
	err := c.BindJSON(&mailBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "POST DATA ERROR"})
		return
	}
	username := c.Query("username")

	incomingPool.RLock()
	mails := incomingPool.m[username]
	incomingPool.RUnlock()

	mails.Lock()
	if mails.m == nil {
		mails.m = make(map[string]Mail)
	}

	mails.m[mailBody.ID] = Mail{
		To:      mailBody.To,
		From:    mailBody.From,
		Content: mailBody.Content,
	}
	mails.Unlock()

	incomingPool.Lock()
	incomingPool.m[username] = mails
	incomingPool.Unlock()

	log.Print(incomingPool.m)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Suc to Add mail.",
	})
	return
}
