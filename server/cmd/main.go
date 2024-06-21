package main

import (
	"fmt"
	"gmail-organizer/cmd/gmail"
	"os"
)

func main () {
	if len(os.Args) < 2 {
		fmt.Println("Please provide an email address to search your inbox for")
		return
	}

	senderAddress := os.Args[1]
	gmail.Main(senderAddress);
}