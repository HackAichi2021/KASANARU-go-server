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
	AccessToken  string  `json:"access_token" binding:"required"`
	RefreshToken string  `json:"refresh_token" binding:"required"`
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

var MatchLend []MatchingForm

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
	fmt.Println("hello")
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

	fmt.Println("form", form)
	fmt.Println("match_lend", MatchLend)
	MatchLend = append(MatchLend, form)
	fmt.Println("match_lend", MatchLend)

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
