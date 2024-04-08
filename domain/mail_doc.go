package domain

type MailDoc struct {
	Values []Value `json:"value"`
}

type Value struct {
	SearchAction string   `json:"@search.action"`
	ID           string   `json:"id"`
	From         string   `json:"from"`
	To           []string `json:"to"`
	Subject      string   `json:"subject"`
	Body         string   `json:"body"`
}

func ConvertToMailDocs(list GmailList) MailDoc {
	values := make([]Value, 0, len(list.GmailList))
	for _, mail := range list.GmailList {
		values = append(values, Value{
			SearchAction: "upload",
			ID:           mail.ID,
			From:         mail.FromAddress,
			To:           mail.ToAddress,
			Subject:      mail.Subject,
			Body:         mail.Body,
		})
	}
	return MailDoc{Values: values}
}

func ConvertToMailDoc(mail Gmail) MailDoc {
	return MailDoc{Values: []Value{
		{
			SearchAction: "upload",
			ID:           mail.ID,
			From:         mail.FromAddress,
			To:           mail.ToAddress,
			Subject:      mail.Subject,
			Body:         mail.Body,
		},
	}}
}
