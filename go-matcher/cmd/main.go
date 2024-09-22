package main

import (
	"fmt"
	"go-matcher/config"
	"go-matcher/services/betting"
	OddsApi "go-matcher/services/fetch"
	"go-matcher/types"
	"log"
	"math"
	"slices"
	"strings"
)

func main() {
	fmt.Println("Setup Rollin")

	payload, err := OddsApi.FetchSportsOdds("soccer_epl", "")
	if err != nil {
		fmt.Println("Error fetching sports odds", err)
	}

	var bestOutcomes []types.BestOutcome
	isFirst := true
	absoluteBestOutcome := types.BestOutcome{}

	// Loop over each sport event in the payload
	for _, value := range payload {
		var betfairMarketsLayOdds types.SportsOddsMarket

		// Find betfair exchange au odds
		for _, bookmaker := range value.Bookmakers {
			if bookmaker.Key == "betfair_ex_au" {
				for i := range bookmaker.Markets {
					if bookmaker.Markets[i].Key == "h2h_lay" {
						betfairMarketsLayOdds = bookmaker.Markets[i]
					}
				}
			}
		}

		// Match up odds from Bookie and Betfair
		for _, bookmaker := range value.Bookmakers {
			bestOutcome, err := handleCheck(betfairMarketsLayOdds, bookmaker)
			if err != nil {
				log.Fatal("unable to find best outcome", err)
			}
			if bestOutcome.ConversionRate == types.ErrConversionRate {
				// Continue to next iteration if we don't have a valid outcome.
				continue
			}

			bestOutcomes = append(bestOutcomes, bestOutcome)
		}
	}

	for _, outcome := range bestOutcomes {
		if isFirst {
			absoluteBestOutcome = types.BestOutcome{
				Bookmaker:         outcome.Bookmaker,
				Betfair:           outcome.Betfair,
				ConversionRate:    outcome.ConversionRate,
				Outcome:           outcome.Outcome,
				Probability:       outcome.Probability,
				BestBackStakeSize: outcome.BestBackStakeSize,
				BestLayStakeSize:  outcome.BestLayStakeSize,
				NetOutcome:        outcome.NetOutcome,
			}
			isFirst = false
			continue
		}

		// Check using absolute 0
		//if math.Abs(outcome.BookieWins) < math.Abs(absoluteBestOutcome.BookieWins) {
		//	absoluteBestOutcome = types.BestOutcome{
		//		Bookmaker:         outcome.Bookmaker,
		//		Betfair:           outcome.Betfair,
		//		ConversionRate:    outcome.ConversionRate,
		//		Outcome:           outcome.Outcome,
		//		Probability:       outcome.Probability,
		//		BestBackStakeSize: outcome.BestBackStakeSize,
		//		BestLayStakeSize:  outcome.BestLayStakeSize,
		//		NetOutcome:        outcome.NetOutcome,
		//	}
		//}
	}

	fmt.Println("absoluteBestOutcome", absoluteBestOutcome)
}

func handleCheck(betfairMarketsLayOdds types.SportsOddsMarket, bookmaker types.SportsOddsBookmaker) (types.BestOutcome, error) {
	bookieStake := 25.0 // default stake val
	minBet := 15.0
	maxBet := 30.0
	step := 2.5
	//betfairStake := 20
	//maxStakeDiff := 20 // Max difference
	isFirst := true
	bestOutcome := types.BestOutcome{}
	betfairCommission := 0.05
	bookies := strings.Split(config.Envs.OddsApiBookies, ",") // Extracts a slice of strings representing each bookie.

	// If bookie doesn't exist
	if !slices.Contains(bookies, bookmaker.Key) {
		return types.BestOutcome{
			ConversionRate: types.ErrConversionRate,
		}, nil
	}

	// Loop over each market for the bookie
	for _, market := range bookmaker.Markets {
		// Loop over each outcome (Win price team 1, Draw, Win Price Team 2)
		for _, outcome := range market.Outcomes {
			// For each outcome, find the equivalent odds outcome on Betfair
			for _, betfairLayOddsOutcomes := range betfairMarketsLayOdds.Outcomes {
				// Matching outcome name means we can use as a comparison.
				if betfairLayOddsOutcomes.Name == outcome.Name {
					backOdds := outcome.Price
					layOdds := betfairLayOddsOutcomes.Price

					// TODO Calculate optimum bet on bookie and betfair
					// Testing
					bestBookieStake, bestExchangeStake, bestOutcomeFloatVal := betting.FindOptimalArbitrageBets(minBet, maxBet, step, backOdds, layOdds, 0.05)

					//fmt.Printf("Optimal bets:\n")
					//fmt.Printf("Platform 1 bestBookieStake: %.2f\n", bestBookieStake)
					//fmt.Printf("Platform 2 bestExchangeStake: %.2f\n", bestExchangeStake)
					//fmt.Printf("Net outcome: %.2f\n", bestOutcomeFloatVal)
					// End Test

					// This is lay odds calculation, hence why the math looks a little strange compared to normal bet profit calculation.
					// Since we are just making stake amount - commission.
					profit := bookieStake - (bookieStake * betfairCommission)

					conversionRate := 100 - (backOdds/layOdds)*profit

					// Filter out favourites (not perfect way), OR matches with close odds
					if backOdds < 5 || (layOdds-backOdds) > 5 {
						continue
					}

					// Filter out events that are likely to succeed.
					probabilityOfOutcome := (1 / backOdds) * 100
					if probabilityOfOutcome > 20 {
						continue
					}

					if isFirst {
						bestOutcome = types.BestOutcome{
							Bookmaker:         bookmaker,
							Betfair:           betfairMarketsLayOdds,
							ConversionRate:    conversionRate,
							Outcome:           outcome,
							Probability:       probabilityOfOutcome,
							BestBackStakeSize: bestBookieStake,
							BestLayStakeSize:  bestExchangeStake,
							NetOutcome:        bestOutcomeFloatVal,
						}
						isFirst = false
						continue
					}

					if math.Abs(bestOutcomeFloatVal) < math.Abs(bestOutcome.NetOutcome) {
						bestOutcome = types.BestOutcome{
							Bookmaker:         bookmaker,
							ConversionRate:    conversionRate,
							Outcome:           outcome,
							Probability:       probabilityOfOutcome,
							BestBackStakeSize: bestBookieStake,
							BestLayStakeSize:  bestExchangeStake,
							NetOutcome:        bestOutcomeFloatVal,
						}
					}
				}
			}
		}
	}

	fmt.Println("Best Outcome ", bestOutcome)
	return bestOutcome, nil
}
