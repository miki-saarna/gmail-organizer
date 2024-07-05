package gmail

import (
	"context"
	"encoding/json"
	"fmt"
	"gmail-organizer/utils"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Client struct {
	*http.Client
}

type UnsubscribeMessage struct {
	sender string
	id string
}

type UnsubscribeErrorList struct {
	Address string
	MessageId string
	MailtoAddress string
	HttpAddress string
	ErrorMessage string
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
					tok = getTokenFromWeb(config)
					saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
					"authorization code: \n%v\n", authURL)

	err := utils.OpenURL(authURL)
	if err != nil {
		fmt.Printf("Failed to automatically open URL %s. Manually copy the link above.", authURL)
	}

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
					log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
					log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
					return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
					log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func InitMessageRemoval(senderAddresses []string) {
	client, _ := main()

	messages, err := client.ListMessagesFromSender(senderAddresses, 500)
	if err != nil {
		log.Fatalf("Could not successfully retrieve emails: %v", err.Error())
	}

	if len(messages) > 0 {
		// apiClient.RemoveMessages(messages)
		
		err := client.BatchPermanentlyDeleteMessages(messages)
		if err != nil {
			log.Fatalf("Could not successfully delete messages: %v", err.Error())
		}
	}
}

func InitTrashListUpdate(senderAddresses []string) {
	client, _ := main()

	trashList, err := client.RetrieveTrashList()
	if err != nil {
		fmt.Printf("error occurred: %v", err.Error())
	}

	for i := 0; i < len(senderAddresses); i++ {
		if _, found := trashList[senderAddresses[i]]; found {
			continue
		}

		filter, err := client.AssignSenderToTrashList(senderAddresses[i])
		if err != nil {
			log.Fatalf("Could not successfully assign senders to trash list: %v\n", err.Error())
		}
		fmt.Printf("Filtered successfully applied: %v\n", *filter)
	}
}

func InitUnsubscribeWithWebDriver(senderAddresses []string) {
	client, _ := main()

	var messages []string

	for _, senderAddress := range senderAddresses {
		retrievedMessages, err := client.ListMessagesFromSender([]string{senderAddress}, 1)
		if err != nil {
			fmt.Printf("Could not successfully retrieve email from %s: %v", senderAddress, err.Error())
		}
		messages = append(messages, retrievedMessages...)
	}

	if len(messages) > 0 {
		err := client.UnsubscribeWithWebDriver(messages)
		if err != nil {
			log.Fatalf("Could not unsubscribe: \n%s", err)
		}
	}
}

func InitUnsubscribe(senderAddresses []string) {
	client, _ := main()

	var messages []UnsubscribeMessage

	for _, senderAddress := range senderAddresses {
		retrievedMessages, err := client.ListMessagesFromSender([]string{senderAddress}, 1)
		if err != nil {
			fmt.Printf("Could not successfully retrieve email from %s: %v", senderAddress, err.Error())
		}

		if len(retrievedMessages) == 1 {
			messages = append(messages, UnsubscribeMessage{senderAddress, retrievedMessages[0]})
		}
	}

	possibleErrList := []UnsubscribeErrorList{}
	errList := []UnsubscribeErrorList{}

	for _, message := range messages {
		data, err := client.GetOriginalMessageById(message.id)
		if err != nil {
			log.Fatalf("error occurred: %v", err.Error())
		}
		headersList := data.Payload.Headers

		var unsubscribeHttpAddress, unsubscribeMailtoAddress string

		for _, header := range headersList {
			if header.Name == "List-Unsubscribe" {
				unsubscribeAddresses := strings.Split(header.Value, ",")
				if len(unsubscribeAddresses) > 1 {
					if unsubscribeAddresses[0][1:8] == "mailto:" {
						unsubscribeMailtoAddress = unsubscribeAddresses[0]
						unsubscribeHttpAddress = unsubscribeAddresses[1]
					} else {
						unsubscribeMailtoAddress = unsubscribeAddresses[1]
						unsubscribeHttpAddress = unsubscribeAddresses[0]
					}
				} else {
					unsubscribeHttpAddress = unsubscribeAddresses[0]
				}
			}
		}

		// fmt.Printf("\nunsubscribeHttpAddress: %s\nunsubscribeMailtoAddress: %s\n", unsubscribeHttpAddress, unsubscribeMailtoAddress)

		if unsubscribeMailtoAddress != "" {
			err := client.UnsubscribeByMailtoAddress(unsubscribeMailtoAddress)
			if err != nil {
				fmt.Printf("error occurred: %s", err.Error())
				e := UnsubscribeErrorList{message.sender, message.id, unsubscribeMailtoAddress, unsubscribeHttpAddress, err.Error()}
				errList = append(errList, e)
			}
		} else if unsubscribeHttpAddress != "" {
			msg, err := client.UnsubscribeByHttpAddress(unsubscribeHttpAddress)
			if err != nil {
				fmt.Printf("error occurred: %s\n", err.Error())
				e := UnsubscribeErrorList{message.sender, message.id, "", unsubscribeHttpAddress, err.Error()}
				errList = append(errList, e)
			} else {
				fmt.Printf("Message body: %v", msg)
				e := UnsubscribeErrorList{message.sender, message.id, "", unsubscribeHttpAddress, fmt.Sprintf("An error may or may not have occurred: %v", msg)}
				possibleErrList = append(possibleErrList, e)
			}
		} else {
			e := UnsubscribeErrorList{message.sender, message.id, "", "", fmt.Errorf("headers of message ID %s does not contain \"List-Unsubscribe\" header", message).Error()}
			errList = append(errList, e)
		}
	}
	
	if len(possibleErrList) > 0 {
		prettyErrList, err := utils.PrettyPrint(possibleErrList)
		if err != nil {
			fmt.Printf("Could not implement prettyPrint on possibleErrList: %s", err.Error())
			fmt.Printf("\n\nError list: %v", possibleErrList)
		} else {
			fmt.Printf("\n\nError list: %v", prettyErrList)
		}
	}
	
	if len(errList) > 0 {
		prettyErrList, err := utils.PrettyPrint(errList)
		if err != nil {
			fmt.Printf("Could not implement prettyPrint on errList: %s", err.Error())
			fmt.Printf("\n\nError list: %v", errList)
		} else {
			fmt.Printf("\n\nError list: %v", prettyErrList)
		}
	}
}

func main() (*Client, *gmail.Service) {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
					log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope, gmail.GmailSettingsBasicScope)
	if err != nil {
					log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
					log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
					log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(r.Labels) == 0 {
					fmt.Println("No labels found.")
					return nil, nil
	}
	fmt.Println("Labels:")
	for _, l := range r.Labels {
					fmt.Printf("- %s\n", l.Name) // fmt.Printf("- %s: %s\n", l.Name, l.Id)
	}

	return &Client{client}, srv
}