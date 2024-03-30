package usecase

import (
	"plum/domain"
)

type ContactService struct {
	ticketService  HubspotService
	slackService   SlackService
	chatGptService ChatGPTService
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

func NewContactService(
	hubspotService HubspotService,
	slackService SlackService,
	chatGptService ChatGPTService,
) *ContactService {
	return &ContactService{
		ticketService:  hubspotService,
		slackService:   slackService,
		chatGptService: chatGptService,
	}
}

func (c *ContactService) RespondContact(contact domain.Contact) error {
	panic("implement me!!")
}
