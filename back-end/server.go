package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	permGameCreation  = 1
	permGameDeletion  = 3
	permPromoteRange1 = 2
	permPromoteRange2 = 3
)

var db *gorm.DB

// var mg *mailgun.Mailgun

func init() {
	var err error
	var updateTable bool

	flagTableUpdate := flag.Bool("table-update", false, "Drops the current tables, and updates the tables.")

	flag.Parse()

	updateTable = *flagTableUpdate

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

	if updateTable {
		// Note, write db backup functionality
		db.DropTableIfExists(&Account{}, &Game{}, &GameAccount{}, &Match{}, &PendingAccount{}, &Team{}, &TeamMember{})
		db.CreateTable(&Account{}, &Game{}, &GameAccount{}, &Match{}, &PendingAccount{}, &Team{}, &TeamMember{})
	}
}

func main() {
	defer db.Close()

	router := gin.Default()

	api := router.Group("/api")
	{
		/*************\
		    No Auth
		\*************/
		api.GET("/game/rankings/:gid", getRankingsForGame)
		api.POST("/login", postLogin)
		api.GET("/match/:mid", getMatch)
		api.GET("/matches/game/:gid", getGameMatches)
		api.GET("/matches/user/:aid", getUserMatches)
		api.GET("/matches/user/:pid/game/:gid", getGameMatchesForUser)
		// api.GET("/matches/team/:tid", func(context *gin.Context) { })
		api.GET("/team/:tid", getTeam)
		api.GET("/user/account/:aid", getUserAccount)
		api.GET("/user/account/:aid/game/:gid", getUserAccountForGame)
		api.GET("/user/confirm/:uuid", getConfirmAccount)
		api.POST("/user/create", postCreateAccount)
		api.GET("/user/rankings/:aid", getRankingsForUserAccount)
		api.GET("/user/ranking/:aid/game/:gid", getRankingsForUserAccountInGame)

		/**********\
		    Auth
		\**********/
		api.DELETE("/user/delete/:aid", assureAuthentication(deleteUser))
		api.DELETE("/game/delete/:gid", assureAuthentication(deleteGame))
		api.DELETE("/match/delete/:mid", assureAuthentication(deleteMatch))
		api.POST("/game/create", assureAuthentication(postCreateGame))
		api.POST("/match/create", assureAuthentication(postCreateMatch))
		api.POST("/user/update", assureAuthentication(postUpdateUser))
		api.PUT("/user/disable/:aid/game/:gid", assureAuthentication(putDisableGameAccount))
		api.PUT("/user/enable/:aid/game/:gid", assureAuthentication(putEnableGameAccount))
		api.PUT("/match/confirm/:mid", assureAuthentication(putMatchConfirm))
		api.PUT("/match/deny/:mid", assureAuthentication(putMatchDeny))
	}

	router.Run(":3800")
}
