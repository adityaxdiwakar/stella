package main

import (
	"github.com/bwmarrin/discordgo"
)

func reactionHandler(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	personalUser, err := s.User("@me")
	if err != nil {
		return
	}

	if e.UserID == personalUser.ID {
		return
	}

	if e.Emoji.ID != "737458650490077196" {
		return
	}

	delMsgObject := removableMessages[e.MessageID]
	if (delMsgObject == RemovableMessageStruct{}) {
		return
	}

	if e.UserID != delMsgObject.AuthorID {
		return
	}

	s.ChannelMessageDelete(delMsgObject.ChannelID, delMsgObject.SentID)
	s.ChannelMessageDelete(delMsgObject.ChannelID, delMsgObject.ReceivedID)
}
