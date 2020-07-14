package main

import (
	"database/sql"
	"errors"
	"log"
	//"github.com/bwmarrin/discordgo"
)

func retrieveTagFromDb(tagName string) (string, error) {
	sqlStatement := `SELECT content FROM tags WHERE id=$1`
	var content string
	row := db.QueryRow(sqlStatement, tagName)
	err := row.Scan(&content)
	if err != nil {
		return "", err
	}
	return content, nil
}

func pushTagToDb(tagName string, tagContent string) error {
	if _, err := retrieveTagFromDb(tagName); err != sql.ErrNoRows {
		return errors.New("Tag already exists")
	}

	sqlStatement := `
        INSERT INTO tags (id, content)
        VALUES ($1, $2)
    `
	_, err := db.Exec(sqlStatement, tagName, tagContent)
	if err != nil {
		return err
	}
	return nil
}
