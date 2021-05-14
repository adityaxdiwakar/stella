package main

import (
	"fmt"

	"github.com/adityaxdiwakar/flux"
	"github.com/bwmarrin/discordgo"
)

func quoteTicker(s *discordgo.Session, m *discordgo.MessageCreate, mSplit []string, count int) {
	if len(mSplit) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a ticker to search")
	}

	// I like mark, bid, bid size, ask, ask size, too,
	// last, last change, percent

	searchResponse, err := fluxS.RequestQuote(flux.QuoteRequestSignature{
		Ticker:      mSplit[1],
		RefreshRate: 300,
		Fields: []flux.QuoteField{
			flux.Bid,
			flux.BidSize,
			flux.Ask,
			flux.AskSize,
			flux.Volume,
			flux.Last,
			flux.LastSize,
			flux.NetChange,
			flux.NetPercentChange,
			flux.Mark,
			flux.MarkChange,
			flux.MarkPercentChange,
		},
	})

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Something went wrong while processing the request!")
		return
	}

	payload := searchResponse.Items[0].Values
	if payload.LAST == 0 {
		if count == 2 {
			s.ChannelMessageSend(m.ChannelID, "Something went wrong while processing the request!")
			return
		} else {
			quoteTicker(s, m, mSplit, count+1)
			return
		}
	}

	var sChange string
	var sPercent string
	if payload.NETCHANGE > 0 {
		sChange = fmt.Sprintf("+%.2f", payload.MARKCHANGE)
		sPercent = fmt.Sprintf("+%.2f%%", payload.MARKPERCENTCHANGE*100)
	} else {
		sChange = fmt.Sprintf("%.2f", payload.MARKCHANGE)
		sPercent = fmt.Sprintf("%.2f%%", payload.MARKPERCENTCHANGE*100)
	}

	responseText := printer.Sprintf(("__Quote Information for %s__\n" +
		"%.2f %s (%s)\n\n" +
		"**Last:** %.2f (x%d)\n" +
		"**Bid:** %.2f (x%d)\n" +
		"**Ask:** %.2f (x%d)\n" +
		"**Volume:** %d"),
		searchResponse.Items[0].Symbol,
		payload.MARK, sChange, sPercent,
		payload.LAST, payload.LASTSIZE,
		payload.BID, payload.BIDSIZE,
		payload.ASK, payload.ASKSIZE,
		payload.VOLUME)
	s.ChannelMessageSend(m.ChannelID, string(responseText))
}
