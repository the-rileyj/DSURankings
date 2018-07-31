package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func getRankingsForGame(context *gin.Context) {
	gid, err := strconv.ParseUint(context.Param("gid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing ID.",
			err.Error(),
		)
		return
	}

	gameAccounts := make([]GameAccount, 0)

	if rdb := db.Preload("Account").Preload("Game").Where(GameAccount{GameID: gid, Enabled: true}).Order("score desc").Find(&gameAccounts); rdb.Error != nil && !rdb.RecordNotFound() {
		if rdb.RecordNotFound() {
			context.JSON(
				200,
				gin.H{
					"data":  gameAccounts,
					"error": false,
				},
			)
		} else {
			errorResponse(
				context,
				"An error occured getting the information, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	responseGameAccounts := make([]ResponseAdvancedGameAccount, 0)
	for _, gameAccount := range gameAccounts {
		responseGameAccounts = append(responseGameAccounts, gameAccount.AdvancedResponse())
	}

	context.JSON(
		200,
		gin.H{
			"data":  responseGameAccounts,
			"error": false,
		},
	)
}

func postLogin(context *gin.Context) {
	defer context.Request.Body.Close()
	var auth Auther

	decoder := json.NewDecoder(context.Request.Body)

	if err := decoder.Decode(&auth); err != nil {
		errorResponse(
			context,
			"Error checking login information, please try again.",
			err.Error(),
		)
		return
	}

	var user Account

	if rdb := db.Where(Account{Email: auth.Email}).First(&user); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"Incorrect login information, please try again.",
				rdb.Error.Error(),
			)
		} else {
			errorResponse(
				context,
				"An error occured getting the login information, please try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(auth.Password)) != nil {
		errorResponse(
			context,
			"Incorrect login information, please try again.",
			"Bad login", // Change to be in sync with response from above part after debugging
		)
		return
	}

	token := getUUID()

	newSession := Session{
		AccountID: user.AccountID,
		UUID:      token,
		CreatedAt: time.Now(),
	}

	if err := db.Create(&newSession).Error; err != nil {
		errorResponse(
			context,
			"An error occured creating the account, try again.",
			err.Error(),
		)
		return
	}

	context.JSON(
		200,
		gin.H{
			"data":  token,
			"error": false,
		},
	)
}

func getMatch(context *gin.Context) {
	mid, err := strconv.ParseUint(context.Param("mid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing match ID.",
			err.Error(),
		)
		return
	}

	var match Match

	// if rdb := db.Preload("LosingTeam").Preload("LosingTeam.TeamMembers").Preload("LosingTeam.TeamMembers")
	// .Preload("LosingTeam.TeamMembers.Accounts").Preload("LosingTeam.TeamMembers.GameAccounts")
	// .Preload("WinningTeam")

	if rdb := db.
		Preload("LosingTeam.TeamMembers.Accounts").Preload("LosingTeam.TeamMembers.GameAccounts").
		Preload("WinningTeam.TeamMembers.Accounts").Preload("WinningTeam.TeamMembers.GameAccounts").
		Preload("Game").Find(&match, mid); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"A match with that ID does not exist.",
				rdb.Error.Error(),
			)
		} else {
			errorResponse(
				context,
				"An error occured getting the information, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	context.JSON(
		200,
		gin.H{
			"data":  match.AdvancedResponse(),
			"error": false,
		},
	)
}

func getGameMatches(context *gin.Context) {
	gid, err := strconv.ParseUint(context.Param("gid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing match ID.",
			err.Error(),
		)
		return
	}

	var matches []Match

	// if rdb := db.Preload("LosingTeam").Preload("LosingTeam.TeamMembers").Preload("LosingTeam.TeamMembers")
	// .Preload("LosingTeam.TeamMembers.Accounts").Preload("LosingTeam.TeamMembers.GameAccounts")
	// .Preload("WinningTeam")

	if rdb := db.
		Preload("LosingTeam.TeamMembers.Accounts").Preload("LosingTeam.TeamMembers.GameAccounts").
		Preload("WinningTeam.TeamMembers.Accounts").Preload("WinningTeam.TeamMembers.GameAccounts").
		Preload("Game").Where(Match{GameID: gid}).Find(&matches); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"Matches for the game with that ID does not exist.",
				rdb.Error.Error(),
			)
		} else {
			errorResponse(
				context,
				"An error occured getting the information, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	matchesResponses := make([]ResponseAdvancedMatch, 0)
	for _, match := range matches {
		matchesResponses = append(matchesResponses, match.AdvancedResponse())
	}

	context.JSON(
		200,
		gin.H{
			"data":  matchesResponses,
			"error": false,
		},
	)
}

func getUserMatches(context *gin.Context) {
	aid, err := strconv.ParseUint(context.Param("aid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing user ID.",
			err.Error(),
		)
		return
	}

	var usersTeamMember []TeamMember

	// if rdb := db.Preload("LosingTeam").Preload("LosingTeam.TeamMembers").Preload("LosingTeam.TeamMembers")
	// .Preload("LosingTeam.TeamMembers.Accounts").Preload("LosingTeam.TeamMembers.GameAccounts")
	// .Preload("WinningTeam")

	if rdb := db.Preload("Match").
		Preload("Match.LosingTeam.TeamMembers.Accounts").Preload("Match.LosingTeam.TeamMembers.GameAccounts").
		Preload("Match.WinningTeam.TeamMembers.Accounts").Preload("Match.WinningTeam.TeamMembers.GameAccounts").
		Preload("Match.Game").Where(TeamMember{AccountID: aid}).Find(&usersTeamMember); rdb.Error != nil {

		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"Matches could not be found for the user with that ID.",
				rdb.Error.Error(),
			)
		} else {
			errorResponse(
				context,
				"An error occured getting the information, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	matchesResponses := make([]ResponseAdvancedMatch, 0)
	for _, userTeamMember := range usersTeamMember {
		matchesResponses = append(matchesResponses, userTeamMember.Match.AdvancedResponse())
	}

	context.JSON(
		200,
		gin.H{
			"data": gin.H{
				"accountID": aid,
				"matches":   matchesResponses,
			},
			"error": false,
		},
	)
}

func getGameMatchesForUser(context *gin.Context) {
	aid, err := strconv.ParseUint(context.Param("aid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing ID.",
			err.Error(),
		)
		return
	}

	gid, err := strconv.ParseUint(context.Param("gid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing ID.",
			err.Error(),
		)
		return
	}

	var usersTeamMember []TeamMember

	// if rdb := db.Preload("LosingTeam").Preload("LosingTeam.TeamMembers").Preload("LosingTeam.TeamMembers")
	// .Preload("LosingTeam.TeamMembers.Accounts").Preload("LosingTeam.TeamMembers.GameAccounts")
	// .Preload("WinningTeam")

	if rdb := db.Preload("Match").
		Preload("Match.LosingTeam.TeamMembers.Accounts").Preload("Match.LosingTeam.TeamMembers.GameAccounts").
		Preload("Match.WinningTeam.TeamMembers.Accounts").Preload("Match.WinningTeam.TeamMembers.GameAccounts").
		Preload("Match.Game").Where(TeamMember{AccountID: aid, GameID: gid}).Find(&usersTeamMember); rdb.Error != nil {

		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"Matches could not be found for the user with that ID.",
				rdb.Error.Error(),
			)
		} else {
			errorResponse(
				context,
				"An error occured getting the information, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	matchesResponses := make([]ResponseAdvancedMatch, 0)
	for _, userTeamMember := range usersTeamMember {
		matchesResponses = append(matchesResponses, userTeamMember.Match.AdvancedResponse())
	}

	context.JSON(
		200,
		gin.H{
			"data": gin.H{
				"accountID": aid,
				"gameID":    gid,
				"matches":   matchesResponses,
			},
			"error": false,
		},
	)
}

func getTeam(context *gin.Context) {
	tid, err := strconv.ParseUint(context.Param("tid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing ID.",
			err.Error(),
		)
		return
	}

	var team Team

	if rdb := db.Preload("TeamMembers").Preload("TeamMembers.Accounts").Preload("TeamMembers.GameAccounts").Find(&team, tid); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"A team with that ID does not exist.",
				rdb.Error.Error(),
			)
		} else {
			errorResponse(
				context,
				"An error occured getting the information, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	context.JSON(
		200,
		gin.H{
			"data":  team.AdvancedResponse(),
			"error": false,
		},
	)
}

func getUserAccount(context *gin.Context) {
	aid, err := strconv.ParseUint(context.Param("aid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing ID.",
			err.Error(),
		)
		return
	}

	var user Account

	if rdb := db.Preload("GameAccounts").Where(Account{AccountID: aid, Enabled: true}).First(&user); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"A user with that ID does not exist.",
				rdb.Error.Error(),
			)
		} else {
			errorResponse(
				context,
				"An error occured getting the information, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	context.JSON(
		200,
		gin.H{
			"data":  user.AdvancedResponse(),
			"error": false,
		},
	)
}

func getUserAccountForGame(context *gin.Context) {
	aid, err := strconv.ParseUint(context.Param("aid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing ID.",
			err.Error(),
		)
		return
	}

	gid, err := strconv.ParseUint(context.Param("gid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing ID.",
			err.Error(),
		)
		return
	}

	var userGameAccount GameAccount

	if rdb := db.Preload("Account").Preload("Game").Where(GameAccount{AccountID: aid, GameID: gid, Enabled: true}).First(&userGameAccount); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"No results were found for the provided user and game.",
				rdb.Error.Error(),
			)
		} else {
			errorResponse(
				context,
				"An error occured getting the information, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	context.JSON(
		200,
		gin.H{
			"data":  userGameAccount.AdvancedResponse(),
			"error": false,
		},
	)
}

