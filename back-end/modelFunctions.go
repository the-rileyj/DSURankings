package main

import (
	"github.com/jinzhu/gorm"
)

func CreateAccount(db *gorm.DB, user *Account) (ResponseAccount, error) {
	err := db.Create(user).Error
	if err != nil {
		return (&Account{}).ResponseAccount(), err
	}
	return user.ResponseAccount(), nil
}

func CreatePendingAccount(db *gorm.DB, pendingUser *PendingAccount) (ResponsePendingAccount, error) {
	var existingUser Account

	if rdb := db.Where(Account{Email: pendingUser.Email}).Or(Account{UserName: pendingUser.UserName}).First(&existingUser); !rdb.RecordNotFound() && rdb.Error != nil {
		return (&PendingAccount{}).ResponseAccount(), rdb.Error
	}

	if existingUser.Email == pendingUser.Email || existingUser.UserName == pendingUser.UserName {
		return (&PendingAccount{}).ResponseAccount(), NewNotUniqueUserError()
	}

	err := db.Create(pendingUser).Error

	if err != nil {
		return (&PendingAccount{}).ResponseAccount(), err
	}

	return pendingUser.ResponseAccount(), nil
}
