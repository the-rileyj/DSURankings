package main

func (a *Account) ResponseAccount() ResponseAccount {
	return ResponseAccount{
		a.AccountID,
		a.Email,
		a.FirstName,
		a.GameAccounts,
		a.LastName,
		a.UserName,
	}
}

func (pa *PendingAccount) ResponseAccount() ResponsePendingAccount {
	return ResponsePendingAccount{
		pa.Email,
		pa.FirstName,
		pa.LastName,
		pa.UserName,
	}
}