func getConfirmAccount(context *gin.Context) {
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
}

func postCreateAccount(context *gin.Context) {
	defer context.Request.Body.Close()

	pendingUser := PendingAccount{}

	decoder := json.NewDecoder(context.Request.Body)

	if err := decoder.Decode(&pendingUser); err != nil {
		errorResponse(
			context,
			"Error recieving registration information, please try again.",
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
	pendingUser.Password = hashPassword(pendingUser.Password)

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
}

func getRankingsForUserAccount(context *gin.Context) {
	aid, err := strconv.ParseUint(context.Param("aid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing ID.",
			err.Error(),
		)
		return
	}

	var user Account

	if rdb := db.Preload("GameAccounts").Preload("GameAccounts.Game").Where(Account{AccountID: aid, Enabled: true}).First(&user); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"A user with that ID does not exist.",
				rdb.Error.Error(),
			)
		} else {
			errorResponse(
				context,
				"An error occured getting the information, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	if len(user.GameAccounts) == 0 {
		context.JSON(
			200,
			gin.H{
				"data":  make([]interface{}, 0),
				"error": false,
			},
		)
		return
	}

	userGameRankings := make([]ResponseGameRanking, 0)

	for _, userGameAccount := range user.GameAccounts {
		if !userGameAccount.Enabled {
			continue
		}

		gameAccounts := make([]GameAccount, 0)

		if rdb := db.Where(GameAccount{GameID: userGameAccount.GameID, Enabled: true}).Order("score desc").Find(&gameAccounts); rdb.Error != nil {
			if rdb.RecordNotFound() {
				userGameRankings = append(userGameRankings, ResponseGameRanking{userGameAccount.Game.Response(), 0, 0})
				continue
			} else {
				errorResponse(
					context,
					"An error occured getting the information, try again.",
					rdb.Error.Error(),
				)
			}
			return
		}

		for index, gameAccount := range gameAccounts {
			if gameAccount.AccountID == aid {
				userGameRankings = append(userGameRankings, ResponseGameRanking{
					userGameAccount.Game.Response(),
					uint64(index + 1),
					uint64(len(gameAccounts)),
				})
			}
		}
	}

	context.JSON(
		200,
		gin.H{
			"data":  userGameRankings,
			"error": false,
		},
	)
}

func getRankingsForUserAccountInGame(context *gin.Context) {
	aid, err := strconv.ParseUint(context.Param("aid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing ID.",
			err.Error(),
		)
		return
	}

	gid, err := strconv.ParseUint(context.Param("gid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing ID.",
			err.Error(),
		)
		return
	}

	gameAccounts := make([]GameAccount, 0)

	if rdb := db.Preload("Game").Where(GameAccount{GameID: gid, Enabled: true}).Order("score desc").Find(&gameAccounts); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"Could not find the rankings for that game.",
				"Game has no rankings.",
			)
		} else {
			errorResponse(
				context,
				"An error occured getting the information, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	for index, gameAccount := range gameAccounts {
		if gameAccount.AccountID == aid {
			context.JSON(
				200,
				gin.H{
					"data": ResponseGameRanking{
						gameAccount.Game.Response(),
						uint64(index + 1),
						uint64(len(gameAccounts)),
					},
					"error": false,
				},
			)
			return
		}
	}

	errorResponse(
		context,
		"Could not find the ranking for the user in that game.",
		"User not in game rankings.",
	)
}

