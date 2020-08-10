package main

import (
	"net/http"
	"time"
)

func main() {
	p("charaboti.com", version(), "started at", config.Address)

	// 静的ファイルの処理
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	//
	// 以下に全てのルートパターンを記述しています
	// ルートハンドラー関数は他のファイルに定義しています

	// index
	mux.HandleFunc("/", index)
	// error
	// mux.HandleFunc("/err", err)

	// route_auth.goで定義されています
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/signup_account", signupAccount)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/authenticate", authenticate)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/mypage", mypage)

	// サーバーの起動処理
	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
