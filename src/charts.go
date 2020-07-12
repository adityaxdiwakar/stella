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

func finvizEquityChartHandler(ticker string, timeframe int8) (string, error) {

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
		queryParams.Technials = "st_c,sch_200p,sma_50,sma_200,sma_20,sma_100,bb_20_2,rsi_b_14,macd_b_12_26_9,stofu_b_14_3_3"
	}

	qStr, _ := query.Values(queryParams)
	chartUrl := fmt.Sprintf("%s?%s", rootUrl, qStr.Encode())

	if finvizCheckContentLength(chartUrl) != nil {
		return "", errors.New("Content Length check failed")
	}

	return chartUrl, nil
}

type FinvizFuturesQueryStruct struct {
	Ticker string `url:"t"`
	Period string `url:"p"`
	Size   string `url:"s"`
}

func finvizFuturesChartHandler(ticker string, timeframe int8) (string, error) {

	// No futures tickers are longer than 4 characters, therefore this is a hard limit
	if len(ticker) > 4 {
		return "", errors.New("Ticker is too long, not moving forward")
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
		return "", errors.New("Content Length check failed")
	}

	return chartUrl, nil
}

func finvizForexChartHandler(ticker string, timeframe int8) (string, error) {
	// No forex tickers are longer than 9 characters, hard limit
	if len(ticker) > 9 {
		return "", errors.New("Ticker is too long, not moving forward")
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
		return "", errors.New("Forex ticker provided is not in list, cannot continue")
	}

	chartUrl := fmt.Sprintf("https://charts.aditya.diwakar.io/fx_image.ashx?%s_%s_l.png", ticker, timeframes[timeframe])

	if finvizCheckContentLength(chartUrl) != nil {
		return "", errors.New("Content length check failed, cannot continue")
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

func finvizChartSender(s *discordgo.Session, m *discordgo.MessageCreate, mSplit []string, isFutures bool, isForex bool) {
	msg, err := s.ChannelMessageSend(m.ChannelID, ":clock1: Fetching your chart, stand by.")
	if err != nil {
		// An error occured that stopped the queue message from being sent, cancel progression
		return
	}

	if len(mSplit[0]) > 1 {
		if _, err := strconv.Atoi(string(mSplit[0][1])); err != nil {
			// The second character of the command is not an integer but the length is 2, therefore invalid
		}
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

	chartsServed += len(tickers)

	var files []*discordgo.File
	var tickerErrorStack []string
	for _, ticker := range tickers {
		var chartUrl string
		var err error

		switch {

		case isFutures:
			chartUrl, err = finvizFuturesChartHandler(ticker, intChartType)

		case isForex:
			chartUrl, err = finvizForexChartHandler(ticker, intChartType)

		default:
			chartUrl, err = finvizEquityChartHandler(ticker, intChartType)

		}

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
