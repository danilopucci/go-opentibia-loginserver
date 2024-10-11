package models

type BanInfo struct {
	Author    string
	Reason    string
	ExpiresAt int64
	IsBanned  bool
}

type AccountInfo struct {
	Id            uint32
	PasswordSHA1  string
	AccountType   uint32
	PremiumEndsAt int64
	Characters    []string
}
