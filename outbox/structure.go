package main

type MailBody struct {
	ID      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}
