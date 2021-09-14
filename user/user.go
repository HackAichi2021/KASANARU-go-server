package user

import (
	"encoding/json"
	"hackaichi2021/auth"
	"hackaichi2021/database"
	"net/http"
)

type Login struct {
	Email    string `json:"email" binding: "required"`
	Password string `json: "password" binding: "required`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
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
		token, _ := auth.CreateToken(form.Email)

		type RegisterResponse struct {
			Status       string `json:"status"`
			Message      string `json:"message"`
			AccessToken  string `json: "access_token"`
			RefreshToken string `json: "refresh_token"`
		}
		response := RegisterResponse{
			Status:       "Success",
			Message:      "User registered successfully",
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		}
		json, _ := json.Marshal(response)

		w.Write(json)

	}

})

// func checkJson(w http.ResponseWriter, r *http.Request) (int, error) {
// 	if r.Header.Get("Content-Type") != "application/json" {
// 		return http.StatusBadRequest, errors.New(JsonParseErr)
// 	}

// 	_, err := strconv.Atoi(r.Header.Get("Content-Length"))
// 	if err != nil {
// 		return http.StatusInternalServerError, errors.New(JsonParseErr)
// 	}

// 	body, err := ioutil.ReadAll(r.Body)
// 	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
// 	if err != nil && err != io.EOF {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return http.StatusInternalServerError, errors.New(JsonParseErr)
// 	}
// 	return 0, nil
// }

// var login = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 	var form database.User
// 	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {

// 	}

// 	if statusCode, err := checkJson(w, r); err != nil {
// 		w.WriteHeader(statusCode)
// 		fmt.Fprintf(w, JsonInvalid)
// 		return
// 	}

// 	var u User
// 	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		fmt.Fprintf(w, AuthenticateFailed)
// 		return
// 	}

// 	if user.UserName != u.UserName || user.Password != u.Password {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		fmt.Fprintf(w, AuthenticateFailed)
// 		return
// 	}
// 	token, err := auth.CreateToken(user.ID)
// 	if err != nil {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		fmt.Fprintf(w, AuthenticateFailed)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(token)
// })
