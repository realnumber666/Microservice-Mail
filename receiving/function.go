package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	ip        = "127.0.0.1"
	outboxURL = "http://" + ip + ":8003"
	inboxURL  = "http://" + ip + ":8005"
)

type SendMailResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func handleAddMails(c *gin.Context) {
	var mailBody MailsBody
	err := c.BindJSON(&mailBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "POST DATA ERROR"})
		return
	}

	for id, mail := range mailBody {
		go PostToInbox(id, mail)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Suc to Add mails.",
	})
	return
}

func PostToInbox(id string, mail Mail) {
	username := strings.Split(mail.To, "@")[0]
	url := inboxURL + "/mail?username=" + username

	mailToSend := MailToSend{
		ID:      id,
		From:    mail.From,
		To:      mail.To,
		Content: mail.Content,
	}
	bytesData, _ := json.Marshal(mailToSend)
	resp, err := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer([]byte(bytesData)))

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	sendMailResp := SendMailResp{}
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &sendMailResp)
	fmt.Printf("%+v \n", sendMailResp)

	if sendMailResp.Code != 200 {
		fmt.Println("Fail to send email, username: ", username, ", id: ", id)
	}
}
