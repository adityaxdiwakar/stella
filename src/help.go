package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func help(s *discordgo.Session, m *discordgo.MessageCreate, mSplit []string) {
	helpMessages := map[string]string{
		"menu": fmt.Sprintf("The following menus are available, use them by" +
			"doing `?help <menu>`\n```charts  | help for equity charts\nfutures | " +
			"help for futures charts\nforex   | help for forex charts\ncompany | help" +
			"for ?div, ?earnings, and ?fun\nlookup  | help for ?name/?search\ntags    | " +
			"help for stella tags```"),
		"charts":  "https://i.imgur.com/T7ywnMv.png",
		"futures": "https://i.imgur.com/Yg4WKol.png",
		"forex":   "https://i.imgur.com/Uo8CZrA.png",
		"company": "https://i.imgur.com/1VI8nOD.png",
		"lookup":  "https://i.imgur.com/UGBE1GB.png",
		"tags":    "https://i.imgur.com/pEbmevO.png",
	}

	var response string
	if len(mSplit) == 1 {
		response = helpMessages["menu"]
	} else {
		response = helpMessages[strings.ToLower(mSplit[1])]
	}

	if response == "" {
		s.ChannelMessageSend(m.ChannelID, "Sorry, that help command is not avaiable, see `?help` for the menu")
		return
	}

	s.ChannelMessageSend(m.ChannelID, response)
}
