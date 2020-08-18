package data

import (
	"database/sql"
	"testing"
)

func SessionDeleteAll() (err error) {
	statement := "delete from sessions"
	_, err = Db.Exec(statement)
	return
}

func UserDeleteAll() (err error) {
	statement := "delete from users"
	_, err = Db.Exec(statement)
	return
}

func Test_UserCreate(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "ユーザーの作成に失敗しました")
	}
	if users[0].Id == 0 {
		t.Errorf("ユーザーIDが存在しません")
	}
	u, err := UserByEmail(users[0].Email)
	if err != nil {
		t.Error(err, "ユーザーが存在しません")
	}
	if users[0].Email != u.Email {
		t.Errorf("取得したユーザーは作成したユーザーと異なっています")
	}
}

func Test_UserDelete(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "ユーザーの作成に失敗しました")
	}
	if err := users[0].Delete(); err != nil {
		t.Error(err, "ユーザーの削除に失敗しました")
	}
	_, err := UserByEmail(users[0].Email)
	if err != sql.ErrNoRows {
		t.Error(err, "ユーザーが削除されていません")
	}
}

func Test_UserUpdate(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "ユーザーの作成に失敗しました")
	}
	users[0].Name = "Random User"
	if err := users[0].Update(); err != nil {
		t.Error(err, "ユーザーの更新に失敗しました")
	}
	u, err := UserByEmail(users[0].Email)
	if err != nil {
		t.Error(err, "ユーザーの取得に失敗しました")
	}
	if u.Name != "Random User" {
		t.Error(err, "ユーザーが更新されていません")
	}
}
func Test_UserByUUID(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "ユーザーの作成に失敗しました")
	}
	u, err := UserByUUID(users[0].Uuid)
	if err != nil {
		t.Error(err, "ユーザーが作成されていません")
	}
	if users[0].Email != u.Email {
		t.Errorf("取得したユーザーは作成したユーザーと異なっています")
	}
}

func Test_Users(t *testing.T) {
	setup()
	for _, user := range users {
		if err := user.Create(); err != nil {
			t.Error(err, "ユーザーの作成に失敗しました")
		}
	}
	u, err := Users()
	if err != nil {
		t.Error(err, "ユーザーの取得に失敗しました")
	}
	if len(u) != 2 {
		t.Error(err, "ユーザーの取得数が異なっています")
	}
	if u[0].Email != users[0].Email {
		t.Error(u[0], users[0], "取得したユーザーが異なっています")
	}
}

func Test_CreateSession(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "ユーザーの作成に失敗しました")
	}
	session, err := users[0].CreateSession()
	if err != nil {
		t.Error(err, "セッションの作成に失敗しました")
	}
	if session.UserId != users[0].Id {
		t.Error("ユーザーとセッションがリンクしていません")
	}
}

func Test_GetSession(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "ユーザーの作成に失敗しました")
	}
	session, err := users[0].CreateSession()
	if err != nil {
		t.Error(err, "セッションの作成に失敗しました")
	}

	s, err := users[0].Session()
	if err != nil {
		t.Error(err, "セッションの取得に失敗しました")
	}
	if s.Id == 0 {
		t.Error("取得されたセッションが存在しません")
	}
	if s.Id != session.Id {
		t.Error("取得されたセッションと異なっています")
	}
}

func Test_checkValidSession(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "ユーザーの作成に失敗しました")
	}
	session, err := users[0].CreateSession()
	if err != nil {
		t.Error(err, "セッションの作成に失敗しました")
	}

	uuid := session.Uuid

	s := Session{Uuid: uuid}
	valid, err := s.Check()
	if err != nil {
		t.Error(err, "セッションのチェックに失敗しました")
	}
	if valid != true {
		t.Error(err, "セッションが無効です")
	}

}

func Test_checkInvalidSession(t *testing.T) {
	setup()
	s := Session{Uuid: "123"}
	valid, err := s.Check()
	if err == nil {
		t.Error(err, "セッションは無効ですが検証されています")
	}
	if valid == true {
		t.Error(err, "誤ってセッションが有効になっています")
	}

}

func Test_DeleteSession(t *testing.T) {
	setup()
	if err := users[0].Create(); err != nil {
		t.Error(err, "ユーザーの作成に失敗しました")
	}
	session, err := users[0].CreateSession()
	if err != nil {
		t.Error(err, "セッションの作成に失敗しました")
	}

	err = session.DeleteByUUID()
	if err != nil {
		t.Error(err, "セッションの削除に失敗しました")
	}
	s := Session{Uuid: session.Uuid}
	valid, err := s.Check()
	if err == nil {
		t.Error(err, "削除したセッションが有効になっています")
	}
	if valid == true {
		t.Error(err, "セッションが削除されていません")
	}
}
