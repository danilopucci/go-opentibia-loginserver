package database

import (
	"database/sql"
	"fmt"
	"go-opentibia-loginserver/models"

	_ "github.com/go-sql-driver/mysql"
)

const DatabaseDriverName = "mysql"

type DatabaseQuery interface {
	GetIpBanInfo(database *sql.DB, ip uint32) (models.BanInfo, error)
	GetAccountInfo(database *sql.DB, accountNumber uint32) (models.AccountInfo, error)
	GetCharactersList(database *sql.DB, accountId uint32) ([]string, error)
}

// placeholders only
type Otx2Query struct{}
type NostalriusQuery struct{}
type DefaultQuery struct{}

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

func GetDatabaseQuery(version string) DatabaseQuery {
	switch version {
	case "tvp":
		return &TvpQuery{}
		// case "otx2":
		// 	return &Otx2Query{}
		// case "nostalrius":
		// 	return &NostalriusQuery{}

		// default:
		// 	return &DefaultQuery{}
	}

	return nil
}
