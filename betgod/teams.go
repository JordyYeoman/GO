package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	allTimeTeamVsTeamStats := getAllTimeTeamVSTeamQuarterStats(allTimeTeamVsTeam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	teamOneIndividualStats := IndividualTeamStats{}
	teamTwoIndividualStats := IndividualTeamStats{}

	// TEAM VS TEAM
	// All time team one stats
	allTimeTeamOneStats := getAllTeamStatsFromDb(b.DB, requestBody.TeamOne)
	allTimeTeamOneQuarterStats := getAllTimeTeamQuarterStats(allTimeTeamOneStats)
	fmt.Println(allTimeTeamOneQuarterStats)
	// Global Team wins at half time and loses game
	allTimeTeamOneWinsSecondQAndLoses := getAllTimeTeamWinsXQuarterAndXOutcome(allTimeTeamOneStats, 2, "WIN", "LOSS")
	// Global Team wins half time and wins game
	allTimeTeamOneWinsSecondQAndWins := getAllTimeTeamWinsXQuarterAndXOutcome(allTimeTeamOneStats, 2, "WIN", "WIN")
	// Global Team loses half time and wins game
	allTimeTeamOneLosesSecondQAndWins := getAllTimeTeamWinsXQuarterAndXOutcome(allTimeTeamOneStats, 2, "LOSS", "WIN")
	// Global Team loses half time and loses game
	allTimeTeamOneLosesSecondQAndLoses := getAllTimeTeamWinsXQuarterAndXOutcome(allTimeTeamOneStats, 2, "LOSS", "LOSS")

	// Totals
	totalTeamOneGamesEver := len(allTimeTeamOneStats)

	// ALL TIME
	// - Global
	teamOneIndividualStats.GlobalStats.AllTime.TeamWinsHalfTimeAndWins = (allTimeTeamOneWinsSecondQAndWins / float64(totalTeamOneGamesEver)) * 100
	teamOneIndividualStats.GlobalStats.AllTime.TeamWinsHalfTimeButLoses = (allTimeTeamOneWinsSecondQAndLoses / float64(totalTeamOneGamesEver)) * 100
	teamOneIndividualStats.GlobalStats.AllTime.TeamLosesHalfTimeAndLoses = (allTimeTeamOneLosesSecondQAndLoses / float64(totalTeamOneGamesEver)) * 100
	teamOneIndividualStats.GlobalStats.AllTime.TeamLosesHalfTimeAndWins = (allTimeTeamOneLosesSecondQAndWins / float64(totalTeamOneGamesEver)) * 100
	// - Global - Quarters
	//teamOneIndividualStats.VersusStats.AllTime.TeamQuarterOneWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterOneWinPercent
	//teamOneIndividualStats.VersusStats.AllTime.TeamQuarterTwoWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterTwoWinPercent
	//teamOneIndividualStats.VersusStats.AllTime.TeamQuarterThreeWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterThreeWinPercent
	//teamOneIndividualStats.VersusStats.AllTime.TeamQuarterFourWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterFourWinPercent

	// - Versus
	teamOneIndividualStats.VersusStats.AllTime.TeamWinsHalfTimeAndWins = allTimeTeamVsTeamStats.TotalTeamOneWinsHalfTimeAndWins
	teamOneIndividualStats.VersusStats.AllTime.TeamWinsHalfTimeButLoses = allTimeTeamVsTeamStats.TotalTeamOneWinsHalfTimeButLoses
	teamOneIndividualStats.VersusStats.AllTime.TeamLosesHalfTimeAndLoses = (allTimeTeamOneLosesSecondQAndLoses / float64(totalTeamOneGamesEver)) * 100
	teamOneIndividualStats.VersusStats.AllTime.TeamLosesHalfTimeAndWins = (allTimeTeamOneLosesSecondQAndWins / float64(totalTeamOneGamesEver)) * 100
	// - Versus - Quarters
	teamOneIndividualStats.VersusStats.AllTime.TeamQuarterOneWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterOneWinPercent
	teamOneIndividualStats.VersusStats.AllTime.TeamQuarterTwoWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterTwoWinPercent
	teamOneIndividualStats.VersusStats.AllTime.TeamQuarterThreeWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterThreeWinPercent
	teamOneIndividualStats.VersusStats.AllTime.TeamQuarterFourWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterFourWinPercent

	// Last 3 Seasons
	//totalGamesTeamOnePlayedLast3Seasons := getTotalGamesPlayedLastXSeasons(allTimeTeamOneStats, 3, 2023)
	//// Team wins at half time and loses game
	//last3SeasonsTeamOneWinsSecondQAndLoses := getXSeasonTeamWinsYQuarterAndZOutcome(allTimeTeamOneStats, 2, "WIN", "LOSS", 3, 2023)
	//// Team wins half time and wins game
	//last3SeasonsTeamOneWinsSecondQAndWins := getXSeasonTeamWinsYQuarterAndZOutcome(allTimeTeamOneStats, 2, "WIN", "WIN", 3, 2023)
	//// Team loses half time and wins game
	//last3SeasonsTeamOneLosesSecondQAndWins := getXSeasonTeamWinsYQuarterAndZOutcome(allTimeTeamOneStats, 2, "LOSS", "WIN", 3, 2023)
	//// Team loses half time and loses game
	//last3SeasonsTeamOneLosesSecondQAndLoses := getXSeasonTeamWinsYQuarterAndZOutcome(allTimeTeamOneStats, 2, "LOSS", "LOSS", 3, 2023)

	//teamOneIndividualStats.G_Last3Years_TeamWinsHalfTimeButLoses = (last3SeasonsTeamOneWinsSecondQAndLoses) / totalGamesTeamOnePlayedLast3Seasons * 100
	//teamOneIndividualStats.G_Last3Years_TeamWinsHalfTimeAndWins = (last3SeasonsTeamOneWinsSecondQAndWins) / totalGamesTeamOnePlayedLast3Seasons * 100
	//teamOneIndividualStats.G_Last3Years_TeamLosesHalfTimeAndLoses = (last3SeasonsTeamOneLosesSecondQAndWins) / totalGamesTeamOnePlayedLast3Seasons * 100
	//teamOneIndividualStats.G_Last3Years_TeamLosesHalfTimeAndWins = (last3SeasonsTeamOneLosesSecondQAndLoses) / totalGamesTeamOnePlayedLast3Seasons * 100
	//teamOneIndividualStats.V_Last3Years_TeamQuarterOneWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterOneWinPercent
	//teamOneIndividualStats.V_Last3Years_TeamWinsHalfTimeAndWins = allTimeTeamVsTeamStats.TotalTeamOneWinsHalfTimeAndWins
	//teamOneIndividualStats.V_Last3Years_TeamWinsHalfTimeAndWins = allTimeTeamVsTeamStats.TotalTeamOneWinsHalfTimeButLoses
	//teamOneIndividualStats.V_Last3Years_TeamQuarterOneWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterOneWinPercent
	//teamOneIndividualStats.V_Last3Years_TeamQuarterTwoWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterTwoWinPercent
	//teamOneIndividualStats.V_Last3Years_TeamQuarterThreeWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterThreeWinPercent
	//teamOneIndividualStats.V_Last3Years_TeamQuarterFourWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterFourWinPercent

	// Last Season
	//totalGamesTeamOnePlayedLastSeason := getTotalGamesPlayedLastXSeasons(allTimeTeamOneStats, 1, 2023)
	//// Team wins at half time and loses game
	//lastSeasonTeamOneWinsSecondQAndLoses := getXSeasonTeamWinsYQuarterAndZOutcome(allTimeTeamOneStats, 2, "WIN", "LOSS", 1, 2023)
	//// Team wins half time and wins game
	//lastSeasonTeamOneWinsSecondQAndWins := getXSeasonTeamWinsYQuarterAndZOutcome(allTimeTeamOneStats, 2, "WIN", "WIN", 1, 2023)
	//// Team loses half time and wins game
	//lastSeasonTeamOneLosesSecondQAndWins := getXSeasonTeamWinsYQuarterAndZOutcome(allTimeTeamOneStats, 2, "LOSS", "WIN", 1, 2023)
	//// Team loses half time and loses game
	//lastSeasonTeamOneLosesSecondQAndLoses := getXSeasonTeamWinsYQuarterAndZOutcome(allTimeTeamOneStats, 2, "LOSS", "LOSS", 1, 2023)

	//teamOneIndividualStats.G_LastSeason_TeamWinsHalfTimeButLoses = (lastSeasonTeamOneWinsSecondQAndLoses / totalGamesTeamOnePlayedLastSeason) * 100
	//teamOneIndividualStats.G_LastSeason_TeamWinsHalfTimeAndWins = (lastSeasonTeamOneWinsSecondQAndWins / totalGamesTeamOnePlayedLastSeason) * 100
	//teamOneIndividualStats.G_LastSeason_TeamLosesHalfTimeAndLoses = (lastSeasonTeamOneLosesSecondQAndWins / totalGamesTeamOnePlayedLastSeason) * 100
	//teamOneIndividualStats.G_LastSeason_TeamLosesHalfTimeAndWins = (lastSeasonTeamOneLosesSecondQAndLoses / totalGamesTeamOnePlayedLastSeason) * 100
	//teamOneIndividualStats.V_LastSeason_TeamQuarterOneWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterOneWinPercent
	//teamOneIndividualStats.V_LastSeason_TeamWinsHalfTimeAndWins = allTimeTeamVsTeamStats.TotalTeamOneWinsHalfTimeAndWins
	//teamOneIndividualStats.V_LastSeason_TeamWinsHalfTimeAndWins = allTimeTeamVsTeamStats.TotalTeamOneWinsHalfTimeButLoses
	//teamOneIndividualStats.V_LastSeason_TeamQuarterOneWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterOneWinPercent
	//teamOneIndividualStats.V_LastSeason_TeamQuarterTwoWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterTwoWinPercent
	//teamOneIndividualStats.V_LastSeason_TeamQuarterThreeWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterThreeWinPercent
	//teamOneIndividualStats.V_LastSeason_TeamQuarterFourWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterFourWinPercent

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
