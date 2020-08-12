package main

import (
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/hide168/charaboti.com/data"
)

func mypage(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		danger(err, "セッションの確認に失敗しました")
		http.Redirect(writer, request, "/err", 302)
	} else {
		user, err := sess.User()
		if err != nil {
			danger(err, "セッションからユーザーを取得出来ませんでした")
			http.Redirect(writer, request, "/err", 302)
		} else {
			generateHTML(writer, user, "layout", "private.navbar", "mypage")
		}
	}
}

func mypageEdit(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		danger(err, "セッションの確認に失敗しました")
		http.Redirect(writer, request, "/err", 302)
	} else {
		user, err := sess.User()
		if err != nil {
			danger(err, "セッションからユーザーを取得出来ませんでした")
			http.Redirect(writer, request, "/err", 302)
		} else {
			generateHTML(writer, user, "layout", "private.navbar", "mypage.edit")
		}
	}
}

func changeProfile(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		danger(err, "セッションの確認に失敗しました")
		http.Redirect(writer, request, "/err", 302)
		return
	}
	user := data.User{
		Uuid: request.FormValue("uuid"),
		Name: request.FormValue("name"),
	}
	if user.Name == "" {
		generateHTML(writer, user, "layout", "private.navbar", "mypage.edit.error")
		return
	}
	file, header, err := request.FormFile("icon")
	if err != nil {
		err = user.ChangeName()
		if err != nil {
			danger(err, "ユーザー名の変更に失敗しました")
			http.Redirect(writer, request, "/err", 302)
			return
		} else {
			http.Redirect(writer, request, "/mypage", 302)
			return
		}
	}
	defer file.Close()
	err = user.ChangeName()
	if err != nil {
		danger(err, "ユーザー名の変更に失敗しました")
		http.Redirect(writer, request, "/err", 302)
		return
	}
	dt, err := ioutil.ReadAll(file)
	if err != nil {
		danger(err, "ファイルの読み込みに失敗しました")
		http.Redirect(writer, request, "/err", 302)
		return
	}
	iconUuid := data.CreateUUID()
	filename := filepath.Join("icons", iconUuid+filepath.Ext(header.Filename))
	err = ioutil.WriteFile(filename, dt, 0777)
	if err != nil {
		danger(err, "ファイルの書き込みに失敗しました")
		http.Redirect(writer, request, "/err", 302)
		return
	}
	filename = "/" + filename
	err = user.ChangeIcon(filename)
	if err != nil {
		danger(err, "アイコンの変更に失敗しました")
		http.Redirect(writer, request, "/err", 302)
		return
	} else {
		http.Redirect(writer, request, "/mypage", 302)
	}
}
