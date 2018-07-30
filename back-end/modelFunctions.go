package main

import (
	"github.com/jinzhu/gorm"
)

func CreateAccount(db *gorm.DB, user *Account) (ResponseBasicAccount, error) {
	err := db.Create(user).Error
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

	err := db.Create(pendingUser).Error

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

	err := db.Create(game).Error

	if err != nil {
		return ResponseGame{}, err
	}

	return game.Response(), nil
}
