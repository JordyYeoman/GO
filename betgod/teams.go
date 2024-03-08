package main

import (
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

type TeamHandler struct {
	DB *sql.DB
}

// By default, Go's JSON decoder only sets fields that start with a capital letter.
type GetAllTimeTeamStatsRequestBody struct {
	TeamName string
}

// Return type for all time team stats
type AllTimeTeamStatsAbbrv struct {
	TeamName             string
	AllTimeWinRate       float64
	TotalGamesPlayed     int
	TotalSeasonsCompared int
}

type GetTeamVsTeamRequestBody struct {
	TeamOne string
	TeamTwo string
}

type GenericTeamStatsAverages struct {
	TeamWinsHalfTimeButLoses      float64
	TeamWinsHalfTimeAndWins       float64
	TeamLosesHalfTimeAndLoses     float64
	TeamLosesHalfTimeAndWins      float64
	TeamQuarterOneWinPercentage   float64
	TeamQuarterTwoWinPercentage   float64
	TeamQuarterThreeWinPercentage float64
	TeamQuarterFourWinPercentage  float64
}

type GenericTeamYearsStats struct {
	AllTime    GenericTeamStatsAverages
	Last5Years GenericTeamStatsAverages
	Last3Years GenericTeamStatsAverages
	LastYear   GenericTeamStatsAverages
}

type IndividualTeamStats struct {
	GlobalStats GenericTeamYearsStats
	VersusStats GenericTeamYearsStats
}

type GetTeamVsTeamResponseBody struct {
	AllTimeTeamWinRate               float64 // Percentage
	AllTimeTeamWinner                string
	Draws                            int
	TeamOne                          string
	TeamOneWins                      int
	TeamOneIndividualStats           IndividualTeamStats
	TeamTwo                          string
	TeamTwoWins                      int
	TeamTwoIndividualStats           IndividualTeamStats
	TotalGamesPlayedAgainstEachOther float64
}

func (b TeamHandler) GetTeamVsTeam(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var requestBody GetTeamVsTeamRequestBody
	// Set up the response body
	var responseBody GetTeamVsTeamResponseBody

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Basic team assignment
	responseBody.TeamOne = requestBody.TeamOne
	responseBody.TeamTwo = requestBody.TeamTwo

	// All Time TeamVsTeam stats average?
	allTimeTeamVsTeam, err := getTeamVsTeamStats(b.DB, requestBody.TeamOne, requestBody.TeamTwo)
	allTimeTeamVsTeamStats := getAllTimeTeamVSTeamQuarterStats(allTimeTeamVsTeam, -1, 2023)
	//lastFiveYearsTeamVsTeamStats := getAllTimeTeamVSTeamQuarterStats(allTimeTeamVsTeam, 5, 2023)
	//lastThreeYearsTeamVsTeamStats := getAllTimeTeamVSTeamQuarterStats(allTimeTeamVsTeam, 3, 2023)
	//lastYearsTeamVsTeamStats := getAllTimeTeamVSTeamQuarterStats(allTimeTeamVsTeam, 1, 2023)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	teamOneIndividualStats := IndividualTeamStats{}
	teamTwoIndividualStats := IndividualTeamStats{}

	// TEAM VS TEAM
	// All time team one stats
	allTimeTeamOneStats := getAllTeamStatsFromDb(b.DB, requestBody.TeamOne)
	allTimeTeamOneQuarterStats := getXTimeTeamQuarterStats(allTimeTeamOneStats, -1, -1)
	lastFiveYearsTeamOneQuarterStats := getXTimeTeamQuarterStats(allTimeTeamOneStats, 5, 2023)
	lastThreeYearsTeamOneQuarterStats := getXTimeTeamQuarterStats(allTimeTeamOneStats, 3, 2023)
	lastYearsTeamOneQuarterStats := getXTimeTeamQuarterStats(allTimeTeamOneStats, 1, 2023)

	// ALL TIME
	// - Global - halves
	teamOneIndividualStats.GlobalStats.AllTime.TeamWinsHalfTimeAndWins = allTimeTeamOneQuarterStats.TeamWinsHalfTimeAndWins
	teamOneIndividualStats.GlobalStats.AllTime.TeamWinsHalfTimeButLoses = allTimeTeamOneQuarterStats.TeamWinsHalfTimeButLoses
	teamOneIndividualStats.GlobalStats.AllTime.TeamLosesHalfTimeAndLoses = allTimeTeamOneQuarterStats.TeamLosesHalfTimeAndLoses
	teamOneIndividualStats.GlobalStats.AllTime.TeamLosesHalfTimeAndWins = allTimeTeamOneQuarterStats.TeamLosesHalfTimeAndWins
	// - Global - Quarters
	teamOneIndividualStats.GlobalStats.AllTime.TeamQuarterOneWinPercentage = allTimeTeamOneQuarterStats.TeamQuarterOneWinPercentage
	teamOneIndividualStats.GlobalStats.AllTime.TeamQuarterTwoWinPercentage = allTimeTeamOneQuarterStats.TeamQuarterTwoWinPercentage
	teamOneIndividualStats.GlobalStats.AllTime.TeamQuarterThreeWinPercentage = allTimeTeamOneQuarterStats.TeamQuarterThreeWinPercentage
	teamOneIndividualStats.GlobalStats.AllTime.TeamQuarterFourWinPercentage = allTimeTeamOneQuarterStats.TeamQuarterFourWinPercentage

	// LAST 5 YEARS
	// - Global - Halves
	teamOneIndividualStats.GlobalStats.Last5Years.TeamWinsHalfTimeAndWins = lastFiveYearsTeamOneQuarterStats.TeamWinsHalfTimeAndWins
	teamOneIndividualStats.GlobalStats.Last5Years.TeamWinsHalfTimeButLoses = lastFiveYearsTeamOneQuarterStats.TeamWinsHalfTimeButLoses
	teamOneIndividualStats.GlobalStats.Last5Years.TeamLosesHalfTimeAndLoses = lastFiveYearsTeamOneQuarterStats.TeamLosesHalfTimeAndLoses
	teamOneIndividualStats.GlobalStats.Last5Years.TeamLosesHalfTimeAndWins = lastFiveYearsTeamOneQuarterStats.TeamLosesHalfTimeAndWins
	// - Global - Quarters
	teamOneIndividualStats.GlobalStats.Last5Years.TeamQuarterOneWinPercentage = lastFiveYearsTeamOneQuarterStats.TeamQuarterOneWinPercentage
	teamOneIndividualStats.GlobalStats.Last5Years.TeamQuarterTwoWinPercentage = lastFiveYearsTeamOneQuarterStats.TeamQuarterTwoWinPercentage
	teamOneIndividualStats.GlobalStats.Last5Years.TeamQuarterThreeWinPercentage = lastFiveYearsTeamOneQuarterStats.TeamQuarterThreeWinPercentage
	teamOneIndividualStats.GlobalStats.Last5Years.TeamQuarterFourWinPercentage = lastFiveYearsTeamOneQuarterStats.TeamQuarterFourWinPercentage

	// LAST 3 YEARS
	// - Global - Halves
	teamOneIndividualStats.GlobalStats.Last3Years.TeamWinsHalfTimeAndWins = lastThreeYearsTeamOneQuarterStats.TeamWinsHalfTimeAndWins
	teamOneIndividualStats.GlobalStats.Last3Years.TeamWinsHalfTimeButLoses = lastThreeYearsTeamOneQuarterStats.TeamWinsHalfTimeButLoses
	teamOneIndividualStats.GlobalStats.Last3Years.TeamLosesHalfTimeAndLoses = lastThreeYearsTeamOneQuarterStats.TeamLosesHalfTimeAndLoses
	teamOneIndividualStats.GlobalStats.Last3Years.TeamLosesHalfTimeAndWins = lastThreeYearsTeamOneQuarterStats.TeamLosesHalfTimeAndWins
	// - Global - Quarters
	teamOneIndividualStats.GlobalStats.Last3Years.TeamQuarterOneWinPercentage = lastThreeYearsTeamOneQuarterStats.TeamQuarterOneWinPercentage
	teamOneIndividualStats.GlobalStats.Last3Years.TeamQuarterTwoWinPercentage = lastThreeYearsTeamOneQuarterStats.TeamQuarterTwoWinPercentage
	teamOneIndividualStats.GlobalStats.Last3Years.TeamQuarterThreeWinPercentage = lastThreeYearsTeamOneQuarterStats.TeamQuarterThreeWinPercentage
	teamOneIndividualStats.GlobalStats.Last3Years.TeamQuarterFourWinPercentage = lastThreeYearsTeamOneQuarterStats.TeamQuarterFourWinPercentage

	// LAST YEARS
	// - Global - Halves
	teamOneIndividualStats.GlobalStats.LastYear.TeamWinsHalfTimeAndWins = lastYearsTeamOneQuarterStats.TeamWinsHalfTimeAndWins
	teamOneIndividualStats.GlobalStats.LastYear.TeamWinsHalfTimeButLoses = lastYearsTeamOneQuarterStats.TeamWinsHalfTimeButLoses
	teamOneIndividualStats.GlobalStats.LastYear.TeamLosesHalfTimeAndLoses = lastYearsTeamOneQuarterStats.TeamLosesHalfTimeAndLoses
	teamOneIndividualStats.GlobalStats.LastYear.TeamLosesHalfTimeAndWins = lastYearsTeamOneQuarterStats.TeamLosesHalfTimeAndWins
	// - Global - Quarters
	teamOneIndividualStats.GlobalStats.LastYear.TeamQuarterOneWinPercentage = lastYearsTeamOneQuarterStats.TeamQuarterOneWinPercentage
	teamOneIndividualStats.GlobalStats.LastYear.TeamQuarterTwoWinPercentage = lastYearsTeamOneQuarterStats.TeamQuarterTwoWinPercentage
	teamOneIndividualStats.GlobalStats.LastYear.TeamQuarterThreeWinPercentage = lastYearsTeamOneQuarterStats.TeamQuarterThreeWinPercentage
	teamOneIndividualStats.GlobalStats.LastYear.TeamQuarterFourWinPercentage = lastYearsTeamOneQuarterStats.TeamQuarterFourWinPercentage

	// All time
	// - Versus
	//teamOneIndividualStats.VersusStats.AllTime.TeamWinsHalfTimeAndWins = allTimeTeamVsTeamStats.TotalTeamOneWinsHalfTimeAndWins
	//teamOneIndividualStats.VersusStats.AllTime.TeamWinsHalfTimeButLoses = allTimeTeamVsTeamStats.TotalTeamOneWinsHalfTimeButLoses
	//teamOneIndividualStats.VersusStats.AllTime.TeamLosesHalfTimeAndLoses = (allTimeTeamOneLosesSecondQAndLoses / float64(totalTeamOneGamesEver)) * 100
	//teamOneIndividualStats.VersusStats.AllTime.TeamLosesHalfTimeAndWins = (allTimeTeamOneLosesSecondQAndWins / float64(totalTeamOneGamesEver)) * 100
	// - Versus - Quarters
	teamOneIndividualStats.VersusStats.AllTime.TeamQuarterOneWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterOneWinPercent
	teamOneIndividualStats.VersusStats.AllTime.TeamQuarterTwoWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterTwoWinPercent
	teamOneIndividualStats.VersusStats.AllTime.TeamQuarterThreeWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterThreeWinPercent
	teamOneIndividualStats.VersusStats.AllTime.TeamQuarterFourWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterFourWinPercent

	// Last 5 years
	// - Versus
	//teamOneIndividualStats.VersusStats.AllTime.TeamWinsHalfTimeAndWins = allTimeTeamVsTeamStats.TotalTeamOneWinsHalfTimeAndWins
	//teamOneIndividualStats.VersusStats.AllTime.TeamWinsHalfTimeButLoses = allTimeTeamVsTeamStats.TotalTeamOneWinsHalfTimeButLoses
	//teamOneIndividualStats.VersusStats.AllTime.TeamLosesHalfTimeAndLoses = (allTimeTeamOneLosesSecondQAndLoses / float64(totalTeamOneGamesEver)) * 100
	//teamOneIndividualStats.VersusStats.AllTime.TeamLosesHalfTimeAndWins = (allTimeTeamOneLosesSecondQAndWins / float64(totalTeamOneGamesEver)) * 100
	// - Versus - Quarters
	//teamOneIndividualStats.VersusStats.AllTime.TeamQuarterOneWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterOneWinPercent
	//teamOneIndividualStats.VersusStats.AllTime.TeamQuarterTwoWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterTwoWinPercent
	//teamOneIndividualStats.VersusStats.AllTime.TeamQuarterThreeWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterThreeWinPercent
	//teamOneIndividualStats.VersusStats.AllTime.TeamQuarterFourWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterFourWinPercent

	// - Versus
	//teamOneIndividualStats.VersusStats.AllTime.TeamWinsHalfTimeAndWins = allTimeTeamVsTeamStats.TotalTeamOneWinsHalfTimeAndWins
	//teamOneIndividualStats.VersusStats.AllTime.TeamWinsHalfTimeButLoses = allTimeTeamVsTeamStats.TotalTeamOneWinsHalfTimeButLoses
	//teamOneIndividualStats.VersusStats.AllTime.TeamLosesHalfTimeAndLoses = (allTimeTeamOneLosesSecondQAndLoses / float64(totalTeamOneGamesEver)) * 100
	//teamOneIndividualStats.VersusStats.AllTime.TeamLosesHalfTimeAndWins = (allTimeTeamOneLosesSecondQAndWins / float64(totalTeamOneGamesEver)) * 100
	// - Versus - Quarters
	//teamOneIndividualStats.VersusStats.AllTime.TeamQuarterOneWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterOneWinPercent
	//teamOneIndividualStats.VersusStats.AllTime.TeamQuarterTwoWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterTwoWinPercent
	//teamOneIndividualStats.VersusStats.AllTime.TeamQuarterThreeWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterThreeWinPercent
	//teamOneIndividualStats.VersusStats.AllTime.TeamQuarterFourWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterFourWinPercent

	// - Versus
	//teamOneIndividualStats.VersusStats.AllTime.TeamWinsHalfTimeAndWins = allTimeTeamVsTeamStats.TotalTeamOneWinsHalfTimeAndWins
	//teamOneIndividualStats.VersusStats.AllTime.TeamWinsHalfTimeButLoses = allTimeTeamVsTeamStats.TotalTeamOneWinsHalfTimeButLoses
	//teamOneIndividualStats.VersusStats.AllTime.TeamLosesHalfTimeAndLoses = (allTimeTeamOneLosesSecondQAndLoses / float64(totalTeamOneGamesEver)) * 100
	//teamOneIndividualStats.VersusStats.AllTime.TeamLosesHalfTimeAndWins = (allTimeTeamOneLosesSecondQAndWins / float64(totalTeamOneGamesEver)) * 100
	// - Versus - Quarters
	//teamOneIndividualStats.VersusStats.AllTime.TeamQuarterOneWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterOneWinPercent
	//teamOneIndividualStats.VersusStats.AllTime.TeamQuarterTwoWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterTwoWinPercent
	//teamOneIndividualStats.VersusStats.AllTime.TeamQuarterThreeWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterThreeWinPercent
	//teamOneIndividualStats.VersusStats.AllTime.TeamQuarterFourWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterFourWinPercent

	// TEAM TWO
	//allTimeTeamTwoStats := getAllTeamStatsFromDb(b.DB, requestBody.TeamTwo)
	// Global Team wins at half time and loses game
	//allTimeTeamTwoWinsSecondQAndLoses := getAllTimeTeamWinsXQuarterAndXOutcome(allTimeTeamTwoStats, 2, "WIN", "LOSS")
	//// Global Team wins half time and wins game
	//allTimeTeamTwoWinsSecondQAndWins := getAllTimeTeamWinsXQuarterAndXOutcome(allTimeTeamTwoStats, 2, "WIN", "WIN")
	//// Global Team loses half time and wins game
	//allTimeTeamTwoLosesSecondQAndWins := getAllTimeTeamWinsXQuarterAndXOutcome(allTimeTeamTwoStats, 2, "LOSS", "WIN")
	//// Global Team loses half time and loses game
	//allTimeTeamTwoLosesSecondQAndLoses := getAllTimeTeamWinsXQuarterAndXOutcome(allTimeTeamTwoStats, 2, "LOSS", "LOSS")

	// Totals
	//totalTeamTwoGamesEver := len(allTimeTeamTwoStats)

	//teamTwoIndividualStats.G_AllTime_TeamWinsHalfTimeButLoses = (float64(len(allTimeTeamTwoWinsSecondQAndLoses)) / float64(totalTeamTwoGamesEver)) * 100
	//teamTwoIndividualStats.G_AllTime_TeamWinsHalfTimeAndWins = (float64(len(allTimeTeamTwoWinsSecondQAndWins)) / float64(totalTeamTwoGamesEver)) * 100
	//teamTwoIndividualStats.G_AllTime_TeamLosesHalfTimeAndWins = (float64(len(allTimeTeamTwoLosesSecondQAndWins)) / float64(totalTeamTwoGamesEver)) * 100
	//teamTwoIndividualStats.G_AllTime_TeamLosesHalfTimeAndLoses = (float64(len(allTimeTeamTwoLosesSecondQAndLoses)) / float64(totalTeamTwoGamesEver)) * 100

	//teamTwoIndividualStats.V_AllTime_TeamQuarterOneWinPercentage = allTimeTeamVsTeamStats.TeamTwoQuarterOneWinPercent
	//teamTwoIndividualStats.V_AllTime_TeamQuarterTwoWinPercentage = allTimeTeamVsTeamStats.TeamTwoQuarterTwoWinPercent
	//teamTwoIndividualStats.V_AllTime_TeamQuarterThreeWinPercentage = allTimeTeamVsTeamStats.TeamTwoQuarterThreeWinPercent
	//teamTwoIndividualStats.V_AllTime_TeamQuarterFourWinPercentage = allTimeTeamVsTeamStats.TeamTwoQuarterFourWinPercent
	//teamTwoIndividualStats.V_AllTime_TeamWinsHalfTimeAndWins = allTimeTeamVsTeamStats.TotalTeamTwoWinsHalfTimeAndWins
	//teamTwoIndividualStats.V_AllTime_TeamWinsHalfTimeButLoses = allTimeTeamVsTeamStats.TotalTeamTwoWinsHalfTimeButLoses

	responseBody.Draws = allTimeTeamVsTeamStats.TotalDraws
	responseBody.TotalGamesPlayedAgainstEachOther = allTimeTeamVsTeamStats.TotalTimesPlayed
	responseBody.TeamOneWins = allTimeTeamVsTeamStats.TotalTeamOneWins
	responseBody.TeamTwoWins = allTimeTeamVsTeamStats.TotalTeamTwoWins

	// Get all time team winner
	if responseBody.TeamOneWins > responseBody.TeamTwoWins {
		responseBody.AllTimeTeamWinner = requestBody.TeamOne
		responseBody.AllTimeTeamWinRate = (float64(responseBody.TeamOneWins) / responseBody.TotalGamesPlayedAgainstEachOther) * 100
	} else if responseBody.TeamOneWins < responseBody.TeamTwoWins {
		responseBody.AllTimeTeamWinner = requestBody.TeamTwo
		responseBody.AllTimeTeamWinRate = (float64(responseBody.TeamTwoWins) / responseBody.TotalGamesPlayedAgainstEachOther) * 100
	} else {
		responseBody.AllTimeTeamWinner = "DRAW" //
		responseBody.AllTimeTeamWinRate = 0
	}

	// Assign individual team stats
	responseBody.TeamOneIndividualStats = teamOneIndividualStats
	responseBody.TeamTwoIndividualStats = teamTwoIndividualStats

	// Finally respond with payload
	respondWithJSON(w, 200, responseBody)
}

func (b TeamHandler) GetAllTimeTeamStats(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var requestBody GetAllTimeTeamStatsRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// No team name supplied
	if requestBody.TeamName == "" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(404) // setting response status code
		w.Write([]byte("Team name required"))
		return
	}

	team := getAllTeamStatsFromDb(b.DB, requestBody.TeamName)

	respondWithJSON(w, 200, team)
}

