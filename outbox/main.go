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
type MailsSameDoamin map[string]Mail // {id: Mail}, all mails in same domain

var pendingPool = struct {
	sync.RWMutex
	m map[string]MailsSameDoamin // {Domain: MailsInDomain}
}{}

type MailsSameFrom map[string]Mail // {id: Mail}, all mails in same from

var fromPool = struct {
	sync.RWMutex
	m map[string]MailsSameFrom // {from: MailsSameFrom}
}{}

func init() {
	pendingPool.m = make(map[string]MailsSameDoamin)
	fromPool.m = make(map[string]MailsSameFrom)
}

func main() {
	r := InitHttpServerRouter()
	err := r.Run(":8003")
	if err != nil {
		log.Panic("Outbox start failed.")
	}
}
