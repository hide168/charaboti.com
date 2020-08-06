package main

import (
	"net/http"
	"sync"

	"github.com/hide168/charaboti.com/data"
)

func signup(writer http.ResponseWriter, request *http.Request) {
	generateHTML(writer, nil, "layout", "signup.default")
}

func signupAccount(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		danger(err, "フォームのパースに失敗しました")
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
	generateHTML(writer, nil, "layout", "signup.complete")
}
