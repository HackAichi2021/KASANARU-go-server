package user

import (
	"encoding/json"
	"fmt"
	"hackaichi2021/auth"
	"hackaichi2021/crypto"
	"hackaichi2021/database"
	"net/http"
	"reflect"
	"strconv"
	"sync"

	jwt "github.com/dgrijalva/jwt-go"
)

type LoginForm struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateForm struct {
	Sex          int    `json:"sex" binding:"required"`
	Game         int    `json:"game" binding:"required"`
	Sport        int    `json:"sport" binding:"required"`
	Book         int    `json:"book" binding:"required"`
	Travel       int    `json:"travel" binding:"required"`
	Internet     int    `json:"internet" binding:"required"`
	Anime        int    `json:"anime" binding:"required"`
	Movie        int    `json:"movie" binding:"required"`
	Music        int    `json:"music" binding:"required"`
	Gourmet      int    `json:"gourmet" binding:"required"`
	Muscle       int    `json:"muscle" binding:"required"`
	Camp         int    `json:"camp" binding:"required"`
	Tv           int    `json:"tv" binding:"required"`
	Cook         int    `json:"cook" binding:"required"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MatchingForm struct {
	Latitude     float64 `json:"latitude" binding:"required"`
	Longitude    float64 `json:"longitude" binding:"required"`
	Lend         int     `json:"lend" binding:"required"`
	AfterArrival int     `json:"after_arrival" binding:"required"`
	Token
}

type Token struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type Matching struct {
	UserId   int `json:"user_id" binding:"required"`
	Info     MatchingForm
	Favorite database.Favorite
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type AuthenticateResponse struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type FavoriteGetReseponse struct {
	Favorite database.Favorite
	Response Response
	UserName string `json:"username"`
	Age      int    `json:"age" binding:"required"`
}

type AIFavoriteForm struct {
	Age1      int `json:"age1" binding:"required"`
	Sex1      int `json:"sex1" binding:"required"`
	Game1     int `json:"game1" binding:"required"`
	Sport1    int `json:"sport1" binding:"required"`
	Book1     int `json:"book1" binding:"required"`
	Travel1   int `json:"travel1" binding:"required"`
	Internet1 int `json:"internet1" binding:"required"`
	Anime1    int `json:"anime1" binding:"required"`
	Movie1    int `json:"movie1" binding:"required"`
	Music1    int `json:"music1" binding:"required"`
	Gourmet1  int `json:"gourmet1" binding:"required"`
	Mucle1    int `json:"mucle1" binding:"required"`
	Camp1     int `json:"camp1" binding:"required"`
	Tv1       int `json:"tv1" binding:"required"`
	Cook1     int `json:"cook1" binding:"required"`
	Age2      int `json:"age2" binding:"required"`
	Sex2      int `json:"sex2" binding:"required"`
	Game2     int `json:"game2" binding:"required"`
	Sport2    int `json:"sport2" binding:"required"`
	Book2     int `json:"book2" binding:"required"`
	Travel2   int `json:"travel2" binding:"required"`
	Internet2 int `json:"internet2" binding:"required"`
	Anime2    int `json:"anime2" binding:"required"`
	Movie2    int `json:"movie2" binding:"required"`
	Music2    int `json:"music2" binding:"required"`
	Gourmet2  int `json:"gourmet2" binding:"required"`
	Mucle2    int `json:"mucle2" binding:"required"`
	Camp2     int `json:"camp2" binding:"required"`
	Tv2       int `json:"tv2" binding:"required"`
	Cook2     int `json:"cook2" binding:"required"`
}

type AIDataForm struct {
	Data []AIFavoriteForm `json:"data"`
}

type AIDataResult struct {
	Result []int `json:"result"`
}

//貸す人へのresponse 借りる人には緯度、経度は0で渡す
type LendResponse struct {
	UserName  string  `json:"user_name" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

type SafeLend struct {
	MatchingSlice [2][]Matching
	NotifiesLend  map[string](chan Matching)
	Mux           sync.Mutex
}

var MatchingGlobal = new(SafeLend)

var Register = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var form database.User
	json.NewDecoder(r.Body).Decode(&form)
	statusCode := database.CreateUser(form)
	w.WriteHeader(statusCode)
	if statusCode != http.StatusCreated {
		response := Response{
			Status:  "Error",
			Message: "Your account is already registered",
		}
		json, _ := json.Marshal(response)

		w.Write(json)

	} else {
		token, _ := auth.CreateTokenByUserIdWithEmail(form.Email)

		response := AuthenticateResponse{
			Status:       "Success",
			Message:      "User registered successfully",
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		}
		json, _ := json.Marshal(response)

		w.Write(json)

	}

})

var Login = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var form LoginForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		fmt.Println(err)
	}
	fmt.Println("form", form)
	user := database.GetOneColumnValueUser("email", form.Email)
	if len(user) == 0 {
		fmt.Println("email error")
		if e := responseError(w, http.StatusUnauthorized); e != nil {
			fmt.Println(e)
		}
		return
	}

	if err := crypto.CompareHashAndPassword(user[0].Password, form.Password); err != nil {
		fmt.Println("password error")
		if e := responseError(w, http.StatusUnauthorized); e != nil {
			fmt.Println(e)
		}
		return
	}

	token, _ := auth.CreateTokenByUserIdWithEmail(form.Email)

	w.WriteHeader(http.StatusOK)
	response := AuthenticateResponse{
		Status:       "Success",
		Message:      "Login successfully",
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	json, _ := json.Marshal(response)

	w.Write(json)
	// if err := responseAuthenticate(w, http.StatusNoContent, token); err != nil {
	// 	fmt.Println("eeerr")
	// }
})

