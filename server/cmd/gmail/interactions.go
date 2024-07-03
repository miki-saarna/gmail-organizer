package gmail

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/tebeka/selenium"
)

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

type FiltersList struct {
	Filter []Filter `json:"filter"`
}

type TrashList map[string]struct{}

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

type messagePayload struct {
	Payload struct{Headers []struct{Name string `json:"name"`; Value string `json:"value"`} `json:"headers"`} `json:"payload"`
}

const baseUrl = "https://gmail.googleapis.com/gmail/v1/users/me"

var (
	chromeDriverPath 	string
	port             	string
	profile 				 	string
	userDataDirectory string
)

func (r *requestCreationError) Error() string {
 return fmt.Sprintf("error creating %v request for url \"%v\": %v", r.method, r.url, r.err)
}

func (r *requestExecutionError) Error() string {
 return fmt.Sprintf("error executing %v request for url \"%v\": %v", r.method, r.url, r.err)
}

func (c *Client) ListMessagesFromSender(senderAddresses []string, maxResults int) ([]string, error) {
	query := senderAddresses[0]
	for i := 1; i < len(senderAddresses); i++ {
		query += fmt.Sprintf(" OR %v", senderAddresses[i])
	}
	encodedQuery := url.QueryEscape(query)
	
	url := fmt.Sprintf("%s/messages?q=maxResults=%d&%s", baseUrl, maxResults, encodedQuery)
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

	fmt.Printf("Number of emails found: %v\n", len(unmarshalledRes.Messages))

	messages := make([]string, len(unmarshalledRes.Messages))
	for idx, message := range unmarshalledRes.Messages {
		messages[idx] = message.Id
	}

	return messages, nil
}

func (c *Client) UnsubscribeWithWebDriver(msgIds []string) error {
	chromeDriverPath 	= os.Getenv("CHROME_DRIVER_PATH")
	port             	= os.Getenv("WEB_DRIVER_PORT")
	userDataDirectory = os.Getenv("CHROME_USER_DATA_DIRECTORY")
	profile 				 	= os.Getenv("CHROME_PROFILE")

	opts := []selenium.ServiceOption{
		selenium.ChromeDriver(chromeDriverPath),
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("error converting port string value to int: %s", err.Error())
	}

	service, err := selenium.NewChromeDriverService(chromeDriverPath, portInt, opts...)
	if err != nil {
		return fmt.Errorf("error starting ChromeDriver server: %s", err.Error())
	}
	defer service.Stop()

	capabilities := selenium.Capabilities{"browserName": "chrome"}

	chromeArgs := []string{
		fmt.Sprintf("user-data-dir=%s", userDataDirectory),
		fmt.Sprintf("profile-directory=%s", profile),
		// "--headless",
	}

	capabilities["goog:chromeOptions"] = map[string]interface{}{
		"args": chromeArgs,
	}

	wd, err := selenium.NewRemote(capabilities, fmt.Sprintf("http://localhost:%d/wd/hub", portInt))
	if err != nil {
		return fmt.Errorf("error connecting to WebDriver server: %s", err.Error())
	}
	defer wd.Quit()

	var failedMsgIds []string

	for _, msgId := range msgIds {
		messageUrl := fmt.Sprintf("https://mail.google.com/mail/u/0/#inbox/%s", msgId)
		if err := wd.Get(messageUrl); err != nil {
			failedMsgIds = append(failedMsgIds, fmt.Sprintf("error navigating to URL for message ID %s: %s", msgId, err.Error()))
			continue
		}

		time.Sleep(1 * time.Second)

		xpath := `//span[contains(@class, 'Ca') and contains(text(), 'Unsubscribe')]`
		elem, err := wd.FindElement(selenium.ByXPATH, xpath)
		if err != nil {
			failedMsgIds = append(failedMsgIds, fmt.Sprintf("error trying to find \"Subscribe\" span element for message ID %s: %s", msgId, err.Error()))
			continue
		}

		if err := elem.Click(); err != nil {
			failedMsgIds = append(failedMsgIds, fmt.Sprintf("error trying to click \"Subscribe\" span element for message ID %s: %s", msgId, err.Error()))
			continue
		}

		xpath = `//button[contains(text(), 'Unsubscribe')]`
		elem, err = wd.FindElement(selenium.ByXPATH, xpath)
		if err != nil {
			failedMsgIds = append(failedMsgIds, fmt.Sprintf("error trying to find \"Subscribe\" button element for message ID %s: %s", msgId, err.Error()))
			continue
		}

		if err := elem.Click(); err != nil {
			failedMsgIds = append(failedMsgIds, fmt.Sprintf("error trying to click \"Subscribe\" button element for message ID %s: %s", msgId, err.Error()))
			continue
		}

		time.Sleep(2 * time.Second)
	}

	if len(failedMsgIds) > 0 {
		errorMsgAsStr := "The following errors occurred:"
		for _, failedMsgId := range failedMsgIds {
			errorMsgAsStr += fmt.Sprintf("\n- %s", failedMsgId)
		}
		return fmt.Errorf("%s", errorMsgAsStr)
	}

	return nil
}

