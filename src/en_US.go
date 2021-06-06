package main

import (
	"fmt"
	"strings"
)

var (
	ErrReutersPkgFailed      = "reuters error: package failed to be requested"
	ErrReutersResponseFailed = "reuters error: non-200/206 response code"
	ErrNoTickerProvided      = "Please provider a ticker to look up"

	LoadingCompanyBio = func(ticker string) string {
		return fmt.Sprintf("Loading company bio for %s, please wait...", strings.ToUpper(ticker))
	}
	CantFindCompanyBio = "The company bio could not be found, try again?"
	ErroredCompanyBio  = "Something went wrong or the ticker does not exist in the database"
)
