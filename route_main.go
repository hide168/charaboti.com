package main

import (
	"net/http"

	"github.com/hide168/charaboti.com/data"
)

func index(writer http.ResponseWriter, request *http.Request) {
	characters, err := data.NewCharacters()
	if err != nil {
		danger(err, "キャラクターの取得に失敗しました")
		http.Redirect(writer, request, "/err", 302)
		return
	}
	_, err = session(writer, request)
	if err != nil {
		generateHTML(writer, &characters, "layout", "public.navbar", "index")
	} else {
		generateHTML(writer, &characters, "layout", "private.navbar", "index")
	}
}

func terms(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, nil, "layout", "public.navbar", "terms")
	} else {
		generateHTML(writer, nil, "layout", "private.navbar", "terms")
	}
}

func privacy(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, nil, "layout", "public.navbar", "privacy")
	} else {
		generateHTML(writer, nil, "layout", "private.navbar", "privacy")
	}
}

func err(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, nil, "layout", "public.navbar", "error")
	} else {
		generateHTML(writer, nil, "layout", "private.navbar", "error")
	}
}
