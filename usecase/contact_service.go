package usecase

import (
	"fmt"
	"plum/domain"
)

type ContactService struct {
	hubspotService HubspotService
	slackService   SlackService
	chatGptService ChatGPTService
	gmailService   GmailService
	aiSearch       AISearch
	db             Db
}

type AISearch interface {
	SearchDocuments(search string) (string, error)
}

type SlackService interface {
	Escalation(contact domain.Contact) error
	Report(contact domain.Contact) error
}

type ChatGPTService interface {
	Generate(contents string, related string, setting domain.ChatgptSetting) (*domain.Generated, error)
}

type HubspotService interface {
	CreateTicket(ticket domain.Ticket) (int, error)
	SearchCompanyByName(companyName string) (int, error)
	SearchContact(email string) (int, error)
	AssociateContactToTicket(ticketId, contactId int) error
	AssociateCompanyToTicket(ticketId, companyId int) error
	GetContact(contractId string) error
}

type GmailService interface {
	CreateDraft(string, string) error
	FollowUpMail(string) error
}

type Db interface {
	GetChatgptSetting() (domain.ChatgptSetting, error)
	UpdateChatgptSetting(domain.ChatgptSetting) error
}

func NewContactService(
	hubspotService HubspotService,
	slackService SlackService,
	chatGptService ChatGPTService,
	gmailService GmailService,
	aiSearch AISearch,
	db Db,
) *ContactService {
	return &ContactService{
		hubspotService: hubspotService,
		slackService:   slackService,
		chatGptService: chatGptService,
		gmailService:   gmailService,
		aiSearch:       aiSearch,
		db:             db,
	}
}

func (c *ContactService) RespondContact(contact domain.Contact) error {
	related, err := c.aiSearch.SearchDocuments(contact.GetContents())
	if err != nil {
		return fmt.Errorf("c.aiSearch.SearchDocuments: %w", err)
	}
	setting, err := c.db.GetChatgptSetting()
	if err != nil {
		return fmt.Errorf("c.db.GetChatgptSetting: %w", err)
	}
	generated, err := c.chatGptService.Generate(contact.GetContents(), related, setting)
	if err != nil {
		return fmt.Errorf("c.chatGptService.Create: %w", err)
	}
	if err := c.gmailService.FollowUpMail(contact.GetEmailAddress()); err != nil {
		return fmt.Errorf("c.gmailService.FollowUpMail: %w", err)
	}
	if !generated.Escalation {
		fmt.Println(generated.Message)
		if err := c.gmailService.CreateDraft(generated.Message, contact.GetEmailAddress()); err != nil {
			return fmt.Errorf("c.gmailService.CreateDraft: %w", err)
		}
	} else {
		if err := c.slackService.Escalation(contact); err != nil {
			return fmt.Errorf("c.slackService.Escalation: %w", err)
		}
	}
	ticketObj := domain.ContactToTicket(contact)
	ticketId, err := c.hubspotService.CreateTicket(ticketObj)
	if err != nil {
		return fmt.Errorf("c.ticketService.CreateTicket: %w", err)
	}
	contactId, err := c.hubspotService.SearchContact(contact.GetEmailAddress())
	if err != nil {
		return fmt.Errorf("c.hubspotService.SearchContact: %w", err)
	}
	if contactId != 0 {
		if err := c.hubspotService.AssociateContactToTicket(ticketId, contactId); err != nil {
			return fmt.Errorf("c.hubspotService.AssociateContactToTicket: %w", err)
		}
	}
	companyId, err := c.hubspotService.SearchCompanyByName(contact.GetCompany())
	if err != nil {
		return fmt.Errorf("c.hubspotService.SearchCompanyByName: %w", err)
	}
	if companyId != 0 {
		if err := c.hubspotService.AssociateCompanyToTicket(companyId, contactId); err != nil {
			return fmt.Errorf("c.hubspotService.AssociateCompanyToTicket: %w", err)
		}
	}
	return nil
}

func (c *ContactService) GetChatgptSetting() (domain.ChatgptSetting, error) {
	return c.db.GetChatgptSetting()
}

func (c *ContactService) UpdateChatgptSetting(setting domain.ChatgptSetting) error {
	return c.db.UpdateChatgptSetting(setting)
}
