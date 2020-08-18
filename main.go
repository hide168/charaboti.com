package main

import (
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
)

func main() {
	p("charaboti.com", version(), "started at", config.Address)

	// 静的ファイルの処理
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(config.Static))
	http.Handle("/static/", http.StripPrefix("/static/", files))
	files = http.FileServer(http.Dir(config.Icons))
	mux.Handle("/icons/", http.StripPrefix("/icons/", files))
	files = http.FileServer(http.Dir(config.Characters))
	mux.Handle("/characters/", http.StripPrefix("/characters/", files))

	//
	// 以下に全てのルートパターンを記述しています
	// ルートハンドラー関数は他のファイルに定義しています

	// route_main.goで定義されています
	http.HandleFunc("/", index)
	mux.HandleFunc("/terms", terms)
	mux.HandleFunc("/privacy", privacy)
	mux.HandleFunc("/err", err)

	// route_auth.goで定義されています
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/signup_account", signupAccount)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/authenticate", authenticate)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/test_login", testLogin)

	// route_mypage.goで定義されています
	mux.HandleFunc("/mypage", mypage)
	mux.HandleFunc("/mypage_edit", mypageEdit)
	mux.HandleFunc("/change_profile", changeProfile)

	// route_character.goで定義されています
	mux.HandleFunc("/character/new", newCharacter)
	mux.HandleFunc("/character/post", postCharacter)
	mux.HandleFunc("/character/list", listCharacter)
	mux.HandleFunc("/character/detail", detailCharacter)
	mux.HandleFunc("/character/delete", deleteCharacter)
	mux.HandleFunc("/character/search", searchCharacter)

	// サーバーの起動処理
	// server := &http.Server{
	// 	Addr:           config.Address,
	// 	Handler:        mux,
	// 	ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
	// 	WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
	// 	MaxHeaderBytes: 1 << 20,
	// }
	// server.ListenAndServe()
	l, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Fatal(err)
	}
	fcgi.Serve(l, nil)
}
