package main

import (
	"fmt"
	"strings"
)

var (
	// Errors
	ErrReutersPkgFailed      = "reuters error: package failed to be requested"
	ErrReutersResponseFailed = "reuters error: non-200/206 response code"
	ErrNoTickerProvided      = "Please provider a ticker to look up"
	ErrNoHead                = "charts error: could not get HEAD of chart image"
	ErrChartSmall            = "charts error: chart too small (assuming empty)"
	ErrTickerOverLength      = "charts error: ignoring ticker that is too long"
	ErrContentLengthFailed   = "charts error: error finding content length"
	ErrTickerNotForexList    = "charts error: cannot continue with forex ticker not in list"
	ErrCannotDownloadFile    = "charts error: could not download file"
	ErrInvalidChart          = "You've requested an invalid chart timeframe"
	ErrTooManyCharts         = "Multi-Charts are limited to 10 per command, please try again with fewer."

	// Errored Parameterized Messages
	ErrInvalidChartTime = func(min, max int) string {
		return fmt.Sprintf("%s, choose between %d and %d.", ErrInvalidChart, min, max)
	}
	ErrCouldNotLoadTickers = func(num int, tickers string) string {
		return fmt.Sprintf("**`%d`** tickers could not be loaded: ``%s``", num, tickers)
	}

	// General Messages
	StellaXReaction    = ":stellax:737458650490077196:"
	CantFindCompanyBio = "The company bio could not be found, try again?"
	ErroredCompanyBio  = "Something went wrong or the ticker does not exist in the database"
	FetchingChart      = ":clock1: Fetching your chart, stand by..."

	// General Parametrized Messages
	LoadingCompanyBio = func(ticker string) string {
		return fmt.Sprintf("Loading company bio for %s, please wait...", strings.ToUpper(ticker))
	}
	ChartWithTimeResponse = func(timeframe string) string {
		return fmt.Sprintf("Here is your %s chart", timeframe)
	}
)
