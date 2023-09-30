package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Chat_id          int
	Tg_name          string
	Eth_address      string
	Private_key      string
	Original_address string
}

func Connect(dbconn string) (*sql.DB, error) {
	sqlDB, err := sql.Open("sqlite3", dbconn)

	if err != nil {
		return nil, err
	}

	_, err = sqlDB.Exec(create)

	if err != nil {
		return nil, err
	}

	return sqlDB, nil

}

func SelectSingleUser(tgNick string, db *sql.DB) (User, error) {

	var user User

	row := db.QueryRow(selectSingleUser, tgNick)

	err := row.Scan(&user.Chat_id, &user.Tg_name, &user.Eth_address, &user.Private_key)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func InsertIntoTable(chat_id int64, tgNick string, eth_address string, privateKey string, db *sql.DB) error {
	stmt, err := db.Prepare(insertIntoTable)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(chat_id, tgNick, eth_address, privateKey)

	if err != nil {
		return err
	}

	return nil
}

const create string = `
  CREATE TABLE IF NOT EXISTS users (
  CHAT_ID INTEGER NOT NULL PRIMARY KEY,
  TG_NAME TEXT NOT NULL,
  ADDRESS TEXT,
  PRIVATE_KEY TEXT UNIQUE NOT NULL
  );`

const selectSingleUser string = `
  Select *
  From users
  Where TG_NAME = ?
  `
const insertIntoTable string = `
INSERT INTO users VALUES ( ?, ?, ?, ?)
`
