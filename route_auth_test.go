package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_Login(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", login)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/login", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
	body := writer.Body.String()
	if strings.Contains(body, "ログイン") == false {
		t.Errorf("Bodyは「ログイン」を含んでいません")
	}
}
