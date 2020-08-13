package main

import (
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/hide168/charaboti.com/data"
)

func newCharacter(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, nil, "layout", "public.navbar", "login.character")
	} else {
		generateHTML(writer, nil, "layout", "private.navbar", "character.new")
	}
}

func postCharacter(writer http.ResponseWriter, request *http.Request) {
	sess, err := session(writer, request)
	if err != nil {
		danger(err, "セッションの確認に失敗しました")
		http.Redirect(writer, request, "/err", 302)
		return
	}
	user, err := sess.User()
	if err != nil {
		danger(err, "セッションからユーザーを取得出来ませんでした")
		http.Redirect(writer, request, "/err", 302)
		return
	}
	file, header, err := request.FormFile("image")
	if err != nil {
		generateHTML(writer, nil, "layout", "private.navbar", "character.new.error")
		return
	}
	defer file.Close()
	dt, err := ioutil.ReadAll(file)
	if err != nil {
		danger(err, "ファイルの読み込みに失敗しました")
		http.Redirect(writer, request, "/err", 302)
		return
	}
	imageUuid := data.CreateUUID()
	filename := filepath.Join("characters", imageUuid+filepath.Ext(header.Filename))
	err = ioutil.WriteFile(filename, dt, 0777)
	if err != nil {
		danger(err, "ファイルの書き込みに失敗しました")
		http.Redirect(writer, request, "/err", 302)
		return
	}
	filename = "/" + filename
	character := data.Character{
		Text:  request.FormValue("text"),
		Image: filename,
	}
	err = character.Create(user.Id)
	if err != nil {
		danger(err, "キャラクターの作成に失敗しました")
		http.Redirect(writer, request, "/err", 302)
		return
	}
	generateHTML(writer, nil, "layout", "private.navbar", "character.new.complete")
}