package domain

import "regexp"

type Form struct {
	CompanyName string `json:"company_name"`
	PhoneNumber string `json:"phone_number"`
	LastName    string `json:"last_name"`
	FirstName   string `json:"first_name"`
	Email       string `json:"email"`
	Contents    string `json:"contents"`
}

func (c *Form) Validation() error {
	if c.Contents == "" {
		return ContentsIsEmptyError
	}
	if c.Email == "" {
		return EmailIsEmptyError
	}
	if !c.isValidEmail() {
		return InvalidEmailError
	}
	return nil
}

func (c *Form) isValidEmail() bool {
	pattern := `^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(c.Email)
}

func (c *Form) GetEmailAddress() string {
	return c.Email
}

func (c *Form) GetContents() string {
	return c.Contents
}