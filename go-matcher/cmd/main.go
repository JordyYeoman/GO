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

	// Do payload comparison
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

		// Compare Lay odds
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
	stake := 25.0 // default stake val
	isFirst := true
	bestOutcome := types.BestOutcome{}
	betfairCommission := 0.05
	bookies := strings.Split(config.Envs.OddsApiBookies, ",")
	minDifference := 4.5

	if !slices.Contains(bookies, bookmaker.Key) {
		return types.BestOutcome{
			ConversionRate: types.ErrConversionRate,
		}, nil
	}

	for _, market := range bookmaker.Markets {
		for _, outcome := range market.Outcomes {
			// For each outcome, find the equivalent
			for _, betfairLayOddsOutcomes := range betfairMarketsLayOdds.Outcomes {
				if betfairLayOddsOutcomes.Name == outcome.Name {
					// Matching name outcomes from betfair and bookie
					backOdds := outcome.Price
					layOdds := betfairLayOddsOutcomes.Price
					risk := layOdds * stake
					// This is lay odds calculation, hence why the math looks a little strange compared to normal bet profit calculation.
					// Since we are just making stake amount - commission.
					profit := stake - (stake * betfairCommission)

					conversionRate := 100 - (backOdds/layOdds)*profit

					// If the lay bet fails:
					bookieWins := ((stake * backOdds) - stake) - risk
					// If the lay bet wins:
					betfairWins := profit - stake // subtract the loss from bookie stake
					// Difference between odds
					diff := layOdds - backOdds
					isDiffHighEnough := diff > minDifference

					// Filter out favourites (not perfect way)
					if backOdds < 5 {
						continue
					}

					// Probability of outcome
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

					//fmt.Println("=========== T ++++++++++++")
					//fmt.Println(math.Abs(bookieWins))
					//fmt.Println(bookieWins)
					//fmt.Println(math.Abs(bestOutcome.BookieWins))
					//fmt.Println(bestOutcome.BookieWins)

					if math.Abs(bookieWins) < math.Abs(bestOutcome.BookieWins) && isDiffHighEnough {
						bestOutcome = types.BestOutcome{
							Bookmaker:      bookmaker,
							BookieWins:     bookieWins,
							BetfairWins:    betfairWins,
							ConversionRate: conversionRate,
							Outcome:        outcome,
						}
					}
				}
			}
		}
	}

	if bestOutcome.BetfairWins == 0 && bestOutcome.BookieWins == 0 {
		return types.BestOutcome{
			ConversionRate: types.ErrConversionRate,
		}, nil
	}

	//fmt.Println("Best Outcome: ", bestOutcome)
	//fmt.Println("Best Bookmaker: ", bestOutcome.Bookmaker)
	//fmt.Println("Best: ", bestOutcome.Outcome)
	return bestOutcome, nil
}
