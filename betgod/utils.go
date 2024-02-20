package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Importing a package for side effects, no direct usages (interface for DB)
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// TeamNames Remember to use uppercase declaration if you want to export
var TeamNames = map[string]bool{
	"Richmond":               true,
	"Carlton":                true,
	"Geelong":                true,
	"Collingwood":            true,
	"Melbourne":              true,
	"Sydney":                 true,
	"Adelaide":               true,
	"Hawthorn":               true,
	"Essendon":               true,
	"St Kilda":               true,
	"Fremantle":              true,
	"Greater Western Sydney": true,
	"Gold Coast":             true,
	"Western Bulldogs":       true,
	"West Coast":             true,
	"Port Adelaide":          true,
	"Brisbane Lions":         true,
	"North Melbourne":        true,
}

func StripDigitsFromString(s string) string {
	var tempVar string

	for _, char := range s {
		if !unicode.IsDigit(char) {
			tempVar += string(char)
		}
	}

	return tempVar
}

func GetMatchData(sliceOfStrings []string) string {
	var tempStr string
	// Part 5 is int + day, we just want the day
	for i, item := range sliceOfStrings {
		if i == 4 {
			tempStr += StripDigitsFromString(item)
		}
		if i > 4 {
			tempStr += " "
			tempStr += item
		}
	}

	return tempStr
}

func createTeamStatsDBTables(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS team_stats (
		id INT AUTO_INCREMENT PRIMARY KEY,
    	team_name TEXT,
		quarter_one_score BIGINT,
		quarter_one_result TEXT,
		quarter_one_data TEXT,
		quarter_two_score BIGINT,
		quarter_two_result TEXT,
		quarter_two_data TEXT,
		quarter_three_score BIGINT,
		quarter_three_data TEXT,
		quarter_three_result TEXT,
		quarter_four_score BIGINT,
		quarter_four_data TEXT,
		quarter_four_result TEXT,
		match_result TEXT,
		match_data TEXT,
		final_score BIGINT
	);`

	_, err := db.Exec(query) // Execute query against DB without returning any rows
	if err != nil {
		log.Fatal(err)
	}
}

func createMatchStatsDBTables(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS match_stats (
		id INT AUTO_INCREMENT PRIMARY KEY,
    	match_id TEXT,
		team_one VARCHAR(255),
		team_two VARCHAR(255),
		winning_team TEXT
	)`

	_, err := db.Exec(query) // Execute query against DB without returning any rows
	if err != nil {
		log.Fatal(err)
	}
}

func insertMatchStats(db *sql.DB, matchStats MatchStats) int {
	query := "INSERT INTO match_stats (match_id, team_one, team_two, winning_team) VALUES (?, ?, ?, ?);"
	result, err := db.Exec(query, matchStats.MatchID, matchStats.TeamOne.TeamName, matchStats.TeamTwo.TeamName, matchStats.WinningTeam)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the last inserted ID
	pk, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return int(pk)
}

func insertTeamStats(db *sql.DB, matchStats MatchStats) int {
	query := "INSERT INTO team_stats (team_name, quarter_one_score, quarter_one_result, quarter_one_data, quarter_two_score, quarter_two_result, quarter_two_data, quarter_three_score, quarter_three_data, quarter_three_result, quarter_four_score, quarter_four_data, quarter_four_result, match_result, match_data, final_score) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
	// TODO: Finish query
	result, err := db.Exec(query, matchStats.MatchID, matchStats.TeamOne.TeamName, matchStats.TeamTwo.TeamName, matchStats.WinningTeam)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the last inserted ID
	pk, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return int(pk)
}

func handleDBConnection() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbUrl := os.Getenv("DB_URL")

	fmt.Println("Connecting to DB:")
	db, dbErr := sql.Open("mysql", dbUrl)
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	if dbErr = db.Ping(); dbErr != nil {
		log.Fatal(dbErr)
	}

	// Upload data per season
	// 1. Create tables for matches + team stats
	//createMatchStatsDBTables(db) // done
	//createTeamStatsDBTables(db) // done
	// 1.5 Create a table for each team name??
	//createTeamDBTables(db) // done

	// Test works
	//var tMatch = MatchStats{
	//	MatchID: "slegjhw;ehjg",
	//	TeamOne: TeamStats{
	//		TeamName: "Full Send",
	//	},
	//	TeamTwo: TeamStats{
	//		TeamName: "Nah bro",
	//	},
	//	WinningTeam: "FULL SEND",
	//}
	//insertMatchStats(db, tMatch) // test

	// 2. Add Match Stats to DB table
	// 3. Add Team Stats

	// For each match, create the match stats and two teams stats entries

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db) // Defer means run this when the wrapping function terminates
}

