package main

import (
	"fmt"
	"go-matcher/config"
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
		fmt.Println("Outcome: ", outcome)
		if isFirst {
			fmt.Println("IS FIRST!!")
			absoluteBestOutcome = types.BestOutcome{
				Bookmaker:      outcome.Bookmaker,
				Betfair:        outcome.Betfair,
				BookieWins:     outcome.BookieWins,
				BetfairWins:    outcome.BetfairWins,
				ConversionRate: outcome.ConversionRate,
				Outcome:        outcome.Outcome,
				Probability:    outcome.Probability,
			}
			isFirst = false
			continue
		}

		// Check using absolute 0
		if math.Abs(outcome.BookieWins) < math.Abs(absoluteBestOutcome.BookieWins) {
			absoluteBestOutcome = types.BestOutcome{
				Bookmaker:      outcome.Bookmaker,
				Betfair:        outcome.Betfair,
				BookieWins:     outcome.BookieWins,
				BetfairWins:    outcome.BetfairWins,
				ConversionRate: outcome.ConversionRate,
				Outcome:        outcome.Outcome,
				Probability:    outcome.Probability,
			}
		}
	}

	fmt.Println("Absolute Best Outcome: ", absoluteBestOutcome)
	fmt.Println("Absolute Best Bookmaker: ", absoluteBestOutcome.Bookmaker)
	fmt.Println("Absolute Best: ", absoluteBestOutcome.Outcome)
	fmt.Println("Absolute Best BetfairWins: ", absoluteBestOutcome.BetfairWins)
	fmt.Println("Absolute Best BookieWins: ", absoluteBestOutcome.BookieWins)

	// Bookie
	fmt.Println("")
}

func handleCheck(betfairMarketsLayOdds types.SportsOddsMarket, bookmaker types.SportsOddsBookmaker) (types.BestOutcome, error) {
	bookieStake := 25.0 // default stake val
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
			// TODO: Calculate odds
			for _, betfairLayOddsOutcomes := range betfairMarketsLayOdds.Outcomes {
				// Matching outcome name means we can use as a comparison.
				if betfairLayOddsOutcomes.Name == outcome.Name {
					backOdds := outcome.Price
					layOdds := betfairLayOddsOutcomes.Price
					risk := layOdds * bookieStake
					// This is lay odds calculation, hence why the math looks a little strange compared to normal bet profit calculation.
					// Since we are just making stake amount - commission.
					profit := bookieStake - (bookieStake * betfairCommission)

					conversionRate := 100 - (backOdds/layOdds)*profit

					// If the lay bet fails:
					bookieWins := ((bookieStake * backOdds) - bookieStake) - risk
					// If the lay bet wins:
					betfairWins := profit - bookieStake // subtract the loss from bookie stake

					// Filter out favourites (not perfect way), OR matches with close odds
					if backOdds < 5 || (layOdds-backOdds) > 5 {
						continue
					}

					// Filter out events that are likely to succeed.
					probabilityOfOutcome := (1 / backOdds) * 100
					if probabilityOfOutcome > 20 {
						continue
					}

					// Ideally, we want both bookieWins + betfairWins to be as low as possible,
					// that way the conversion process loss is at a minimum.

					//fmt.Println("============ Conversion Percent ============")
					//fmt.Println("Conversion Rate: ", conversionRate)
					//fmt.Println("bookieWins: ", bookieWins)
					//fmt.Println("Betfair Wins: ", betfairWins)
					//fmt.Println("Market: ", market.Key)
					//fmt.Println("Bookie: ", bookmaker.Key)
					//fmt.Println("Bookie Market Outcome Key: ", outcome.Name)
					//fmt.Println("Betfair Market Outcome Key:", betfairLayOddsOutcomes.Name)
					//fmt.Println("back odds (Bookie): ", backOdds)
					//fmt.Println("Lay odds (Betfair): ", layOdds)

					if isFirst {
						//fmt.Println("======== Is First =========")
						bestOutcome = types.BestOutcome{
							Bookmaker:      bookmaker,
							Betfair:        betfairMarketsLayOdds,
							BookieWins:     bookieWins,
							BetfairWins:    betfairWins,
							ConversionRate: conversionRate,
							Outcome:        outcome,
							Probability:    probabilityOfOutcome,
						}
						isFirst = false
						continue
					}

					if math.Abs(bookieWins) < math.Abs(bestOutcome.BookieWins) {
						bestOutcome = types.BestOutcome{
							Bookmaker:      bookmaker,
							BookieWins:     bookieWins,
							BetfairWins:    betfairWins,
							ConversionRate: conversionRate,
							Outcome:        outcome,
							Probability:    probabilityOfOutcome,
						}
					}
				}
			}
		}
	}

	// If both are still set to 0, we can assume a failure has occurred somewhere.
	if bestOutcome.BetfairWins == 0 && bestOutcome.BookieWins == 0 {
		return types.BestOutcome{
			ConversionRate: types.ErrConversionRate,
		}, nil
	}

	return bestOutcome, nil
}
