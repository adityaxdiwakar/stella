package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"strings"

	"github.com/adityaxdiwakar/flux"
	"github.com/bwmarrin/discordgo"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

func quoteTicker(s *discordgo.Session, m *discordgo.MessageCreate, mSplit []string, count int) {
	if !fluxS.Established {
		s.ChannelMessageSend(m.ChannelID, "Sorry, Stella is not connected to TDAmeritrade at this time")
		return
	}

	if len(mSplit) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a ticker to search")
		return
	}

	tickers, tickerMap := removeDuplicates(mSplit[1:])
	if len(tickers) > 10 {
		s.ChannelMessageSend(m.ChannelID, "Please only put up to 10 symbols to quote")
		return
	}

	quoteResponse, err := fluxS.RequestQuote(flux.QuoteRequestSignature{
		Ticker:      strings.Join(tickers, ","),
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
	if err != nil || len(quoteResponse.Items) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Could not quote any of the tickers requested")
		return
	}

	photoFile, err := os.Open(fmt.Sprintf("assets/images/quotes/Quote %dx.png", len(tickers)))
	if err != nil {
		fmt.Println(err)
		return
	}

	i, _, err := image.Decode(photoFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	img := i.(*image.NRGBA)

	ttf, err := ioutil.ReadFile("assets/fonts/OpenSans-Regular.ttf")
	if err != nil {
		fmt.Println(err)
		return
	}

	font, err := truetype.Parse(ttf)

	c := freetype.NewContext()
	size := 24.0

	c.SetDPI(72)
	c.SetFont(font)
	c.SetClip(img.Bounds())
	c.SetFontSize(size)
	c.SetDst(img)
	heights := []int{69, 121, 173, 225, 277, 329, 381, 433, 485, 537}

	errorCount := 0
	for i, quote := range quoteResponse.Items {
		addRow(quote, heights[i], c, font)
		tickerMap[quote.Symbol] = false
		errorCount = i
	}

	for key, val := range tickerMap {
		if val {
			addLabel(key, c, font, heights[errorCount+1], 75, image.White, 24.0, 140)
			addLabel("Could not load data for this ticker, try again?", c, font, heights[errorCount+1],
				549, image.White, 24.0, 500)
			errorCount += 1
		}
	}

	buff := new(bytes.Buffer)
	png.Encode(buff, img)

	file := discordgo.File{
		Name:   "file.png",
		Reader: bytes.NewReader(buff.Bytes()),
	}

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{&file},
	})

}

func addRow(quote flux.QuoteItem, height int, c *freetype.Context, font *truetype.Font) {
	payload := quote.Values

	last := payload.MARK
	change := payload.MARKCHANGE
	changePercent := payload.MARKPERCENTCHANGE
	if change == 0 && payload.NETCHANGE != 0 {
		change = payload.NETCHANGE
		last = payload.LAST
		changePercent = payload.NETCHANGEPERCENT
	}

	fields := []string{
		quote.Symbol,
		printer.Sprintf("%.2f", last),
		printer.Sprintf("%.2f (%.2f%%)", change, changePercent*100),
		printer.Sprintf("%.2f", payload.BID),
		printer.Sprintf("%.2f", payload.ASK),
		printer.Sprintf("%d", payload.VOLUME),
	}
	offsets := []int{75, 220, 394, 568, 698, 875}
	primaryColor := getColor(change)
	colors := []image.Image{image.White, primaryColor, primaryColor, primaryColor, primaryColor, image.White}
	widths := []int{140, 500, 500, 500, 500, 500}

	for i, fieldStr := range fields {
		addLabel(fieldStr, c, font, height, offsets[i], colors[i], 24.0, widths[i])
	}
}

func getColor(data float64) image.Image {
	if data > 0 {
		return tdaGreen
	} else if data < 0 {
		return tdaRed
	}
	return image.White
}

func addLabel(str string, c *freetype.Context, font *truetype.Font, height int, offset int, color image.Image, size float64, space int) {
	width, nSize := calculateWidth(str, c, font, size, space)
	pt := freetype.Pt(offset-(width/2).Round(), height+int(c.PointToFixed(size)>>6))
	c.SetFontSize(nSize)
	c.SetSrc(color)
	c.DrawString(str, pt)

}

func calculateWidth(str string, c *freetype.Context, font *truetype.Font, size float64, space int) (fixed.Int26_6, float64) {
	face := truetype.NewFace(font, &truetype.Options{
		Size: size,
	})

	twidth := c.PointToFixed(0)
	for _, ch := range str {
		awidth, _ := face.GlyphAdvance(rune(ch))
		twidth += awidth
	}

	if twidth.Floor() > (space - 15) {
		return calculateWidth(str, c, font, float64(space-15)/float64(twidth.Floor())*float64(size), space)
	}

	return twidth, size
}

func removeDuplicates(dupeSlice []string) ([]string, map[string]bool) {
	keys := make(map[string]bool)
	output := []string{}
	for _, entry := range dupeSlice {
		if _, value := keys[entry]; !value {
			entryUpper := strings.ToUpper(entry)
			keys[entryUpper] = true
			output = append(output, entryUpper)
		}
	}

	return output, keys
}
