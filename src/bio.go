package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ReutersResponse struct {
	Ts         int        `json:"ts"`
	MarketData MarketData `json:"market_data"`
	Ric        string     `json:"ric"`
	Status     Status     `json:"status"`
}
type LocalName struct {
	Lang string `json:"lang"`
	Name string `json:"name"`
}
type SigDevs struct {
	DevelopmentID string `json:"development_id"`
	LastUpdate    string `json:"last_update"`
	Headline      string `json:"headline"`
	Description   string `json:"description"`
}
type Phone struct {
	Type             string `json:"type"`
	CountryPhoneCode string `json:"country_phone_code"`
	CityAreaCode     string `json:"city_area_code"`
	Number           string `json:"number"`
}
type Officers struct {
	Name  string `json:"name"`
	Rank  int    `json:"rank"`
	Title string `json:"title"`
	Age   string `json:"age"`
	Since int    `json:"since"`
}
type Recommendation struct {
	UnverifiedMean          float64 `json:"unverified_mean"`
	PreliminaryMean         float64 `json:"preliminary_mean"`
	Mean                    float64 `json:"mean"`
	High                    int     `json:"high"`
	Low                     int     `json:"low"`
	NumberOfRecommendations int     `json:"number_of_recommendations"`
}
type Data struct {
	FiscalYear int     `json:"fiscal_year"`
	Value      float64 `json:"value"`
	Estimate   bool    `json:"estimate"`
}
type EpsPerYear struct {
	Currency string `json:"currency"`
	Data     []Data `json:"data"`
}
type RevenuePerYear struct {
	Currency string `json:"currency"`
	Data     []Data `json:"data"`
}
type MarketData struct {
	Ric                          string         `json:"ric"`
	ExchangeName                 string         `json:"exchange_name"`
	Last                         string         `json:"last"`
	LastTime                     string         `json:"last_time"`
	NetChange                    string         `json:"net_change"`
	Currency                     string         `json:"currency"`
	PercentChange                string         `json:"percent_change"`
	Modified                     string         `json:"modified"`
	Volume                       string         `json:"volume"`
	DayHigh                      string         `json:"day_high"`
	DayLow                       string         `json:"day_low"`
	FiftytwoWkHigh               string         `json:"fiftytwo_wk_high"`
	FiftytwoWkLow                string         `json:"fiftytwo_wk_low"`
	PrevDayClose                 string         `json:"prev_day_close"`
	Open                         string         `json:"open"`
	CompanyName                  string         `json:"company_name"`
	FundamentalExchangeName      string         `json:"fundamental_exchange_name"`
	LocalName                    []LocalName    `json:"local_name"`
	MarketCap                    string         `json:"market_cap"`
	ShareVolume3M                string         `json:"share_volume_3m"`
	Beta                         string         `json:"beta"`
	EpsExclExtraTtm              string         `json:"eps_excl_extra_ttm"`
	PeExclExtraTtm               string         `json:"pe_excl_extra_ttm"`
	PsAnnual                     string         `json:"ps_annual"`
	PsTtm                        string         `json:"ps_ttm"`
	PcfShareTtm                  string         `json:"pcf_share_ttm"`
	PbAnnual                     string         `json:"pb_annual"`
	PbQuarterly                  string         `json:"pb_quarterly"`
	DividendYieldIndicatedAnnual string         `json:"dividend_yield_indicated_annual"`
	LtDebtEquityAnnual           string         `json:"lt_debt_equity_annual"`
	TotalDebtEquityAnnual        string         `json:"total_debt_equity_annual"`
	LtDebtEquityQuarterly        string         `json:"lt_debt_equity_quarterly"`
	TotalDebtEquityQuarterly     string         `json:"total_debt_equity_quarterly"`
	SharesOut                    string         `json:"shares_out"`
	RoeTtm                       string         `json:"roe_ttm"`
	RoiTtm                       string         `json:"roi_ttm"`
	SigDevs                      []SigDevs      `json:"sig_devs"`
	About                        string         `json:"about"`
	AboutJp                      string         `json:"about_jp"`
	Website                      string         `json:"website"`
	StreetAddress                []string       `json:"street_address"`
	City                         string         `json:"city"`
	State                        string         `json:"state"`
	PostalCode                   string         `json:"postal_code"`
	Country                      string         `json:"country"`
	Phone                        Phone          `json:"phone"`
	Sector                       string         `json:"sector"`
	Industry                     string         `json:"industry"`
	ForwardPE                    string         `json:"forward_PE"`
	Officers                     []Officers     `json:"officers"`
	Recommendation               Recommendation `json:"recommendation"`
	EpsPerYear                   EpsPerYear     `json:"eps_per_year"`
	RevenuePerYear               RevenuePerYear `json:"revenue_per_year"`
	NextEvent                    interface{}    `json:"next_event"`
}
type Status struct {
	Code int `json:"code"`
}

