package main

import "log"

// {MailAddress: NetworkAddress}
var addressMap map[string]string

func init() {
	addressMap = make(map[string]string)
}

func main() {
	r := InitHttpServerRouter()
	err := r.Run(":8000")
	if err != nil {
		log.Panic("Blue book start failed.")
	}
}