package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//IMPLEMENT CHECKS FOR STRING LENGTHS

// type User struct {
// 	Username string `json: username`
// 	Password string `json: password`
// }

// var router *gin.Engine

func init() {
	accounts := []Account{
		Account{Email: "wow@test.com", GlobalPermissions: 3, UserName: "wow", Password: "p"},
		Account{Email: "nice@test.com", GlobalPermissions: 3, UserName: "nice", Password: "p"},
		Account{Email: "good@test.com", GlobalPermissions: 3, UserName: "good", Password: "p"},
		Account{Email: "ok@test.com", GlobalPermissions: 3, UserName: "ok", Password: "p"},
	}

	db, err := gorm.Open(
		"postgres",
		"host=0.0.0.0 port=5432 user=postgres password=postgrespassword sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.AutoMigrate(&Account{}, &Game{}, &GameAccount{}, &Match{}, &PendingAccount{}, &Team{}, &TeamMember{})

	for _, account := range accounts {
		var existing_account Account
		if err := db.Where("email = ?", account.Email).First(&existing_account).Error; err != nil {
			fmt.Println(err)
			db.Create(&account)
			db.Where("email = ?", account.Email).First(&existing_account)
		}
		// else {
		// 	//fmt.Println(existing_account)
		// }

		// if db.NewRecord(account) {
		// 	db.Create(&account)
		// }

		gameAccount := GameAccount{AccountID: existing_account.AccountID, GamePermissions: 3, GameID: 2, Score: 100}
		var existing_gameaccount GameAccount

		if err := db.Where("game_id = ? and account_id = ?", gameAccount.GameID, existing_account.ID).First(&existing_gameaccount).Error; err != nil {
			db.Create(&gameAccount)
			fmt.Println("Noice")
		} else {
			fmt.Println(existing_gameaccount)
		}
	}

	log.Fatal("NICE")
}

func main() {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/", func(context *gin.Context) {

		})
	}

	// router.Run()
}
