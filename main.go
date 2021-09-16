package main

import (
	"log"
	"net/http"
	"os"

	"hackaichi2021/database"
	api_user "hackaichi2021/user"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/cors"
)

func main() {

	database.GormConnect()
	r := mux.NewRouter()
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
