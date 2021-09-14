package database

import (
	"fmt"
	"hackaichi2021/crypto"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	Id        int64  `json:"exampleId" gorm:"primaryKey"`
	UserName  string `json:"username" binding:"required" gorm:"type:varchar(255);not null"`
	Email     string `json:"email" binding:"required" gorm:"type:varchar(255);not null"`
	Password  string `json:"password" binding:"required" gorm:"type:varchar(1024);not null"`
	Age       int64  `json:"age" binding:"required" gorm:"not null"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime"`
	DeletedAt int64  `json:"deletedAt"`
}

func GormConnect() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: strings.Join([]string{
			"host=" + os.Getenv("HOST"),
			"dbname=" + os.Getenv("DB_NAME"),
			"user=" + os.Getenv("DB_USER"),
			"password=" + os.Getenv("DB_PASSWORD"),
			"port=" + os.Getenv("DB_PORT"),
			"sslmode=" + os.Getenv("SSLMODE"),
		}, " "),
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            false,
		Logger:                 logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("success")
	}

	return db
}

func CreateUser(u User) error {
	u.Password, _ = crypto.PasswordEncrypt(u.Password)

	db_conn := GormConnect()
	db, err := db_conn.DB()
	if err != nil {
		return err
	}
	defer db.Close()

	db_conn.Create(&u)

	if err != nil {
		return err
	}

	return nil
}
