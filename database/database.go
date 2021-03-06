package database

import (
	"fmt"
	"hackaichi2021/crypto"
	"net/http"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	Id        int    `json:"exampleId" gorm:"primaryKey"`
	UserName  string `json:"username" binding:"required" gorm:"type:varchar(255);not null"`
	Email     string `json:"email" binding:"required" gorm:"type:varchar(255);not null"`
	Password  string `json:"password" binding:"required" gorm:"type:varchar(1024);not null"`
	Age       int    `json:"age" binding:"required" gorm:"not null"`
	UpdatedAt int64  `json:"updatedAt" gorm:"autoUpdateTime"`
	CreatedAt int64  `json:"createdAt" gorm:"autoCreateTime"`
	DeletedAt int64  `json:"deletedAt"`
}

type Feedback struct {
	Categories string `json:"categories"`
	Star       int64  `json:"star"`
}

type Favorite struct {
	UserId   int `json:"id" gorm:"unique;not null"`
	Age      int `json:"age" binding:"required"`
	Sex      int `json:"sex" binding:"required"`
	Game     int `json:"game" binding:"required"`
	Sport    int `json:"sport" binding:"required"`
	Book     int `json:"book" binding:"required"`
	Travel   int `json:"travel" binding:"required"`
	Internet int `json:"internet" binding:"required"`
	Anime    int `json:"anime" binding:"required"`
	Movie    int `json:"movie" binding:"required"`
	Music    int `json:"music" binding:"required"`
	Gourmet  int `json:"gourmet" binding:"required"`
	Muscle   int `json:"muscle" binding:"required"`
	Camp     int `json:"camp" binding:"required"`
	Tv       int `json:"tv" binding:"required"`
	Cook     int `json:"cook" binding:"required"`
}

func GormConnect() *gorm.DB {

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

func CreateUser(u User) int {
	u.Password, _ = crypto.PasswordEncrypt(u.Password)

	db_conn := GormConnect()
	db, err := db_conn.DB()
	if err != nil {
		return http.StatusInternalServerError
	}
	defer db.Close()

	if result := db_conn.Create(&u); result.Error != nil {
		return http.StatusBadRequest
	}

	return http.StatusCreated
}

func GetIdByEmail(email string) []User {
	db_conn := GormConnect()
	db, err := db_conn.DB()
	if err != nil {
		return nil
	}
	defer db.Close()

	item := []User{}
	db_conn.Find(&item, "email=?", email)
	fmt.Println("item", item)
	return item
}

func GetUserByUserId(id int) []User {
	db_conn := GormConnect()
	db, err := db_conn.DB()
	if err != nil {
		return nil
	}
	defer db.Close()

	item := []User{}
	db_conn.Find(&item, "id=?", id)
	fmt.Println("item", item)
	return item
}

func GetOneColumnValueUser(column string, email string) []User {
	db_conn := GormConnect()
	db, err := db_conn.DB()
	if err != nil {
		return nil
	}
	defer db.Close()

	item := []User{}
	db_conn.Find(&item, column+"=?", email)
	fmt.Println("item", item)
	return item

}

func InsertOrUpdateFavorite(item Favorite) error {
	db_conn := GormConnect()
	db, err := db_conn.DB()
	if err != nil {
		return nil
	}
	defer db.Close()

	var u []Favorite
	db_conn.Find(&u, "user_id=?", item.UserId)
	fmt.Println("u", u)
	if len(u) == 0 {
		if result := db_conn.Create(&item); result.Error != nil {
			return result.Error
		}
	} else {
		// db_conn.Save(item)
		result := db_conn.Model(&Favorite{}).Where("user_id=?", u[0].UserId).Updates(&item)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func GetFavorite(id int) []Favorite {
	db_conn := GormConnect()
	db, err := db_conn.DB()
	if err != nil {
		return nil
	}
	defer db.Close()

	var item []Favorite
	db_conn.Find(&item, "user_id=?", id)
	if len(item) == 1 {
		return item
	}
	return nil
}

func InsertFeedback(item Feedback) error {
	db_conn := GormConnect()
	db, err := db_conn.DB()
	if err != nil {
		return nil
	}
	defer db.Close()

	if result := db_conn.Table("feedbacks").Create(&item); result.Error != nil {
		return result.Error
	}
	return nil
}
