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
	G_TeamOneLosesHalfTimeButWins    float64
	V_TeamOneWinsHalfTimeButLoses    float64 // Versus TeamTwo Percentage
	V_TeamOneLosesHalfTimeButWins    float64 // Versus TeamTwo Percentage
	TeamTwo                          string
	G_TeamTwoWinsHalfTimeButLoses    float64 // Versus Any Team Percentage
	G_TeamTwoLosesHalfTimeButWins    float64
	V_TeamTwoWinsHalfTimeButLoses    float64 // Versus TeamTwo Percentage
	V_TeamTwoLosesHalfTimeButWins    float64 // Versus TeamTwo Percentage
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tempTeamOneTally := 0.00
	tempTeamTwoTally := 0.00
	totalTimesPlayedEachOther := 0.00

	// 1. Get raining champ in team v team match-ups
	for _, team := range allTimeTeamVsTeam {
		totalTimesPlayedEachOther++
		if team.WinningTeam == requestBody.TeamOne {
			tempTeamOneTally++
		} else if team.WinningTeam == requestBody.TeamTwo {
			tempTeamTwoTally++
		}
	}

	if tempTeamOneTally > tempTeamTwoTally {
		responseBody.AllTimeTeamWinRate = (tempTeamOneTally / totalTimesPlayedEachOther) * 100
		responseBody.AllTimeTeamWinner = requestBody.TeamOne
	} else if tempTeamOneTally < tempTeamTwoTally {
		responseBody.AllTimeTeamWinRate = (tempTeamTwoTally / totalTimesPlayedEachOther) * 100
		responseBody.AllTimeTeamWinner = requestBody.TeamTwo
	} else {
		responseBody.AllTimeTeamWinRate = 50.00
		responseBody.AllTimeTeamWinner = "DRAW"
	}

	responseBody.TotalGamesPlayedAgainstEachOther = totalTimesPlayedEachOther

	// TEAM ONE
	// All time team one stats
	allTimeTeamOneStats := getAllTeamStatsFromDb(b.DB, requestBody.TeamOne)
	fmt.Println()
	fmt.Println(allTimeTeamOneStats)
	fmt.Println()
	// Team wins at half time and loses game
	allTimeTeamOneWinsSecondQAndLoses := getAllTimeTeamWinsXQuarterAndXOutcome(allTimeTeamOneStats, 2, "WIN", "LOSS")
	// Team wins half time and wins game
	allTimeTeamOneWinsSecondQAndWins := getAllTimeTeamWinsXQuarterAndXOutcome(allTimeTeamOneStats, 2, "WIN", "WIN")

	// Totals - loses/wins halves but X
	totalTeamOneGamesEver := len(allTimeTeamOneStats)
	fmt.Println()
	fmt.Println("Total Games team one EVER")
	fmt.Println(totalTeamOneGamesEver)
	fmt.Println()
	fmt.Println("teamOneWinsSecondQ")
	fmt.Println(allTimeTeamOneWinsSecondQAndLoses)
	fmt.Println()
	fmt.Println("LosesSecondQ")
	fmt.Println(allTimeTeamOneWinsSecondQAndLoses)
	fmt.Println()
	t1WinsHalfButLoses := (len(allTimeTeamOneWinsSecondQAndLoses) / totalTeamOneGamesEver) * 100
	t1LosesHalfButWins := (len(allTimeTeamOneWinsSecondQAndWins) / totalTeamOneGamesEver) * 100

	fmt.Println()
	fmt.Println(t1WinsHalfButLoses)
	fmt.Println()
	fmt.Println(t1LosesHalfButWins)
	fmt.Println()
	responseBody.G_TeamOneWinsHalfTimeButLoses = float64(t1WinsHalfButLoses)
	responseBody.G_TeamOneLosesHalfTimeButWins = float64(t1LosesHalfButWins)

	//V_TeamOneWinsHalfTimeButLoses    float64 // Versus TeamTwo Percentage
	//V_TeamOneLosesHalfTimeButWins

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
