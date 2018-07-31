package main

func (a *Account) APIAccount() APIAccount {
	return APIAccount{
		a.AccountID,
		a.Email,
		a.FirstName,
		a.GlobalPermissions,
		a.LastName,
		a.UserName,
	}
}

func (a *Account) AdvancedResponse() ResponseAdvancedAccount {
	return ResponseAdvancedAccount{
		a.GameAccounts,
		a.BasicResponse(),
	}
}

func (a *Account) BasicResponse() ResponseBasicAccount {
	return ResponseBasicAccount{
		a.AccountID,
		a.Email,
		a.FirstName,
		a.LastName,
		a.UserName,
	}
}

func (pa *PendingAccount) Response() ResponsePendingAccount {
	return ResponsePendingAccount{
		pa.Email,
		pa.FirstName,
		pa.LastName,
		pa.UserName,
	}
}

func (g *Game) Response() ResponseGame {
	return ResponseGame{
		g.Colors,
		g.GameID,
		g.GameName,
	}
}

func (ga *GameAccount) AdvancedResponse() ResponseAdvancedGameAccount {
	return ResponseAdvancedGameAccount{
		ga.BasicResponse(),
		ga.Game.Response(),
	}
}

func (ga *GameAccount) BasicResponse() ResponseBasicGameAccount {
	return ResponseBasicGameAccount{
		ga.Account.BasicResponse(),
		ga.GameID,
		ga.Score,
	}
}

func (m *Match) AdvancedResponse() ResponseAdvancedMatch {
	return ResponseAdvancedMatch{
		m.Account.BasicResponse(),
		m.Game.Response(),
		m.LosingTeam.AdvancedResponse(),
		m.WinningTeam.AdvancedResponse(),
		m.MatchID,
		m.Confirmed,
		m.GameID,
		m.LosingTeamID,
		m.MatchTime,
		m.WinningTeamID,
	}
}

func (t *Team) AdvancedResponse() ResponseAdvancedTeam {
	rt := ResponseAdvancedTeam{
		MatchID:     t.MatchID,
		TeamID:      t.TeamID,
		TeamMembers: make([]ResponseAdvancedTeamMember, 0),
	}
	for _, teamMember := range t.TeamMembers {
		rt.TeamMembers = append(rt.TeamMembers, teamMember.AdvancedResponse())
	}
	return rt
}

func (t *Team) BasicResponse() ResponseBasicTeam {
	rt := ResponseBasicTeam{
		MatchID:     t.MatchID,
		TeamID:      t.TeamID,
		TeamMembers: make([]ResponseBasicTeamMember, 0),
	}
	for _, teamMember := range t.TeamMembers {
		rt.TeamMembers = append(rt.TeamMembers, teamMember.BasicResponse())
	}
	return rt
}

func (tm *TeamMember) AdvancedResponse() ResponseAdvancedTeamMember {
	return ResponseAdvancedTeamMember{
		tm.Account.BasicResponse(),
		tm.GameAccount.BasicResponse(),
		tm.BasicResponse(),
	}
}

func (tm *TeamMember) BasicResponse() ResponseBasicTeamMember {
	return ResponseBasicTeamMember{
		tm.AccountID,
		tm.Confirmed,
		tm.GameID,
		tm.MatchID,
		tm.TeamID,
		tm.TeamMembers,
	}
}