func exchangeUrl(ticker, exchange string) string {
	return fmt.Sprintf("https://www.reuters.com/companies/api/getFetchCompanyProfile/%s%s", ticker, exchange)
}

var ErrRequestPkgFailed = errors.New(ErrReutersPkgFailed)
var ErrResponseFailed = errors.New(ErrReutersResponseFailed)

func getReutersData(ticker, exchange string) (*ReutersResponse, error) {
	req, err := http.NewRequest("GET", exchangeUrl(ticker, exchange), nil)
	if err != nil {
		return nil, ErrRequestPkgFailed
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4181.8 Safari/537.36")
	req.Header.Set("Host", "www.reuters.com")

	response := ReutersResponse{}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, ErrResponseFailed
	} else if resp.StatusCode == 206 {
		response := ReutersResponse{}
		json.NewDecoder(resp.Body).Decode(&response)

		if response.MarketData.CompanyName != "" && response.MarketData.About != "" {
			return &response, nil
		}

		return nil, ErrResponseFailed
	} else if resp.StatusCode != 200 {
		return nil, ErrResponseFailed
	}

	json.NewDecoder(resp.Body).Decode(&response)

	return &response, nil
}

func repeatReutersRequest(ticker, exchange string, ch chan ReutersResponse) {
	for i := 0; i < 3; i++ {
		r, err := getReutersData(ticker, exchange)
		if err != nil {
			time.Sleep(3000 * time.Millisecond)
			continue
		}
		ch <- *r
	}
}

func reutersBio(s *discordgo.Session, m *discordgo.MessageCreate, mSplit []string) {
	if len(mSplit) < 2 {
		s.ChannelMessageSend(m.ChannelID, ErrNoTickerProvided)
		return
	}

	ticker := mSplit[1]

	loadingMsg, _ := s.ChannelMessageSend(m.ChannelID, LoadingCompanyBio(ticker))

	reutersChannel := make(chan ReutersResponse)

	for _, exchange := range []string{".O", ".N", ""} {
		go repeatReutersRequest(ticker, exchange, reutersChannel)
	}

	select {

	case reutersData := <-reutersChannel:
		companyName := reutersData.MarketData.CompanyName
		aboutCompany := reutersData.MarketData.About

		if companyName == "" || aboutCompany == "" {
			s.ChannelMessageEdit(loadingMsg.ChannelID, loadingMsg.ID, CantFindCompanyBio)
			return
		}

		embed := &discordgo.MessageEmbed{
			Title:       companyName,
			Description: aboutCompany,
		}

		emptyContent := ""
		s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Content: &emptyContent,
			Embed:   embed,
			ID:      loadingMsg.ID,
			Channel: loadingMsg.ChannelID,
		})

		return

	case <-time.After(10 * time.Second):
		s.ChannelMessageEdit(loadingMsg.ChannelID, loadingMsg.ID, ErroredCompanyBio)
		return
	}

}
