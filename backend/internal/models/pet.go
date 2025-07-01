package models

import "time"

type Pet struct {
	id          uint
	name        string
	species     string
	breed       string
	isAdopted   bool
	birthDate   time.Time
	adoptDate   time.Time
	description string
	adoptUser   User
	crtDate     time.Time
	uptDate     time.Time
}

type SimplifiedPet struct {
	id        uint
	name      string
	species   string
	breed     string
	isAdopted bool
	adoptUser User
}
