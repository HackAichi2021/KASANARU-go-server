package main

import (
	"bytes"
	renameJson "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"hackaichi2021/database"
	"hackaichi2021/user"
	api_user "hackaichi2021/user"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/cors"
)

const (
	DISTANCE_FROM_USERS = 150 // マッチング相手との距離
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
	r.Handle("/api/user/feedback/post", api_user.FeedbackPost).Methods("POST")
	go monitor()

	c := cors.Default()
	muxWithMiddlewares := http.TimeoutHandler(r, time.Second*20, "Timeout!")
	chain := alice.New(c.Handler, logHandler).Then(muxWithMiddlewares)

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

func monitor() {
	for {
		time.Sleep(3 * time.Second) // 3秒待つ
		if len(api_user.MatchingGlobal.MatchingSlice[0]) > 0 {
			fmt.Println("ok")
			fmt.Println("len", len(api_user.MatchingGlobal.MatchingSlice[0]), len(api_user.MatchingGlobal.MatchingSlice[1]))
			api_user.MatchingGlobal.Mux.Lock()
			var maxValue float64
			var maxIndex int = -1
			var match api_user.Matching
			for i, v := range api_user.MatchingGlobal.MatchingSlice[1] {
				fmt.Println("inside")
				if distance(api_user.MatchingGlobal.MatchingSlice[0][0].Info.Latitude,
					api_user.MatchingGlobal.MatchingSlice[0][0].Info.Longitude, v.Info.Latitude, v.Info.Longitude, "M") < DISTANCE_FROM_USERS {
					tmp := []api_user.AIFavoriteForm{
						{
							Age1:      api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Age,
							Sex1:      api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Sex,
							Game1:     api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Game,
							Sport1:    api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Sport,
							Book1:     api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Book,
							Travel1:   api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Travel,
							Internet1: api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Internet,
							Anime1:    api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Anime,
							Movie1:    api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Movie,
							Music1:    api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Music,
							Gourmet1:  api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Gourmet,
							Mucle1:    api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Muscle,
							Camp1:     api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Camp,
							Tv1:       api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Tv,
							Cook1:     api_user.MatchingGlobal.MatchingSlice[0][0].Favorite.Cook,
							Age2:      v.Favorite.Age,
							Sex2:      v.Favorite.Sex,
							Game2:     v.Favorite.Game,
							Sport2:    v.Favorite.Sport,
							Book2:     v.Favorite.Book,
							Travel2:   v.Favorite.Travel,
							Internet2: v.Favorite.Internet,
							Anime2:    v.Favorite.Anime,
							Movie2:    v.Favorite.Movie,
							Music2:    v.Favorite.Music,
							Gourmet2:  v.Favorite.Gourmet,
							Mucle2:    v.Favorite.Muscle,
							Camp2:     v.Favorite.Camp,
							Tv2:       v.Favorite.Tv,
							Cook2:     v.Favorite.Cook,
						},
					}
					item := api_user.AIDataForm{
						Data: tmp,
					}
					f, err := HttpPost(os.Getenv("URL"), os.Getenv("AUTHENTICATION"), item)
					fmt.Println("f", f)
					if err != nil {
						fmt.Println("err", err)
					}
					if maxValue < f {
						maxValue = f
						match = v
						maxIndex = i
					}
				}

				if maxIndex != -1 {
					fmt.Println("貸す側への返答", match)
					api_user.MatchingGlobal.NotifiesLend[api_user.MatchingGlobal.MatchingSlice[0][0].Info.AccessToken] <- match
					delete(api_user.MatchingGlobal.NotifiesLend, api_user.MatchingGlobal.MatchingSlice[0][0].Info.AccessToken)

					fmt.Println("借りる側への返答", api_user.MatchingGlobal.MatchingSlice[0][0])
					api_user.MatchingGlobal.NotifiesLend[match.Info.AccessToken] <- api_user.MatchingGlobal.MatchingSlice[0][0]
					delete(api_user.MatchingGlobal.NotifiesLend, match.Info.AccessToken)

					api_user.MatchingGlobal.MatchingSlice[0] = unset(api_user.MatchingGlobal.MatchingSlice[0], 0)
					api_user.MatchingGlobal.MatchingSlice[1] = unset(api_user.MatchingGlobal.MatchingSlice[1], maxIndex)
				} else {
					if len(api_user.MatchingGlobal.MatchingSlice[0]) >= 2 {
						tmp := api_user.MatchingGlobal.MatchingSlice[0][0]
						api_user.MatchingGlobal.MatchingSlice[0] = api_user.MatchingGlobal.MatchingSlice[0][1:]
						api_user.MatchingGlobal.MatchingSlice[0] = append(api_user.MatchingGlobal.MatchingSlice[0][1:], tmp)
					}
				}
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

func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		} else if unit[0] == "M" {
			dist = dist * 1.609344
			dist /= 1000
		}
	}

	return dist
}

func HttpPost(url string, token string, json api_user.AIDataForm) (float64, error) {
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

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r := regexp.MustCompile(`\d+\.\d+`)
	u := r.FindAllStringSubmatch(string(body), 1)
	fmt.Println("u", u[0][0])
	a, _ := strconv.ParseFloat(u[0][0], 64)

	fmt.Println("result", result)
	return a, nil
}