func deleteUser(context *gin.Context, apiAccount APIAccount) {
	// var userAccount Account

	// if rdb := db.First(&userAccount, apiAccount.AccountID); rdb.Error != nil {
	// 	if rdb.RecordNotFound() {
	// 		errorResponse(
	// 			context,
	// 			"Could not find an account with that ID to delete.",
	// 			"Account doesn't exist.",
	// 		)
	// 	} else {
	// 		errorResponse(
	// 			context,
	// 			"An error occured getting the information, try again.",
	// 			rdb.Error.Error(),
	// 		)
	// 	}
	// 	return
	// }

	// userAccount.Enabled = false

	// var userGameAccounts []GameAccount

	// if rdb := db.Where(GameAccount{AccountID: apiAccount.AccountID}).Find(&userGameAccounts); !rdb.RecordNotFound() && rdb.Error != nil {
	// 	errorResponse(
	// 		context,
	// 		"An error occured getting the game accounts information for deletion, try again.",
	// 		rdb.Error.Error(),
	// 	)
	// 	return
	// }

	db.Model(&Account{}).Where(Account{AccountID: apiAccount.AccountID}).Update("enabled", false)

	db.Where(Session{AccountID: apiAccount.AccountID}).Delete(Session{})

	db.Model(&GameAccount{}).Where(GameAccount{AccountID: apiAccount.AccountID}).Update("enabled", false)

	// db.Save(userAccount)

	context.JSON(
		200,
		gin.H{
			"data":  true,
			"error": false,
		},
	)
}

