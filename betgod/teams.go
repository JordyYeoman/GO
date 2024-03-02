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

type GetTeamVsTeamResponseBody struct {
	AllTimeTeamWinRate                 float64 // Percentage
	AllTimeTeamWinner                  string
	Draws                              int
	TeamOne                            string
	TeamOneWins                        int
	G_TeamOneWinsHalfTimeButLoses      float64 // Versus Any Team Percentage
	G_TeamOneWinsHalfTimeAndWins       float64 // Versus Any Team Percentage
	V_TeamOneWinsHalfTimeButLoses      float64 // Versus TeamTwo Percentage
	V_TeamOneWinsHalfTimeAndWins       float64 // Versus TeamTwo Percentage
	V_TeamOneQuarterOneWinPercentage   float64
	V_TeamOneQuarterTwoWinPercentage   float64
	V_TeamOneQuarterThreeWinPercentage float64
	V_TeamOneQuarterFourWinPercentage  float64
	TeamTwo                            string
	TeamTwoWins                        int
	G_TeamTwoWinsHalfTimeButLoses      float64 // Versus Any Team Percentage
	G_TeamTwoWinsHalfTimeAndWins       float64
	V_TeamTwoWinsHalfTimeButLoses      float64 // Versus TeamTwo Percentage
	V_TeamTwoWinsHalfTimeAndWins       float64 // Versus TeamTwo Percentage
	TotalGamesPlayedAgainstEachOther   float64
	V_TeamTwoQuarterOneWinPercentage   float64
	V_TeamTwoQuarterTwoWinPercentage   float64
	V_TeamTwoQuarterThreeWinPercentage float64
	V_TeamTwoQuarterFourWinPercentage  float64
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

	// TEAM ONE
	// All time team one stats
	allTimeTeamOneStats := getAllTeamStatsFromDb(b.DB, requestBody.TeamOne)
	// Global Team wins at half time and loses game
	allTimeTeamOneWinsSecondQAndLoses := getAllTimeTeamWinsXQuarterAndXOutcome(allTimeTeamOneStats, 2, "WIN", "LOSS")
	// Global Team wins half time and wins game
	allTimeTeamOneWinsSecondQAndWins := getAllTimeTeamWinsXQuarterAndXOutcome(allTimeTeamOneStats, 2, "WIN", "WIN")
	// Totals
	totalTeamOneGamesEver := len(allTimeTeamOneStats)
	t1WinsHalfButLoses := (float64(len(allTimeTeamOneWinsSecondQAndLoses)) / float64(totalTeamOneGamesEver)) * 100
	t1WinsHalfTimeAndWins := (float64(len(allTimeTeamOneWinsSecondQAndWins)) / float64(totalTeamOneGamesEver)) * 100

	// TEAM TWO
	allTimeTeamTwoStats := getAllTeamStatsFromDb(b.DB, requestBody.TeamTwo)
	// Global Team wins at half time and loses game
	allTimeTeamTwoWinsSecondQAndLoses := getAllTimeTeamWinsXQuarterAndXOutcome(allTimeTeamTwoStats, 2, "WIN", "LOSS")
	// Global Team wins half time and wins game
	allTimeTeamTwoWinsSecondQAndWins := getAllTimeTeamWinsXQuarterAndXOutcome(allTimeTeamTwoStats, 2, "WIN", "WIN")
	// Totals
	totalTeamTwoGamesEver := len(allTimeTeamTwoStats)
	t2WinsHalfButLoses := (float64(len(allTimeTeamTwoWinsSecondQAndLoses)) / float64(totalTeamTwoGamesEver)) * 100
	t2WinsHalfTimeAndWins := (float64(len(allTimeTeamTwoWinsSecondQAndWins)) / float64(totalTeamTwoGamesEver)) * 100

	// TEAM VS TEAM
	responseBody.V_TeamOneWinsHalfTimeAndWins = allTimeTeamVsTeamStats.TotalTeamOneWinsHalfTimeAndWins
	responseBody.V_TeamOneWinsHalfTimeButLoses = allTimeTeamVsTeamStats.TotalTeamOneWinsHalfTimeButLoses
	responseBody.G_TeamOneWinsHalfTimeButLoses = t1WinsHalfButLoses
	responseBody.G_TeamOneWinsHalfTimeAndWins = t1WinsHalfTimeAndWins
	responseBody.V_TeamOneQuarterOneWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterOneWinPercent
	responseBody.V_TeamOneQuarterTwoWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterTwoWinPercent
	responseBody.V_TeamOneQuarterThreeWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterThreeWinPercent
	responseBody.V_TeamOneQuarterFourWinPercentage = allTimeTeamVsTeamStats.TeamOneQuarterFourWinPercent
	responseBody.TeamOneWins = allTimeTeamVsTeamStats.TotalTeamOneWins

	responseBody.G_TeamTwoWinsHalfTimeButLoses = t2WinsHalfButLoses
	responseBody.G_TeamTwoWinsHalfTimeAndWins = t2WinsHalfTimeAndWins
	responseBody.V_TeamTwoQuarterOneWinPercentage = allTimeTeamVsTeamStats.TeamTwoQuarterOneWinPercent
	responseBody.V_TeamTwoQuarterTwoWinPercentage = allTimeTeamVsTeamStats.TeamTwoQuarterTwoWinPercent
	responseBody.V_TeamTwoQuarterThreeWinPercentage = allTimeTeamVsTeamStats.TeamTwoQuarterThreeWinPercent
	responseBody.V_TeamTwoQuarterFourWinPercentage = allTimeTeamVsTeamStats.TeamTwoQuarterFourWinPercent
	responseBody.TeamTwoWins = allTimeTeamVsTeamStats.TotalTeamTwoWins
	responseBody.V_TeamTwoWinsHalfTimeAndWins = allTimeTeamVsTeamStats.TotalTeamTwoWinsHalfTimeAndWins
	responseBody.V_TeamTwoWinsHalfTimeButLoses = allTimeTeamVsTeamStats.TotalTeamTwoWinsHalfTimeButLoses

	responseBody.Draws = allTimeTeamVsTeamStats.TotalDraws
	responseBody.TotalGamesPlayedAgainstEachOther = allTimeTeamVsTeamStats.TotalTimesPlayed

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
