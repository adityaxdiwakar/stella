package main

import (
	"database/sql"
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
	sqlStatement := `
        INSERT INTO tags (id, content)
        VALUES ($1, $2)`

	switch _, err := retrieveTagFromDb(tagName); err {
	case nil:
		sqlStatement = `
            UPDATE tags
            SET content = $2
            WHERE id = $1
        `

	case sql.ErrNoRows:
		break

	default:
		return err

	}

	_, err := db.Exec(sqlStatement, tagName, tagContent)
	if err != nil {
		return err
	}
	return nil
}

func deleteTagFromDb(tagName string) error {
	sqlStatement := `DELETE FROM tags WHERE ID = $1`

	_, err := retrieveTagFromDb(tagName)
	if err == sql.ErrNoRows {
		return sql.ErrNoRows
	} else if err != nil {
		return err
	}

	_, err = db.Exec(sqlStatement, tagName)
	if err != nil {
		return err
	}
	return nil
}
