package data

import "time"

type Character struct {
	Id        int
	Uuid      string
	Text      string
	UserId    int
	Image     string
	CreatedAt time.Time
}

func (character *Character) Create(userId int) (err error) {
	statement := "insert into characters (uuid, text, user_id, image, created_at) values (?, ?, ?, ?, ?)"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(CreateUUID(), character.Text, userId, character.Image, time.Now())
	return
}
