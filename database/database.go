package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/benjamin-larsen/NoctesChat-WebSocket/models"
)

var DB *sql.DB

type SQLQueryInterface interface {
	QueryRow(query string, args ...any) *sql.Row
}

func HasUserToken(token models.UserToken, tx SQLQueryInterface) (bool, error) {
	success := false

	err := tx.QueryRow(
		"SELECT 1 FROM user_tokens WHERE user_id = ? AND key_hash = ? FOR SHARE;",
		token.UserId,
		token.Token[:],
	).Scan(&success)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return success, nil
}

func InitDB() {
	db, err := sql.Open("mysql", os.Getenv("db_conn"))

	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	DB = db;
	log.Print("Connected succesfully to Database")
}