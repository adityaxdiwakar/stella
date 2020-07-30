package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/google/go-querystring/query"
)

type RemovableMessageStruct struct {
	SentID     string
	ReceivedID string
	ChannelID  string
	AuthorID   string
}

type FinvizEquityQueryStruct struct {
	Ticker    string `url:"t"`
	Type      string `url:"ty"`
	Technials string `url:"ta"`
	Period    string `url:"p"`
	Size      string `url:"s"`
}

func finvizCheckContentLength(chartUrl string) error {
	res, err := http.Head(chartUrl)
	if err != nil {
		return errors.New("Error getting HEAD of chart image")
	}

	if res.ContentLength < 7500 {
		return errors.New("Chart is too small, most likely empty")
	}

	return nil

}

func finvizEquityChartHandler(ticker string, timeframe int8) (string, string, error) {

	if len(ticker) > 8 {
		return "", "", errors.New("Ticker is too long, not moving forward")
	}

	// Timeframes for translations to Finviz-Query language
	timeframes := []string{"i1", "i3", "i5", "i15", "i30", "d", "w", "m"}

	// Determine chart type and remove TA on charts that are not >= 6
	binaryTA := int(0)
	if timeframe < 6 {
		binaryTA = 1
	}

	queryParams := FinvizEquityQueryStruct{
		Ticker:    ticker,
		Type:      "c",
		Technials: strconv.Itoa(binaryTA),
		Period:    timeframes[timeframe],
		Size:      "l",
	}

	rootUrl := "https://charts.aditya.diwakar.io/elite/chart.ashx"
	if timeframe > 5 {
		rootUrl = "https://charts.aditya.diwakar.io/chart.ashx"
	}

	if timeframe == 5 {
		queryParams.Technials = "st_c,sch_200p,sma_20,sma_50,sma_200,rsi_b_14,macd_b_12_26_9"
	}

	qStr, _ := query.Values(queryParams)
	chartUrl := fmt.Sprintf("%s?%s", rootUrl, qStr.Encode())

	if finvizCheckContentLength(chartUrl) != nil {
		return "", "", errors.New("Content Length check failed")
	}

	return chartUrl, timeframes[timeframe], nil
}

type FinvizFuturesQueryStruct struct {
	Ticker string `url:"t"`
	Period string `url:"p"`
	Size   string `url:"s"`
}

func finvizFuturesChartHandler(ticker string, timeframe int8) (string, string, error) {

	// No futures tickers are longer than 4 characters, therefore this is a hard limit
	if len(ticker) > 4 {
		return "", "", errors.New("Ticker is too long, not moving forward")
	}

	// Timeframes for translating to Finviz-Query language
	timeframes := []string{"m5", "h1", "d1", "w1", "m1"}

	queryParams := FinvizEquityQueryStruct{
		Ticker: ticker,
		Period: timeframes[timeframe],
		Size:   "l",
	}

	rootUrl := "https://charts.aditya.diwakar.io/fut_chart.ashx"

	qStr, _ := query.Values(queryParams)
	chartUrl := fmt.Sprintf("%s?%s", rootUrl, qStr.Encode())

	if finvizCheckContentLength(chartUrl) != nil {
		return "", "", errors.New("Content Length check failed")
	}

	return chartUrl, timeframes[timeframe], nil
}

func finvizForexChartHandler(ticker string, timeframe int8) (string, string, error) {
	// No forex tickers are longer than 9 characters, hard limit
	if len(ticker) > 9 {
		return "", "", errors.New("Ticker is too long, not moving forward")
	}

	ticker = strings.ToUpper(ticker)

	timeframes := []string{"m5", "h1", "d1", "w1", "mo"}
	currencies := []string{"EURUSD", "GBPUSD", "USDJPY", "USDCAD", "USDCHF", "AUDUSD", "NZDUSD", "EURGBP", "GBPJPY", "BTCUSD"}

	isValid := false
	for _, currency := range currencies {
		if currency == ticker {
			isValid = true
		}
	}
	if !isValid {
		return "", "", errors.New("Forex ticker provided is not in list, cannot continue")
	}

	chartUrl := fmt.Sprintf("https://charts.aditya.diwakar.io/fx_image.ashx?%s_%s_l.png", ticker, timeframes[timeframe])

	if finvizCheckContentLength(chartUrl) != nil {
		return "", "", errors.New("Content length check failed, cannot continue")
	}

	return chartUrl, timeframes[timeframe], nil
}

