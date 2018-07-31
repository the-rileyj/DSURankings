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
		db.DropTableIfExists(&Account{}, &Game{}, &GameAccount{}, &Match{}, &PendingAccount{}, &Session{}, &Team{}, &TeamMember{})
		db.CreateTable(&Account{}, &Game{}, &GameAccount{}, &Match{}, &PendingAccount{}, &Session{}, &Team{}, &TeamMember{})
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
		api.GET("/game/rankings/:gid", getRankingsForGame)                       // Check
		api.POST("/login", postLogin)                                            // Check
		api.GET("/match/:mid", getMatch)                                         // Not Check
		api.GET("/matches/game/:gid", getGameMatches)                            // Not Check
		api.GET("/matches/user/:aid", getUserMatches)                            // Not Check
		api.GET("/matches/user/:aid/game/:gid", getGameMatchesForUser)           // Not Check
		api.GET("/team/:tid", getTeam)                                           // Not Check
		api.GET("/user/account/:aid", getUserAccount)                            // Check
		api.GET("/user/account/:aid/game/:gid", getUserAccountForGame)           // Check
		api.GET("/user/confirm/:uuid", getConfirmAccount)                        // Check
		api.POST("/user/create", postCreateAccount)                              // Check
		api.GET("/user/rankings/:aid", getRankingsForUserAccount)                // Check
		api.GET("/user/ranking/:aid/game/:gid", getRankingsForUserAccountInGame) // Check
		// api.GET("/matches/team/:tid", func(context *gin.Context) { }) // Not Check

		/**********\
		    Auth
		\**********/
		api.DELETE("/user/delete", assureAuthentication(deleteUser))                         // Check
		api.DELETE("/game/delete/:gid", assureAuthentication(deleteGame))                    // Not Check
		api.DELETE("/match/delete/:mid", assureAuthentication(deleteMatch))                  // Not Check
		api.POST("/game/create", assureAuthentication(postCreateGame))                       // Check
		api.POST("/match/create", assureAuthentication(postCreateMatch))                     // Not Check
		api.POST("/user/update", assureAuthentication(postUpdateUser))                       // Not Check FINISH
		api.PUT("/game/account/create/:gid", assureAuthentication(putCreateGameAccount))     // Check
		api.PUT("/user/disable/:aid/game/:gid", assureAuthentication(putDisableGameAccount)) // Not Check FINISH
		api.PUT("/user/enable/:aid/game/:gid", assureAuthentication(putEnableGameAccount))   // Not Check FINISH
		api.PUT("/match/confirm/:mid", assureAuthentication(putMatchConfirm))                // Not Check
		api.PUT("/match/deny/:mid", assureAuthentication(putMatchDeny))                      // Not Check
	}

	router.Run(":3800")
}
