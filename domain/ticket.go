package domain

type Ticket struct {
	Subject string
	Content string
}

func ContactToTicket(contact Contact) Ticket {
	return Ticket{
		Subject: "お問い合わせがありました",
		Content: contact.GetContents(),
	}
}
