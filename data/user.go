package data

import (
	"regexp"
	"sync"
	"time"
)

type User struct {
	Id              int
	Uuid            string
	Name            string
	Email           string
	Password        string
	ConfirmPassword string
	CreatedAt       time.Time
}

func (user *User) CheckName(wg *sync.WaitGroup, ch chan string, check chan bool) {
	if user.Name != "" {
		ch <- "signup.valid.name"
		check <- true
	} else {
		ch <- "signup.invalid.name"
		check <- false
	}
	wg.Done()
}

func (user *User) CheckEmail(wg *sync.WaitGroup, ch chan string, check chan bool) {
	match, _ := regexp.MatchString("^[0-9a-z_./?-]+@([0-9a-z-]+.)+[0-9a-z-]+$", user.Email)
	if match == true {
		ch <- "signup.valid.email"
		check <- true
	} else {
		ch <- "signup.invalid.email"
		check <- false
	}
	wg.Done()
}

func (user *User) CheckPassword(wg *sync.WaitGroup, ch chan string, check chan bool) {
	match, _ := regexp.MatchString("[A-Za-z0-9]{8,}", user.Password)
	if match != true {
		user.ConfirmPassword = ""
		ch <- "signup.invalid.password"
		ch <- "signup.none.confirm-password"
		check <- false
	} else if user.Password != user.ConfirmPassword {
		ch <- "signup.valid.password"
		ch <- "signup.invalid.confirm-password"
		check <- false
	} else {
		ch <- "signup.valid.password"
		ch <- "signup.valid.confirm-password"
		check <- true
	}
	wg.Done()
}
