package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

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

	var gameAccounts []GameAccount

	if rdb := db.Preload("Accounts").Where(GameAccount{GameID: gid, Enabled: true}).Order("score desc").Find(&gameAccounts); rdb.Error != nil && !rdb.RecordNotFound() {
		// if rdb.RecordNotFound() {
		// 	errorResponse(
		// 		context,
		// 		"None for this game.", // Redo when return behavior is define, this should eventually return an empty list
		// 		rdb.Error.Error(),
		// 	)
		// } else {
		errorResponse(
			context,
			"An error occured getting the information, try again.",
			rdb.Error.Error(),
		)

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

	if rdb := db.Preload("GameAccounts").First(&user, aid); rdb.Error != nil {
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
	// Like regi profile but with information pertaining specifically to game
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

	if rdb := db.Preload("GameAccounts").Preload("GameAccounts.Game").First(&user, aid); rdb.Error != nil {
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

	if rdb := db.Where(GameAccount{GameID: gid}).Order("score desc").Find(&gameAccounts); rdb.Error != nil {
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

	var userAccount Account

	if rdb := db.First(&userAccount, apiAccount.AccountID); rdb.Error != nil {
		if rdb.RecordNotFound() {
			errorResponse(
				context,
				"Could not find an account with that ID to delete.",
				"Account doesn't exist.",
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

	db.Delete(userAccount)
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
}

func postCreateGame(context *gin.Context, apiAccount APIAccount) {
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

	context.JSON(
		200,
		gin.H{
			"data":  gameResponse,
			"error": false,
		},
	)
}

func postCreateMatch(context *gin.Context, apiAccount APIAccount) {

}

func postUpdateUser(context *gin.Context, apiAccount APIAccount) {

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
				"An error occured getting the information, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	userGameAccount.Enabled = false

	db.Save(userGameAccount)
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
				"An error occured getting the information, try again.",
				rdb.Error.Error(),
			)
		}
		return
	}

	userGameAccount.Enabled = true

	db.Save(userGameAccount)
}

func putMatchConfirm(context *gin.Context, apiAccount APIAccount) {
	// Implement score calculation on all confirms
}

func putMatchDeny(context *gin.Context, apiAccount APIAccount) {
	// Take down match if deny
}
