package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	// r.Handle("/api/user/matching", api_user.Match).Methods("POST")
	r.Handle("/api/user/favorite/get", api_user.FavoriteGet).Methods("POST")
	go monitor()

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

type AIRequest struct {
}

func monitor() {
	// a := api_user.LendResponse{
	// 	UserName:  "test",
	// 	Latitude:  134.31,
	// 	Longitude: 24.13,
	// }
	for {
		time.Sleep(1 * time.Second) // 1秒待つ
		maxValue := 0
		maxIndex := 0
		fmt.Println(maxIndex)
		if len(api_user.MatchingSlice[0]) > 0 {
			for i, _ := range api_user.MatchingSlice[1] {

				apiValue := 100
				if maxValue < apiValue {
					maxValue = apiValue
					maxIndex = i
				}
			}
			// api_user.NotifiesLend[api_user.MatchingSlice[0][0].AccessToken] <- a
			// delete(api_user.NotifiesLend, api_user.MatchingSlice[0][0].AccessToken)
			// api_user.MatchingSlice[0] = api_user.MatchingSlice[0][1:]
		}

		// fmt.Println("slice", api_user.MatchingSlice)
	}
}
