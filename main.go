package main

import (
	"bytes"
	renameJson "encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"time"

	"hackaichi2021/database"
	"hackaichi2021/user"
	api_user "hackaichi2021/user"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/cors"
)

const (
	EQUATORIAL_RADIUS = 6378137.0            // 赤道半径 GRS80
	POLAR_RADIUS      = 6356752.314          // 極半径 GRS80
	ECCENTRICITY      = 0.081819191042815790 // 第一離心率 GRS80
)

type Coodinate struct {
	Latitude  float64
	Longitude float64
}

func main() {

	api_user.MatchingGlobal.NotifiesLend = map[string](chan user.Matching){}
	database.GormConnect()
	r := mux.NewRouter()
	r.Handle("/api/user/register", api_user.Register).Methods("POST")
	r.Handle("/api/user/login", api_user.Login).Methods("POST")
	r.Handle("/api/user/update", api_user.Update).Methods("POST")
	r.Handle("/api/user/matching", api_user.Match).Methods("POST")
	r.Handle("/api/user/favorite/get", api_user.FavoriteGet).Methods("POST")
	r.Handle("/curl", a).Methods("GET")
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

var a = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	out, _ := exec.Command("curl", "google.com").Output()
	fmt.Println(string(out))
})

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
		// maxValue := 0
		// maxIndex := 0
		// fmt.Println(api_user.MatchingGlobal, "test")
		if len(api_user.MatchingGlobal.MatchingSlice[0]) > 0 {
			fmt.Println("ok")
			api_user.MatchingGlobal.Mux.Lock()
			for _, v := range api_user.MatchingGlobal.MatchingSlice[1] {
				fmt.Println("inside")
				// tmp := []api_user.AIFavoriteForm{
				// 	{
				// 		Age1:      api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Age,
				// 		Sex1:      api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Sex,
				// 		Game1:     api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Game,
				// 		Sport1:    api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Sport,
				// 		Book1:     api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Book,
				// 		Travel1:   api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Travel,
				// 		Internet1: api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Internet,
				// 		Anime1:    api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Anime,
				// 		Movie1:    api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Movie,
				// 		Music1:    api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Music,
				// 		Gourmet1:  api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Gourmet,
				// 		Mucle1:    api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Muscle,
				// 		Camp1:     api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Camp,
				// 		Tv1:       api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Tv,
				// 		Cook1:     api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Cook,
				// 		Age2:      v.Favorite.Age,
				// 		Sex2:      v.Favorite.Sex,
				// 		Game2:     v.Favorite.Game,
				// 		Sport2:    v.Favorite.Sport,
				// 		Book2:     v.Favorite.Book,
				// 		Travel2:   v.Favorite.Travel,
				// 		Internet2: v.Favorite.Internet,
				// 		Anime2:    v.Favorite.Anime,
				// 		Movie2:    v.Favorite.Movie,
				// 		Music2:    v.Favorite.Music,
				// 		Gourmet2:  v.Favorite.Gourmet,
				// 		Mucle2:    v.Favorite.Muscle,
				// 		Camp2:     v.Favorite.Camp,
				// 		Tv2:       v.Favorite.Tv,
				// 		Cook2:     v.Favorite.Cook,
				// 	},
				// }
				// item := api_user.AIDataForm{
				// 	Data: tmp,
				// }
				// _, err := HttpPost(os.Getenv("URL"), os.Getenv("AUTHENTICATION"), item)
				// if err != nil {
				// fmt.Println("err", err)
				// }
				// apiValue := 100
				// if maxValue < apiValue {
				// maxValue = apiValue
				// maxIndex = i
				api_user.MatchingGlobal.NotifiesLend[api_user.MatchingGlobal.MatchingSlice[0][0].Info.AccessToken] <- v
				delete(api_user.MatchingGlobal.NotifiesLend, api_user.MatchingGlobal.MatchingSlice[0][0].Info.AccessToken)

				api_user.MatchingGlobal.NotifiesLend[v.Info.AccessToken] <- api_user.MatchingGlobal.MatchingSlice[0][0]
				delete(api_user.MatchingGlobal.NotifiesLend, v.Info.AccessToken)

				api_user.MatchingGlobal.MatchingSlice[0] = unset(api_user.MatchingGlobal.MatchingSlice[0], 0)
				api_user.MatchingGlobal.MatchingSlice[1] = unset(api_user.MatchingGlobal.MatchingSlice[1], 0)
				break
			}
			api_user.MatchingGlobal.Mux.Unlock()
		}
	}
	// fmt.Println("slice", api_user.MatchingSlice)
}

func unset(s []api_user.Matching, i int) []api_user.Matching {
	if i >= len(s) {
		return s
	}
	return append(s[:i], s[i+1:]...)
}

func hubenyDistance(src Coodinate, dst Coodinate) float64 {
	dx := degree2radian(dst.Longitude - src.Longitude)
	dy := degree2radian(dst.Latitude - src.Latitude)
	my := degree2radian((src.Latitude + dst.Latitude) / 2)

	W := math.Sqrt(1 - (Power2(ECCENTRICITY) * Power2(math.Sin(my)))) // 卯酉線曲率半径の分母
	m_numer := EQUATORIAL_RADIUS * (1 - Power2(ECCENTRICITY))         // 子午線曲率半径の分子

	M := m_numer / math.Pow(W, 3) // 子午線曲率半径
	N := EQUATORIAL_RADIUS / W    // 卯酉線曲率半径

	d := math.Sqrt(Power2(dy*M) + Power2(dx*N*math.Cos(my)))

	return d
}

func degree2radian(x float64) float64 {
	return x * math.Pi / 180
}

func Power2(x float64) float64 {
	return math.Pow(x, 2)
}

func HttpPost(url string, token string, json api_user.AIDataForm) (int, error) {
	// url = "http://localhost:8000"
	fmt.Println("url", url, "token", token)
	s, err := renameJson.Marshal(json)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(s),
	)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	// req.Header.Set("Accept", "*/*")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer resp.Body.Close()

	var result []api_user.AIDataResult
	fmt.Println("Body", resp.Body)
	if err := renameJson.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println(err)

		return 0, nil
	}
	fmt.Println("result", result)
	return 0, nil
	// return result[0], nil
}