func (b TeamHandler) GetAllTimeTeamAbbrvStats(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var requestBody GetAllTimeTeamStatsRequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// No team name supplied
	if requestBody.TeamName == "" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(404) // setting response status code
		w.Write([]byte("Team name required"))
		return
	}

	teamStats := getAllTeamStatsFromDb(b.DB, requestBody.TeamName)

	allTimeTeamStats := AllTimeTeamStatsAbbrv{}
	allTimeTeamStats.TeamName = requestBody.TeamName
	seasonsIteratedOver := map[int]bool{}

	totalWins := 0.0
	totalPlayed := 0

	// All time win rate
	for _, t := range teamStats {
		totalPlayed++
		// Add season to map
		seasonsIteratedOver[t.Season] = true
		if t.MatchResult == "WIN" {
			totalWins++
		}
	}

	allTimeTeamStats.TotalGamesPlayed = totalPlayed
	allTimeTeamStats.AllTimeWinRate = (totalWins / float64(totalPlayed)) * 100
	allTimeTeamStats.TotalSeasonsCompared = len(seasonsIteratedOver)

	respondWithJSON(w, 200, allTimeTeamStats)
}

func (b TeamHandler) List(w http.ResponseWriter, r *http.Request) {
	var teams []string

	for t := range TeamNames {
		teams = append(teams, t)
	}

	respondWithJSON(w, 200, teams)
}

func TeamRoutes(db *sql.DB) chi.Router {
	r := chi.NewRouter()
	teamHandler := TeamHandler{
		DB: db,
	}
	r.Get("/allTimeTeamStats", teamHandler.GetAllTimeTeamStats)
	r.Get("/teamVsTeam", teamHandler.GetTeamVsTeam)
	r.Get("/allTimeTeamStatsAbbrv", teamHandler.GetAllTimeTeamAbbrvStats)
	r.Get("/list", teamHandler.List)
	return r
}
