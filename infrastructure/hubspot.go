package infrastructure

import "plum/domain"

type Hubspot struct {
}

func NewHubspot() *Hubspot {
	return &Hubspot{}
}

/*
CreateTicket
Hubspotにチケットを新たに作成する。
*/
func (h *Hubspot) CreateTicket(ticket domain.Ticket) error {
	return nil
}
