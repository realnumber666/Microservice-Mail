package main

type MailsBody map[string]Mail

type MailToSend struct {
	ID      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}
