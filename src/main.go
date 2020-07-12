package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var startTime time.Time
var chartsServed int
var messagesSeen int64
var ctx = context.Background()
var rdb *redis.Client

func init() {
	startTime = time.Now()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rand.Seed(time.Now().Unix()) // global pseudo random generator

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Could not make connection with Redis")
	}
}

func uptime() string {
	return time.Since(startTime).Round(time.Second).String()
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
	messagesSeen += 1
	rdb.Incr(ctx, "stats.msgs.seen")

	// Check if the prefix is mentioned by using strings.HasPrefix
	if !strings.HasPrefix(m.Content, os.Getenv("PREFIX")) {
		return
	}

	// If the prefix is present, remove the prefix for later handling
	m.Content = m.Content[len(os.Getenv("PREFIX")):]
	mSplit := strings.Split(m.Content, " ")

	switch {

	case mSplit[0] == "ping":
		ping(s, m)

	case strings.HasPrefix(mSplit[0], "c"):
		finvizChartSender(s, m, mSplit, false, false)

	case strings.HasPrefix(mSplit[0], "f"):
		finvizChartSender(s, m, mSplit, true, false)

	case strings.HasPrefix(mSplit[0], "x"):
		finvizChartSender(s, m, mSplit, false, true)

	case mSplit[0] == "8ball":
		eightballSend(s, m)

	case mSplit[0] == "v":
		stellaVersion(s, m)

	}
}

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func stellaVersion(s *discordgo.Session, m *discordgo.MessageCreate) {
	var title string
	if os.Getenv("ENV") == "DEV" {
		title = "About Stella [Dev]"
	} else {
		title = "About Stella"
	}

	lifetimeMessagesSeen, msgsErr := rdb.Get(ctx, "stats.msgs.seen").Result()
	lifetimeChartsServed, chartsErr := rdb.Get(ctx, "stats.charts.served").Result()

	if msgsErr != nil || chartsErr != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not get the lifetime statistics, try again later.")
	}

	embed := &discordgo.MessageEmbed{
		Color:       0x00cd6e,
		Title:       title,
		Description: "Discord Bot for Financial Markets",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: "Status",
				Value: fmt.Sprintf("%s\n%s\n%s\n%s",
					fmt.Sprintf("**Messages Seen**: %d", messagesSeen),
					fmt.Sprintf("**Charts Served:** %d", chartsServed),
					fmt.Sprintf("**Uptime**: %s", uptime()),
					fmt.Sprintf("**Version**: v0.3"),
				),
			},
			&discordgo.MessageEmbedField{
				Name: "Lifetime Statistics",
				Value: fmt.Sprintf("%s\n%s",
					fmt.Sprintf("**Messages Seen**: %s", lifetimeMessagesSeen),
					fmt.Sprintf("**Charts Served**: %s", lifetimeChartsServed),
				),
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "http://img.aditya.diwakar.io/stellaLogo.png",
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Made with ❤️ by Aditya Diwakar",
		},
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
