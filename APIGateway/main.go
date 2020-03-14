package main

import "log"

type Mail struct {
	ID      string
	To      string
	From    string
	Content string
}

func main() {
	r := InitHttpServerRouter()
	err := r.Run(":8001")
	if err != nil {
		log.Panic("API Gateway start failed.")
	}
}
