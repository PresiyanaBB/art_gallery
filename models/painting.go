package model

import (
	"time"
)

type Painting struct {
	ID                string
	Title             string
	Description       string
	MIMEType          string
	Data              []byte
	Author            User
	DateOfPublication time.Time
	Width             int
	Height            int
	Genre             Genre
	Price             float64
}
