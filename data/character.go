package data

import "time"

type Character struct {
	Id        int
	Uuid      string
	Name      string
	Text      string
	UserId    int
	Image     string
	CreatedAt time.Time
}

func (character *Character) Create(userId int) (err error) {
	statement := "insert into characters (uuid, name, text, user_id, image, created_at) values (?, ?, ?, ?, ?, ?)"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	var name string
	if character.Name == "" {
		name := "名無し"
	} else {
		name := character.Name
	}
	_, err = stmt.Exec(CreateUUID(), name, character.Text, userId, character.Image, time.Now())
	return
}

func Characters() (characters []Character, err error) {
	rows, err := Db.Query("SELECT id, uuid, name, text, user_id, image, created_at FROM characters ORDER BY created_at DESC")
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

func (character *Character) User() (user User) {
	user = User{}
	Db.QueryRow("SELECT id, uuid, name, email, password, icon, created_at FROM users WHERE id = ?", character.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.Icon, &user.CreatedAt)
	return
}

func CharacterByUUID(uuid string) (character Character, err error) {
	character = Character{}
	err = Db.QueryRow("SELECT id, uuid, name, text, user_id, image, created_at FROM characters WHERE uuid = ?", uuid).
		Scan(&character.Id, &character.Uuid, &character.Name, &character.Text, &character.UserId, &character.Image, &character.CreatedAt)
	return
}

func Search(word string) (characters []Character, err error) {
	rows, err := Db.Query("SELECT id, uuid, name, text, user_id, image, created_at FROM characters WHERE name LIKE = '%?%'", word)
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
