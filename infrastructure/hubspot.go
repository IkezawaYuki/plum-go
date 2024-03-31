package infrastructure

import (
	"errors"
	"fmt"
	"github.com/belong-inc/go-hubspot"
	"plum/domain"
)

type Hubspot struct {
	client *hubspot.Client
}

func NewHubspot(token string) *Hubspot {
	client, err := hubspot.NewClient(hubspot.SetPrivateAppToken(token))
	if err != nil {
		panic(err)
	}
	return &Hubspot{client: client}
}

/*
CreateTicket
Hubspotにチケットを新たに作成する。
*/
func (h *Hubspot) CreateTicket(ticket domain.Ticket) error {
	return nil
}

func (h *Hubspot) GetContact(contractId string) error {
	fmt.Println("GetContact is invoked")
	res, err := h.client.CRM.Contact.Get(contractId, &hubspot.Contact{}, nil)
	if err != nil {
		return err
	}
	contact, ok := res.Properties.(*hubspot.Contact)
	if !ok {
		return errors.New("unable to assert type")
	}

	// Use contact fields.
	fmt.Println(contact.FirstName)
	fmt.Println(contact.Message)
	fmt.Println(contact.Website)
	return nil
}
