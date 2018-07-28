package main

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ResponseAccount struct {
	AccountID    uint   `json:"accountID"`
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	GameAccounts []GameAccount
	LastName     string `json:"lastName"`
	UserName     string `json:"userName"`
}

type ResponsePendingAccount struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserName  string `json:"userName"`
}

type Account struct {
	AccountID         uint          `gorm:"primary_key"`
	Email             string        `gorm:"unique;not null"`
	FirstName         string        `gorm:"not null"`
	GameAccounts      []GameAccount `gorm:"foreignkey:AccountID;association_foreignkey:AccountID"`
	GlobalPermissions uint          `gorm:"default:0"`
	LastName          string        `gorm:"not null"`
	UserName          string        `gorm:"unique;not null"`
	Password          string
}

type Game struct {
	Colors       string        //This may change
	GameAccounts []GameAccount `gorm:"foreignkey:GameID;association_foreignkey:GameID"`
	GameID       uint          `gorm:"primary_key"`
	GameName     string        `gorm:"unique"`
	Matches      []Match       `gorm:"foreignkey:GameID;association_foreignkey:GameID"`
}

type GameAccount struct {
	Account         Account      `gorm:"foreignkey:AccountID;association_foreignkey:AccountID"`
	AccountID       uint         `gorm:"unique_index:idx_game"`
	AccountTeams    []TeamMember `gorm:"foreignkey:AccountID;association_foreignkey:AccountID"`
	Games           Game         `gorm:"foreignkey:GameID;association_foreignkey:GameID"`
	GameID          uint         `gorm:"unique_index:idx_game"`
	GamePermissions uint         `gorm:"default:0"`
	Score           uint
}

type Match struct {
	Game          Game `gorm:"foreignkey:GameID;association_foreignkey:GameID"`
	LosingTeam    Team `gorm:"foreignkey:TeamID;association_foreignkey:LosingTeamID"`
	MatchID       uint `gorm:"primary_key"`
	WinningTeam   Team `gorm:"foreignkey:TeamID;association_foreignkey:WinningTeamID"`
	GameID        uint
	LosingTeamID  uint
	MatchTime     time.Time
	WinningTeamID uint
}

type PendingAccount struct {
	Email     string `gorm:"unique;not null"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	UserName  string `gorm:"unique;not null"`
	UUID      string `gorm:"primary_key"`
	Password  string
}

type Team struct {
	TeamID      uint         `gorm:"primary_key"`
	TeamMembers []TeamMember `gorm:"foreignkey:TeamID;association_foreignkey:TeamID"`
}

// Note TeamMembers to indcate that there is more than one person on the team
type TeamMember struct {
	AccountID   uint
	GameID      uint
	Team        Team `gorm:"foreignkey:TeamID;association_foreignkey:TeamID"`
	TeamID      uint
	TeamMembers bool
}
