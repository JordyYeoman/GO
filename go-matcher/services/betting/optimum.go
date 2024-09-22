package betting

import (
	"math"
)

// CalculateArbitrageOutcome calculates the outcome considering exchange commission
func CalculateArbitrageOutcome(bookieStake, exchangeStake, bookieOdds, exchangeOdds, exchangeCommission float64) float64 {
	exchangeProfit := exchangeStake * (exchangeOdds - 1)
	bookieProfit := bookieStake * (bookieOdds - 1)

	// Apply commission to exchange profit
	exchangeProfitAfterCommission := exchangeProfit * (1 - exchangeCommission)

	return exchangeProfitAfterCommission - bookieProfit
}

// FindOptimalArbitrageBets finds optimal bets considering exchange commission
func FindOptimalArbitrageBets(minBet, maxBet, step, bookieOdds, exchangeOdds, exchangeCommission float64) (float64, float64, float64) {
	bestBookieStake, bestExchangeStake, bestOutcome := 0.0, 0.0, math.Inf(1)

	for bookieStake := minBet; bookieStake <= maxBet; bookieStake += step {
		for exchangeStake := minBet; exchangeStake <= maxBet; exchangeStake += step {
			outcome := CalculateArbitrageOutcome(bookieStake, exchangeStake, bookieOdds, exchangeOdds, exchangeCommission)
			if math.Abs(outcome) < math.Abs(bestOutcome) {
				bestBookieStake, bestExchangeStake, bestOutcome = bookieStake, exchangeStake, outcome
			}
		}
	}

	return bestBookieStake, bestExchangeStake, bestOutcome
}
