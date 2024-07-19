package models

import "time"

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Employee struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	BirthDate time.Time `json:"birth_date"`
}

type Subscription struct {
	IdFrom int `json:"id_from"`
	IdTo   int `json:"id_to"`
}

type SubscriptionRequest struct {
	IdTo int `json:"id_to"`
}

type SignUpRequest struct {
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	BirthDate time.Time `json:"birth_date"`
}
