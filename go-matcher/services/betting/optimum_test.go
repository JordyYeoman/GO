package betting

import (
	"math"
	"testing"
)

func TestCalculateOutcome(t *testing.T) {
	tests := []struct {
		name            string
		x, o            float64
		yOdds, pOdds    float64
		expectedOutcome float64
	}{
		{"Equal bets, equal odds", 10, 10, 2.0, 2.0, 0},
		{"Equal bets, different odds", 10, 10, 2.0, 1.8, 2},
		{"Different bets, equal odds", 15, 20, 2.0, 2.0, -5},
		{"Different bets, different odds", 15, 20, 2.0, 1.8, -1}, // Updated expected outcome
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outcome := CalculateOutcome(tt.x, tt.o, tt.yOdds, tt.pOdds)
			if math.Abs(outcome-tt.expectedOutcome) > 1e-9 {
				t.Errorf("CalculateOutcome(%f, %f, %f, %f) = %f; want %f",
					tt.x, tt.o, tt.yOdds, tt.pOdds, outcome, tt.expectedOutcome)
			}
		})
	}
}

func TestFindOptimalBets(t *testing.T) {
	tests := []struct {
		name                       string
		minBet, maxBet, step       float64
		yOdds, pOdds               float64
		expectedXRange             [2]float64
		expectedORange             [2]float64
		maxExpectedOutcomeAbsValue float64
	}{
		{
			name:   "Simple case",
			minBet: 15, maxBet: 30, step: 0.1,
			yOdds: 2.0, pOdds: 1.8,
			expectedXRange:             [2]float64{15.0, 15.3},
			expectedORange:             [2]float64{18.8, 19.2},
			maxExpectedOutcomeAbsValue: 0.1,
		},
		{
			name:   "Equal odds case",
			minBet: 15, maxBet: 30, step: 0.1,
			yOdds: 2.0, pOdds: 2.0,
			expectedXRange:             [2]float64{15.0, 15.1},
			expectedORange:             [2]float64{15.0, 15.1},
			maxExpectedOutcomeAbsValue: 0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, o, outcome := FindOptimalBets(tt.minBet, tt.maxBet, tt.step, tt.yOdds, tt.pOdds)

			if x < tt.expectedXRange[0] || x > tt.expectedXRange[1] {
				t.Errorf("Expected X between %f and %f, but got %f", tt.expectedXRange[0], tt.expectedXRange[1], x)
			}
			if o < tt.expectedORange[0] || o > tt.expectedORange[1] {
				t.Errorf("Expected O between %f and %f, but got %f", tt.expectedORange[0], tt.expectedORange[1], o)
			}
			if math.Abs(outcome) > tt.maxExpectedOutcomeAbsValue {
				t.Errorf("Expected absolute outcome to be at most %f, but got %f", tt.maxExpectedOutcomeAbsValue, math.Abs(outcome))
			}
		})
	}
}
