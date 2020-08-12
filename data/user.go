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

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

func (user *User) CheckName(wg *sync.WaitGroup, ch chan string, check chan bool) {
	var count int
	err := Db.QueryRow("select count(*) from users where name = ?", user.Name).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
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
	var count int
	err := Db.QueryRow("select count(*) from users where email = ?", user.Email).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	match, _ := regexp.MatchString("^[0-9a-z_./?-]+@([0-9a-z-]+.)+[0-9a-z-]+$", user.Email)
	if count > 0 {
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
	statement := "insert into users (uuid, name, email, password, icon, created_at) values (?, ?, ?, ?, ?, ?)"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(CreateUUID(), user.Name, user.Email, Encrypt(user.Password), "default.jpg", time.Now())
	return
}

func UserByEmail(email string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE email = ?", email).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

func (user *User) CreateSession() (session Session, err error) {
	statement := "insert into sessions (uuid, email, user_id, created_at) values (?, ?, ?, ?)"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	uuid := CreateUUID()
	_, err = stmt.Exec(uuid, user.Email, user.Id, time.Now())
	if err != nil {
		return
	}
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = ?", uuid).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

func (session *Session) Check() (valid bool, err error) {
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = ?", session.Uuid).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	if err != nil {
		valid = false
		return
	}
	if session.Id != 0 {
		valid = true
	}
	return
}

func (session *Session) DeleteByUUID() (err error) {
	statement := "delete from sessions where uuid = ?"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(session.Uuid)
	return
}

func (session *Session) User() (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = ?", session.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

func (user *User) ChangeName() (err error) {
	statement := "update users set name = ? where uuid = ?"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Name, user.Uuid)
	return
}

func (user *User) ChangeIcon(icon string) (err error) {
	statement := "update users set icon = ? where uuid = ?"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(icon, user.Uuid)
	return
}
