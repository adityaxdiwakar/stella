package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	dg, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("Error creating Discord Session due to:", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "ping" {
		ping(s, m)
	} else if m.Content == "queryChart" {
		finvizChartHandler("aapl", 5)
	}
}

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