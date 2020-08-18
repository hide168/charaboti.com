package data

import "testing"

func CharacterDeleteAll() (err error) {
	statement := "delete from characters"
	_, err = Db.Exec(statement)
	if err != nil {
		return
	}
	return
}

func Test_CreateCharacter(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "ユーザーの作成に失敗しました")
	}
	characters[0].Create(users[0].Id)
	conv := users[0].Characters()
	if conv[0].UserId != users[0].Id {
		t.Error("ユーザーとキャラクターがリンクされていません")
	}
	if err := users[1].Create(); err != nil {
		t.Error(err, "ユーザーの作成に失敗しました")
	}
	characters[1].Create(users[1].Id)
	conv = users[1].Characters()
	if conv[1].Name != "名無し" {
		t.Error("キャラクター名空白時の「名無し」が機能していません")
	}

}
