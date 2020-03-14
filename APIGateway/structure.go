package main

type MailBody struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

type CommonResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type MailResp struct {
	Code    int `json:"code"`
	Message MailBody
}

type MailsSameUser map[string]MailBody // {id: Mail}, all mails in same domain

type MailsSameUserResp struct {
	Code    int           `json:"code"`
	Message MailsSameUser `json:"message"`
}

type FromPool map[string]MailBody // {Domain: MailsInDomain}
type FromPoolResp struct {
	Code    int      `json:"code"`
	Message FromPool `json:"message"`
}

type MailsSameDoamin map[string]MailBody    // {id: Mail}, all mails in same domain
type PendingPool map[string]MailsSameDoamin // {Domain: MailsInDomain}
type PendingPoolResp struct {
	Code    int         `json:"code"`
	Message PendingPool `json:"message"`
}
