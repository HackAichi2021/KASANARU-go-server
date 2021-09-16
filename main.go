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
	"github.com/justinas/alice"
	"github.com/rs/cors"
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

	database.GormConnect()
	r := mux.NewRouter()
	r.Handle("/", public)
	r.Handle("/public", public)
	r.Handle("/private", auth.JwtMiddleware.Handler(private))
	r.Handle("/auth", auth.GetTokenHandler)
	r.Handle("/login", login)
	r.Handle("/api/user/register", api_user.Register).Methods("POST")
	r.Handle("/api/user/login", api_user.Login).Methods("POST")
	r.Handle("/api/user/update", api_user.Update).Methods("POST")
	r.Handle("/api/user/matching", api_user.Match).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "application/json"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	chain := alice.New(c.Handler, logHandler).Then(r)

	//サーバー起動
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), chain); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func logHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method: %v; URL: %v; Protocol: %v", r.Method, r.URL, r.Proto)
		h.ServeHTTP(w, r)
	})
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
	token, err := auth.CreateTokenByUserIdWithEmail(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, AuthenticateFailed)
		return
	}
	json.NewEncoder(w).Encode(token)
})
