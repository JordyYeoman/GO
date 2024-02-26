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

type GetTeamVsTeamRequestBody struct {
	TeamOne string
	TeamTwo string
}

type GetTeamVsTeamResponseBody struct {
	AllTimeTeamWinRate               float64 // Percentage
	AllTimeTeamWinner                string
	TeamOne                          string
	G_TeamOneWinsHalfTimeButLoses    float64 // Versus Any Team Percentage
	G_TeamOneWinsHalfTimeAndWins     float64 // Versus Any Team Percentage
	V_TeamOneWinsHalfTimeButLoses    float64 // Versus TeamTwo Percentage
	V_TeamOneWinsHalfTimeAndWins     float64 // Versus TeamTwo Percentage
	TeamTwo                          string
	G_TeamTwoWinsHalfTimeButLoses    float64 // Versus Any Team Percentage
	G_TeamTwoWinsHalfTimeAndWins     float64
	V_TeamTwoWinsHalfTimeButLoses    float64 // Versus TeamTwo Percentage
	V_TeamTwoWinsHalfTimeAndWins     float64 // Versus TeamTwo Percentage
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

	//tempTeamOneTally := 0.00
	//tempTeamTwoTally := 0.00
	totalTimesPlayedEachOther := 0.00

	fmt.Println(allTimeTeamVsTeam)

	//if tempTeamOneTally > tempTeamTwoTally {
	//	responseBody.AllTimeTeamWinRate = (tempTeamOneTally / totalTimesPlayedEachOther) * 100
	//	responseBody.AllTimeTeamWinner = requestBody.TeamOne
	//} else if tempTeamOneTally < tempTeamTwoTally {
	//	responseBody.AllTimeTeamWinRate = (tempTeamTwoTally / totalTimesPlayedEachOther) * 100
	//	responseBody.AllTimeTeamWinner = requestBody.TeamTwo
	//} else {
	//	responseBody.AllTimeTeamWinRate = 50.00
	//	responseBody.AllTimeTeamWinner = "DRAW"
	//}

	responseBody.TotalGamesPlayedAgainstEachOther = totalTimesPlayedEachOther

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
	responseBody.V_TeamTwoWinsHalfTimeAndWins = allTimeTeamVsTeamStats.TotalTeamTwoWinsHalfTimeAndWins
	responseBody.V_TeamTwoWinsHalfTimeButLoses = allTimeTeamVsTeamStats.TotalTeamTwoWinsHalfTimeButLoses

	responseBody.G_TeamOneWinsHalfTimeButLoses = t1WinsHalfButLoses
	responseBody.G_TeamOneWinsHalfTimeAndWins = t1WinsHalfTimeAndWins
	responseBody.G_TeamTwoWinsHalfTimeButLoses = t2WinsHalfButLoses
	responseBody.G_TeamTwoWinsHalfTimeAndWins = t2WinsHalfTimeAndWins

	// TODO: TEAM TWO
	// ....
	// ...

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

	testTeam := getAllTeamStatsFromDb(b.DB, requestBody.TeamName)

	respondWithJSON(w, 200, testTeam)
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
	r.Get("/list", teamHandler.List)
	return r
}
