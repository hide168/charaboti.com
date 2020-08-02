package main

import "net/http"

func index(writer http.ResponseWriter, request *http.Request) {
	generateHTML(writer, nil, "layout", "index")
}
