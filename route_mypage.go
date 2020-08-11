package main

import (
	"log"
	"net/http"

	"github.com/hide168/charaboti.com/data"
)

func mypage(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/err", 302)
	} else {
		user, err := sess.User()
		if err != nil {
			danger(err, "セッションからユーザーを取得出来ませんでした。")
			http.Redirect(writer, request, "/err", 302)
		} else {
			generateHTML(writer, user, "layout", "private.navbar", "mypage")
		}
	}
}

func mypageEdit(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		http.Redirect(writer, request, "/err", 302)
	} else {
		user, err := sess.User()
		if err != nil {
			danger(err, "セッションからユーザーを取得出来ませんでした。")
			http.Redirect(writer, request, "/err", 302)
		} else {
			generateHTML(writer, user, "layout", "private.navbar", "mypage.edit")
		}
	}
}

func changeProfile(writer http.ResponseWriter, request *http.Request) {
	// sess, err := session(writer, request)
	// if err != nil {
	// 	http.Redirect(writer, request, "/err", 302)
	// 	return
	// }
	err := request.ParseForm()
	if err != nil {
		danger(err, "フォームのパースに失敗しました")
		http.Redirect(writer, request, "/err", 302)
		return
	}
	user := data.User{
		Uuid: request.PostFormValue("uuid"),
		Name: request.PostFormValue("name"),
	}
	if user.Name == "" {
		log.Print(user.Name)
		log.Print(user.Uuid)
		generateHTML(writer, user, "layout", "private.navbar", "mypage.edit.error")
		return
	}
	file, _, err := request.FormFile("icon")
	if err != nil {
		http.Redirect(writer, request, "/err", 302)
		return
	}
	defer file.Close()
	http.Redirect(writer, request, "/", 302)
	// data, err := ioutil.ReadAll(file)
}
