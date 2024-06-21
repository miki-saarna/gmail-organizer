package gmail

import (
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

var baseUrl string = "https://gmail.googleapis.com/gmail/v1/users/me/messages"

func ListMessagesFromSender(client *http.Client, sender string) ([]MessageObj, error) {
	url := fmt.Sprintf("%v?q=from:%v",baseUrl, sender)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errMessage := fmt.Errorf("Error creating GET request for URL \"%v\": %v", url, err.Error())
		fmt.Println(errMessage)
		return nil, errMessage
	}

	resp, err := client.Do(req)
	if err != nil {
		errMessage := fmt.Errorf("Error executing GET request for URL \"%v\": %v", url, err.Error())
		fmt.Println(errMessage)
		return nil, errMessage
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
    body, _ := io.ReadAll(resp.Body)
    errMessage := fmt.Errorf("HTTP status: %v\nResponse body: %v", resp.Status, string(body))
		fmt.Println(errMessage)
    return nil, fmt.Errorf("HTTP status %v for URL %v", resp.Status, url)
}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errMessage := fmt.Errorf("Error reading response body for URL %v: %v", url, err.Error())
		fmt.Println(errMessage)
		return nil, errMessage
	}

	var apiResp ApiResp

	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		errMessage := fmt.Errorf("Error unmarshalling JSON response for URL %v: %v", url, err.Error())
		fmt.Println(errMessage)
		return nil, errMessage
	}

	fmt.Printf("apiResp: %v", len(apiResp.Messages))
	return apiResp.Messages, nil


}