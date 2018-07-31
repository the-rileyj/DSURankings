package main

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type APIAccount struct {
	AccountID         uint64 `json:"accountID"`
	Email             string `json:"email"`
	FirstName         string `json:"firstName"`
	GlobalPermissions uint64 `json:"globalPermissions"`
	LastName          string `json:"lastName"`
	UserName          string `json:"userName"`
}

type ApiAuther struct {
	Token string `json:"token"`
}

type Auther struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Session struct {
	Account   Account `gorm:"foreignkey:AccountID;association_foreignkey:AccountID"`
	UUID      string  `gorm:"primary_key"`
	AccountID uint64
	CreatedAt time.Time
}

type RequestMatch struct {
	Losers    []uint64 `json:"losers"`
	Winners   []uint64 `json:"winners"`
	GameID    uint64
	MatchTime time.Time
}

type ResponseAdvancedAccount struct {
	GameAccounts []GameAccount
	ResponseBasicAccount
}

type ResponseBasicAccount struct {
	AccountID uint64 `json:"accountID"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserName  string `json:"userName"`
}

type ResponseGame struct {
	Colors   string `json:"colors"`
	GameID   uint64 `json:"gameID"`
	GameName string `json:"gameName"`
}

type ResponseAdvancedGameAccount struct {
	ResponseBasicGameAccount `json:"gameAccount"`
	ResponseGame             `json:"game"`
}

type ResponseBasicGameAccount struct {
	Account ResponseBasicAccount `json:"account"`
	GameID  uint64               `json:"gameID"`
	Score   uint64               `json:"score"`
}

type ResponseAdvancedMatch struct {
	Account        ResponseBasicAccount `json:"initiator"`
	Game           ResponseGame         `json:"game"`
	LosingTeam     ResponseAdvancedTeam `json:"losingTeam"`
	WinningTeam    ResponseAdvancedTeam `json:"winningTeam"`
	MatchID        uint64               `json:"matchID"`
	Confirmed      bool                 `json:"confirmed"`
	GameID         uint64               `json:"gameID"`
	LosingTeamID   uint64               `json:"losingTeamID"`
	MatchStartTime time.Time            `json:"matchStartTime"`
	WinningTeamID  uint64               `json:"winningTeamID"`
}

type ResponsePendingAccount struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	UserName  string `json:"userName"`
}

type ResponseGameRanking struct {
	Game  ResponseGame `json:"game"`
	Rank  uint64       `json:"rank"`
	Total uint64       `json:"total"`
}

type ResponseAdvancedTeam struct {
	TeamMembers []ResponseAdvancedTeamMember `json:"teamMembers"`
	MatchID     uint64                       `json:"matchID"`
	TeamID      uint64                       `json:"teamID"`
}

type ResponseBasicTeam struct {
	TeamMembers []ResponseBasicTeamMember `json:"teamMembers"`
	MatchID     uint64                    `json:"matchID"`
	TeamID      uint64                    `json:"teamID"`
}

type ResponseAdvancedTeamMember struct {
	ResponseBasicAccount
	ResponseBasicGameAccount
	ResponseBasicTeamMember
}

type ResponseBasicTeamMember struct {
	AccountID   uint64 `json:"accountID"`
	Confirmed   bool   `json:"confirmed"`
	GameID      uint64 `json:"gameID"`
	MatchID     uint64 `json:"matchID"`
	TeamID      uint64 `json:"teamID"`
	TeamMembers bool   `json:"teamMembers"`
}

/* GLOBAL PERMISSIONS SCALE
===========================
0 - None at all
1 - Game Creation
2 - Modification of User Permisions When Account.GlobalPermissions <= 1 and assignable to range 0:1
3 - Modification of User Permisions When Account.GlobalPermissions <= 2 and assignable to range 0:2 and Game deletion
*/

type Account struct {
	GameAccounts      []GameAccount `gorm:"foreignkey:AccountID;association_foreignkey:AccountID"`
	TeamMembers       []TeamMember  `gorm:"foreignkey:AccountID;association_foreignkey:AccountID"`
	AccountID         uint64        `gorm:"primary_key"`
	Email             string        `gorm:"unique;not null"`
	FirstName         string        `gorm:"not null"`
	GlobalPermissions uint64        `gorm:"default:0"`
	LastName          string        `gorm:"not null"`
	UserName          string        `gorm:"unique;not null"`
	Enabled           bool          `gorm:"default:true"`
	Password          string
}

type Game struct {
	GameAccounts []GameAccount `gorm:"foreignkey:GameID;association_foreignkey:GameID"`
	Matches      []Match       `gorm:"foreignkey:GameID;association_foreignkey:GameID"`
	GameID       uint64        `gorm:"primary_key"`
	GameName     string        `gorm:"unique" json:"gameName"`
	AccountID    uint64        // The person who originally created the game
	Colors       string        ` json:"colors"` // This may change
}

type GameAccount struct {
	Account         Account      `gorm:"foreignkey:AccountID;association_foreignkey:AccountID"`
	AccountTeams    []TeamMember `gorm:"foreignkey:AccountID;association_foreignkey:AccountID"`
	Game            Game         `gorm:"foreignkey:GameID;association_foreignkey:GameID"`
	AccountID       uint64       `gorm:"unique_index:idx_game"`
	GameID          uint64       `gorm:"unique_index:idx_game"`
	GamePermissions uint64       `gorm:"default:0"`
	Score           uint64       `gorm:"default:300"`
	Enabled         bool         `gorm:"default:true"`
}

type Match struct {
	Account       Account `gorm:"foreignkey:AccountID;association_foreignkey:AccountID"` // The initiator of the match
	Game          Game    `gorm:"foreignkey:GameID;association_foreignkey:GameID"`
	LosingTeam    *Team   `gorm:"foreignkey:TeamID;association_foreignkey:LosingTeamID"`
	WinningTeam   *Team   `gorm:"foreignkey:TeamID;association_foreignkey:WinningTeamID"`
	Confirmed     bool    `gorm:"default:false"`
	MatchID       uint64  `gorm:"primary_key"`
	AccountID     uint64
	GameID        uint64
	LosingTeamID  uint64
	MatchTime     time.Time
	WinningTeamID uint64
	// Confirmations []MatchConfirmation `gorm:"foreignkey:MatchID;association_foreignkey:MatchID"`
}

// type MatchConfirmation struct {
// Match          Match `gorm:"foreignkey:MatchID;association_foreignkey:MatchID"`
// 	LosingTeam    Team `gorm:"foreignkey:TeamID;association_foreignkey:LosingTeamID"`
// 	WinningTeam   Team `gorm:"foreignkey:TeamID;association_foreignkey:WinningTeamID"`
// 	MatchID       uint64
// 	Confirm       bool
// 	GameID        uint
// 	LosingTeamID  uint
// 	MatchTime     time.Time
// 	WinningTeamID uint
// }

type PendingAccount struct {
	TeamMembers []TeamMember `gorm:"foreignkey:MatchID;association_foreignkey:MatchID"`
	Email       string       `gorm:"unique;not null"`
	FirstName   string       `gorm:"not null"`
	LastName    string       `gorm:"not null"`
	UserName    string       `gorm:"unique;not null"`
	UUID        string       `gorm:"primary_key"`
	Password    string
}

type Team struct {
	Match       Match        `gorm:"foreignkey:MatchID;association_foreignkey:MatchID"`
	Single      bool         `gorm:"default:false"`
	TeamMembers []TeamMember `gorm:"foreignkey:TeamID;association_foreignkey:TeamID"`
	TeamID      uint64       `gorm:"primary_key"`
	MatchID     uint64
}

// Note TeamMembers to indcate that there is more than one person on the team
type TeamMember struct {
	Account     Account     `gorm:"foreignkey:AccountID;association_foreignkey:AccountID"`
	GameAccount GameAccount `gorm:"foreignkey:AccountID;association_foreignkey:AccountID"`
	Match       Match       `gorm:"foreignkey:MatchID;association_foreignkey:MatchID"`
	Team        Team        `gorm:"foreignkey:TeamID;association_foreignkey:TeamID"`
	Confirmed   bool        `gorm:"default:false"`
	Deny        bool        `gorm:"default:false"`
	TeamMembers bool        `gorm:"default:true"`
	AccountID   uint64
	GameID      uint64
	MatchID     uint64
	TeamID      uint64
	Winner      bool
}
