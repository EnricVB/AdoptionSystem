package models

import (
	"time"
)

type User struct {
	id            uint
	name          string
	surname       string
	email         string
	address       string
	password      string
	twoFactorAuth uint
	failedLogins  uint
	isBlocked     bool
	crtDate       time.Time
	uptDate       time.Time
}

type NonValidatedUser struct {
	id      uint
	name    string
	surname string
	email   string
	address string
	crtDate time.Time
	uptDate time.Time
}

type SimplifiedUser struct {
	id      uint
	name    string
	surname string
	email   string
	address string
}
