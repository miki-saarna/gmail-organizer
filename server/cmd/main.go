package main

import (
	"encoding/json"
	"fmt"
	"gmail-organizer/cmd/gmail"
	"gmail-organizer/utils"
	"io"
	"log"
	"os"
)

const (
	deletion = "Delete emails"
	updateSpam =  "Update SPAM list"
	exit =  "Exit"
)

type deletionList []string


func main () {
	var deletionList deletionList

	jsonFile, err := os.Open("deletionList.json") // relative path?
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

	options := utils.Options{deletion, updateSpam, exit}
	selectedOption, err := options.SelectOption()
	if err != nil {
		log.Fatalf("There was an error selecting an option: %v", err)
	}

	if selectedOption == deletion {
		var senderBulletPointList string
		for _, senderAddress := range deletionList {
			senderBulletPointList += fmt.Sprintf("\n- %v", senderAddress)
		}

		confirmationMsg := utils.ConfirmationMsg(fmt.Sprintf("Are you sure you would like to permanently delete all emails from the following senders: %v?", senderBulletPointList))
		isConfirmed, err := confirmationMsg.AskForConfirmation()
		if err != nil {
			log.Fatalf("There was an error selection an option: %v", err)
		} else if (isConfirmed) {
			gmail.InitMessageRemoval(deletionList);
		}
	} else if selectedOption == updateSpam {
		// create update SPAM list func
	}

	// if len(os.Args) < 2 {
	// 	fmt.Println("Please provide an email address to search your inbox for")
	// 	return
	// }
	// senderAddress := os.Args[1]
}