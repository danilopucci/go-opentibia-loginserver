package database

import (
	"database/sql"
	"fmt"
	"go-opentibia-loginserver/models"
)

type TvpQuery struct{}

func (q *TvpQuery) GetIpBanInfo(database *sql.DB, ip uint32) (models.BanInfo, error) {
	var banInfo models.BanInfo
	statement := fmt.Sprintf("SELECT `reason`, `expires_at`, `banned_by` FROM `ip_bans` WHERE `ip` = %d", ip)

	err := database.QueryRow(statement).Scan(&banInfo.Reason, &banInfo.ExpiresAt, &banInfo.Author)
	if err != nil {
		banInfo.IsBanned = false
		if err != sql.ErrNoRows {
			return banInfo, err
		}
		return banInfo, nil
	}

	banInfo.IsBanned = true
	return banInfo, nil
}

func (q *TvpQuery) GetAccountInfo(database *sql.DB, accountNumber uint32) (models.AccountInfo, error) {
	var accountInfo models.AccountInfo
	statement := fmt.Sprintf("SELECT `id`, `password`, `type`, `premium_ends_at` FROM `accounts` WHERE `id` = %d", accountNumber)

	err := database.QueryRow(statement).Scan(&accountInfo.Id, &accountInfo.PasswordSHA1, &accountInfo.AccountType, &accountInfo.PremiumEndsAt)
	if err != nil && err != sql.ErrNoRows {
		return accountInfo, err
	}

	return accountInfo, nil
}

func (q *TvpQuery) GetCharactersList(database *sql.DB, accountId uint32) ([]string, error) {
	var characterList []string

	statement := fmt.Sprintf("SELECT `name` FROM `players` WHERE `account_id` = %d AND `deletion` = 0 ORDER BY `name` ASC", accountId)

	rows, err := database.Query(statement)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var characterName string

		err := rows.Scan(&characterName)
		if err != nil {
			return nil, err
		}

		characterList = append(characterList, characterName)
	}

	return characterList, nil
}
