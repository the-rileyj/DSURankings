package main

import (
	"github.com/jinzhu/gorm"
)

func CreateAccount(db *gorm.DB, user *Account) (ResponseBasicAccount, error) {
	err := db.Create(&user).Error
	if err != nil {
		return ResponseBasicAccount{}, err
	}
	return user.BasicResponse(), nil
}

func CreatePendingAccount(db *gorm.DB, pendingUser *PendingAccount) (ResponsePendingAccount, error) {
	var existingUser Account

	if rdb := db.Where(Account{Email: pendingUser.Email}).Or(Account{UserName: pendingUser.UserName}).First(&existingUser); !rdb.RecordNotFound() && rdb.Error != nil {
		return ResponsePendingAccount{}, rdb.Error
	}

	if existingUser.Email == pendingUser.Email || existingUser.UserName == pendingUser.UserName {
		return ResponsePendingAccount{}, NewNotUniqueUserError()
	}

	err := db.Create(&pendingUser).Error

	if err != nil {
		return ResponsePendingAccount{}, err
	}

	return pendingUser.Response(), nil
}

func CreateGame(db *gorm.DB, game *Game) (ResponseGame, error) {
	var existingGame Game

	if rdb := db.Where(Game{GameName: game.GameName}).First(&existingGame); !rdb.RecordNotFound() && rdb.Error != nil {
		return ResponseGame{}, rdb.Error
	}

	if existingGame.GameName == game.GameName {
		return ResponseGame{}, NewNotUniqueGameError()
	}

	err := db.Create(&game).Error

	if err != nil {
		return ResponseGame{}, err
	}

	return game.Response(), nil
}

func CreateMatch(db *gorm.DB, match *Match, losers, winners *[]uint64) (ResponseAdvancedMatch, error) {
	// Remake with Rollback later http://doc.gorm.io/advanced.html#transactions
	err := db.Create(&match).Error
	if err != nil {
		return ResponseAdvancedMatch{}, err
	}

	losingTeam := Team{MatchID: match.MatchID, Single: len(*losers) == 1}
	winningTeam := Team{MatchID: match.MatchID, Single: len(*winners) == 1}

	err = db.Create(&losingTeam).Error
	if err != nil {
		db.Delete(&match)
		return ResponseAdvancedMatch{}, err
	}

	err = db.Create(&winningTeam).Error
	if err != nil {
		db.Delete(&match)
		db.Delete(&losingTeam)
		return ResponseAdvancedMatch{}, err
	}

	match.LosingTeamID = losingTeam.TeamID
	match.WinningTeamID = winningTeam.TeamID

	db.Update(&match)

	teamMembers := make([]TeamMember, 0)

	if len(*losers) == 1 {
		newTeamMember := TeamMember{AccountID: (*losers)[0], GameID: match.GameID, MatchID: match.MatchID, TeamID: match.LosingTeamID, TeamMembers: false, Winner: false}
		teamMembers = append(teamMembers, newTeamMember)
		err = db.Create(&newTeamMember).Error

		if err != nil {
			db.Delete(&match)
			db.Delete(&losingTeam)
			db.Delete(&winningTeam)
			return ResponseAdvancedMatch{}, err
		}
	} else {
		for _, loserID := range *losers {
			newTeamMember := TeamMember{AccountID: loserID, GameID: match.GameID, MatchID: match.MatchID, TeamID: match.LosingTeamID, Winner: false}
			teamMembers = append(teamMembers, newTeamMember)
			err = db.Create(&newTeamMember).Error

			if err != nil {
				db.Delete(&match)
				db.Delete(&losingTeam)
				db.Delete(&winningTeam)
				for _, deleteTeamMember := range teamMembers {
					db.Delete(&deleteTeamMember)
				}
				return ResponseAdvancedMatch{}, err
			}
		}
	}

	if len(*winners) == 1 {
		newTeamMember := TeamMember{AccountID: (*winners)[0], GameID: match.GameID, MatchID: match.MatchID, TeamID: match.LosingTeamID, TeamMembers: false, Winner: true}
		teamMembers = append(teamMembers, newTeamMember)
		err = db.Create(&newTeamMember).Error

		if err != nil {
			db.Delete(&match)
			db.Delete(&losingTeam)
			db.Delete(&winningTeam)
			for _, deleteTeamMember := range teamMembers {
				db.Delete(&deleteTeamMember)
			}
			return ResponseAdvancedMatch{}, err
		}
	} else {
		for _, winnerID := range *winners {
			newTeamMember := TeamMember{AccountID: winnerID, GameID: match.GameID, MatchID: match.MatchID, TeamID: match.LosingTeamID, Winner: true}
			teamMembers = append(teamMembers, newTeamMember)
			err = db.Create(&newTeamMember).Error

			if err != nil {
				db.Delete(&match)
				db.Delete(&losingTeam)
				db.Delete(&winningTeam)
				for _, deleteTeamMember := range teamMembers {
					db.Delete(&deleteTeamMember)
				}
				return ResponseAdvancedMatch{}, err
			}
		}
	}

	if rdb := db.
		Preload("LosingTeam.TeamMembers.Accounts").Preload("LosingTeam.TeamMembers.GameAccounts").
		Preload("WinningTeam.TeamMembers.Accounts").Preload("WinningTeam.TeamMembers.GameAccounts").
		Preload("Game").Find(&match, match.MatchID); rdb.Error != nil {
		db.Delete(&match)
		db.Delete(&losingTeam)
		db.Delete(&winningTeam)
		for _, deleteTeamMember := range teamMembers {
			db.Delete(&deleteTeamMember)
		}
		return ResponseAdvancedMatch{}, rdb.Error
	}

	return match.AdvancedResponse(), nil
}
