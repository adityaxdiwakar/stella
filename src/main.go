package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/text/message"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var startTime time.Time
var chartsServed int
var messagesSeen int64
var ctx = context.Background()
var rdb *redis.Client
var db *sql.DB
var printer *message.Printer

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "stella"
)

func init() {
	// initialize the global starttime, for uptime calculations
	startTime = time.Now()

	// load the dotenv file for environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// global pseudo random generator
	rand.Seed(time.Now().Unix())

	// establish connection with Redis DB
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	// ping rdb to test, use context for the situation
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Could not make connection with Redis")
	}

	pSqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", pSqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// establish english printer
	printer = message.NewPrinter(message.MatchLanguage("en"))
}

func uptime() string {
	return time.Since(startTime).Round(time.Second).String()
}

func main() {
	defer db.Close()

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

	case mSplit[0] == "addtag":
		addTag(s, m, mSplit)

	case mSplit[0] == "tag":
		retrieveTag(s, m, mSplit)

	case mSplit[0] == "deltag":
		deleteTag(s, m, mSplit)

	case mSplit[0] == "showtags":
		showTags(s, m)

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

	lifetimeMessagesSeenStr, msgsErr := rdb.Get(ctx, "stats.msgs.seen").Result()
	lifetimeChartsServedStr, chartsErr := rdb.Get(ctx, "stats.charts.served").Result()

	lifetimeMessagesSeen, msgsErr := strconv.Atoi(lifetimeMessagesSeenStr)
	lifetimeChartsServed, chartsErr := strconv.Atoi(lifetimeChartsServedStr)

	if msgsErr != nil || chartsErr != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not get the lifetime statistics, try again later.")
		return
	}

	embed := &discordgo.MessageEmbed{
		Color:       0x00cd6e,
		Title:       title,
		Description: "Discord Bot for Financial Markets",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: "Status",
				Value: fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
					printer.Sprintf("Messages Seen: **%d**", messagesSeen),
					printer.Sprintf("Charts Served: **%d**", chartsServed),
					printer.Sprintf("Uptime: **%s**", uptime()),
					printer.Sprintf("Version: **v0.43**"),
					printer.Sprintf("Heartbeat Latency: **%dms**", s.HeartbeatLatency().Milliseconds()),
				),
			},
			&discordgo.MessageEmbedField{
				Name: "Lifetime Statistics",
				Value: fmt.Sprintf("%s\n%s",
					printer.Sprintf("Messages Seen: **%d**", lifetimeMessagesSeen),
					printer.Sprintf("Charts Served: **%d**", lifetimeChartsServed),
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
