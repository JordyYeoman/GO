package main

import "testing"

func TestFindCorrectTeamName(t *testing.T) {
	testStr := "Greater Western Sydney  0.3   4.4   6.9   8.9  57Collingwood won by 1 pt [Match stats]\n"

	result := FindCorrectTeamName(testStr)
	expected := "Greater Western Sydney"

	if result != expected {
		t.Errorf("Test case 1 failed. Expected team name of %v, got %v", expected, result)
	}
}

// {Greater Western Sydney 3  0.3 28  4.4 45 6.9  57 8.9  LOSS Collingwood won by 1 pt [Match stats] 57}

func TestExtractMatchStats(t *testing.T) {
	// Test case 1: Valid game URL
	testHtmlText := "Collingwood  2.2   2.6   7.7  8.10  58Fri 22-Sep-2023 7:50 PM Att: 97,665 Venue: M.C.G.\nGreater Western Sydney  0.3   4.4   6.9   8.9  57Collingwood won by 1 pt [Match stats]"

	expected := MatchStats{
		TeamOne: TeamStats{
			TeamName:           "Collingwood",
			QuarterOneScore:    14,
			QuarterOneData:     "2.2",
			QuarterOneResult:   "WIN",
			QuarterTwoScore:    18,
			QuarterTwoData:     "2.6",
			QuarterTwoResult:   "LOSS",
			QuarterThreeScore:  49,
			QuarterThreeData:   "7.7",
			QuarterThreeResult: "WIN",
			QuarterFourScore:   58,
			QuarterFourData:    "8.10",
			QuarterFourResult:  "WIN",
			MatchResult:        "WIN",
			MatchData:          "Fri 22-Sep-2023 7:50 PM Att: 97,665 Venue: M.C.G.",
			FinalScore:         58,
		},
		TeamTwo: TeamStats{
			TeamName:           "Greater Western Sydney",
			QuarterOneScore:    3,
			QuarterOneData:     "0.3",
			QuarterOneResult:   "LOSS",
			QuarterTwoScore:    28,
			QuarterTwoData:     "4.4",
			QuarterTwoResult:   "WIN",
			QuarterThreeScore:  45,
			QuarterThreeData:   "6.9",
			QuarterThreeResult: "LOSS",
			QuarterFourScore:   57,
			QuarterFourData:    "8.9",
			QuarterFourResult:  "LOSS",
			MatchResult:        "LOSS",
			MatchData:          "Collingwood won by 1 pt [Match stats]",
			FinalScore:         57,
		},
	}

	result := ExtractMatchStats(testHtmlText)

	if result.TeamOne.TeamName != expected.TeamOne.TeamName {
		t.Errorf("Test case 1 failed. Expected team one name of %v, got %v", expected.TeamOne.TeamName, result.TeamOne.TeamName)
	}

	if result.TeamOne != expected.TeamOne {
		t.Errorf("Test case 1 failed. Expected %v, got %v", expected.TeamOne, result.TeamOne)
	}

	if result.TeamTwo != expected.TeamTwo {
		t.Errorf("Test case 1 failed. Expected %v, got %v", expected.TeamTwo, result.TeamTwo)
	}

}
