package main

type TeamStats struct {
	TeamName           string // Carlton / Collingwood etc etc
	QuarterOneScore    int
	QuarterOneResult   string
	QuarterOneData     string
	QuarterTwoScore    int
	QuarterTwoResult   string
	QuarterTwoData     string
	QuarterThreeScore  int
	QuarterThreeData   string
	QuarterThreeResult string
	QuarterFourScore   int
	QuarterFourData    string
	QuarterFourResult  string
	MatchResult        string
	MatchData          string // String containing info about venue / game time / data etc - might be useful in the future :D
	FinalScore         int
}

type MatchStats struct {
	MatchID     string
	TeamOne     TeamStats
	TeamTwo     TeamStats
	WinningTeam string
	Season      string
}

type TeamStatsWithMatchId struct {
	MatchID            string
	TeamName           string
	QuarterOneScore    int
	QuarterOneResult   string
	QuarterOneData     string
	QuarterTwoScore    int
	QuarterTwoResult   string
	QuarterTwoData     string
	QuarterThreeScore  int
	QuarterThreeData   string
	QuarterThreeResult string
	QuarterFourScore   int
	QuarterFourData    string
	QuarterFourResult  string
	MatchResult        string
	MatchData          string
	FinalScore         int
}
