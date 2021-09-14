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
	"hackaichi2021/database"
	api_user "hackaichi2021/user"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	JsonParseErr       = "invalid json"
	JsonInvalid        = "Invalid json"
	AuthenticateFailed = "AuthenticateFailed"
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
	err := godotenv.Load()
	if err != nil {
	}

	database.GormConnect()
	os.Setenv("PORT", "8080")

	r := mux.NewRouter()
	r.Handle("/", public)
	r.Handle("/public", public)
	r.Handle("/private", auth.JwtMiddleware.Handler(private))
	r.Handle("/auth", auth.GetTokenHandler)
	r.Handle("/login", login)
	r.Handle("/register", api_user.Register)

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
		fmt.Fprintf(w, JsonInvalid)
		return
	}

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, AuthenticateFailed)
		return
	}

	if user.UserName != u.UserName || user.Password != u.Password {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, AuthenticateFailed)
		return
	}
	token, err := auth.CreateToken(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, AuthenticateFailed)
		return
	}
	json.NewEncoder(w).Encode(token)
})
