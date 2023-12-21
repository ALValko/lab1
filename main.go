package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
)

var (
	mx sync.RWMutex
	db = make(map[string]string)
)

type IndexHandler struct {
}

func (h IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result := false
	if r.Method == http.MethodPost {
		mx.RLock()
		defer mx.RUnlock()
		password, ok := db[r.FormValue("username")]

		result = ok && password == r.FormValue("password")
	}

	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, result)
}

func init() {
	mx.Lock()
	defer mx.Unlock()
	db["user"] = "pwd"
	db["admin"] = "admin"
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", IndexHandler{})

	log.Printf("Сервер запущен на :8000\n")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
