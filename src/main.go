package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/adityaxdiwakar/flux"
	tda "github.com/adityaxdiwakar/tda-go"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"golang.org/x/text/message"
)

var startTime time.Time
var chartsServed int
var messagesSeen int64
var ctx = context.Background()
var rdb *redis.Client
var db *sql.DB
var printer *message.Printer
var tds tda.Session
var fluxS *flux.Session
var tickerChannels []string
var conf tomlConfig
var removableMessages map[string]RemovableMessageStruct

var stellaHttpClient = &http.Client{Timeout: 10 * time.Second}

func init() {
	if _, err := toml.DecodeFile("config/config.toml", &conf); err != nil {
		log.Fatalf("error: could not parse configuration: %v\n", err)
	}

	// initialize the global starttime, for uptime calculations
	startTime = time.Now()

	// global pseudo random generator
	rand.Seed(time.Now().Unix())

	// establish connection with Redis DB
	rdb = redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Address,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})
	// ping rdb to test, use context for the situation
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Could not make connection with Redis")
	}

	pSqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s "+
		"sslmode=disable", conf.Database.Host, conf.Database.Port,
		conf.Database.User, conf.Database.Password, conf.Database.DBName)

	db, err = sql.Open("postgres", pSqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// establish english printer
	printer = message.NewPrinter(message.MatchLanguage(conf.Language))

	// intitialize tda lib
	tds = tda.Session{
		Refresh:     conf.TDAmeritrade.RefreshToken,
		ConsumerKey: conf.TDAmeritrade.ConsumerKey,
		RootUrl:     "https://api.tdameritrade.com/v1",
	}
	tds.InitSession()

	fluxS, err = flux.New(tds)
	if err != nil {
		log.Fatal(err)
	}

	fluxS.Open()

	removableMessages = make(map[string]RemovableMessageStruct)
}

func uptime() string {
	return time.Since(startTime).Round(time.Second).String()
}

func main() {
	defer db.Close()

	tickerPtr := flag.Bool("ticker", true, "ticker boolean")
	flag.Parse()

	dg, err := discordgo.New(fmt.Sprintf("Bot %s", conf.DiscordConfig.Token))
	if err != nil {
		fmt.Println("Error creating Discord Session due to:", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(reactionHandler)

	if *tickerPtr {
		go channelTicker(dg)
		go playingTicker(dg)
	} else {
		fmt.Println("Non-Default Boot: Ticker Feature Disabled")
	}

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	fmt.Println(`

	  Stella is now loaded.


          //       **       //
        //////     **     //////
        //////     **     //////
        //////   ******   //////
        //////   ******   //////
        //////   ******   //////
          //     ******     //
          //       **       //
                   **

    `)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	fmt.Println("Interrupt received, terminating Stella.")

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	messagesSeen += 1
	rdb.Incr(ctx, "stats.msgs.seen")

	// Check if the prefix is mentioned by using strings.HasPrefix
	if !strings.HasPrefix(m.Content, conf.DiscordConfig.Prefix) {
		return
	}

	// If the prefix is present, remove the prefix for later handling
	m.Content = m.Content[len(conf.DiscordConfig.Prefix):]
	mSplit := strings.Split(m.Content, " ")

	switch {

	case mSplit[0] == "ping":
		ping(s, m)

	case mSplit[0] == "name":
		sendCompanyName(s, m, mSplit)

	case mSplit[0] == "fun":
		sendFundamentals(s, m, mSplit)

	case mSplit[0] == "div":
		sendDividends(s, m, mSplit)

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

	case mSplit[0] == "search":
		searchTicker(s, m, mSplit)

	case mSplit[0] == "help":
		help(s, m, mSplit)

	case mSplit[0] == "bio":
		reutersBio(s, m, mSplit)

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
	if conf.DiscordConfig.Env == "dev" {
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
			{
				Name: "Status",
				Value: fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
					printer.Sprintf("Messages Seen: **%d**", messagesSeen),
					printer.Sprintf("Charts Served: **%d**", chartsServed),
					printer.Sprintf("Uptime: **%s**", uptime()),
					printer.Sprintf("Version: **v0.93.1**"),
					printer.Sprintf("Heartbeat Latency: **%dms**", s.HeartbeatLatency().Milliseconds()),
				),
			},
			{
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
