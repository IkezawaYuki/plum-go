package domain

type Contact interface {
	GetEmailAddress() string
	GetContents() string
	GetCompany() string
}
