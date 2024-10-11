package database

import (
	"database/sql"
	"fmt"
	"go-opentibia-loginserver/models"

	_ "github.com/go-sql-driver/mysql"
)

const DatabaseDriverName = "mysql"

func CreateDatabaseConnection(user string, password string, host string, port int, databaseName string) (*sql.DB, error) {
	dsn := generateConnectionString(user, password, host, port, databaseName)

	db, err := sql.Open(DatabaseDriverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("[createDatabaseConnection] - error opening database connection: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("[createDatabaseConnection] - error checking database connection: %s", err)
	}

	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		return nil, fmt.Errorf("[createDatabaseConnection] - error querying database version: %s", err)
	}

	fmt.Printf("[createDatabaseConnection] - connected to %s database (%s); Version: %s\n", DatabaseDriverName, databaseName, version)

	return db, nil
}

func generateConnectionString(user string, password string, host string, port int, databaseName string) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		user,
		password,
		host,
		port,
		databaseName,
	)
}

func GetIpBanInfo(database *sql.DB, ip uint32) (models.BanInfo, error) {
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

func GetAccountInfo(database *sql.DB, accountNumber uint32) (models.AccountInfo, error) {
	var accountInfo models.AccountInfo
	statement := fmt.Sprintf("SELECT `id`, `password`, `type`, `premium_ends_at` FROM `accounts` WHERE `id` = %d", accountNumber)

	err := database.QueryRow(statement).Scan(&accountInfo.Id, &accountInfo.PasswordSHA1, &accountInfo.AccountType, &accountInfo.PremiumEndsAt)
	if err != nil && err != sql.ErrNoRows {
		return accountInfo, err
	}

	return accountInfo, nil
}

func GetCharactersList(database *sql.DB, accountId uint32) ([]string, error) {
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
