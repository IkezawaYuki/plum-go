package infrastructure

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"plum/domain"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("TestMain is invoked")
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println(os.Getenv("HUBSPOT_ACCESS_TOKEN"))
	code := m.Run()
	os.Exit(code)
}

func TestHubspot_GetContact(t *testing.T) {
	hubspot := NewHubspot(os.Getenv("HUBSPOT_ACCESS_TOKEN"))
	if err := hubspot.GetContact("6695090231"); err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func TestHubspot_CreateTicket(t *testing.T) {
	hubspot := NewHubspot(os.Getenv("HUBSPOT_ACCESS_TOKEN"))
	ticket := domain.Ticket{
		Subject: "TEST TEST_IKEZAWA",
		Content: "TEST TEST TEST",
		OwnerId: 415434072,
	}
	if _, err := hubspot.CreateTicket(ticket); err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func TestHubspot_SearchContact(t *testing.T) {
	hubspot := NewHubspot(os.Getenv("HUBSPOT_ACCESS_TOKEN"))
	contactId, err := hubspot.SearchContact("yuki.ikezawa@strategy-drive.jp")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	fmt.Println(contactId)

	_, err = hubspot.SearchContact("yddddduki.ikezawa@strategy-drive.jp")
	if errors.Is(err, domain.HubspotNoResultsError) == false {
		t.Fatalf("%s", err.Error())
	} else {
		fmt.Println(err)
	}
}

func TestHubspot_CreateAndAssociate(t *testing.T) {
	hubspot := NewHubspot(os.Getenv("HUBSPOT_ACCESS_TOKEN"))
	ticket := domain.Ticket{
		Subject: "TEST TEST TEST",
		Content: "TEST TEST TEST",
		OwnerId: 415434072,
	}
	ticketId, err := hubspot.CreateTicket(ticket)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	fmt.Println(ticketId)
	contactId, err := hubspot.SearchContact("yuki.ikezawa@strategy-drive.jp")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	fmt.Println(contactId)
	err = hubspot.AssociateContactToTicket(ticketId, contactId)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	companyId, err := hubspot.SearchCompanyByName("鼓動")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	fmt.Println(companyId)
	err = hubspot.AssociateCompanyToTicket(ticketId, companyId)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func TestHubspot_SearchCompanyByName(t *testing.T) {
	hubspot := NewHubspot(os.Getenv("HUBSPOT_ACCESS_TOKEN"))
	_id, err := hubspot.SearchCompanyByName("鼓動")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	fmt.Println(_id)

	_, err = hubspot.SearchCompanyByName("wahaha本舗")
	if errors.Is(err, domain.HubspotNoResultsError) == false {
		t.Fatalf("%s", err.Error())
	} else {
		fmt.Println(err)
	}
}
