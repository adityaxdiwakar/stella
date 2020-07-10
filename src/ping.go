package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Collect current timestamp for comparision
	now := time.Now()

	// Send the message to Discord
	msg, _ := s.ChannelMessageSend(m.ChannelID, "Pong!")

	// Calculate the timestamp from the snowflake
	timestamp, _ := discordgo.SnowflakeTimestamp(msg.ID)

	// Find the difference in the API timestamp and local timestamp
	diff := int32(timestamp.Sub(now).Seconds() * 1000)

	// Send the ping message into the respective channel
	msg, _ = s.ChannelMessageEdit(msg.ChannelID, msg.ID, fmt.Sprintf(":ping_pong: WS Roundtrip: %dms!", diff))
}
