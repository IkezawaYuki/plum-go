package usecase

import (
	"plum/domain"
)

type ContactService struct {
	ticketService  HubspotService
	slackService   SlackService
	chatGptService ChatGPTService
	gmailService   GmailService
	aiSearch       AISearch
}

type AISearch interface {
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
	CreateDraft(string, string) error
}

func NewContactService(
	hubspotService HubspotService,
	slackService SlackService,
	chatGptService ChatGPTService,
	gmailService GmailService,
	aiSearch AISearch,
) *ContactService {
	return &ContactService{
		ticketService:  hubspotService,
		slackService:   slackService,
		chatGptService: chatGptService,
		gmailService:   gmailService,
		aiSearch:       aiSearch,
	}
}

func (c *ContactService) RespondContact(contact domain.Contact) error {
	panic("implement me!!")
}

func (c *ContactService) GmailToHubspot(mail domain.Gmail) error {
	panic("implement me!!")
}

func (c *ContactService) GmailToAiSearch(mailList domain.GmailList) error {
	panic("implement me!!")
}