func deleteGame(context *gin.Context, apiAccount APIAccount) {
	if apiAccount.GlobalPermissions < permGameDeletion {
		errorResponse(
			context,
			"Insufficient permissions to perform operation.",
			"User does not have permissions.",
		)
		return
	}

	gid, err := strconv.ParseUint(context.Param("gid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing ID.",
			err.Error(),
		)
		return
	}

	var game Game

	if rdb := db.Where(Game{GameID: gid}).First(&game); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"Could not find a game with the specified game ID to delete.",
				"Game doesn't exist.",
			)
		} else {
			errorResponse(
				context,
				"An error occured getting the information, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	db.Delete(game)

	context.JSON(
		200,
		gin.H{
			"data":  true,
			"error": false,
		},
	)
}

func deleteMatch(context *gin.Context, apiAccount APIAccount) {
	mid, err := strconv.ParseUint(context.Param("mid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing ID.",
			err.Error(),
		)
		return
	}

	var match Match

	if rdb := db.Where(Match{AccountID: apiAccount.AccountID, MatchID: mid}).First(&match); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"Could not find a match with the specified match and user ID to delete.",
				"Game account doesn't exist.",
			)
		} else {
			errorResponse(
				context,
				"An error occured getting the information, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	db.Delete(match)

	context.JSON(
		200,
		gin.H{
			"data":  true,
			"error": false,
		},
	)
}

func postCreateGame(context *gin.Context, apiAccount APIAccount) {
	defer context.Request.Body.Close()

	if apiAccount.GlobalPermissions < permGameCreation {
		errorResponse(
			context,
			"Insufficient permissions to perform operation.",
			"User does not have permissions.",
		)
		return
	}
	var game Game

	decoder := json.NewDecoder(context.Request.Body)

	if err := decoder.Decode(&game); err != nil {
		errorResponse(
			context,
			"Error recieving game information, please try again.",
			err.Error(),
		)
		return
	}

	if game.GameName == "" {
		errorResponse(
			context,
			"Name of the game cannot be blank.",
			"Blank name field.",
		)
		return
	}

	if len(game.GameName) > 100 {
		errorResponse(
			context,
			"Name of the game is too long (over 100 characters).",
			"Blank name field.",
		)
		return
	}

	if len(game.Colors) > 50 { //IMPROVE WITH REGEX IN FUTURE
		context.JSON(
			400,
			gin.H{
				"error": true,
				"msg":   "Color scheme is too long.",
			},
		)
		return
	}

	newGame := Game{
		Colors:   game.Colors,
		GameName: game.GameName,
	}

	gameResponse, err := CreateGame(db, &newGame)

	if err != nil {
		if terr, ok := err.(*NotUniqueGameError); ok {
			errorResponse(
				context,
				terr.Error(),
				err.Error(),
			)
		} else {
			errorResponse(
				context,
				"Error creating game, try again.",
				err.Error(),
			)
		}
		return
	}

	newGameAccount := GameAccount{
		AccountID:       apiAccount.AccountID,
		GameID:          newGame.GameID,
		GamePermissions: 5,
	}

	if err := db.Create(&newGameAccount).Error; err != nil {
		errorResponse(
			context,
			"An error occured creating the account, try again.",
			err.Error(),
		)
		return
	}

	context.JSON(
		200,
		gin.H{
			"data":  gameResponse,
			"error": false,
		},
	)
}

func postCreateMatch(context *gin.Context, apiAccount APIAccount) {
	defer context.Request.Body.Close()

	var match RequestMatch

	decoder := json.NewDecoder(context.Request.Body)

	if err := decoder.Decode(&match); err != nil {
		errorResponse(
			context,
			"Error recieving match information, please try again.",
			err.Error(),
		)
		return
	}

	if len(match.Losers) == 0 {
		errorResponse(
			context,
			"The list of losers cannot be blank.",
			"Blank losers list.",
		)
		return
	}

	if len(match.Winners) == 0 {
		errorResponse(
			context,
			"The list of winners cannot be blank.",
			"Blank winners list.",
		)
		return
	}

	for _, loserID := range match.Losers {
		for _, winnerID := range match.Winners {
			if loserID == winnerID {
				errorResponse(
					context,
					"A winner cannot also be a loser.",
					"ID in both losers and winners list.",
				)
				return
			}
		}
	}

	tdb := db.Where(GameAccount{GameID: match.GameID})
	allIDs := append(match.Losers, match.Winners...)
	gameAccounts := make([]GameAccount, 0)

	for _, accountID := range allIDs {
		tdb.Or(GameAccount{AccountID: accountID, GameID: match.GameID})
	}

	if rdb := tdb.Find(&gameAccounts); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"Could not find the users with accounts for that game, make sure everyone is in the rankings for that game and that game exists.",
				"An account isn't in game rankings.",
			)
		} else {
			errorResponse(
				context,
				"An error occured getting the information, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	if len(allIDs) != len(gameAccounts) {
		errorResponse(
			context,
			"Could not find all the users with accounts for that game, make sure everyone is in the rankings for that game.",
			"Missing game accounts for match.",
		)
		return
	}

	newMatch := Match{
		AccountID: apiAccount.AccountID,
		GameID:    match.GameID,
		MatchTime: match.MatchTime,
	}

	matchResponse, err := CreateMatch(db, &newMatch, &match.Losers, &match.Winners)

	if err != nil {
		errorResponse(
			context,
			"Error creating match, try again.",
			err.Error(),
		)
		return
	}

	context.JSON(
		200,
		gin.H{
			"data":  matchResponse,
			"error": false,
		},
	)
}

