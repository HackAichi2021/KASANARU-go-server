package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"hackaichi2021/auth"

	"github.com/gorilla/mux"
)

type post struct {
	Title string `json:"title"`
	Tag   string `json:"tag"`
	URL   string `json:"url"`
}

func main() {
	r := mux.NewRouter()
	r.Handle("/", public)
	r.Handle("/public", public)
	r.Handle("/private", auth.JwtMiddleware.Handler(private))
	r.Handle("/auth", auth.GetTokenHandler)

	//サーバー起動
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

var public = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	post := &post{
		Title: "test",
		Tag:   "test",
		URL:   "test",
	}
	json.NewEncoder(w).Encode(post)
})

var private = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	post := &post{
		Title: "test",
		Tag:   "test",
		URL:   "test",
	}
	json.NewEncoder(w).Encode(post)
})