func (c *Client) GetOriginalMessageById(msgId string) (*messagePayload, error) {
	url := fmt.Sprintf("%v/messages/%s?format=%s", baseUrl, msgId, "full")
	reqMethod := "GET"

	req, err := http.NewRequest(reqMethod, url, nil)
	if err != nil {
		return nil, &requestCreationError{reqMethod, url, err.Error()}
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, &requestCreationError{reqMethod, url, err.Error()}
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("there was an error reading the response's body: %v", err.Error())
	}

	var randomInterface messagePayload
	err = json.Unmarshal(body, &randomInterface)
	if err != nil {
		return nil, fmt.Errorf("there was an error unmarshalling the data: %v", err.Error())
	}

	return &randomInterface, nil
}

func (c *Client) UnsubscribeByHttpAddress(httpAddress string) (string, error) {
	method := "POST"
	address := httpAddress[1:len(httpAddress) - 1]
	// fmt.Printf("%v\n", address)

	req, err := http.NewRequest(method, address, nil)
	if err != nil {
		return "", &requestCreationError{method, address, err.Error()}
	}

	res, err := c.Do(req)
	if err != nil {
		return "", &requestExecutionError{method, address, err.Error()}
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
    return "", fmt.Errorf("HTTP status %v for url %v", res.Status, address)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %s", err.Error())
	}
	
	var dataInterface interface{}
	err = json.Unmarshal(data, &dataInterface)
	if err != nil {
		return string(data), nil
	}

	return "May need to manually open link", nil
}

func (c *Client) UnsubscribeByMailtoAddress(mailtoAddress string) error {
	address := mailtoAddress[8:len(mailtoAddress) - 1]
	method := "POST"
	url := fmt.Sprintf("%v/messages/send", baseUrl)

	rawEmail := fmt.Sprintf("From: mikitosaarna@gmail.com\r\nTo: %s\r\nSubject: Unsubscribe Request\r\n\r\nPlease unsubscribe me from this mailing list.", address)
	encodedRawEmail := base64.URLEncoding.EncodeToString([]byte(rawEmail))
	body := map[string]string{"raw": encodedRawEmail}
	
	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling body: %s", err.Error())
	}
	reader := bytes.NewReader(marshalledBody)

	req, err := http.NewRequest(method, url , reader)
	if err != nil {
		return &requestCreationError{method, url, err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = c.Do(req)
	if err != nil {
		return &requestExecutionError{method, url, err.Error()}
	}

	fmt.Printf("Successfully sent unsubscribe request email to: %s", address)
	return nil
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

func (c *Client) RetrieveTrashList() (TrashList, error) {
	filters, err := c.RetrieveAllFilters()
	if err != nil {
		return nil, fmt.Errorf("error retrieving list of filters: %v", err.Error())
	}

	trashList := make(TrashList)

	for _, filter := range filters {
		for _, labelId := range filter.Action.AddLabelIds {
			if labelId == "TRASH" {
				trashList[filter.Criteria.From] = struct{}{}
				continue
			}
		}
	}

	return trashList, nil

}

func (c *Client) RetrieveAllFilters() ([]Filter, error) {
	url := fmt.Sprintf("%v/settings/filters", baseUrl)
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
		return nil, fmt.Errorf("HTTP status %v for url %v", res.Status, url)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body for url %v: %v", url, err.Error())
	}

	var filterRes FiltersList
	err = json.Unmarshal(body, &filterRes)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall response.body: %v", err.Error())
	}

	return filterRes.Filter, nil
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
