package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func handleAddMail(c *gin.Context) {
	var mailBody MailBody
	err := c.BindJSON(&mailBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "POST DATA ERROR"})
		return
	}

	domain := strings.Split(mailBody.To, "@")[1]

	pendingPool.RLock()
	mails := pendingPool.m[domain]
	pendingPool.RUnlock()

	if mails == nil {
		mails = make(map[string]Mail)
	}

	mails[mailBody.ID] = Mail{
		To:      mailBody.To,
		From:    mailBody.From,
		Content: mailBody.Content,
	}

	pendingPool.Lock()
	pendingPool.m[domain] = mails
	pendingPool.Unlock()

	// Generate fromPool
	from := strings.Split(mailBody.From, "@")[0]

	fromPool.RLock()
	fromMails := fromPool.m[from]
	fromPool.RUnlock()

	if fromMails == nil {
		fromMails = make(map[string]Mail)
	}

	fromMails[mailBody.ID] = Mail{
		To:      mailBody.To,
		From:    mailBody.From,
		Content: mailBody.Content,
	}

	fromPool.Lock()
	fromPool.m[from] = fromMails
	fromPool.Unlock()

	log.Print(pendingPool.m)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Suc to Add mail.",
	})
	return
}

func handleReturnAndClearMails(c *gin.Context) {
	pendingPool.RLock()
	data := pendingPool.m
	pendingPool.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": data,
	})

	// clear pending pool
	pendingPool.Lock()
	pendingPool.m = make(map[string]MailsSameDoamin)
	pendingPool.Unlock()

	fromPool.Lock()
	fromPool.m = make(map[string]MailsSameFrom)
	fromPool.Unlock()

	return
}

func handleReturnMails(c *gin.Context) {
	username := c.Param("username")

	fromPool.RLock()
	data := fromPool.m[username]
	fromPool.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": data,
	})

	return
}
