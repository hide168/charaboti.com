package main

import (
	"net/http"
	"sync"

	"github.com/hide168/charaboti.com/data"
)

func signup(writer http.ResponseWriter, request *http.Request) {
	generateHTML(writer, nil, "layout", "public.navbar", "signup.default")
}

func signupAccount(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		danger(err, "フォームのパースに失敗しました")
		http.Redirect(writer, request, "/err", 302)
		return
	}
	user := data.User{
		Name:            request.PostFormValue("name"),
		Email:           request.PostFormValue("email"),
		Password:        request.PostFormValue("password"),
		ConfirmPassword: request.PostFormValue("confirm-password"),
	}
	var wg sync.WaitGroup
	var templates []string
	ch := make(chan string, 4)
	check := make(chan bool, 3)
	for i := 0; i < 3; i++ {
		wg.Add(1)
	}
	go user.CheckName(&wg, ch, check)
	go user.CheckEmail(&wg, ch, check)
	go user.CheckPassword(&wg, ch, check)
	wg.Wait()
	close(ch)
	close(check)
	for i := range ch {
		templates = append(templates, i)
	}
	templates = append(templates, "signup.layout")
	for i := range check {
		if i == false {
			generateHTML(writer, user, templates...)
			return
		}
	}
	if err := user.Create(); err != nil {
		danger(err, "ユーザーの作成に失敗しました")
		generateHTML(writer, nil, "layout", "public.navbar", "signup.error")
		return
	}
	generateHTML(writer, nil, "layout", "public.navbar", "signup.complete")
}

func login(writer http.ResponseWriter, request *http.Request) {
	generateHTML(writer, nil, "layout", "public.navbar", "login.default")
}

func authenticate(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	user, err := data.UserByEmail(request.PostFormValue("email"))
	if err != nil {
		danger(err, "ユーザーが見つかりません")
	}
	if user.Password == data.Encrypt(request.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			danger(err, "セッションの生成に失敗しました")
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(writer, &cookie)
		http.Redirect(writer, request, "/", 302)
	} else {
		generateHTML(writer, nil, "layout", "public.navbar", "login.error")
	}
}

func logout(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("_cookie")
	if err != http.ErrNoCookie {
		warning(err, "Cookieの取得に失敗しました")
		session := data.Session{Uuid: cookie.Value}
		session.DeleteByUUID()
	}
	http.Redirect(writer, request, "/", 302)
}
