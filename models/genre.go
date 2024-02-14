package model

type Genre struct {
	ID   string
	Name string
}

type GenreTypes int

// Enumeration of genres
const (
	Undefined GenreTypes = iota
	Portrait
	Landscape
	StillLife
	History
	Abstract
	Surrealism
	Impressionism
	Expressionism
	Realism
	Baroque
)

// GenreNames maps genre constants to their string representations
var GenreTypesString = map[GenreTypes]string{
	Undefined:     "",
	Portrait:      "Portrait",
	Landscape:     "Landscape",
	StillLife:     "StillLife",
	History:       "History",
	Abstract:      "Abstract",
	Surrealism:    "Surrealism",
	Impressionism: "Impressionism",
	Expressionism: "Expressionism",
	Realism:       "Realism",
	Baroque:       "Baroque",
}
