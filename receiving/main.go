package main

import "log"

type Mail struct {
	To      string
	From    string
	Content string
}
type MailsSameDoamin map[string]Mail // {id: Mail}, all mails in same domain

var incomingPool map[string]MailsSameDoamin // {username: MailsInDomain}

func init() {
	incomingPool = make(map[string]MailsSameDoamin)
}

func main() {
	r := InitHttpServerRouter()
	err := r.Run(":8002")
	if err != nil {
		log.Panic("Receiving start failed.")
	}
}
