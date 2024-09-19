package types

import "time"

// SingleTon
var ErrConversionRate = 717172.1234

type SportsOddsPayload struct {
	Id           string                `json:"id"`
	SportKey     string                `json:"sport_key"`
	SportTitle   string                `json:"sport_title"`
	CommenceTime time.Time             `json:"commence_time"`
	HomeTeam     string                `json:"home_team"`
	AwayTeam     string                `json:"away_team"`
	Bookmakers   []SportsOddsBookmaker `json:"bookmakers"`
}

type SportsOddsBookmaker struct {
	Key        string             `json:"key"`
	Title      string             `json:"title"`
	LastUpdate time.Time          `json:"last_update"`
	Markets    []SportsOddsMarket `json:"markets"`
}

type SportsOddsMarket struct {
	Key        string    `json:"key"`
	LastUpdate time.Time `json:"last_update"`
	Outcomes   []Outcome `json:"outcomes"`
}

type Outcome struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type BestOutcome struct {
	Bookmaker         SportsOddsBookmaker
	Betfair           SportsOddsMarket
	BookieWins        float64
	BetfairWins       float64
	ConversionRate    float64
	Outcome           Outcome
	Probability       float64
	BestBackStakeSize float64
	BestLayStakeSize  float64
	NetOutcome        float64
}