func finvizChartUrlDownloader(Url string, ticker string) (discordgo.File, error) {
	resp, err := http.Get(Url)
	if err != nil {
		return discordgo.File{}, errors.New("Could not download file, for some reason")
	}

	return discordgo.File{Name: fmt.Sprintf("%s.png", ticker), Reader: resp.Body}, nil
}

func finvizChartTimeframeTranslator(timeframe string) string {
	translationMap := map[string]string{
		"i1":  "1-minute intraday",
		"i3":  "3-minute intraday",
		"i5":  "5-minute intraday",
		"i15": "15-minute intraday",
		"i30": "30-minute intraday",
		"d":   "daily",
		"w":   "weekly",
		"m":   "monthly",
		"m5":  "5-minute",
		"h1":  "hourly",
		"d1":  "daily",
		"w1":  "weekly",
		"m1":  "monthly",
		"mo":  "monthly",
	}

	return translationMap[timeframe]
}

type channelChartObject struct {
	Ticker string
	File   *discordgo.File
	Error  error
}

func finvizDownloadToChannel(chartUrl, ticker string, ch chan channelChartObject) {
	file, err := finvizChartUrlDownloader(chartUrl, ticker)
	ch <- channelChartObject{
		Ticker: ticker,
		File:   &file,
		Error:  err,
	}
}

type channelURLObject struct {
	Ticker           string
	TimeframeMessage string
	ChartUrl         string
	Error            error
}

func finvizUrlToChannel(ticker string, marketType int, intChartType int8, ch chan channelURLObject) {

	var chartUrl string
	var timeframeMessage string
	var err error

	switch marketType {

	case 1:
		chartUrl, timeframeMessage, err = finvizFuturesChartHandler(ticker, intChartType)

	case 2:
		chartUrl, timeframeMessage, err = finvizForexChartHandler(ticker, intChartType)

	default:
		chartUrl, timeframeMessage, err = finvizEquityChartHandler(ticker, intChartType)

	}

	ch <- channelURLObject{
		Ticker:           ticker,
		TimeframeMessage: timeframeMessage,
		ChartUrl:         chartUrl,
		Error:            err,
	}
}

func boolToInt(a bool) int {
	r := 0
	if a {
		r = 1
	}
	return r
}

