package user

import (
	"encoding/json"
	"hackaichi2021/auth"
	"hackaichi2021/database"
	"net/http"
)

var Register = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var form database.User
	json.NewDecoder(r.Body).Decode(&form)
	statusCode := database.CreateUser(form)
	w.WriteHeader(statusCode)
	if statusCode != http.StatusCreated {
		type Response struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}

		response := Response{
			Status:  "Error",
			Message: "Your account is already registered",
		}
		json, _ := json.Marshal(response)

		w.Write(json)

	} else {
		token, _ := auth.CreateToken(form.Email)

		type Response struct {
			Status       string `json:"status"`
			Message      string `json:"message"`
			AccessToken  string `json: "access_token"`
			RefreshToken string `json: "refresh_token"`
		}
		response := Response{
			Status:       "Success",
			Message:      "User registered successfully",
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		}
		json, _ := json.Marshal(response)

		w.Write(json)

	}

})
