package domain

import "regexp"

type Form struct {
	Company   string `json:"company"`
	Phone     string `json:"phone"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Content   string `json:"contents"`
}

func NewForm(company, phone, lastName, firstName, email, content string) Form {
	return Form{
		Company:   company,
		Phone:     phone,
		LastName:  lastName,
		FirstName: firstName,
		Email:     email,
		Content:   content,
	}
}

func (c *Form) Validation() error {
	if c.Content == "" {
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
	return c.Content
}
