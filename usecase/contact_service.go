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
}

type AISearch interface {
}

type SlackService interface {
	Escalation(contact domain.Contact) error
	Report(contact domain.Contact) error
}

type ChatGPTService interface {
	Create(string) (*domain.Generated, error)
}

type HubspotService interface {
	CreateTicket(ticket domain.Ticket) error
}

type GmailService interface {
	CreateDraft(string, string) error
	FollowUpMail(string) error
}

func NewContactService(
	hubspotService HubspotService,
	slackService SlackService,
	chatGptService ChatGPTService,
	gmailService GmailService,
	aiSearch AISearch,
) *ContactService {
	return &ContactService{
		hubspotService: hubspotService,
		slackService:   slackService,
		chatGptService: chatGptService,
		gmailService:   gmailService,
		aiSearch:       aiSearch,
	}
}

func (c *ContactService) RespondContact(contact domain.Contact) error {
	generated, err := c.chatGptService.Create(contact.GetContents())
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
		if err := c.slackService.Report(contact); err != nil {
			return fmt.Errorf("c.slackService.Report: %w", err)
		}
	} else {
		if err := c.slackService.Escalation(contact); err != nil {
			return fmt.Errorf("c.slackService.Escalation: %w", err)
		}
	}

	ticket := domain.ContactToTicket(contact)
	if err := c.hubspotService.CreateTicket(ticket); err != nil {
		return fmt.Errorf("c.ticketService.CreateTicket: %w", err)
	}

	return nil
}
