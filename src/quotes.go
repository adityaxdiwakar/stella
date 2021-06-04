package main

import (
	"encoding/json"
	"image"
	"os"
	"strings"

	"github.com/adityaxdiwakar/flux"
	"github.com/bwmarrin/discordgo"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

func quoteTicker(s *discordgo.Session, m *discordgo.MessageCreate, mSplit []string, count int) {
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

	quoteResponses := []*flux.QuoteStoredCache{}
	for {
		if len(erroredTickers)+len(quoteResponses) == len(tickers) {
			break
		}

		quoteResponse := <-quoteChannel
		quoteResponses = append(quoteResponses, quoteResponse)
	}
	json.NewEncoder(os.Stdout).Encode(quoteResponses)

	/*
		photoFile, err := os.Open("Quote 1x.png")
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

		start := time.Now()
		c := freetype.NewContext()
		size := 24.0

		greenColor := color.RGBA{66, 169, 43, 255}
		redColor := color.RGBA{255, 51, 51, 255}

		var primaryColor image.Image
		if payload.MARKCHANGE > 0 {
			primaryColor = image.NewUniform(greenColor)
		} else if payload.MARKCHANGE < 0 {
			primaryColor = image.NewUniform(redColor)
		} else {
			primaryColor = image.White
		}

		c.SetDPI(72)
		c.SetFont(font)
		c.SetClip(img.Bounds())
		c.SetSrc(image.White)
		c.SetFontSize(size)
		c.SetDst(img)

		// add ticker
		addLabel(tickerStr, c, font, 65, 75, image.White, 24.0, 140)

		// add mark
		addLabel(markStr, c, font, 65, 220, primaryColor, 24.0, 500)

		// mark change
		addLabel(markChStr, c, font, 65, 394, primaryColor, 24.0, 500)

		// bid string
		bidStr := printer.Sprintf("%.2f", payload.BID)
		addLabel(bidStr, c, font, 65, 568, primaryColor, 24.0, 500)

		// ask string
		addLabel(askStr, c, font, 65, 698, primaryColor, 24.0, 500)

		// volume
		addLabel(volumeStr, c, font, 65, 875, image.White, 24.0, 500)

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
	*/
}

func addRow(quote *flux.QuoteStoredCache, height int, primaryColor image.Image, c *freetype.Context, font *truetype.Font) {
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
	colors := []image.Image{image.White, primaryColor, primaryColor, primaryColor, image.White}
	widths := []int{140, 500, 500, 500, 500, 500}

	for i, fieldStr := range fields {
		addLabel(fieldStr, c, font, 65, offsets[i], colors[i], 24.0, widths[i])
	}
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
