package user

import (
	"encoding/json"
	"fmt"
	"hackaichi2021/database"
	"net/http"
)

var Register = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var form database.User
	json.NewDecoder(r.Body).Decode(&form)
	fmt.Println("call api")
	database.CreateUser(form)
	json.NewEncoder(w).Encode(form)

})
