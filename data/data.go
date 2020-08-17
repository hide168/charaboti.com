package data

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hide168/charaboti.com/data"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("mysql", "root:golang@tcp(mysql-container:3306)/mysql?parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		log.Fatal(err)
	}
	user := data.User{
		Name:     "テストユーザー",
		Email:    "test@mail.com",
		Password: "testuser",
	}
	var count int
	err := Db.QueryRow("select count(*) from users where name = ?", user.Name).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		if err := user.Create(); err != nil {
			log.Fatal(err)
		}
	}
	return
}

func CreateUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("UUIDの生成に失敗しました", err)
	}
	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}
