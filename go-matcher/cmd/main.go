package main

import (
	"fmt"
	"go-matcher/config"
	OddsApi "go-matcher/services/fetch"
	"go-matcher/types"
	"slices"
	"strings"
)

func main() {
	fmt.Println("Setup Rollin")

	payload, err := OddsApi.FetchSportsOdds("soccer_epl", "")
	if err != nil {
		fmt.Println("Error fetching sports odds", err)
	}

	// Do payload comparison
	for index, value := range payload {
		//
		fmt.Println("Index: ", index)

		//
		var betfairMarketsLayOdds types.SportsOddsMarket

		// Find betfair exchange au odds
		for bookmakerIndex, bookmaker := range value.Bookmakers {
			fmt.Println("Bookmaker index: ", bookmakerIndex)

			if bookmaker.Key == "betfair_ex_au" {
				for i := range bookmaker.Markets {
					if bookmaker.Markets[i].Key == "h2h_lay" {
						betfairMarketsLayOdds = bookmaker.Markets[i]
					}
				}
			}
		}

		// Compare Lay odds
		for _, bookmaker := range value.Bookmakers {
			handleCheck(betfairMarketsLayOdds, bookmaker)
		}
	}
}

func handleCheck(betfairMarketsLayOdds types.SportsOddsMarket, bookmaker types.SportsOddsBookmaker) {
	stake := 10.0 // default stake val
	betfairCommission := 0.05
	bookies := strings.Split(config.Envs.OddsApiBookies, ",")

	if !slices.Contains(bookies, bookmaker.Key) {
		return
	}

	for _, market := range bookmaker.Markets {
		for _, outcome := range market.Outcomes {
			// For each outcome, find the equivalent
			for _, betfairLayOddsOutcomes := range betfairMarketsLayOdds.Outcomes {
				if betfairLayOddsOutcomes.Name == outcome.Name {
					// Matching name outcomes from betfair and bookie
					backOdds := outcome.Price
					layOdds := betfairLayOddsOutcomes.Price
					// This is lay odds calculation, hence why the math looks a little strange compared to normal bet profit calculation.
					// Since we are just making stake amount - commission.
					profit := stake - (stake * betfairCommission)

					conversionRate := 100 - (backOdds/layOdds)*profit
					fmt.Println("============ Conversion Percent ============")
					fmt.Println("Conversion Rate: ", conversionRate)
					fmt.Println("Market: ", market.Key)
					fmt.Println("Bookie: ", bookmaker.Key)
					fmt.Println("Bookie Market Outcome Key: ", outcome.Name)
					fmt.Println("Betfair Market Outcome Key:", betfairLayOddsOutcomes.Name)
					fmt.Println("back odds (Bookie): ", backOdds)
					fmt.Println("Lay odds (Betfair): ", layOdds)
				}
			}
		}
	}

	// Do calculation

	// Get back odds

	// Get lay odds

	// Betfair commission

	//
}