func GetFinalScore(str string) int {
	var tempScore string
	for _, char := range str {
		if unicode.IsDigit(char) {
			tempScore += string(char)
		} else {
			break // Stop iteration if non-digit character encountered
		}
	}

	s, err := strconv.Atoi(tempScore)
	if err != nil {
		fmt.Println("Can't convert this to an int!")
	}

	return s
}

func (s *TeamStats) SetQuarterScore(quarter, score int) {
	switch quarter {
	case 1:
		s.QuarterOneScore += score
	case 2:
		s.QuarterTwoScore += score
	case 3:
		s.QuarterThreeScore += score
	case 4:
		s.QuarterFourScore += score
	}
}

func updateQuarterResult(teamOne, teamTwo *TeamStats) {
	updateQuarter(&teamOne.QuarterOneResult, &teamTwo.QuarterOneResult, teamOne.QuarterOneScore, teamTwo.QuarterOneScore)
	updateQuarter(&teamOne.QuarterTwoResult, &teamTwo.QuarterTwoResult, teamOne.QuarterTwoScore, teamTwo.QuarterTwoScore)
	updateQuarter(&teamOne.QuarterThreeResult, &teamTwo.QuarterThreeResult, teamOne.QuarterThreeScore, teamTwo.QuarterThreeScore)
	updateQuarter(&teamOne.QuarterFourResult, &teamTwo.QuarterFourResult, teamOne.QuarterFourScore, teamTwo.QuarterFourScore)
}

func updateQuarter(teamOneResult, teamTwoResult *string, teamOneScore, teamTwoScore int) {
	if teamOneScore > teamTwoScore {
		*teamOneResult = "WIN"
		*teamTwoResult = "LOSS"
	} else if teamOneScore < teamTwoScore {
		*teamOneResult = "LOSS"
		*teamTwoResult = "WIN"
	} else {
		*teamOneResult = "DRAW"
		*teamTwoResult = "DRAW"
	}
}

func RemoveTeamName(line, team string) string {
	return strings.TrimPrefix(line, team)
}

func ExtractTeamStats(line, team string) TeamStats {
	var stats TeamStats
	stats.TeamName = team
	parts := strings.Fields(line) // Split the line by spaces
	endOfTeamScoresInStringSplit := 4

	// Final Score
	stats.FinalScore = GetFinalScore(parts[endOfTeamScoresInStringSplit])

	// Match data
	stats.MatchData = GetMatchData(parts)

	// Quarters
	for i := 0; i < endOfTeamScoresInStringSplit; i++ {
		score := parts[i]
		scoreParts := strings.Split(score, ".")
		if len(scoreParts) != 2 {
			fmt.Println("Invalid score format:", score)
			continue
		}
		score1, err1 := strconv.Atoi(scoreParts[0])
		score2, err2 := strconv.Atoi(scoreParts[1])
		if err1 != nil || err2 != nil {
			fmt.Println("Error converting score to int:", score)
			continue
		}
		quarter := i + 1 // Quarter 1 corresponds to index 1, Quarter 2 to index 2, and so on
		switch quarter {
		case 1:
			stats.QuarterOneData = score
			stats.QuarterOneScore = score1*6 + score2
		case 2:
			stats.QuarterTwoData = score
			stats.QuarterTwoScore = score1*6 + score2
		case 3:
			stats.QuarterThreeData = score
			stats.QuarterThreeScore = score1*6 + score2
		case 4:
			stats.QuarterFourData = score
			stats.QuarterFourScore = score1*6 + score2
		}
	}

	//fmt.Println("==========")
	//fmt.Println(stats)
	//fmt.Println("==========")
	return stats
}

func FindCorrectTeamName(str string) string {
	// We need to check ALL team names to find best match in substring
	var foundTeamNames []string
	var correctTeamName = ""

	// Find all matching team names
	for team := range TeamNames {
		if strings.Contains(str, team) {
			foundTeamNames = append(foundTeamNames, team)
		}
	}

	// If a longer team name exists in slice, use that instead.
	// EG 'Sydney' and 'Greater Western Sydney', return the longer string
	for _, t := range foundTeamNames {
		if len(t) > len(correctTeamName) {
			correctTeamName = t
		}
	}

	return correctTeamName
}
