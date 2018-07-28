package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

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

var db *gorm.DB

// var mg *mailgun.Mailgun

func init() {
	var err error

	// accounts := []Account{
	// 	Account{Email: "wow@test.com", GlobalPermissions: 3, UserName: "wow", Password: "p"},
	// 	Account{Email: "nice@test.com", GlobalPermissions: 3, UserName: "nice", Password: "p"},
	// 	Account{Email: "good@test.com", GlobalPermissions: 3, UserName: "good", Password: "p"},
	// 	Account{Email: "ok@test.com", GlobalPermissions: 3, UserName: "ok", Password: "p"},
	// }

	db, err = gorm.Open(
		"postgres",
		"host=0.0.0.0 port=5432 user=postgres password=postgrespassword sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Account{}, &Game{}, &GameAccount{}, &Match{}, &PendingAccount{}, &Team{}, &TeamMember{})

	// mg = mailgun.NewMailgun("mail.therileyjohnson.com", private, public)
	// for _, account := range accounts {
	// 	var existing_account Account
	// 	if err := db.Where("email = ?", account.Email).First(&existing_account).Error; err != nil {
	// 		fmt.Println(err)
	// 		db.Create(&account)
	// 		db.Where("email = ?", account.Email).First(&existing_account)
	// 	}
	// 	// else {
	// 	// 	//fmt.Println(existing_account)
	// 	// }

	// 	// if db.NewRecord(account) {
	// 	// 	db.Create(&account)
	// 	// }

	// 	gameAccount := GameAccount{AccountID: existing_account.AccountID, GamePermissions: 3, GameID: 2, Score: 100}
	// 	var existing_gameaccount GameAccount

	// 	if err := db.Where("game_id = ? and account_id = ?", gameAccount.GameID, existing_account.ID).First(&existing_gameaccount).Error; err != nil {
	// 		db.Create(&gameAccount)
	// 		fmt.Println("Noice")
	// 	} else {
	// 		fmt.Println(existing_gameaccount)
	// 	}
	// }
}

func main() {
	defer db.Close()

	router := gin.Default()

	api := router.Group("/api")
	{
		//No Auth
		api.GET("/game/rankings/:gid", func(context *gin.Context) {

		})
		api.GET("/match/:mid", func(context *gin.Context) {

		})
		api.GET("/matches/game/:gid", func(context *gin.Context) {

		})
		api.GET("/matches/user/:pid", func(context *gin.Context) {

		})
		api.GET("/matches/user/:pid/game/:gid", func(context *gin.Context) {

		})
		api.GET("/matches/team/:tid", func(context *gin.Context) {

		})
		api.GET("/team/:tid", func(context *gin.Context) {

		})
		api.GET("/user/account/:aid", func(context *gin.Context) {

		})
		api.GET("/user/account/:aid/game/:gid", func(context *gin.Context) {

		})
		api.GET("/user/confirm/:uuid", func(context *gin.Context) {
			uuid := context.Param("uuid")

			pendingUser := PendingAccount{}

			if err := db.Where(PendingAccount{UUID: uuid}).First(&pendingUser).Error; err != nil {
				errorResponse(
					context,
					"Sorry, an error occurred when confirming your account, please check that you have the correct link and try again.",
					err.Error(),
				)
			}

			newUser := Account{
				Email:     pendingUser.Email,
				FirstName: pendingUser.FirstName,
				LastName:  pendingUser.LastName,
				UserName:  pendingUser.UserName,
				Password:  pendingUser.Password,
			}

			account, err := CreateAccount(db, &newUser)

			if err != nil {
				errorResponse(
					context,
					"Sorry, an error occurred when confirming your account, please check that you have the correct link and try again.",
					err.Error(),
				)
				return
			}

			db.Delete(pendingUser)

			context.JSON(
				200,
				gin.H{
					"data":  account,
					"error": false,
				},
			)
		})
		api.POST("/user/create", func(context *gin.Context) {
			pendingUser := PendingAccount{}

			decoder := json.NewDecoder(context.Request.Body)

			if err := decoder.Decode(&pendingUser); err != nil {
				errorResponse(
					context,
					"Error recieving account information, please try again.",
					err.Error(),
				)
				return
			}

			if pendingUser.FirstName == "" || pendingUser.LastName == "" || pendingUser.Password == "" || pendingUser.UserName == "" {
				errorResponse(
					context,
					"None of the fields can be blank.",
					"Blank field.",
				)
				return
			}

			matchEmail := false
			emailDomains := []string{"trojans.dsu.edu", "pluto.dsu.edu", "dsu.edu"}
			pendingUser.Email = strings.ToLower(pendingUser.Email)

			for _, emailRegex := range emailDomains {
				r := regexp.MustCompile(fmt.Sprintf(`^[A-Za-z0-9][A-Za-z0-9_\+\.]*@%s$`, emailRegex))
				if r.Match([]byte(pendingUser.Email)) {
					matchEmail = true
				}
			}

			if !matchEmail {
				errorResponse(
					context,
					"Email is invalid, must be a valid email with either a trojans.dsu.edu, pluto.dsu.edu, or dsu.edu domain.",
					"Email no match.",
				)
				return
			}

			if len(pendingUser.Email) > 150 {
				context.JSON(
					400,
					gin.H{
						"error": true,
						"msg":   "Email is too long (over 150 characters).",
					},
				)
				return
			}

			if len(pendingUser.UserName) > 25 {
				context.JSON(
					400,
					gin.H{
						"error": true,
						"msg":   "Username is too long (over 25 characters).",
					},
				)
				return
			}

			pendingUser.UUID = getUUID()

			pendingAccount, err := CreatePendingAccount(db, &pendingUser)

			if err != nil {
				if terr, ok := err.(*NotUniqueUserError); ok {
					errorResponse(
						context,
						terr.Error(),
						err.Error(),
					)
				} else {
					errorResponse(
						context,
						"Error creating account, try again.",
						err.Error(),
					)
				}
				return
			}

			// Reference https://github.com/the-rileyj/DSU_Chess/blob/master/chess.go
			// _, _, err = mg.Send(mailgun.NewMessage("robot@mail.therileyjohnson.com", "Registration", fmt.Sprintf("Click %s/%s to confirm your email!", URL, pendingAccount.UUID), email))

			context.JSON(
				200,
				gin.H{
					"data":  pendingAccount,
					"error": false,
				},
			)
		})
		api.GET("/user/ranking/:pid", func(context *gin.Context) {

		})
		api.GET("/user/ranking/:pid/game/:gid", func(context *gin.Context) {

		})

		// Auth
		api.POST("/user/create/game/:gid", func(context *gin.Context) {

		})
		api.DELETE("/user/delete/:aid/game/:gid", func(context *gin.Context) {

		})
		api.DELETE("/user/delete/:aid", func(context *gin.Context) {

		})
		api.PUT("/user/update/:aid", func(context *gin.Context) {

		})
		api.POST("/match/create", func(context *gin.Context) {

		})
		api.DELETE("/match/delete/:mid", func(context *gin.Context) {

		})
		api.PUT("/match/update/:mid", func(context *gin.Context) {

		})
		api.POST("/game/create", func(context *gin.Context) {

		})
		api.POST("/game/delete/:gid", func(context *gin.Context) {

		})
	}

	router.Run(":3800")
}