var Update = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var form UpdateForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		response := Response{
			Status:  "Error",
			Message: "Update failed",
		}
		json, _ := json.Marshal(response)

		w.Write(json)
		return
	}

	fmt.Println("form", form)
	tokenString := form.AccessToken
	claims := jwt.MapClaims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("SIGNINGKEY"), nil
	})

	fmt.Println(token)

	c, ok := claims["user_id"].(float64)
	var id int
	if ok {
		id = int(c)
	}

	user := database.GetOneColumnValueUser("id", strconv.Itoa(id))
	if len(user) > 0 {
		favorite := database.Favorite{
			UserId:   id,
			Age:      user[0].Age,
			Sex:      form.Sex,
			Game:     form.Game,
			Sport:    form.Sport,
			Book:     form.Book,
			Travel:   form.Travel,
			Internet: form.Internet,
			Anime:    form.Anime,
			Movie:    form.Movie,
			Music:    form.Music,
			Gourmet:  form.Gourmet,
			Muscle:   form.Muscle,
			Camp:     form.Camp,
			Tv:       form.Tv,
			Cook:     form.Cook,
		}
		if err := database.InsertOrUpdateFavorite(favorite); err != nil {
			fmt.Println(err)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		response := Response{
			Status:  "Error",
			Message: "Update failed",
		}
		json, _ := json.Marshal(response)

		w.Write(json)
		return
	}
	fmt.Println("uudisaj")

	w.WriteHeader(http.StatusNoContent)

})

var Match = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var form MatchingForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		response := Response{
			Status:  "Error",
			Message: "Match failed",
		}
		json, _ := json.Marshal(response)
		w.Write(json)
		return
	}

	fmt.Println(form.AccessToken)
	tmp := database.GetFavorite(decodeUserIdFromAccessToken(form.AccessToken))
	fmt.Println("tuuka")
	var item Matching
	if len(tmp) > 0 {
		item = Matching{
			UserId:   decodeUserIdFromAccessToken(form.AccessToken),
			Info:     form,
			Favorite: tmp[0],
		}
	} else {
		item = Matching{
			UserId: decodeUserIdFromAccessToken(form.AccessToken),
			Info:   form,
		}
	}
	MatchingGlobal.Mux.Lock()
	fmt.Println("rend", form.Lend)
	MatchingGlobal.MatchingSlice[form.Lend] = append(MatchingGlobal.MatchingSlice[form.Lend], item)
	MatchingGlobal.Mux.Unlock()

	MatchingGlobal.Mux.Lock()
	MatchingGlobal.NotifiesLend[form.AccessToken] = make(chan Matching)
	MatchingGlobal.Mux.Unlock()

	fmt.Println("wait")
	b := <-MatchingGlobal.NotifiesLend[form.AccessToken]
	fmt.Println("answer", b)

	fmt.Println("id", decodeUserIdFromAccessToken(form.AccessToken))
	w.WriteHeader(http.StatusOK)
	jsonStr, _ := json.Marshal(b)
	w.Write(jsonStr)

})

var FavoriteGet = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var form UpdateForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		response := Response{
			Status:  "Error",
			Message: "Favorite Get failed",
		}
		json, _ := json.Marshal(response)
		w.Write(json)
		return
	}

	favorite := database.GetFavorite(decodeUserIdFromAccessToken(form.AccessToken))
	user := database.GetUserByUserId(decodeUserIdFromAccessToken(form.AccessToken))
	fmt.Println(favorite)
	if len(favorite) > 0 && len(user) > 0 {
		w.WriteHeader(http.StatusOK)
		response := FavoriteGetReseponse{
			Favorite: favorite[0],
			Response: Response{
				Status:  "Success",
				Message: "Favorite get successfully",
			},
			UserName: user[0].UserName,
			Age:      user[0].Age,
		}
		fmt.Println("response", response)
		json, _ := json.Marshal(response)

		w.Write(json)
	} else if len(user) > 0 {
		w.WriteHeader(http.StatusOK)
		response := FavoriteGetReseponse{
			Response: Response{
				Status:  "Success",
				Message: "Favorite get successfully",
			},
			UserName: user[0].UserName,
			Age:      user[0].Age,
		}
		json, _ := json.Marshal(response)

		w.Write(json)
	}

})

func P(t interface{}) {
	fmt.Println(reflect.TypeOf(t))
}

func responseError(w http.ResponseWriter, statusCode int) error {
	response := Response{
		Status:  "Error",
		Message: http.StatusText(statusCode),
	}

	json, err := json.Marshal(response)
	if err != nil {
		return err
	}

	w.WriteHeader(statusCode)
	w.Write(json)
	return nil
}

func responseAuthenticate(w http.ResponseWriter, statusCode int, token *auth.TokenDetails) error {
	w.WriteHeader(statusCode)
	response := AuthenticateResponse{
		Status:       "Success",
		Message:      http.StatusText(statusCode),
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	fmt.Println(response)
	json, err := json.Marshal(response)
	if err != nil {
		fmt.Println("aaa")
		return err
	}

	w.Write(json)
	return nil
}

func decodeUserIdFromAccessToken(tokenString string) int {
	claims := jwt.MapClaims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("SIGNINGKEY"), nil
	})

	fmt.Println(token)

	c, ok := claims["user_id"].(float64)
	var id int
	if ok {
		id = int(c)
	}
	return id
}
