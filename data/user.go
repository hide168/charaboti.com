package data

import (
	"log"
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
	log.Print(user.Name)
	if user.Name != "" {
		log.Print(user.Email)
		ch <- "signup.valid.name"
		check <- true
		log.Print("successA")
	} else {
		ch <- "signup.invalid.name"
		check <- false
		log.Print("successC")
	}
	wg.Done()
	log.Print("successD")
}

func (user *User) CheckEmail(wg *sync.WaitGroup, ch chan string, check chan bool) {
	log.Print("successE")
	match, _ := regexp.MatchString("^[0-9a-z_./?-]+@([0-9a-z-]+.)+[0-9a-z-]+$", user.Email)
	log.Print("successF")
	if match == true {
		ch <- "signup.valid.email"
		check <- true
		log.Print("successG")
	} else {
		ch <- "signup.invalid.email"
		check <- false
		log.Print("successH")
	}
	wg.Done()
	log.Print("successI")
}

func (user *User) CheckPassword(wg *sync.WaitGroup, ch chan string, check chan bool) {
	log.Print("successJ")
	match, _ := regexp.MatchString("[A-Za-z0-9]{8,}", user.Password)
	log.Print("successK")
	if match != true {
		user.ConfirmPassword = ""
		ch <- "signup.invalid.password"
		ch <- "signup.none.confirm-password"
		check <- false
		log.Print("successL")
	} else if user.Password == user.ConfirmPassword {
		ch <- "signup.valid.password"
		ch <- "signup.invalid.confirm-password"
		check <- false
		log.Print("successM")
	} else {
		ch <- "signup.valid.password"
		ch <- "signup.valid.confirm-password"
		check <- true
		log.Print("successN")
	}
	wg.Done()
	log.Print("successO")
}
