package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Importing a package for side effects, no direct usages (interface for DB)
	"github.com/gocolly/colly"
	"github.com/google/uuid"
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

func connectToDB() *sql.DB {
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

	return db
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

func updateQuarterResult(teamOne, teamTwo *TeamStatsWithMatchId) {
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

// Scraping code:

type AFLSeasonList struct {
	seasonLink string
	seasonYear string
}

func scrapePageData() {
	fmt.Println("System Online and Ready Sir")

	// Generate season data
	var aflSeasonList []AFLSeasonList
	totalSeasons := 30 // Total amount of seasons to record
	lastSeason := 2023 // Season we want to start counting back from

	for i := 0; i < totalSeasons; i++ {
		var season AFLSeasonList
		// Convert lastSeason - i to string
		seasonYear := strconv.Itoa(lastSeason - i)

		// Concatenate the URL parts into a slice of strings
		urlParts := []string{"https://afltables.com/afl/seas/", seasonYear, ".html"}

		// Join the URL parts with an empty separator
		url := strings.Join(urlParts, "")

		season.seasonLink = url
		season.seasonYear = seasonYear
		// Append the URL to aflSeasonList
		aflSeasonList = append(aflSeasonList, season)
	}

	// Create large slice of slices of matches
	var pageData [][]MatchStats
	//Loop over each page link and create dataset
	for _, season := range aflSeasonList {
		p, err := getPageStats(season.seasonLink, season.seasonYear)
		// P1
		// Subroutines should only 'handle' the error if it can recover from it,
		// Return the error ^up if it can.
		// Bubble this bad boi

		// P2 - Graceful error handling.
		// Never log AND return the error.
		// What else can you do with the err??
		if err != nil {
			log.Fatal(err)
			//log.WithError(err).Warn("Getting page stats is DEAD")
		}

		pageData = append(pageData, p)
	}

	fmt.Println(pageData)
}

// TODO:
// Return this out somewhere??!?!
// Every () should return an error.
func ExtractMatchStats(gameURL string, season string) (MatchStats, error) {
	// Struct to contain full match data
	var MatchResult = MatchStats{
		Season:  season,
		MatchID: uuid.New().String(),
	}
	teamOneSet := false

	lines := strings.Split(gameURL, "\n")

	for _, line := range lines {
		// Extract team name and find the actual team name in map
		tempLine := strings.Fields(line)

		if len(tempLine) < 5 {
			continue
		}

		tempStrC := strings.Join(tempLine[:5], " ")
		// Find which team to use
		teamToUse := FindCorrectTeamName(tempStrC)

		if teamToUse != "" {
			//fmt.Println("Found team:", teamToUse)

			// Slice team name from string
			adjustedLine := RemoveTeamName(line, teamToUse)
			stats := ExtractTeamStats(adjustedLine, teamToUse)

			if !teamOneSet {
				MatchResult.TeamOne = TeamStatsWithMatchId{
					MatchID:            MatchResult.MatchID,
					TeamName:           stats.TeamName,
					QuarterOneScore:    stats.QuarterOneScore,
					QuarterOneResult:   stats.QuarterOneResult,
					QuarterOneData:     stats.QuarterOneData,
					QuarterTwoScore:    stats.QuarterTwoScore,
					QuarterTwoResult:   stats.QuarterTwoResult,
					QuarterTwoData:     stats.QuarterTwoData,
					QuarterThreeScore:  stats.QuarterThreeScore,
					QuarterThreeResult: stats.QuarterThreeResult,
					QuarterThreeData:   stats.QuarterThreeData,
					QuarterFourScore:   stats.QuarterFourScore,
					QuarterFourResult:  stats.QuarterFourResult,
					QuarterFourData:    stats.QuarterFourData,
					MatchResult:        stats.MatchResult,
					MatchData:          stats.MatchData,
					FinalScore:         stats.FinalScore,
				}
				teamOneSet = true
				continue
			}

			MatchResult.TeamTwo = TeamStatsWithMatchId{
				MatchID:            MatchResult.MatchID,
				TeamName:           stats.TeamName,
				QuarterOneScore:    stats.QuarterOneScore,
				QuarterOneResult:   stats.QuarterOneResult,
				QuarterOneData:     stats.QuarterOneData,
				QuarterTwoScore:    stats.QuarterTwoScore,
				QuarterTwoResult:   stats.QuarterTwoResult,
				QuarterTwoData:     stats.QuarterTwoData,
				QuarterThreeScore:  stats.QuarterThreeScore,
				QuarterThreeResult: stats.QuarterThreeResult,
				QuarterThreeData:   stats.QuarterThreeData,
				QuarterFourScore:   stats.QuarterFourScore,
				QuarterFourResult:  stats.QuarterFourResult,
				QuarterFourData:    stats.QuarterFourData,
				MatchResult:        stats.MatchResult,
				MatchData:          stats.MatchData,
				FinalScore:         stats.FinalScore,
			}
			break
		}
	}

	// Find match winner
	tempTeamOneOutcome := MatchResult.TeamOne.FinalScore
	tempTeamTwoOutcome := MatchResult.TeamTwo.FinalScore
	// Set Quarter Results for each team (Needed when doing large single team analysis)
	if MatchResult.TeamOne.QuarterOneScore > MatchResult.TeamTwo.QuarterOneScore {
		MatchResult.TeamOne.QuarterOneResult = "WIN"
		MatchResult.TeamTwo.QuarterOneResult = "LOSS"
	} else if MatchResult.TeamOne.QuarterOneScore < MatchResult.TeamTwo.QuarterOneScore {
		MatchResult.TeamOne.QuarterOneResult = "LOSS"
		MatchResult.TeamTwo.QuarterOneResult = "WIN"
	} else {
		MatchResult.TeamOne.QuarterOneResult = "DRAW"
		MatchResult.TeamTwo.QuarterOneResult = "DRAW"
	}
	updateQuarterResult(&MatchResult.TeamOne, &MatchResult.TeamTwo)

	if tempTeamOneOutcome > tempTeamTwoOutcome {
		MatchResult.TeamOne.MatchResult = "WIN"
		MatchResult.TeamTwo.MatchResult = "LOSS"
		MatchResult.WinningTeam = MatchResult.TeamOne.TeamName
	} else if tempTeamOneOutcome < tempTeamTwoOutcome {
		MatchResult.TeamOne.MatchResult = "LOSS"
		MatchResult.TeamTwo.MatchResult = "WIN"
		MatchResult.WinningTeam = MatchResult.TeamTwo.TeamName
	} else {
		// It's a draw
		MatchResult.TeamOne.MatchResult = "DRAW"
		MatchResult.TeamTwo.MatchResult = "DRAW"
	}

	// If no team, return nothing
	if MatchResult.TeamTwo.TeamName == "" {
		return MatchStats{}, nil
	}

	return MatchResult, nil
}

func getPageStats(url string, year string) ([]MatchStats, error) {
	fmt.Println("Scraping: ")
	fmt.Println(url)
	endOfRelevantPage := false // Exiting before finals to ease scraping, can come back and add into data.

	var sliceOMatchStats []MatchStats

	c := colly.NewCollector()
	var err error
	c.OnHTML("table", func(e *colly.HTMLElement) {
		var matchStats MatchStats
		if endOfRelevantPage { // When we reach the final ladder 'year + season'
			return
		}

		// TODO: Error checking here for tables that aren't match stats

		if strings.Contains(e.Text, "Ladder") {
			return
		}

		if strings.Contains(e.Text, year+" Ladder") {
			endOfRelevantPage = true
		}

		// Every 2nd table on the page has the data we require
		// Ignore round number + we start at round 1.
		//fmt.Println(e.Text)
		matchStats, err = ExtractMatchStats(e.Text, year)
		if err != nil {
			// Handle error
			return
		}
		// handle err for above method.

		// Only add match stats if team names exist
		if matchStats.TeamOne.TeamName != "" {
			sliceOMatchStats = append(sliceOMatchStats, matchStats)
		}
	})

	if err != nil {
		return nil, err
	}

	if err := c.Visit(url); err != nil {
		return nil, err
	}

	return sliceOMatchStats, nil
}
