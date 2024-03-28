package domain

type Ticket struct {
	Properties TicketProperties `json:"properties"`
}

type TicketProperties struct {
	HsPipeline      string `json:"hs_pipeline"`
	HsPipelineStage string `json:"hs_pipeline_stage"`
	Subject         string `json:"subject"`
	Content         string `json:"content"`
}

func NewTicket(subject string, content string) Ticket {
	return Ticket{
		Properties: TicketProperties{
			HsPipeline:      "0",
			HsPipelineStage: "1",
			Subject:         subject,
			Content:         content,
		},
	}
}
