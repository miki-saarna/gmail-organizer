package main

import (
	"encoding/json"
	"fmt"
	"gmail-organizer/cmd/gmail"
	"gmail-organizer/utils"
	"io"
	"os"
)

type deletionList []string

func main () {
	var deletionList deletionList

	jsonFile, err := os.Open("deletionList.json") // needs absolute path?
	if err != nil {
		fmt.Printf("could not open json list: %v", err.Error())
		return
	}
	defer jsonFile.Close()

	data, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Printf("could not read json list: %v", err.Error())
		return
	}

	err = json.Unmarshal(data, &deletionList)
	if err != nil {
		fmt.Printf("could not unmarshal json list: %v", err.Error())
		return
	}
	// fmt.Printf("unmarshalled: %v", deletionList)

	var senderBulletPointList string
	for _, senderAddress := range deletionList {
		senderBulletPointList += fmt.Sprintf("\n- %v", senderAddress)
	}

	isConfirmed := utils.AskForConfirmation(fmt.Sprintf("Are you sure you would like to permanently delete all emails from the following senders: %v?", senderBulletPointList))
	if (isConfirmed) {
		gmail.Main(deletionList);
	}

	// if len(os.Args) < 2 {
	// 	fmt.Println("Please provide an email address to search your inbox for")
	// 	return
	// }
	// senderAddress := os.Args[1]
}