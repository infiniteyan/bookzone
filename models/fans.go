package models

type Fans struct {
	Id 			int
	MemberId 	int
	FansId 		int	`orm:"index"`
}

type FansData struct {
	MemberId 	int
	Nickname 	string
	Avatar 		string
	Account 	string
}