package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	tda "github.com/adityaxdiwakar/tda-go"
	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/copier"
)

type ShortListFundamentals struct {
	PeRatio           json.Number `name:"Price-Earnings Ratio"`
	PbRatio           json.Number `name:"Price-Book Ratio"`
	TotalDebtToEquity json.Number `name:"Debt-Equity Ratio"`
	PegRatio          json.Number `name:"PEG Ratio"`
	Beta              json.Number `name:"Beta"`
	DividendYield     json.Number `name:"Div Yield"`
	ReturnOnEquity    json.Number `name:"Equity Return"`
	MarketCap         json.Number `name:"Market Cap"`
}

func getRequiredFundamentals(ticker string) (*ShortListFundamentals, *string, *string, error) {
	allFundamentals, err := tds.GetInstrumentFundamentals(ticker)
	if err != nil {
		return nil, nil, nil, err
	}

	var shortList ShortListFundamentals
	copier.Copy(&shortList, &allFundamentals.Fundamental)

	return &shortList, &allFundamentals.Symbol, &allFundamentals.Description, nil
}

func sendFundamentals(s *discordgo.Session, m *discordgo.MessageCreate, mSplit []string) {
	if len(mSplit) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a ticker to look up fundamentals for!")
		return
	}

	ticker := mSplit[1]

	shortList, tdaTicker, tdaDescription, err := getRequiredFundamentals(ticker)

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

	t := reflect.TypeOf(ShortListFundamentals{})
	v := reflect.ValueOf(*shortList)
	c := v.Type()

	var embedFields []*discordgo.MessageEmbedField
	for i := 0; i < v.NumField(); i++ {
		f, _ := t.FieldByName(c.Field(i).Name)
		tagName, _ := f.Tag.Lookup("name")
		if v.Field(i).String() == "0.0" {
			continue
		}

		tagValue := v.Field(i).String()
		if c.Field(i).Name == "MarketCap" {
			cap, err := strconv.ParseFloat(tagValue, 64)
			if err == nil {
				tagValue = fmt.Sprintf("$%.2fB", float64(cap)/1000)
			}
		} else if c.Field(i).Name == "DividendYield" {
			tagValue = fmt.Sprintf("%s%%", tagValue)
		}

		embedFields = append(embedFields, &discordgo.MessageEmbedField{
			Name:   tagName,
			Value:  tagValue,
			Inline: true,
		})
	}

	embed := &discordgo.MessageEmbed{
		Title:  fmt.Sprintf("[%s] %s", *tdaTicker, *tdaDescription),
		Fields: embedFields,
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
