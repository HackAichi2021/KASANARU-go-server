package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
	}

	fmt.Println("HOST", os.Getenv("HOST"))
	fmt.Println("USER", os.Getenv("DB_USER"))
	fmt.Println("PASSWORD", os.Getenv("PASSWORD"))
	fmt.Println("DB_NAME", os.Getenv("DB_NAME"))
	fmt.Println("PORT", os.Getenv("PORT"))
	fmt.Println("SSLMODE", os.Getenv("SSLMODE"))
	fmt.Println("TIME_ZONE", os.Getenv("TIME_ZONE"))

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("HOST"), os.Getenv("DB_USER"), os.Getenv("PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("PORT"), os.Getenv("SSLMODE"), os.Getenv("TIME_ZONE"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("success")
	}
	fmt.Println(db)
}
