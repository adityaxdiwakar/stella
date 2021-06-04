package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/adityaxdiwakar/flux"
	"github.com/bwmarrin/discordgo"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

func quoteTicker(s *discordgo.Session, m *discordgo.MessageCreate, mSplit []string, count int) {
	start := time.Now()
	if len(mSplit) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a ticker to search")
		return
	}

	erroredTickers := []string{}
	tickers := mSplit[1:]
	if len(tickers) > 5 {
		s.ChannelMessageSend(m.ChannelID, "Please only put up to 5 symbols to quote")
		return
	}

	quoteChannel := make(chan *flux.QuoteStoredCache, 5)

	go func(tickers []string, quoteChan chan *flux.QuoteStoredCache) {
		for _, ticker := range tickers {
			searchResponse, err := fluxS.RequestQuote(flux.QuoteRequestSignature{
				Ticker:      ticker,
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
			if err != nil || len(searchResponse.Items) == 0 {
				erroredTickers = append(erroredTickers, strings.ToUpper(ticker))
				continue
			}

			quoteChan <- searchResponse
		}
	}(tickers, quoteChannel)

	photoFile, err := os.Open(fmt.Sprintf("Quote %dx.png", len(tickers)))
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

	ttf, err := ioutil.ReadFile("OpenSans-Regular.ttf")
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
	quoteResponses := []*flux.QuoteStoredCache{}
	heights := []int{69, 121, 173, 225, 277}
	for {
		if len(erroredTickers)+len(quoteResponses) == len(tickers) {
			break
		}

		quoteResponse := <-quoteChannel
		quoteResponses = append(quoteResponses, quoteResponse)
		addRow(quoteResponse, heights[len(quoteResponses)-1], c, font)
	}
	json.NewEncoder(os.Stdout).Encode(quoteResponses)

	buff := new(bytes.Buffer)
	png.Encode(buff, img)

	file := discordgo.File{
		Name:   "file.png",
		Reader: bytes.NewReader(buff.Bytes()),
	}

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{&file},
	})

	fmt.Println(time.Now().Sub(start))
}

func addRow(quote *flux.QuoteStoredCache, height int, c *freetype.Context, font *truetype.Font) {
	payload := quote.Items[0].Values
	fields := []string{
		quote.Items[0].Symbol,
		printer.Sprintf("%.2f", payload.MARK),
		printer.Sprintf("%.2f (%.2f%%)", payload.MARKCHANGE, payload.MARKPERCENTCHANGE*100),
		printer.Sprintf("%.2f", payload.BID),
		printer.Sprintf("%.2f", payload.ASK),
		printer.Sprintf("%d", payload.VOLUME),
	}
	offsets := []int{75, 220, 394, 568, 698, 875}
	primaryColor := getColor(payload.MARKCHANGE)
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
