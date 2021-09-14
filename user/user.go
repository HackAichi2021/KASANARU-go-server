package user

import (
	"encoding/json"
	"hackaichi2021/database"
	"net/http"
	"strconv"
)

var Register = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var form database.User
	json.NewDecoder(r.Body).Decode(&form)
	statusCode := database.CreateUser(form)
	w.WriteHeader(statusCode)
	w.Write([]byte(strconv.Itoa(statusCode) + " - " + http.StatusText(statusCode)))

})