func finvizChartSender(s *discordgo.Session, m *discordgo.MessageCreate, mSplit []string, isFutures bool, isForex bool) {
	if len(mSplit[0]) > 1 {
		if _, err := strconv.Atoi(string(mSplit[0][1])); err != nil {
			// The second character of the command is not an integer but the length is 2, therefore invalid
			return
		}
	}

	msg, err := s.ChannelMessageSend(m.ChannelID, ":clock1: Fetching your chart, stand by.")
	if err != nil {
		// An error occured that stopped the queue message from being sent, cancel progression
		return
	}

	intChartType := int8(-1)
	if len(mSplit[0]) == 1 || unicode.IsLetter(rune(mSplit[0][1])) {
		if isFutures || isForex {
			intChartType = 1
		} else {
			intChartType = 5
		} // Default values
	} else {
		uncastedChartType, _ := strconv.Atoi(string(mSplit[0][1]))
		intChartType = int8(uncastedChartType)
	}

	if (intChartType > 5 || intChartType < 1) && (isFutures || isForex) {
		s.ChannelMessageEdit(msg.ChannelID, msg.ID, "You've requested an invalid chart timeframe, choose between 1 and 5.")
		return
	} else if (intChartType > 7 || intChartType < 0) && !isFutures {
		s.ChannelMessageEdit(msg.ChannelID, msg.ID, "You've requested an invalid chart timeframe, choose between 0 and 7.")
		return
	}

	if isFutures || isForex {
		intChartType--
	}

	tickers := unique(mSplit[1:])

	if len(tickers) > 10 {
		s.ChannelMessageEdit(msg.ChannelID, msg.ID, "Multiple charts are limited to 10 per command unfortunately, please try again with fewer than 10 tickers")
	}

	for i, t := range tickers {
		if t == "rty" {
			tickers[i] = "er2"
		}

		tickers[i] = strings.ReplaceAll(tickers[i], ".", "-")
	}

	chartsServed += len(tickers)
	rdb.IncrBy(ctx, "stats.charts.served", int64(len(tickers)))

	var timeframeMessage string
	var files []*discordgo.File
	var tickerErrorStack []string
	var urlStack []string
	urlChannel := make(chan channelURLObject, len(tickers))
	tickerChannel := make(chan channelChartObject, len(tickers))
	for _, ticker := range tickers {
		// override tickers (a/b shares and russell-2k, see #5)

		marketType := boolToInt(isFutures) + 2*boolToInt(isForex)

		go finvizUrlToChannel(ticker, marketType, intChartType, urlChannel)
	}

	for {
		if len(tickerErrorStack)+len(urlStack) == len(tickers) {
			break
		}

		chanURL := <-urlChannel
		if chanURL.Error != nil {
			tickerErrorStack = append(tickerErrorStack, chanURL.Ticker)
		} else {
			timeframeMessage = chanURL.TimeframeMessage
			urlStack = append(urlStack, chanURL.ChartUrl)
			go finvizDownloadToChannel(chanURL.ChartUrl, chanURL.Ticker, tickerChannel)
		}
	}

	fileMap := make(map[string]*discordgo.File)

	for {
		if len(tickerErrorStack)+len(fileMap) == len(tickers) {
			break
		}

		chanChart := <-tickerChannel
		if chanChart.Error != nil {
			tickerErrorStack = append(tickerErrorStack, chanChart.Ticker)
		} else {
			fileMap[chanChart.Ticker] = chanChart.File
		}
	}

	for _, erroredTicker := range tickerErrorStack {
		for index, t := range tickers {
			if t == erroredTicker {
				tickers = append(tickers[:index], tickers[index+1:]...)
				break
			}
		}
	}

	for _, t := range tickers {
		files = append(files, fileMap[t])
	}

	messageStack := []*discordgo.Message{}
	s.ChannelMessageDelete(msg.ChannelID, msg.ID)
	if len(tickerErrorStack) != len(tickers) {
		// Delete the interrim message, since it cannot be edited w/ files
		msg, _ := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: fmt.Sprintf("Here is your %s chart", finvizChartTimeframeTranslator(timeframeMessage)),
			Files:   files,
		})
		messageStack = append(messageStack, msg)
		removableMessages[msg.ID] = RemovableMessageStruct{
			SentID:     msg.ID,
			ReceivedID: m.ID,
			ChannelID:  msg.ChannelID,
			AuthorID:   m.Author.ID,
		}
	}

	if len(tickerErrorStack) > 0 {
		joinedTickerStack := strings.Join(tickerErrorStack, ", ")
		errorTickersMsg := fmt.Sprintf("**`%d`** tickers could not be loaded: ``%s``", len(tickerErrorStack), joinedTickerStack)
		msg, _ := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: errorTickersMsg,
		})
		messageStack = append(messageStack, msg)
		removableMessages[msg.ID] = RemovableMessageStruct{
			SentID:     msg.ID,
			ReceivedID: m.ID,
			ChannelID:  msg.ChannelID,
			AuthorID:   m.Author.ID,
		}
	}

	for _, message := range messageStack {
		s.MessageReactionAdd(message.ChannelID, message.ID, ":stellax:737458650490077196")
	}

	time.Sleep(90 * time.Second)

	for _, message := range messageStack {
		personalUser, err := s.User("@me")
		if err != nil {
			continue
		}
		personalID := personalUser.ID
		err = s.MessageReactionRemove(message.ChannelID, message.ID, ":stellax:737458650490077196", personalID)
		if err != nil {
			log.Println(err)
		}
	}

}
