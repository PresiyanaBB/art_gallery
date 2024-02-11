package model

import (
	"html/template"
	"time"
)

type Painting struct {
	ID                string
	Title             string
	Description       string
	Src               template.URL
	Author            User
	DateOfPublication time.Time
	Width             int
	Height            int
	Genre             Genre
	Price             float64
}
