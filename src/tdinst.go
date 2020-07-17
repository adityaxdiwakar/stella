package main

import (
	"encoding/json"
	//	tda "github.com/adityaxdiwakar/tda-go"
	"github.com/adityaxdiwakar/tda-go"
	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/copier"
)

type ShortListFundamentals struct {
	PeRatio           json.Number
	PbRatio           json.Number
	TotalDebtToEquity json.Number
	PegRatio          json.Number
	Beta              json.Number
	DividendYield     json.Number
	ReturnOnEquity    json.Number
	MarketCap         json.Number
}

func getRequiredFundamentals(ticker string) (*ShortListFundamentals, error) {
	allFundamentals, err := tds.GetInstrumentFundamentals(ticker)
	if err != nil {
		return nil, err
	}

	var shortList ShortListFundamentals
	copier.Copy(&shortList, &allFundamentals.Fundamental)

	return &shortList, nil
}

func sendFundamentals(s *discordgo.Session, m *discordgo.MessageCreate, mSplit []string) {
	if len(mSplit) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a ticker to look up fundamentals for!")
		return
	}

	ticker := mSplit[1]

	shortList, err := getRequiredFundamentals(ticker)

	switch err {
	case nil:
		break

	case tda.FundamentalsEmpty:
		s.ChannelMessageSend(m.ChannelID, "That ticker could not be looked up unfortunately, try again.")
		return

	default:
		s.ChannelMessageSend(m.ChannelID, "Something went wrong while processing the fundamentals, try again later.")
		return
	}

}
