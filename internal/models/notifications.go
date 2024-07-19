package models

type Notification struct {
	Email              string `json:"email"`
	NameBirthdayPerson string `json:"name"`
	Count              int    `json:"count"`
}

type BirthdayPerson struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	EmailAddress string `json:"email_address"`
}
