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

func insertMatchStats(db *sql.DB, matchStats MatchStats) int {
	query := "INSERT INTO match_stats (match_id, team_one, team_two, winning_team, season) VALUES (?, ?, ?, ?, ?);"
	result, err := db.Exec(query, matchStats.MatchID, matchStats.TeamOne.TeamName, matchStats.TeamTwo.TeamName, matchStats.WinningTeam, matchStats.Season)
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

func insertTeamStats(db *sql.DB, teamStats TeamStatsWithMatchId) int {
	query := "INSERT INTO team_stats (match_id, team_name, quarter_one_score, quarter_one_result, quarter_one_data, quarter_two_score, quarter_two_result, quarter_two_data, quarter_three_score, quarter_three_data, quarter_three_result, quarter_four_score, quarter_four_data, quarter_four_result, match_result, match_data, final_score) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
	// TODO: Finish query
	result, err := db.Exec(query, teamStats.MatchID, teamStats.TeamName, teamStats.QuarterOneScore, teamStats.QuarterOneResult, teamStats.QuarterOneData, teamStats.QuarterTwoScore, teamStats.QuarterTwoResult, teamStats.QuarterTwoData, teamStats.QuarterThreeScore, teamStats.QuarterThreeResult, teamStats.QuarterThreeData, teamStats.QuarterFourScore, teamStats.QuarterFourResult, teamStats.QuarterFourData, teamStats.MatchResult, teamStats.MatchData, teamStats.FinalScore)
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

func handleDBConnection(seasons [][]MatchStats) {
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

	// For each match, create the match stats and two teams stats entries
	// Remember to add season number also...
	for _, matches := range seasons {
		for _, match := range matches {
			// Create a team entry for each team
			var teamOne = TeamStatsWithMatchId{
				MatchID:            match.MatchID,
				TeamName:           match.TeamOne.TeamName,
				QuarterOneScore:    match.TeamOne.QuarterOneScore,
				QuarterOneData:     match.TeamOne.QuarterOneData,
				QuarterOneResult:   match.TeamOne.QuarterOneResult,
				QuarterTwoScore:    match.TeamOne.QuarterTwoScore,
				QuarterTwoData:     match.TeamOne.QuarterTwoData,
				QuarterTwoResult:   match.TeamOne.QuarterTwoResult,
				QuarterThreeScore:  match.TeamOne.QuarterThreeScore,
				QuarterThreeData:   match.TeamOne.QuarterThreeData,
				QuarterThreeResult: match.TeamOne.QuarterThreeResult,
				QuarterFourScore:   match.TeamOne.QuarterFourScore,
				QuarterFourData:    match.TeamOne.QuarterFourData,
				QuarterFourResult:  match.TeamOne.QuarterFourResult,
				MatchResult:        match.TeamOne.MatchResult,
				MatchData:          match.TeamOne.MatchData,
				FinalScore:         match.TeamOne.FinalScore,
			}

			var teamTwo = TeamStatsWithMatchId{
				MatchID:            match.MatchID,
				TeamName:           match.TeamTwo.TeamName,
				QuarterOneScore:    match.TeamTwo.QuarterOneScore,
				QuarterOneData:     match.TeamTwo.QuarterOneData,
				QuarterOneResult:   match.TeamTwo.QuarterOneResult,
				QuarterTwoScore:    match.TeamTwo.QuarterTwoScore,
				QuarterTwoData:     match.TeamTwo.QuarterTwoData,
				QuarterTwoResult:   match.TeamTwo.QuarterTwoResult,
				QuarterThreeScore:  match.TeamTwo.QuarterThreeScore,
				QuarterThreeData:   match.TeamTwo.QuarterThreeData,
				QuarterThreeResult: match.TeamTwo.QuarterThreeResult,
				QuarterFourScore:   match.TeamTwo.QuarterFourScore,
				QuarterFourData:    match.TeamTwo.QuarterFourData,
				QuarterFourResult:  match.TeamTwo.QuarterFourResult,
				MatchResult:        match.TeamTwo.MatchResult,
				MatchData:          match.TeamTwo.MatchData,
				FinalScore:         match.TeamTwo.FinalScore,
			}
			// Create separate team entries
			insertTeamStats(db, teamOne)
			insertTeamStats(db, teamTwo)
			// Create match entry
			insertMatchStats(db, match)
		}
	}

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
