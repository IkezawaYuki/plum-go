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

type HubspotCompany struct {
	Total   int `json:"total"`
	Results []struct {
		ID         string `json:"id"`
		Properties struct {
			Createdate         time.Time `json:"createdate"`
			Domain             any       `json:"domain"`
			HsLastmodifieddate time.Time `json:"hs_lastmodifieddate"`
			HsObjectID         string    `json:"hs_object_id"`
			Industry           any       `json:"industry"`
			Name               string    `json:"name"`
		} `json:"properties"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
		Archived  bool      `json:"archived"`
	} `json:"results"`
}

func (h *Hubspot) SearchCompanyByName(companyName string) (int, error) {
	apiKey := h.token
	companyId := 0
	// リクエストヘッダー
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + apiKey,
	}

	// 検索クエリ
	data := map[string]interface{}{
		"filterGroups": []map[string]interface{}{
			{
				"filters": []map[string]interface{}{
					{
						"propertyName": "name",
						"operator":     "EQ",
						"value":        companyName,
					},
				},
			},
		},
		"properties": []string{"name", "domain", "industry"},
	}

	// JSONにエンコード
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return companyId, err
	}

	// HubSpotの会社検索エンドポイント
	url := "https://api.hubapi.com/crm/v3/objects/companies/search"

	// POSTリクエストを作成
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return companyId, err
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
		return companyId, err
	}
	defer resp.Body.Close()

	// レスポンスを読み取る
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return companyId, err
	}

	// ステータスコードに応じてレスポンスを表示
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Search successful.")
		fmt.Printf("%s\n", body)
		var hubspotCompany HubspotCompany
		err := json.Unmarshal(body, &hubspotCompany)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
		}
		if hubspotCompany.Total == 0 {
			return companyId, nil
		}
		result, err := strconv.Atoi(hubspotCompany.Results[0].ID)
		if err != nil {
			fmt.Println("Error converting hubspot ID to int:", err)
		}
		companyId = result
	} else {
		fmt.Printf("Error: %d\n%s\n", resp.StatusCode, body)
	}
	return companyId, nil
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
		if hubspotContact.Total == 0 {
			return contactId, nil
		}
		result, err := strconv.Atoi(hubspotContact.Results[0].ID)
		if err != nil {
			return contactId, err
		}
		contactId = result
	} else {
		fmt.Printf("Error: %d\n%s\n", resp.StatusCode, body)
	}
	return contactId, nil
}

func (h *Hubspot) AssociateContactToTicket(ticketId, contactId int) error {
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

func (h *Hubspot) AssociateCompanyToTicket(ticketId, companyId int) error {
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
					"id": companyId,
				},
				"to": map[string]interface{}{
					"id": ticketId,
				},
				"type": "company_to_ticket",
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
	url := "https://api.hubapi.com/crm/v3/associations/company/ticket/batch/create"

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
