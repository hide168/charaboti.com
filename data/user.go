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
	statement := "select count(*) from users where name=?"
	stmt, _ := Db.Prepare(statement)
	defer stmt.Close()
	result, _ := stmt.Exec(user.Name)
	if result > 0 {
		ch <- "signup.duplicate.name"
		check <- false
	} else if user.Name != "" {
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
	statement := "select count(*) from users where email=?"
	stmt, _ := Db.Prepare(statement)
	defer stmt.Close()
	result, _ := stmt.Exec(user.Email)
	if result > 0 {
		ch <- "signup.duplicate.email"
		check <- false
	} else if match == true {
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

func (user *User) Create() (err error) {
	statement := "insert into users (uuid, name, email, password, created_at) values (?, ?, ?, ?, ?)"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(createUUID(), user.Name, user.Email, Encrypt(user.Password), time.Now())
	return
}
