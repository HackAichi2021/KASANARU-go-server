package user

import (
	"encoding/json"
	"fmt"
	"hackaichi2021/auth"
	"hackaichi2021/crypto"
	"hackaichi2021/database"
	"net/http"
)

type LoginForm struct {
	Email    string `json:"email" binding: "required"`
	Password string `json: "password" binding: "required`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type AuthenticateResponse struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	AccessToken  string `json: "access_token"`
	RefreshToken string `json: "refresh_token"`
}

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

	responseAuthenticate(w, http.StatusNoContent, token)

	// response := AuthenticateResponse{
	// 	Status:       "Success",
	// 	Message:      "User registered successfully",
	// 	AccessToken:  token.AccessToken,
	// 	RefreshToken: token.RefreshToken,
	// }
	// json, _ := json.Marshal(response)

	// w.Write(json)
	// if err := responseAuthenticate(w, http.StatusNoContent, token); err != nil {
	// 	fmt.Println("eeerr")
	// }
})

func responseError(w http.ResponseWriter, statusCode int) error {
	w.WriteHeader(statusCode)
	response := Response{
		Status:  "Error",
		Message: http.StatusText(statusCode),
	}

	json, err := json.Marshal(response)
	if err != nil {
		return err
	}

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
