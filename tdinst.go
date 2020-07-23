package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

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

type ShortListDiviends struct {
	DividendYield     json.Number `name:"Yield"`
	DividendPayAmount json.Number `name:"Dividend ($)"`
	DividendDate      string      `name:"Ex-Div Date"`
	DividendPayDate   string      `name:"Pay Date"`
}

func createEmbed(t reflect.Type, v reflect.Value) []*discordgo.MessageEmbedField {
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
		} else if c.Field(i).Name == "DividendPayAmount" {
			tagValue = fmt.Sprintf("$%s", tagValue)
		} else if strings.HasPrefix(c.Field(i).Name, "Div") && strings.HasSuffix(c.Field(i).Name, "Date") {
			tagValue = strings.Split(tagValue, " ")[0]
		}

		embedFields = append(embedFields, &discordgo.MessageEmbedField{
			Name:   tagName,
			Value:  tagValue,
			Inline: true,
		})
	}
	return embedFields
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

func getRequiredDividends(ticker string) (*ShortListDiviends, *string, *string, error) {
	allFundamentals, err := tds.GetInstrumentFundamentals(ticker)
	if err != nil {
		return nil, nil, nil, err
	}

	var shortList ShortListDiviends
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

	embed := &discordgo.MessageEmbed{
		Title:  fmt.Sprintf("[%s] %s Fundamentals", *tdaTicker, *tdaDescription),
		Fields: createEmbed(reflect.TypeOf(ShortListFundamentals{}), reflect.ValueOf(*shortList)),
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func sendDividends(s *discordgo.Session, m *discordgo.MessageCreate, mSplit []string) {
	if len(mSplit) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a ticker to look up dividend information for!")
		return
	}

	ticker := mSplit[1]

	shortList, tdaTicker, tdaDescription, err := getRequiredDividends(ticker)

	switch err {
	case nil:
		break

	case tda.FundamentalsEmpty:
		s.ChannelMessageSend(m.ChannelID, "That ticker could not be looked up unfortunately, try again.")
		return

	default:
		s.ChannelMessageSend(m.ChannelID, "Something went wrong while processing the dividend info, try again later.")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:  fmt.Sprintf("[%s] %s Dividend", *tdaTicker, *tdaDescription),
		Fields: createEmbed(reflect.TypeOf(ShortListDiviends{}), reflect.ValueOf(*shortList)),
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func sendCompanyName(s *discordgo.Session, m *discordgo.MessageCreate, mSplit []string) {
	if len(mSplit) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a ticker to look up dividend information for!")
		return
	}

	ticker := mSplit[1]
	_, tdaTicker, tdaDescription, err := getRequiredFundamentals(ticker)

	switch err {
	case nil:
		break

	case tda.FundamentalsEmpty:
		s.ChannelMessageSend(m.ChannelID, "That ticker could not be looked up unfortunately, try again.")
		return

	default:
		s.ChannelMessageSend(m.ChannelID, "Something went wrong while processing the company name, try again later.")
		return
	}

	messageToSend := fmt.Sprintf("**`%s`** is %s", *tdaTicker, *tdaDescription)
	s.ChannelMessageSend(m.ChannelID, messageToSend)
}
