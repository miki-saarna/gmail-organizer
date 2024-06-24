package gmail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var baseUrl string = "https://gmail.googleapis.com/gmail/v1/users/me"

type criteria struct {
	From string `json:"from"`
}

type action struct {
  AddLabelIds []string `json:"addLabelIds"`
}

type Filter struct {
  Criteria criteria `json:"criteria"`
  Action action `json:"action"`
}

type messageObj struct {
	Id string `json:"id"`
}

type unmarshalledRes struct {
	Messages []messageObj `json:"messages"`
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

func (c *Client) ListMessagesFromSender(sender string) ([]string, error) {
	url := fmt.Sprintf("%v/messages?q=from:%v",baseUrl, sender)
	reqMethod := "GET"

	req, err := http.NewRequest(reqMethod, url, nil)
	if err != nil {
		return nil, &requestCreationError{reqMethod, url, err.Error()}
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, &requestExecutionError{reqMethod, url, err.Error()}
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
    // body, _ := io.ReadAll(res.Body)
    return nil, fmt.Errorf("HTTP status %v for url %v", res.Status, url)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body for url %v: %v", url, err.Error())
	}

	var unmarshalledRes unmarshalledRes

	err = json.Unmarshal(body, &unmarshalledRes)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON response for url %v: %v", url, err.Error())
	}

	fmt.Printf("number of emails from sender %s found: %v\n", sender, len(unmarshalledRes.Messages))

	messages := make([]string, len(unmarshalledRes.Messages))
	for idx, message := range unmarshalledRes.Messages {
		messages[idx] = message.Id
	}

	return messages, nil
}

func (c *Client) RemoveMessages(messages []messageObj) {
	for _, message := range messages {
		url := fmt.Sprintf("%s/messages/%v/trash", baseUrl, message.Id)

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

		var unmarshalledMessageObj messageObj

		err = json.Unmarshal(body, &unmarshalledMessageObj)
		if err != nil {
			fmt.Printf("error unmarshalling JSON response from url %s: %v", url, err.Error())
			return
		}

		fmt.Printf("removed email with id: %v\n", unmarshalledMessageObj.Id)
	}
}

func (c *Client) BatchPermanentlyDeleteMessages(messageIds []string) (error) {
	url := fmt.Sprintf("%v/messages/batchDelete", baseUrl)
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

func (c *Client) AssignSenderToTrashList(sender string) (*Filter, error) {
	url := fmt.Sprintf("%v/settings/filters", baseUrl)
	reqMethod := "POST"
	reqBody := Filter{
		criteria{sender},
		action{[]string{"TRASH"}},
	}
	// reqBody := map[string]interface{}{
	// 	"criteria": map[string]interface{}{
	// 			"from": sender,
	// 	},
	// 	"action": map[string]interface{}{
	// 			"addLabelIds": []string{"TRASH"},
	// 	},
	// }

	// marshal/json-ize body
	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("there was an error marshalizing %v", reqBody)
	}
	reader := bytes.NewReader((data))

	req, err := http.NewRequest(reqMethod, url, reader)
	if err != nil {
		return nil, &requestCreationError{reqMethod, url, err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")

	// init HTTP request
	res, err := c.Do(req)
	if err != nil {
		return nil, &requestExecutionError{reqMethod, url, err.Error()}
	}
	defer res.Body.Close()

	// res.StatusCode validation necessary?
	if res.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("HTTP status %v for url %v", res.Status, url)
	}


	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("there was an error reading body: %v", res.Body)
	}

	unmarshalledFilter := Filter{}

	// unmarshal res.body
	err = json.Unmarshal(body, &unmarshalledFilter)
	if err != nil {
		return nil, fmt.Errorf("there was an error unmarshalling body: %v", err.Error())
	}

	// return res.body
	return &unmarshalledFilter, nil
}
