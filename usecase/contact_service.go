package usecase

import (
	"plum/domain"
)

type ContactService struct {
	ticketService  HubspotService
	slackService   SlackService
	chatGptService ChatGPTService
	gmailService   GmailService
}

type SlackService interface {
	SendMessage(webhookUrl, msg string) error
}

type ChatGPTService interface {
	Create(string) (string, error)
}

type HubspotService interface {
	CreateTicket(ticket domain.Ticket) error
}

type GmailService interface {
	CreateDraft(string) error
	Crawling() error
}

func NewContactService(
	hubspotService HubspotService,
	slackService SlackService,
	chatGptService ChatGPTService,
	gmailService GmailService,
) *ContactService {
	return &ContactService{
		ticketService:  hubspotService,
		slackService:   slackService,
		chatGptService: chatGptService,
		gmailService:   gmailService,
	}
}

func (c *ContactService) RespondContact(contact domain.Contact) error {
	panic("implement me!!")
}

func (c *ContactService) GmailToHubspot(mail domain.Mail) error {
	panic("implement me!!")
}
