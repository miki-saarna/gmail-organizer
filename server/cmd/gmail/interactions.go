package gmail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)


type MessageObj struct {
	Id string `json:"id"`
}

type ApiResp struct {
	Messages []MessageObj `json:"messages"`
}

type TrashApiResp struct {
	Message MessageObj `json:"message"`
}

type batchDeleteBody struct {
	Ids []string `json:"ids"`
}

var baseUrl string = "https://gmail.googleapis.com/gmail/v1/users/me/messages"

func ListMessagesFromSender(client *http.Client, sender string) ([]MessageObj, error) {
	url := fmt.Sprintf("%v?q=from:%v",baseUrl, sender)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errMessage := fmt.Errorf("error creating get request for url \"%v\": %v", url, err.Error())
		fmt.Println(errMessage)
		return nil, errMessage
	}

	resp, err := client.Do(req)
	if err != nil {
		errMessage := fmt.Errorf("error executing get request for url \"%v\": %v", url, err.Error())
		fmt.Println(errMessage)
		return nil, errMessage
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
    body, _ := io.ReadAll(resp.Body)
    errMessage := fmt.Errorf("HTTP status: %v\nResponse body: %v", resp.Status, string(body))
		fmt.Println(errMessage)
    return nil, fmt.Errorf("HTTP status %v for url %v", resp.Status, url)
}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errMessage := fmt.Errorf("error reading response body for url %v: %v", url, err.Error())
		fmt.Println(errMessage)
		return nil, errMessage
	}

	var apiResp ApiResp

	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		errMessage := fmt.Errorf("error unmarshalling JSON response for url %v: %v", url, err.Error())
		fmt.Println(errMessage)
		return nil, errMessage
	}

	fmt.Printf("number of emails from sender %s found: %v\n", sender, len(apiResp.Messages))
	return apiResp.Messages, nil
}

func RemoveMessages(client *http.Client, messages []MessageObj) {
	for _, message := range messages {
		url := fmt.Sprintf("%s/%v/trash", baseUrl, message.Id)

		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			fmt.Printf("error creating post request for url %s: %v", url, err.Error())
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("error executing post request for url %s: %v", url, err.Error())
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("HTTP status: %v\nResponse body: %s", resp.Status, string(body))
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("error reading response body for url %s: %v", url, err.Error())
			return
		}

		var trashApiResp MessageObj

		err = json.Unmarshal(body, &trashApiResp)
		if err != nil {
			fmt.Printf("error unmarshalling JSOn response from url %s: %v", url, err.Error())
			return
		}

		fmt.Printf("removed email with id: %v\n", trashApiResp.Id)
	}
}

func BatchPermanentlyDeleteMessages (client *http.Client, messageIds []string) {
	url := fmt.Sprintf("%v/batchDelete", baseUrl)
	reqBody := batchDeleteBody{ Ids: messageIds }

	data, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal(err)
	}
	reader := bytes.NewReader((data))

	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Printf("an error occurred: %v", err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("an error occurred: %v", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("an error occurred, status code: %v", resp.StatusCode)
		return
	}

	fmt.Printf("removed emails with id:")
	for _, id := range messageIds {
		fmt.Printf("\n-%s", id)
	}

}