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
	Icon            string
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

	_, err = stmt.Exec(CreateUUID(), user.Name, user.Email, Encrypt(user.Password), "/icons/default.jpg", time.Now())
	return
}

func (user *User) Delete() (err error) {
	statement := "delete from users where id = ?"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)
	return
}

func (user *User) Update() (err error) {
	statement := "update users set name = ?, email = ? where id = ?"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Name, user.Email, user.Id)
	return
}

func UserByEmail(email string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE email = ?", email).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

func UserByUUID(uuid string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = ?", uuid).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

func Users() (users []User, err error) {
	rows, err := Db.Query("SELECT id, uuid, name, email, password, created_at FROM users")
	if err != nil {
		return
	}
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt); err != nil {
			return
		}
		users = append(users, user)
	}
	rows.Close()
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

func (user *User) Session() (session Session, err error) {
	session = Session{}
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE user_id = ?", user.Id).
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
	err = Db.QueryRow("SELECT id, uuid, name, email, icon, created_at FROM users WHERE id = ?", session.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Icon, &user.CreatedAt)
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

func (user *User) Characters() (characters []Character) {
	rows, err := Db.Query("SELECT id, uuid, name, text, user_id, image, created_at FROM characters WHERE user_id = ? ORDER BY created_at DESC", user.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		conv := Character{}
		if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Name, &conv.Text, &conv.UserId, &conv.Image, &conv.CreatedAt); err != nil {
			return
		}
		characters = append(characters, conv)
	}
	rows.Close()
	return
}
