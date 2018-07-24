package main

import (
	"time"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Account struct {
	gorm.Model
	AccountID         uint          `gorm:"primary_key"`
	Email             string        `gorm:"unique;not null"`
	GameAccounts      []GameAccount `gorm:"foreignkey:AccountID;association_foreignkey:AccountID"`
	GlobalPermissions uint          `gorm:"default:0"`
	UserName          string        `gorm:"unique;not null"`
	Password          string
}

type Game struct {
	gorm.Model
	Colors       string        //This may change
	GameAccounts []GameAccount `gorm:"foreignkey:GameID;association_foreignkey:GameID"`
	GameID       uint          `gorm:"primary_key"`
	GameName     string        `gorm:"unique"`
	Matches      []Match       `gorm:"foreignkey:GameID;association_foreignkey:GameID"`
}

type GameAccount struct {
	gorm.Model
	Account         Account      `gorm:"foreignkey:AccountID;association_foreignkey:AccountID"`
	AccountID       uint         `gorm:"unique_index:idx_game"`
	AccountTeams    []TeamMember `gorm:"foreignkey:AccountID;association_foreignkey:AccountID"`
	Games           Game         `gorm:"foreignkey:GameID;association_foreignkey:GameID"`
	GameID          uint         `gorm:"unique_index:idx_game"`
	GamePermissions uint         `gorm:"default:0"`
	Score           uint
}

type Match struct {
	gorm.Model
	Game          Game `gorm:"foreignkey:GameID;association_foreignkey:GameID"`
	LosingTeam    Team `gorm:"foreignkey:TeamID;association_foreignkey:LosingTeamID"`
	MatchID       uint `gorm:"primary_key"`
	WinningTeam   Team `gorm:"foreignkey:TeamID;association_foreignkey:WinningTeamID"`
	GameID        uint
	LosingTeamID  uint
	MatchTime     time.Time
	WinningTeamID uint
}

// Need logic to check that Email and Username don't exist in normal 'accounts' table
type PendingAccount struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	UserName string `gorm:"unique;not null"`
	UUID     string `gorm:"unique;not null"`
	Password string
}

type Team struct {
	gorm.Model
	TeamID      uint         `gorm:"primary_key"`
	TeamMembers []TeamMember `gorm:"foreignkey:TeamID;association_foreignkey:TeamID"`
}

// Note TeamMembers to indcate that there is more than one person on the team
type TeamMember struct {
	gorm.Model
	AccountID   uint
	GameID      uint
	Team        Team `gorm:"foreignkey:TeamID;association_foreignkey:TeamID"`
	TeamID      uint `gorm:"unique"`
	TeamMembers bool
}
