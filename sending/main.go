package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	ip          = "127.0.0.1"
	outboxURL   = "http://" + ip + ":8003"
	inboxURL    = "http://" + ip + ":8005"
	blubBookURL = "http://" + ip + ":8000"
)

type Mail struct {
	To      string
	From    string
	Content string
}
type MailsSameDoamin map[string]Mail        // {id: Mail}, all mails in same domain
type PendingPool map[string]MailsSameDoamin // {Domain: MailsInDomain}
type PendingPoolResp struct {
	Code    int         `json:"code"`
	Message PendingPool `json:"message"`
}
type NetworkAddressResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type SendMailResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type MailBody struct {
	ID      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

var GPendingPool PendingPool

func main() {
	GPendingPool = make(map[string]MailsSameDoamin)
	pullPendingMail()
	c := time.Tick(10 * time.Second)
	for {
		<-c
		go pullPendingMail()
	}
}

func pullPendingMail() {
	// 先调用发件箱接口获得所有待发送邮件
	url := outboxURL + "/pendingMail"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	pendingPoolResp := PendingPoolResp{}
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &pendingPoolResp)
	fmt.Printf("%+v \n", pendingPoolResp)

	if pendingPoolResp.Code == 200 {
		message := pendingPoolResp.Message
		for k, v := range message {
			GPendingPool[k] = v
		}

		// 每个domain调用blue book的接口获得网络ip
		for domain, mails := range GPendingPool {
			url = blubBookURL + "/address/" + domain

			resp, err := http.Get(url)

			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			networkAddressResp := NetworkAddressResp{}
			body, err := ioutil.ReadAll(resp.Body)
			json.Unmarshal([]byte(body), &networkAddressResp)
			fmt.Printf("%+v \n", networkAddressResp)
			if networkAddressResp.Code != 200 {
				log.Print("Network address doesn't exist")
				continue
			}
			address := "http://" + networkAddressResp.Message + "/mails"

			// 向对应的address发送该domain的所有邮件
			bytesData, _ := json.Marshal(mails)
			resp, err = http.Post(address, "application/json;charset=utf-8", bytes.NewBuffer([]byte(bytesData)))

			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			sendMailResp := SendMailResp{}
			body, err = ioutil.ReadAll(resp.Body)
			json.Unmarshal([]byte(body), &sendMailResp)
			fmt.Printf("%+v \n", sendMailResp)

			if sendMailResp.Code != 200 {
				fmt.Println("Fail to send email, domain:", domain)
				continue
			}

			delete(GPendingPool, domain)
			fmt.Println(GPendingPool)
		}
	}

}