func postUpdateUser(context *gin.Context, apiAccount APIAccount) {
	//TODO
}

func putCreateGameAccount(context *gin.Context, apiAccount APIAccount) {
	gid, err := strconv.ParseUint(context.Param("gid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing ID.",
			err.Error(),
		)
		return
	}

	var userGameAccount GameAccount

	if rdb := db.Where(GameAccount{AccountID: apiAccount.AccountID, GameID: gid}).First(&userGameAccount); !rdb.RecordNotFound() || (rdb.Error != nil && !rdb.RecordNotFound()) {
		if !rdb.RecordNotFound() {
			errorResponse(
				context,
				"Game account for the user with that ID already exists.",
				"Game account already exist.",
			)
		} else {
			errorResponse(
				context,
				"An error occured creating the account, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	newGameAccount := GameAccount{
		AccountID: apiAccount.AccountID,
		GameID:    gid,
	}

	if err := db.Create(&newGameAccount).Error; err != nil {
		errorResponse(
			context,
			"An error occured creating the account, try again.",
			err.Error(),
		)
		return
	}

	context.JSON(
		200,
		gin.H{
			"data":  true,
			"error": false,
		},
	)
}

func putDisableGameAccount(context *gin.Context, apiAccount APIAccount) {
	aid, err := strconv.ParseUint(context.Param("aid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing ID.",
			err.Error(),
		)
		return
	}

	// Assure that the user deleting the specified account is the same user who is authenticated
	if aid != apiAccount.AccountID {
		errorResponse(
			context,
			"You do not have the permisions for performing that action.",
			err.Error(),
		)
		return
	}

	gid, err := strconv.ParseUint(context.Param("gid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing ID.",
			err.Error(),
		)
		return
	}

	var userGameAccount GameAccount

	if rdb := db.Where(GameAccount{AccountID: apiAccount.AccountID, GameID: gid}).First(&userGameAccount); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"Could not find a game account with that ID to delete.",
				"Game account doesn't exist.",
			)
		} else {
			errorResponse(
				context,
				"An error occured deleting the account, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	userGameAccount.Enabled = false

	if err := db.Save(userGameAccount).Error; err != nil {
		errorResponse(
			context,
			"An error occured disabling the account, try again.",
			err.Error(),
		)
		return
	}

	context.JSON(
		200,
		gin.H{
			"data":  true,
			"error": false,
		},
	)
}

func putEnableGameAccount(context *gin.Context, apiAccount APIAccount) {
	aid, err := strconv.ParseUint(context.Param("aid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing ID.",
			err.Error(),
		)
		return
	}

	// Assure that the user deleting the specified account is the same user who is authenticated
	if aid != apiAccount.AccountID {
		errorResponse(
			context,
			"You do not have the permisions for performing that action.",
			err.Error(),
		)
		return
	}

	gid, err := strconv.ParseUint(context.Param("gid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing ID.",
			err.Error(),
		)
		return
	}

	var userGameAccount GameAccount

	if rdb := db.Where(GameAccount{AccountID: apiAccount.AccountID, GameID: gid}).First(&userGameAccount); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"Could not find a game account with that ID to delete.",
				"Game account doesn't exist.",
			)
		} else {
			errorResponse(
				context,
				"An error occured re-enabling the account, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	userGameAccount.Enabled = true

	if err := db.Save(userGameAccount).Error; err != nil {
		errorResponse(
			context,
			"An error occured re-enabling the account, try again.",
			err.Error(),
		)
		return
	}

	context.JSON(
		200,
		gin.H{
			"data":  true,
			"error": false,
		},
	)
}

func putMatchConfirm(context *gin.Context, apiAccount APIAccount) {
	// Implement score calculation on loser confirm >= 50% and (winner confirm == 100% or (len(winners) > 1 and winner confirm >= 2))
	// True indicates update,
	// False means match has already been confirmed
	mid, err := strconv.ParseUint(context.Param("mid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing match ID.",
			err.Error(),
		)
		return
	}

	var userTeamMember TeamMember

	if rdb := db.Where(TeamMember{AccountID: apiAccount.AccountID, MatchID: mid}).First(&userTeamMember); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"Could not find a game account for that user associated with that match.",
				"Game account doesn't exist for match.",
			)
		} else {
			errorResponse(
				context,
				"An error occured confirming the match, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	var match Match

	if rdb := db.
		Preload("LosingTeam.TeamMembers.Accounts").Preload("LosingTeam.TeamMembers.GameAccounts").
		Preload("WinningTeam.TeamMembers.Accounts").Preload("WinningTeam.TeamMembers.GameAccounts").
		Preload("Game").Find(&match, mid); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"Could not find a match with that ID to confirm.",
				"Match with ID doesn't exist.",
			)
		} else {
			errorResponse(
				context,
				"An error occured confirming the match, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	if match.Confirmed {
		errorResponse(
			context,
			"Match has already been confirmed.",
			"Match already confirmed",
		)
		return
	}

	userTeamMember.Confirmed = true
	userTeamMember.Deny = false

	if err := db.Save(&userTeamMember).Error; err != nil {
		errorResponse(
			context,
			"An error occured confirming the match, try again.",
			err.Error(),
		)
		return
	}

	var loserConfirms, loserDenies, winnerConfirms, winnerDenies int
	for _, loser := range match.LosingTeam.TeamMembers {
		if loser.Confirmed {
			if loser.Deny {
				loserDenies++
			} else {
				loserConfirms++
			}
		}
	}

	for _, winner := range match.LosingTeam.TeamMembers {
		if winner.Confirmed {
			if winner.Deny {
				winnerDenies++
			} else {
				winnerConfirms++
			}
		}
	}

	if userTeamMember.Winner {
		winnerConfirms++
	} else {
		loserConfirms++
	}

	// Implement score calculation on loser confirm >= 50% and (winner confirm == 100% or (len(winners) > 1 and winner confirm >= 2))
	if loserConfirms >= int(len(match.LosingTeam.TeamMembers)/2) && (winnerConfirms == len(match.WinningTeam.TeamMembers) || winnerConfirms >= 2) {
		//Confirm Match and update scores here
	}

	context.JSON(
		200,
		gin.H{
			"data":  true,
			"error": false,
		},
	)
}

func putMatchDeny(context *gin.Context, apiAccount APIAccount) {
	// Take down match if loser deny >= 50% or winner deny > 0%
	// True indicates update,
	// False means match has already been confirmed
	mid, err := strconv.ParseUint(context.Param("mid"), 10, 32)

	if err != nil {
		errorResponse(
			context,
			"Error parsing match ID.",
			err.Error(),
		)
		return
	}

	var userTeamMember TeamMember

	if rdb := db.Where(TeamMember{AccountID: apiAccount.AccountID, MatchID: mid}).First(&userTeamMember); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"Could not find a game account for that user associated with that match.",
				"Game account doesn't exist for match.",
			)
		} else {
			errorResponse(
				context,
				"An error occured confirming the match, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	var match Match

	if rdb := db.
		Preload("LosingTeam.TeamMembers.Accounts").Preload("LosingTeam.TeamMembers.GameAccounts").
		Preload("WinningTeam.TeamMembers.Accounts").Preload("WinningTeam.TeamMembers.GameAccounts").
		Preload("Game").Find(&match, mid); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"Could not find a match with that ID to confirm.",
				"Match with ID doesn't exist.",
			)
		} else {
			errorResponse(
				context,
				"An error occured confirming the match, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	if match.Confirmed {
		errorResponse(
			context,
			"Match has already been confirmed.",
			"Match already confirmed",
		)
		return
	}

	userTeamMember.Confirmed = true
	userTeamMember.Deny = true

	if err := db.Save(&userTeamMember).Error; err != nil {
		errorResponse(
			context,
			"An error occured confirming the match, try again.",
			err.Error(),
		)
		return
	}

	var loserConfirms, loserDenies, winnerConfirms, winnerDenies int
	for _, loser := range match.LosingTeam.TeamMembers {
		if loser.Confirmed {
			if loser.Deny {
				loserDenies++
			} else {
				loserConfirms++
			}
		}
	}

	for _, winner := range match.LosingTeam.TeamMembers {
		if winner.Confirmed {
			if winner.Deny {
				winnerDenies++
			} else {
				winnerConfirms++
			}
		}
	}

	if userTeamMember.Winner {
		winnerDenies++
	} else {
		loserDenies++
	}

	// Take down match if loser deny >= 50% or winner deny > 0%
	if loserDenies >= int(len(match.LosingTeam.TeamMembers)/2) || winnerDenies > 0 {
		if rdb := db.Where(Match{MatchID: mid}).Delete(Match{}); rdb.Error != nil {
			if rdb.RecordNotFound() {
				errorResponse(
					context,
					"Could not find a Match with that ID to delete.",
					"Match with MID doesn't exist.",
				)
			} else {
				errorResponse(
					context,
					"An error occured deleting that match, try again.",
					rdb.Error.Error(),
				)
			}
			return
		}

		if rdb := db.Where(Team{MatchID: mid}).Delete(Team{}); rdb.Error != nil {
			if rdb.RecordNotFound() {
				errorResponse(
					context,
					"Could not find a team with that ID to delete.",
					"Team with MID doesn't exist.",
				)
			} else {
				errorResponse(
					context,
					"An error occured deleting that team, try again.",
					rdb.Error.Error(),
				)
			}
			return
		}

		if rdb := db.Where(TeamMember{MatchID: mid}).Delete(TeamMember{}); rdb.Error != nil {
			if rdb.RecordNotFound() {
				errorResponse(
					context,
					"Could not any team members with that ID to delete.",
					"TeamMembers with MID doesn't exist.",
				)
			} else {
				errorResponse(
					context,
					"An error occured deleting the TeamMembers, try again.",
					rdb.Error.Error(),
				)
			}
			return
		}

		context.JSON(
			200,
			gin.H{
				"data":  false,
				"error": false,
			},
		)
		return
	}

	context.JSON(
		200,
		gin.H{
			"data":  true,
			"error": false,
		},
	)
}
