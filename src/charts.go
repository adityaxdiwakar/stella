package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/google/go-querystring/query"
)

type FinvizQueryStruct struct {
	Ticker    string `url:"t"`
	Type      string `url:"ty"`
	Technials string `url:"ta"`
	Period    string `url:"p"`
	Size      string `url:"s"`
}

func finvizChartHandler(ticker string, timeframe int8) (string, error) {

	if len(ticker) > 8 {
		return "", errors.New("Ticker is too long, not moving forward")
	}

	// Timeframes for translations to Finviz-Query language
	timeframes := []string{"i1", "i3", "i5", "i15", "i30", "d", "w", "m"}

	// Determine chart type and remove TA on charts that are not >= 6
	binaryTA := int(0)
	if timeframe < 6 {
		binaryTA = 1
	}

	queryParams := FinvizQueryStruct{
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
		queryParams.Technials = "st_c,sch_200p,sma_50,sma_200,sma_20,sma_100,bb_20_2,rsi_b_14,macd_b_12_26_9,stofu_b_14_3_3"
	}

	qStr, _ := query.Values(queryParams)
	chartUrl := fmt.Sprintf("%s?%s", rootUrl, qStr.Encode())

	res, err := http.Head(chartUrl)
	if err != nil {
		return "", errors.New("Error getting HEAD of chart image")
	}

	if res.ContentLength < 7500 {
		return "", errors.New("Chart is too small, most likely empty")
	}

	return chartUrl, nil

}

func finvizChartUrlDownloader(Url string) (discordgo.File, error) {
	resp, err := http.Get(Url)
	if err != nil {
		return discordgo.File{}, errors.New("Could not download file, for some reason")
	}

	return discordgo.File{Name: "chart.png", Reader: resp.Body}, nil
}

func finvizChartSender(s *discordgo.Session, m *discordgo.MessageCreate, mSplit []string) {
	msg, err := s.ChannelMessageSend(m.ChannelID, ":clock1: Fetching your chart, stand by.")
	if err != nil {
		// An error occured that stopped the queue message from being sent, cancel progression
		return
	}

	intChartType := int8(-1)
	if len(mSplit[0]) == 1 || unicode.IsLetter(rune(mSplit[0][1])) {
		intChartType = 5
	} else {
		uncastedChartType, _ := strconv.Atoi(string(mSplit[0][1]))
		intChartType = int8(uncastedChartType)
	}

	tickers := unique(mSplit[1:])

	if len(tickers) == 1 {
		chartUrl, err := finvizChartHandler(tickers[0], intChartType)
		if err != nil {
			errorMessage := fmt.Sprintf("`%s` could not be found, could it be a typo?", tickers[0])
			s.ChannelMessageSend(m.ChannelID, errorMessage)
		} else {
			messageEmbed := discordgo.MessageEmbed{Image: &discordgo.MessageEmbedImage{URL: chartUrl}}
			msgContent := ""
			s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				Content: &msgContent,
				Embed:   &messageEmbed,
				ID:      msg.ID,
				Channel: msg.ChannelID,
			})
		}
	} else {
		var files []*discordgo.File
		var tickerErrorStack []string
		for _, ticker := range tickers {
			chartUrl, err := finvizChartHandler(ticker, intChartType)
			if err != nil {
				tickerErrorStack = append(tickerErrorStack, ticker)
			} else {
				file, err := finvizChartUrlDownloader(chartUrl)
				if err != nil {
					tickerErrorStack = append(tickerErrorStack, ticker)
				} else {
					files = append(files, &file)
				}
			}
		}

		// Delete the interrim message, since it cannot be edited w/ files
		s.ChannelMessageDelete(msg.ChannelID, msg.ID)
		for _, file := range files {
			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Files: []*discordgo.File{file},
			})
		}
		if len(tickerErrorStack) > 0 {
			joinedTickerStack := strings.Join(tickerErrorStack, ", ")
			errorTickersMsg := fmt.Sprintf("**`%d`** tickers could not be loaded: ``%s``", len(tickerErrorStack), joinedTickerStack)
			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: errorTickersMsg,
			})
		}
	}
}
