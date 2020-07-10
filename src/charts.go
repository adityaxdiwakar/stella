package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/google/go-querystring/query"
)

type FinvizQueryStruct struct {
	Ticker    string `url:"t"`
	Type      string `url:"ty"`
	Technials string `url:"ta"`
	Period    string `url:"p"`
	Size      string `url:"s"`
}

func finvizChartHandler(ticker string, timeframe int8) string {

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

	rootUrl := "https://elite.finviz.com/chart.ashx"
	if timeframe > 5 {
		rootUrl = "https://finviz.com/chart.ashx"
	}

	if timeframe == 5 {
		queryParams.Technials = "st_c,sch_200p,sma_50,sma_200,sma_20,sma_100,bb_20_2,rsi_b_14,macd_b_12_26_9,stofu_b_14_3_3"
	}

	qStr, _ := query.Values(queryParams)
	chartUrl := fmt.Sprintf("%s?%s", rootUrl, qStr.Encode())

	return "nil"
}
