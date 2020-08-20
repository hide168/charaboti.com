package main

import (
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
)

func main() {
	p("charaboti.com", version(), "started at tcp 127.0.0.1:9000")

	// 静的ファイルの処理
	// mux := http.NewServeMux()
	files := http.FileServer(http.Dir(config.Static))
	http.Handle("/static/", http.StripPrefix("/static/", files))
	files = http.FileServer(http.Dir(config.Icons))
	http.Handle("/icons/", http.StripPrefix("/icons/", files))
	files = http.FileServer(http.Dir(config.Characters))
	http.Handle("/characters/", http.StripPrefix("/characters/", files))

	//
	// 以下に全てのルートパターンを記述しています
	// ルートハンドラー関数は他のファイルに定義しています

	// route_main.goで定義されています
	http.HandleFunc("/", index)
	http.HandleFunc("/terms", terms)
	http.HandleFunc("/privacy", privacy)
	http.HandleFunc("/err", err)

	// route_auth.goで定義されています
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/signup_account", signupAccount)
	http.HandleFunc("/login", login)
	http.HandleFunc("/authenticate", authenticate)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/test_login", testLogin)

	// route_mypage.goで定義されています
	http.HandleFunc("/mypage", mypage)
	http.HandleFunc("/mypage_edit", mypageEdit)
	http.HandleFunc("/change_profile", changeProfile)

	// route_character.goで定義されています
	http.HandleFunc("/character/new", newCharacter)
	http.HandleFunc("/character/post", postCharacter)
	http.HandleFunc("/character/list", listCharacter)
	http.HandleFunc("/character/detail", detailCharacter)
	http.HandleFunc("/character/delete", deleteCharacter)
	http.HandleFunc("/character/search", searchCharacter)

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
