package main

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
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

func listTagsFromDb() ([]string, error) {
	rows, err := db.Query(`SELECT id FROM tags`)
	if err != nil {
		return []string{""}, err
	}

	var tags []string
	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			return []string{""}, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func addTag(s *discordgo.Session, m *discordgo.MessageCreate, mSplit []string) {
	if len(mSplit) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a tag name and tag content for the tag to be added")
		return
	}

	tagName := mSplit[1]
	tagContent := strings.Join(mSplit[2:], " ")

	err := pushTagToDb(tagName, tagContent)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "There was an error adding your tag to the database, if the issue persist, contact aditya")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The tag **%s** was successfully added to the database!", tagName))
}

func retrieveTag(s *discordgo.Session, m *discordgo.MessageCreate, mSplit []string) {
	if len(mSplit) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a tag name to be retrieved from the tag database")
		return
	}

	tagName := mSplit[1]

	switch content, err := retrieveTagFromDb(tagName); err {
	case nil:
		s.ChannelMessageSend(m.ChannelID, content)

	case sql.ErrNoRows:
		s.ChannelMessageSend(m.ChannelID, "The tag you requested could not be found, try again")

	default:
		s.ChannelMessageSend(m.ChannelID, "An error occured when requesting the tag from the database, if the issue persists "+
			"contact aditya")

	}

}

func deleteTag(s *discordgo.Session, m *discordgo.MessageCreate, mSplit []string) {
	if len(mSplit) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a tag name to be deleted from the tag database")
		return
	}

	tagName := mSplit[1]

	switch err := deleteTagFromDb(tagName); err {

	case nil:
		s.ChannelMessageSend(m.ChannelID, "Success! I've sent the tag into the nearest black hole.")

	case sql.ErrNoRows:
		s.ChannelMessageSend(m.ChannelID, "I can't delete something that doesn't exist, I'm not a magician!")

	default:
		s.ChannelMessageSend(m.ChannelID, "Something went wrong with deleting your tag, if the issue persists, contact aditya.")

	}
}

func showTags(s *discordgo.Session, m *discordgo.MessageCreate) {
	tags, err := listTagsFromDb()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "An error occured when attempting to load the list of tags, if the issue persists: "+
			"contact aditya")
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Here you are the available tags: ``%s`` "+
		"Please only test in #bot-spam", strings.Join(tags, ", ")))
}
