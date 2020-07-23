package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func channelTicker(s *discordgo.Session) {
	for range time.Tick(300 * time.Second) {
		price, change, percentage, err := getFuturesData()
		if err != nil {
			log.Println("[ticker] Could not retrieve futures data for profile status")
			continue
		}
		message := fmt.Sprintf("%.2f (%s, %s%%)", *price, *change, *percentage)
		for _, channelID := range []string{"703080609358020608", "709860290694742098"} {
			_, err := s.ChannelEdit(channelID, message)
			if err != nil {
				log.Printf("[ticker] Could not change channel title for %s due to: %v\n", channelID, err)
			}
		}
	}
}

func playingTicker(s *discordgo.Session) {
	for range time.Tick(15 * time.Second) {
		price, _, percentage, err := getFuturesData()
		if err != nil {
			log.Println("[ticker] Could not retrieve futures data for profile status")
			continue
		}
		message := fmt.Sprintf("%.2f (%s%%)", *price, *percentage)
		s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Game: &discordgo.Game{
				Name: message,
				Type: discordgo.GameTypeWatching,
			},
		})
	}
}

type FuturesPayload struct {
	Code    int `json:"code"`
	Payload struct {
		Timestamp     int64 `json:"timestamp"`
		ContractID    int   `json:"contract_id"`
		SessionVolume int   `json:"session_volume"`
		OpenInterest  int   `json:"open_interest"`
		SessionPrices struct {
			Open       float64 `json:"open"`
			High       float64 `json:"high"`
			Settlement float64 `json:"settlement"`
			Low        float64 `json:"low"`
		} `json:"session_prices"`
		Depth struct {
			Bid struct {
				Price float64 `json:"price"`
				Size  int     `json:"size"`
			} `json:"bid"`
			Ask struct {
				Price float64 `json:"price"`
				Size  int     `json:"size"`
			} `json:"ask"`
		} `json:"depth"`
		Trade struct {
			Price float64 `json:"price"`
			Size  int     `json:"size"`
		} `json:"trade"`
	} `json:"payload"`
}

func getFuturesData() (*float64, *string, *string, error) {
	res, err := stellaHttpClient.Get("https://md.aditya.diwakar.io/recent/")
	if err != nil {
		return nil, nil, nil, err
	}

	defer res.Body.Close()

	httpPayload := FuturesPayload{}
	err = json.NewDecoder(res.Body).Decode(&httpPayload)
	if err != nil {
		return nil, nil, nil, err
	}

	payload := httpPayload.Payload

	price := payload.Trade.Price
	numericalChange := price - payload.SessionPrices.Settlement
	numericalPercentage := numericalChange / payload.SessionPrices.Settlement * 100
	var change string
	if numericalChange >= 0 {
		change = fmt.Sprintf("+%.2f", numericalChange)
	} else {
		change = fmt.Sprintf("%.2f", numericalChange)
	}

	var percentage string
	if numericalPercentage >= 0 {
		percentage = fmt.Sprintf("+%.2f", numericalPercentage)
	} else {
		percentage = fmt.Sprintf("%.2f", numericalPercentage)
	}

	return &price, &change, &percentage, nil
}
