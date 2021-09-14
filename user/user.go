package user

import (
	"encoding/json"
	"hackaichi2021/database"
	"net/http"
)

var Register = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var form database.User
	json.NewDecoder(r.Body).Decode(&form)
	database.CreateUser(form)
	json.NewEncoder(w).Encode(form)

})
