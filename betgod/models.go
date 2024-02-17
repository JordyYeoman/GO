package main

type TeamStats struct {
	TeamName          string // Carlton / Collingwood etc etc
	QuarterOneScore   int
	QuarterOneData    string
	QuarterTwoScore   int
	QuarterTwoData    string
	QuarterThreeScore int
	QuarterThreeData  string
	QuarterFourScore  int
	QuarterFourData   string
	MatchResult       string
	MatchData         string // String containing info about venue / game time / data etc - might be useful in the future :D
	FinalScore        int
}

type MatchStats struct {
	MatchID     string
	TeamOne     TeamStats
	TeamTwo     TeamStats
	WinningTeam string
}
