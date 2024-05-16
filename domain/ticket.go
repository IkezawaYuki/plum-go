package domain

type Ticket struct {
	Subject string
	Content string
	OwnerId int
}

var users = map[string]int{
	"yuki_ikezawa":  415434072,
	"yui_sakaguchi": 446121051,
}

func ContactToTicket(contact Contact) Ticket {
	return Ticket{
		Subject: "[PLUM] お問い合わせがありました",
		Content: contact.GetContents(),
		OwnerId: users["yuki_ikezawa"],
	}
}
