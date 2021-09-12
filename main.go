package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"hackaichi2021/auth"

	"github.com/gorilla/mux"
)

var (
	JsonParseErr = "invalid json"
)

type post struct {
	Title string `json:"title"`
	Tag   string `json:"tag"`
	URL   string `json:"url"`
}
type User struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

//A sample use
var user = User{
	ID:       "hogehoge",
	UserName: "username",
	Password: "password",
}

func main() {
	r := mux.NewRouter()
	r.Handle("/", public)
	r.Handle("/public", public)
	r.Handle("/private", auth.JwtMiddleware.Handler(private))
	r.Handle("/auth", auth.GetTokenHandler)
	r.Handle("/login", login)

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

func checkJson(w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return http.StatusBadRequest, errors.New(JsonParseErr)
	}

	_, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		return http.StatusInternalServerError, errors.New(JsonParseErr)
	}

	body, err := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusInternalServerError)
		return http.StatusInternalServerError, errors.New(JsonParseErr)
	}
	return 0, nil
}

var login = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if statusCode, err := checkJson(w, r); err != nil {
		w.WriteHeader(statusCode)
		fmt.Fprintf(w, "Invalid json")
		return
	}

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Please provide valid login details111")
		return
	}

	if user.UserName != u.UserName || user.Password != u.Password {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Please provide valid login details")
		return
	}
	token, err := auth.CreateToken(user.ID, user.UserName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "invalid login parameter")
		return
	}
	json.NewEncoder(w).Encode(token)
})
