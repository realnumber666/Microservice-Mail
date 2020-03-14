package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	ip        = "127.0.0.1"
	outboxURL = "http://" + ip + ":8003"
	inboxURL  = "http://" + ip + ":8005"
)

func handleSendMail(c *gin.Context) {
	// 从结构体中解析出原始mail
	var mailBody MailBody
	err := c.BindJSON(&mailBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "POST DATA ERROR"})
		return
	}

	to := mailBody.To
	from := mailBody.From
	// 校验邮件地址是否正确
	check_to := checkEmail(to)
	check_from := checkEmail(from)

	if !(check_to && check_from) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Irregular email addresses"})
		return
	}

	// 校验是否是自己发给自己的
	if to == from {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Can't send it to yourself"})
		return
	}

	// Generate id by md5(timestamp+from+content)
	currentTime := time.Now().Unix()
	timeStr := strconv.FormatInt(currentTime, 10)

	id := createIdByMD5(timeStr + mailBody.From + mailBody.Content)
	// 封装成处理后的mail结构体
	mail := Mail{
		ID:      id,
		To:      mailBody.To,
		From:    mailBody.From,
		Content: mailBody.Content,
	}

	// 调用outbox的接口
	bytesData, _ := json.Marshal(mail)
	resp, err := http.Post(outboxURL+"/mail", "application/json;charset=utf-8", bytes.NewBuffer([]byte(bytesData)))

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	sendMailResp := CommonResp{}
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &sendMailResp)
	fmt.Printf("%+v \n", sendMailResp)

	if sendMailResp.Code != 200 {
		fmt.Println("Fail to send email, id: ", id)
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Fail to send email, id: " + id})
		return
	}

	// 根据返回码返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": mail,
	})
	return
}

func handleGetMailList(c *gin.Context) {
	// 获取username参数
	username := c.Param("username")

	// 调用inbox的接口获得邮件列表
	resp, err := http.Get(inboxURL + "/mails/" + username)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	mailsSameUserResp := MailsSameUserResp{}
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &mailsSameUserResp)
	fmt.Printf("%+v \n", mailsSameUserResp)

	mailsSameUser := mailsSameUserResp.Message
	if mailsSameUser == nil {
		mailsSameUser = MailsSameUser{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": mailsSameUser,
	})
	return
}

func handleGetMail(c *gin.Context) {
	// 获取username和id参数
	username := c.Query("username")
	id := c.Query("id")

	// 调用inbox接口获得邮件列表
	resp, err := http.Get(inboxURL + "/mail?username=" + username + "&id=" + id)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	mailResp := MailResp{}
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &mailResp)
	fmt.Printf("%+v \n", mailResp)

	if mailResp.Code != 200 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Fail to get email, id: " + id,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": mailResp.Message,
	})
	return
}

func handleDeleteMail(c *gin.Context) {
	// 获取username和id
	username := c.Query("username")
	id := c.Query("id")

	// 调用inbox接口删除该邮件
	req, _ := http.NewRequest("DELETE", inboxURL+"/mail?username="+username+"&id="+id, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	deleteResp := CommonResp{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &deleteResp)
	fmt.Printf("%+v \n", deleteResp)

	if deleteResp.Code != 200 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "Fail to delete email, id: " + id,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": deleteResp.Message,
	})
	return
}

func handleGetOutboxMailList(c *gin.Context) {
	// 获取username参数
	username := c.Param("username")

	// 调用outbox的接口获得邮件列表
	url := outboxURL + "/userPendingMail/" + username
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fromPoolResp := FromPoolResp{}
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &fromPoolResp)
	fmt.Printf("%+v \n", fromPoolResp)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": fromPoolResp.Message,
	})
	return
}

func checkEmail(email string) (b bool) {
	if m, _ := regexp.MatchString("^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+", email); !m {
		return false
	}
	return true
}

func createIdByMD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
