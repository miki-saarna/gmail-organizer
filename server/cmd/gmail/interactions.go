package gmail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

type requestCreationError struct {
	method, url, err string
}

type requestExecutionError struct {
	method, url, err string
}

func (r *requestCreationError) Error() string {
 return fmt.Sprintf("error creating %v request for url \"%v\": %v", r.method, r.url, r.err)
}

func (r *requestExecutionError) Error() string {
 return fmt.Sprintf("error executing %v request for url \"%v\": %v", r.method, r.url, r.err)
}

var baseUrl string = "https://gmail.googleapis.com/gmail/v1/users/me/messages"

func (c *Client) ListMessagesFromSender(sender string) ([]MessageObj, error) {
	url := fmt.Sprintf("%v?q=from:%v",baseUrl, sender)
	reqMethod := "GET"

	req, err := http.NewRequest(reqMethod, url, nil)
	if err != nil {
		return nil, &requestCreationError{reqMethod, url, err.Error()}
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, &requestExecutionError{reqMethod, url, err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
    // body, _ := io.ReadAll(resp.Body)
    return nil, fmt.Errorf("HTTP status %v for url %v", resp.Status, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body for url %v: %v", url, err.Error())
	}

	var apiResp ApiResp

	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON response for url %v: %v", url, err.Error())
	}

	fmt.Printf("number of emails from sender %s found: %v\n", sender, len(apiResp.Messages))
	return apiResp.Messages, nil
}

func (c *Client) RemoveMessages(messages []MessageObj) {
	for _, message := range messages {
		url := fmt.Sprintf("%s/%v/trash", baseUrl, message.Id)

		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			fmt.Printf("error creating post request for url %s: %v", url, err.Error())
			return
		}

		resp, err := c.Do(req)
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
			fmt.Printf("error unmarshalling JSON response from url %s: %v", url, err.Error())
			return
		}

		fmt.Printf("removed email with id: %v\n", trashApiResp.Id)
	}
}

func (c *Client) BatchPermanentlyDeleteMessages(messageIds []string) (error) {
	url := fmt.Sprintf("%v/batchDelete", baseUrl)
	reqBody := batchDeleteBody{ Ids: messageIds }
	reqMethod := "POST"

	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("error unmsharlling reqBody: %v\nreqBody: %v", err.Error(), reqBody)
	}
	reader := bytes.NewReader((data))

	req, err := http.NewRequest(reqMethod, url, reader)
	if err != nil {
		return &requestCreationError{reqMethod, url, err.Error()}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return &requestCreationError{reqMethod, url, err.Error()}
	}
	defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	fmt.Printf("an error occurred, status code: %v", resp.StatusCode)
	// 	return
	// }

	fmt.Printf("Permanently deleted emails with ID:")
	for _, id := range messageIds {
		fmt.Printf("\n-%s", id)
	}

	return nil
}