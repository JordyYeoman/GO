package betting

import (
	"math"
)

func FindOptimalBets(minBet, maxBet, step, yOdds, pOdds float64) (float64, float64, float64) {
	bestX, bestO, bestOutcome := 0.0, 0.0, math.Inf(1)

	for x := minBet; x <= maxBet; x += step {
		for o := minBet; o <= maxBet; o += step {
			outcome := CalculateOutcome(x, o, yOdds, pOdds)
			if math.Abs(outcome) < math.Abs(bestOutcome) {
				bestX, bestO, bestOutcome = x, o, outcome
			}
		}
	}

	return bestX, bestO, bestOutcome
}

func CalculateOutcome(x, o, yOdds, pOdds float64) float64 {
	return x*(yOdds-1) - o*(pOdds-1)
}
