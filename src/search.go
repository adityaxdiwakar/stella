package main

import (
	"fmt"
	"strings"

	"github.com/adityaxdiwakar/flux"
	"github.com/bwmarrin/discordgo"
)

func searchTicker(s *discordgo.Session, m *discordgo.MessageCreate, mSplit []string) {
	if len(mSplit) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a query!")
		return
	}

	query := strings.Join(mSplit[1:], " ")

	searchResponse, err := fluxS.RequestSearch(flux.SearchRequestSignature{
		Pattern: query,
		Limit:   5,
	})

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Something went wrong while processing your query!")
		return
	}

	validResponses := 0
	response := fmt.Sprintf("Relevant responses to your query `%s`:", query)
	for _, instrument := range searchResponse.Value.Instruments {
		if strings.HasPrefix(instrument.Symbol, "0") {
			continue
		}

		if instrument.SourceType != "NYSE" && instrument.SourceType != "NASDAQ" {
			continue
		}

		response += fmt.Sprintf("\n- (**`%s`**) %s", instrument.Symbol, instrument.Description)
		validResponses++
	}

	if validResponses == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No tickers available for your query: `%s`", query))
		return
	}

	s.ChannelMessageSend(m.ChannelID, response)
}
