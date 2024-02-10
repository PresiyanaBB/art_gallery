package model

import "time"

type User struct {
	ID                 string
	FirstName          string
	LastName           string
	DateOfRegistration time.Time
	Email              string
	Password           string
}
