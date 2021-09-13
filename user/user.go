package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
	}

	env := os.Getenv("ENV")
	DB := os.Getenv("DB")
	DBUser := os.Getenv("DB_USER")
	DBPass := os.Getenv("DB_PASS")

	fmt.Println(env)    // local
	fmt.Println(DB)     // test_db
	fmt.Println(DBUser) // root
	fmt.Println(DBPass) // pass1234
}
