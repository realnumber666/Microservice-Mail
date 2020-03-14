package main

import (
	"log"
	"sync"
)

type Mail struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Content string `json:"content"`
}
type MailsSameUser struct {
	sync.RWMutex
	m map[string]Mail // {id: Mail}, all mails in same domain
}

var incomingPool = struct {
	sync.RWMutex
	m map[string]MailsSameUser // {username: MailsInDomain}
}{}

func init() {
	incomingPool.m = make(map[string]MailsSameUser)
}

func main() {
	r := InitHttpServerRouter()
	err := r.Run(":8005")
	if err != nil {
		log.Panic("Inbox start failed.")
	}
}
