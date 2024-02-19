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
