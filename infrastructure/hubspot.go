package infrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/belong-inc/go-hubspot"
	"io"
	"net/http"
	"plum/domain"
	"plum/logger"
	"strconv"
	"time"
)

type Hubspot struct {
	token  string
	client *hubspot.Client
}

func NewHubspot(token string) *Hubspot {
	client, err := hubspot.NewClient(hubspot.SetPrivateAppToken(token))
	if err != nil {
		panic(err)
	}
	return &Hubspot{client: client, token: token}
}

/*
CreateTicket
Hubspotにチケットを新たに作成する。
*/
func (h *Hubspot) CreateTicket(ticket domain.Ticket) (int, error) {
	createdId := 0
	logger.Logger.Info("CreateTicket is invoked")
	properties := make(map[string]interface{})
	properties["hs_pipeline"] = "0"
	properties["hs_pipeline_stage"] = "1"
	properties["subject"] = ticket.Subject
	properties["content"] = ticket.Content
	properties["hubspot_owner_id"] = ticket.OwnerId
	var ticketReq = &hubspot.CrmTicketCreateRequest{
		Properties: properties,
	}
	createTicker, err := h.client.CRM.Tickets.Create(ticketReq)
	if err != nil {
		return createdId, fmt.Errorf("create ticket is failed: %w", err)
	}
	fmt.Println(createTicker)
	strId := createTicker.Id
	createdId, err = strconv.Atoi(fmt.Sprintf("%v", strId))
	return createdId, nil
}

type HubspotContact struct {
	Results []struct {
		Archived   bool      `json:"archived"`
		CreatedAt  time.Time `json:"createdAt"`
		ID         string    `json:"id"`
		Properties struct {
			Createdate       time.Time `json:"createdate"`
			Email            string    `json:"email"`
			Firstname        string    `json:"firstname"`
			HsObjectID       string    `json:"hs_object_id"`
			Lastmodifieddate time.Time `json:"lastmodifieddate"`
			Lastname         string    `json:"lastname"`
		} `json:"properties"`
		UpdatedAt time.Time `json:"updatedAt"`
	} `json:"results"`
	Total int `json:"total"`
}

func (h *Hubspot) SearchContact(email string) (int, error) {
	// HubSpot APIキー
	apiKey := h.token
	contactId := 0

	// リクエストヘッダー
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + apiKey,
	}

	// 検索クエリ
	data := map[string]interface{}{
		"filterGroups": []interface{}{
			map[string]interface{}{
				"filters": []interface{}{
					map[string]interface{}{
						"propertyName": "email",
						"operator":     "EQ",
						"value":        email,
					},
				},
			},
		},
		"properties": []string{"email", "firstname", "lastname"},
	}

	// JSONにエンコード
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return contactId, err
	}

	// HubSpotの検索エンドポイント
	url := "https://api.hubapi.com/crm/v3/objects/contacts/search"

	// POSTリクエストを作成
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return contactId, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// HTTPクライアントを使用してリクエストを送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return contactId, err
	}
	defer resp.Body.Close()

	// レスポンスを読み取る
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return contactId, err
	}

	// ステータスコードに応じてレスポンスを表示
	if resp.StatusCode == http.StatusOK {
		var hubspotContact HubspotContact
		if err := json.Unmarshal(body, &hubspotContact); err != nil {
			fmt.Println("Error decoding response JSON:", err)
			return contactId, err
		}
		fmt.Println(hubspotContact)
		contactId, err = strconv.Atoi(hubspotContact.Results[0].ID)
		if err != nil {
			return contactId, err
		}
	} else {
		fmt.Printf("Error: %d\n%s\n", resp.StatusCode, body)
	}
	return contactId, nil
}

func (h *Hubspot) Associate(ticketId, contactId int) error {
	// HubSpot APIキー
	apiKey := h.token

	// リクエストヘッダー
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + apiKey,
	}

	// 関連付けのペイロード
	data := map[string]interface{}{
		"inputs": []map[string]interface{}{
			{
				"from": map[string]interface{}{
					"id": contactId,
				},
				"to": map[string]interface{}{
					"id": ticketId,
				},
				"type": "contact_to_ticket",
			},
		},
	}
	fmt.Println(data)

	// JSONにエンコード
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return err
	}

	// HubSpotの関連付けエンドポイント
	url := "https://api.hubapi.com/crm/v3/associations/contact/ticket/batch/create"

	// POSTリクエストを作成
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	// ヘッダーを設定
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// HTTPクライアントを使用してリクエストを送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	// レスポンスを読み取る
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	// ステータスコードに応じてレスポンスを表示
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Association created successfully.")
		fmt.Printf("%s\n", body)
	} else {
		fmt.Printf("Error: %d\n%s\n", resp.StatusCode, body)
	}

	return nil
}

func (h *Hubspot) GetContact(contractId string) error {
	logger.Logger.Info("GetContact is invoked")
	res, err := h.client.CRM.Contact.Get(contractId, &hubspot.Contact{}, nil)
	if err != nil {
		return err
	}
	contact, ok := res.Properties.(*hubspot.Contact)
	if !ok {
		return errors.New("unable to assert type")
	}

	fmt.Println(contact.FirstName)
	fmt.Println(contact.Message)
	fmt.Println(contact.Website)
	return nil
}

func (h *Hubspot) DispatchEvent() error {
	return nil
}
