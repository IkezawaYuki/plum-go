package infrastructure

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
	"mime"
	"net/http"
	"os"
)

type GmailService struct {
	service *gmail.Service
}

func NewGmailService(credential []byte, tokenFilePath string) *GmailService {
	config, err := google.ConfigFromJSON(credential, gmail.GmailComposeScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config, tokenFilePath)
	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	return &GmailService{service: srv}
}

func getClient(config *oauth2.Config, tokenFilePath string) *http.Client {
	tok, err := tokenFromFile(tokenFilePath)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFilePath, tok)
	}
	return config.Client(context.Background(), tok)
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer func() {
		_ = f.Close()
	}()
	_ = json.NewEncoder(f).Encode(token)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

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

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

var content = `===================================

本メールは自動返信メールとなっております

===================================



この度はお問い合わせありがとうございます。


お問い合わせ内容を確認してご連絡を差し上げます。


引き続き宜しくお願い致します。`

func (g *GmailService) FollowUpMail(toAddress string) error {
	var message gmail.Message
	subject := "お問い合わせありがとうございます"
	encodedSubject := mime.QEncoding.Encode("utf-8", subject)
	messageStr := []byte(
		"From: 'me'\r\n" +
			"To: " + toAddress + " \r\n" +
			"Subject: " + encodedSubject + " \r\n\r\n" +
			content)
	message.Raw = base64.URLEncoding.EncodeToString(messageStr)

	sendCall, err := g.service.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return fmt.Errorf("send mail is failed: %w", err)
	}
	fmt.Println(sendCall.Id)
	return nil
}

func (g *GmailService) CreateDraft(contents string, toAddress string) error {
	var message gmail.Message
	subject := "お問い合わせありがとうございます"
	encodedSubject := mime.QEncoding.Encode("utf-8", subject)
	messageStr := []byte(
		"From: 'me'\r\n" +
			"To: " + toAddress + " \r\n" +
			"Subject: " + encodedSubject + " \r\n\r\n" +
			contents)
	message.Raw = base64.URLEncoding.EncodeToString(messageStr)
	draft := &gmail.Draft{
		Message: &message,
	}
	draft, err := g.service.Users.Drafts.Create("me", draft).Do()
	if err != nil {
		return fmt.Errorf("create draft mail is failed: %w", err)
	}
	return nil
}
