package main

import (
	"fmt"
	"log"
	"time"

	"github.com/adityaxdiwakar/flux"
	"github.com/bwmarrin/discordgo"
)

func channelTicker(s *discordgo.Session) {
	for range time.Tick(307 * time.Second) {
		price, change, percentage, err := getFuturesData()
		if err != nil {
			log.Println("[ticker] Could not retrieve futures data for profile status")
			fmt.Println(err)
			continue
		}
		message := fmt.Sprintf("%.2f (%s, %s)", *price, *change, *percentage)

		// set for every relevant server
		for _, channelID := range []string{
			"703080609358020608",
			"709860290694742098",
		} {
			_, err := s.ChannelEdit(channelID, message)
			if err != nil {
				log.Printf("[ticker] Could not change channel title for %s due to: %v\n", channelID, err)
			}
		}
	}
}

func playingTicker(s *discordgo.Session) {
	for range time.Tick(13 * time.Second) {
		price, _, percentage, err := getFuturesData()
		if err != nil {
			log.Println("[ticker] Could not retrieve futures data for profile status")
			continue
		}
		message := fmt.Sprintf("%.2f (%s)", *price, *percentage)
		s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Game: &discordgo.Game{
				Name: message,
				Type: discordgo.GameTypeWatching,
			},
		})
	}
}

func getFuturesData() (*float64, *string, *string, error) {
	sig := flux.QuoteRequestSignature{
		Ticker:      "/ES",
		RefreshRate: 300,
		Fields: []flux.QuoteField{
			flux.Mark,
			flux.MarkChange,
			flux.MarkPercentChange,
		},
	}

	payload, err := fluxS.RequestQuote(sig)
	if err != nil {
		return nil, nil, nil, err
	}

	price := payload.Items[0].Values.MARK
	change := payload.Items[0].Values.MARKCHANGE
	percent := payload.Items[0].Values.MARKPERCENTCHANGE * 100

	var sChange string
	var sPercent string
	if change > 0 {
		sChange = fmt.Sprintf("+%.2f", change)
		sPercent = fmt.Sprintf("+%.2f%%", percent)
	} else {
		sChange = fmt.Sprintf("%.2f", change)
		sPercent = fmt.Sprintf("%.2f%%", percent)
	}

	return &price, &sChange, &sPercent, nil
}
