package main

import "net/http"

func index(writer http.ResponseWriter, request *http.Request) {
	_, err := session(writer, request)
	if err != nil {
		generateHTML(writer, nil, "layout", "public.navbar", "index")
	} else {
		generateHTML(writer, nil, "layout", "private.navbar", "index")
	}
}

func err(writer http.ResponseWriter, request *http.Request) {
	generateHTML(writer, nil, "layout", "public.navbar", "error")
}
