package main

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// FYI - All queries run over 30 seasons of scraped data.
func main() {
	fmt.Println("Connect to DB and analyse data")

	// 1. Create queries to question teams quarter performance
	// 2. Create queries to question team vs team quarter performance
	// 3. Create query to get half / fulltime result average of teams
	// 4. Create query to question half / fulltime result of 2 specific teams

	// Connect to DB
	db := connectToDB()
	teamName := "Collingwood"

	allTeamStats := getAllTeamStats(db, teamName)
	var quarterOneWins []TeamStatsWithMatchId
	var quarterOneLosses []TeamStatsWithMatchId
	var quarterOneDraws []TeamStatsWithMatchId

	for _, team := range allTeamStats {
		if team.QuarterOneResult == "LOSS" {
			quarterOneLosses = append(quarterOneLosses, team)
		} else if team.QuarterOneResult == "WIN" {
			quarterOneWins = append(quarterOneWins, team)
		} else {
			quarterOneDraws = append(quarterOneDraws, team)
		}
	}

	fmt.Println()
	fmt.Printf("%+v Quarter One Wins %+v", teamName, len(quarterOneWins))
	fmt.Println()
	fmt.Printf("%+v Quarter One Losses %+v", teamName, len(quarterOneLosses))
	fmt.Println()
	fmt.Printf("%+v Quarter One Draws %+v", teamName, len(quarterOneDraws))
	fmt.Println()

	// Disconnect DB
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db) // Defer means run this when the wrapping function terminates
}

func getAllTeamStats(db *sql.DB, teamName string) []TeamStatsWithMatchId {
	var data = []TeamStatsWithMatchId{}

	// Placeholder values to hold query data
	var match_id string
	var team_name string
	var quarter_one_score int
	var quarter_one_result string
	var quarter_one_data string
	var quarter_two_score int
	var quarter_two_result string
	var quarter_two_data string
	var quarter_three_score int
	var quarter_three_data string
	var quarter_three_result string
	var quarter_four_score int
	var quarter_four_data string
	var quarter_four_result string
	var match_result string
	var match_data string
	var final_score int

	rows, err := db.Query("SELECT match_id, team_name, quarter_one_score, quarter_one_result, quarter_one_data, quarter_two_score, quarter_two_result, quarter_two_data, quarter_three_score, quarter_three_data, quarter_three_result, quarter_four_score, quarter_four_data, quarter_four_result, match_result, match_data, final_score from team_stats WHERE team_name = ?", teamName)
	if err != nil {
		log.WithError(err).Warn("Error querying db")
	}

	defer rows.Close()
	// to SCAN db vals
	// var name string
	// var available bool
	// var price float64

	for rows.Next() {
		err := rows.Scan(&match_id, &team_name, &quarter_one_score, &quarter_one_result, &quarter_one_data, &quarter_two_score, &quarter_two_result, &quarter_two_data, &quarter_three_score, &quarter_three_data, &quarter_three_result, &quarter_four_score, &quarter_four_data, &quarter_four_result, &match_result, &match_data, &final_score)
		if err != nil {
			log.WithError(err).Warn("Error mapping db query to TeamStatsWithMatchId")
		}
		data = append(data, TeamStatsWithMatchId{match_id, team_name, quarter_one_score, quarter_one_result, quarter_one_data, quarter_two_score, quarter_two_result, quarter_two_data, quarter_three_score, quarter_three_data, quarter_three_result, quarter_four_score, quarter_four_data, quarter_four_result, match_result, match_data, final_score})
	}

	fmt.Println(data)
	return data
}
