package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Importing a package for side effects, no direct usages (interface for DB)
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
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
	query := "INSERT INTO team_stats (match_id, team_name, quarter_one_score, quarter_one_result, quarter_one_data, quarter_two_score, quarter_two_result, quarter_two_data, quarter_three_score, quarter_three_data, quarter_three_result, quarter_four_score, quarter_four_data, quarter_four_result, match_result, match_data, final_score, season) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"

	result, err := db.Exec(query, teamStats.MatchID, teamStats.TeamName, teamStats.QuarterOneScore, teamStats.QuarterOneResult, teamStats.QuarterOneData, teamStats.QuarterTwoScore, teamStats.QuarterTwoResult, teamStats.QuarterTwoData, teamStats.QuarterThreeScore, teamStats.QuarterThreeResult, teamStats.QuarterThreeData, teamStats.QuarterFourScore, teamStats.QuarterFourResult, teamStats.QuarterFourData, teamStats.MatchResult, teamStats.MatchData, teamStats.FinalScore, teamStats.Season)
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

func GetQuarterResult(team TeamStatsWithMatchId, quarter int) string {
	switch quarter {
	case 1:
		return team.QuarterOneResult
	case 2:
		return team.QuarterTwoResult
	case 3:
		return team.QuarterThreeResult
	case 4:
		return team.QuarterFourResult
	default:
		return "" // Handle invalid quarter
	}
}

// Function to fetch team stats for a given match_id and team
func GetTeamStats(db *sql.DB, matchID string, teamName string) (TeamStatsWithMatchId, error) {
	var teamStats TeamStatsWithMatchId

	// Query team stats
	err := db.QueryRow("SELECT match_id, team_name, quarter_one_score, quarter_one_result, quarter_one_data, quarter_two_score, quarter_two_result, quarter_two_data, quarter_three_score, quarter_three_result, quarter_three_data, quarter_four_score, quarter_four_result, quarter_four_data, match_result, match_data, final_score, season FROM team_stats WHERE match_id = ? AND team_name = ?", matchID, teamName).Scan(&teamStats.MatchID, &teamStats.TeamName, &teamStats.QuarterOneScore, &teamStats.QuarterOneResult, &teamStats.QuarterOneData, &teamStats.QuarterTwoScore, &teamStats.QuarterTwoResult, &teamStats.QuarterTwoData, &teamStats.QuarterThreeScore, &teamStats.QuarterThreeResult, &teamStats.QuarterThreeData, &teamStats.QuarterFourScore, &teamStats.QuarterFourResult, &teamStats.QuarterFourData, &teamStats.MatchResult, &teamStats.MatchData, &teamStats.FinalScore, &teamStats.Season)
	if err != nil {
		return teamStats, err
	}

	return teamStats, nil
}

// Used for weighting higher levels of more recent years
func getOnlyLast5YearsOfTeamStats(allTeamStats []TeamStatsWithMatchId, startingSeasonToCountBackFrom int) []TeamStatsWithMatchId {
	var data []TeamStatsWithMatchId

	for _, team := range allTeamStats {
		if team.Season > startingSeasonToCountBackFrom-5 {
			data = append(data, team)
		}
	}

	return data
}

func getAllTeamStatsFromDb(db *sql.DB, teamName string) []TeamStatsWithMatchId {
	var data []TeamStatsWithMatchId

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
	var season int

	rows, err := db.Query("SELECT match_id, team_name, quarter_one_score, quarter_one_result, quarter_one_data, quarter_two_score, quarter_two_result, quarter_two_data, quarter_three_score, quarter_three_data, quarter_three_result, quarter_four_score, quarter_four_data, quarter_four_result, match_result, match_data, final_score, season from team_stats WHERE team_name = ?", teamName)
	if err != nil {
		log.WithError(err).Warn("Error querying db")
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&match_id, &team_name, &quarter_one_score, &quarter_one_result, &quarter_one_data, &quarter_two_score, &quarter_two_result, &quarter_two_data, &quarter_three_score, &quarter_three_data, &quarter_three_result, &quarter_four_score, &quarter_four_data, &quarter_four_result, &match_result, &match_data, &final_score, &season)
		if err != nil {
			log.WithError(err).Warn("Error mapping db query to TeamStatsWithMatchId")
		}
		data = append(data, TeamStatsWithMatchId{match_id, team_name, quarter_one_score, quarter_one_result, quarter_one_data, quarter_two_score, quarter_two_result, quarter_two_data, quarter_three_score, quarter_three_data, quarter_three_result, quarter_four_score, quarter_four_data, quarter_four_result, match_result, match_data, final_score, season})
	}

	//fmt.Println(data)
	return data
}

func getAllTimeTeamWinsXQuarterAndXOutcome(teamStatsList []TeamStatsWithMatchId, quarter int, quarterResult string, matchResult string) []TeamStatsWithMatchId {
	var filteredTeamList []TeamStatsWithMatchId

	for _, team := range teamStatsList {
		if GetQuarterResult(team, quarter) == quarterResult && team.MatchResult == matchResult {
			filteredTeamList = append(filteredTeamList, team)
		}
	}

	return filteredTeamList
}

// Return every time teamOne plays teamTwo
func getTeamVsTeamStats(db *sql.DB, teamOne string, teamTwo string) ([]MatchStats, error) {
	var data []MatchStats

	var (
		match_id     string
		team_one     string
		team_two     string
		winning_team string
		season       string
	)

	rows, err := db.Query("SELECT match_id, team_one, team_two, winning_team, season from match_stats WHERE (team_one = ? AND team_two = ?) OR (team_two = ? AND team_one = ?)", teamOne, teamTwo, teamOne, teamTwo)
	if err != nil {
		log.WithError(err).Warn("Error querying db")
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	for rows.Next() {
		err := rows.Scan(&match_id, &team_one, &team_two, &winning_team, &season)
		if err != nil {
			return nil, err
		}
		// Query team stats for each team in the match
		teamOneStats, err := GetTeamStats(db, match_id, teamOne)
		if err != nil {
			return nil, err
		}
		teamTwoStats, err := GetTeamStats(db, match_id, teamTwo)
		if err != nil {
			return nil, err
		}
		// Construct MatchStats struct
		matchData := MatchStats{
			MatchID:     match_id,
			TeamOne:     teamOneStats,
			TeamTwo:     teamTwoStats,
			WinningTeam: winning_team,
			Season:      season,
		}
		data = append(data, matchData)
	}

	//fmt.Println(data)
	return data, nil
}

// First iteration of DB stats querying
func fakeMain() {
	//// Connect to DB
	db := connectToDB()
	teamOne := "Collingwood"

	// Get all time team stats
	allTimeTeamStats := getAllTeamStatsFromDb(db, teamOne)
	// all team stats for only last 5 years
	allTeamStatsOfLast5Years := getOnlyLast5YearsOfTeamStats(allTimeTeamStats, 2022)
	// Team wins at half time and loses game
	allTimeTeamWinsSecondQAndLoses := getAllTimeTeamWinsXQuarterAndXOutcome(allTimeTeamStats, 2, "WIN", "LOSS")
	// Team wins half time and wins game
	allTimeTeamWinsSecondQAndWins := getAllTimeTeamWinsXQuarterAndXOutcome(allTimeTeamStats, 2, "WIN", "WIN")

	// The search for MaxProb = Outcome of event and estimated odds.
	// If odds are below bookie odds, then we know we have a high likelihood of +EV bets.
	fmt.Println()
	fmt.Printf("Collingwood wins half time and loses over 30 years: %+v", len(allTimeTeamWinsSecondQAndLoses))
	fmt.Println()
	fmt.Printf("Collingwood wins half time and wins over 30 years: %+v", len(allTimeTeamWinsSecondQAndWins))
	fmt.Println()
	fmt.Printf("Total collingwood games over last  30 years: %+v", len(allTimeTeamStats))
	fmt.Println()
	fmt.Printf("All time team stats of last 5 years count: %+v", len(allTeamStatsOfLast5Years))
	fmt.Println()
	fmt.Println()

	// Disconnect DB
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.WithError(err).Warn("Failed to disconnect DB")
		}
	}(db) // Defer means run this when the wrapping function terminates
}

// Generic method to handle responses requiring json - all of em? lel
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// marshall the payload into a JSON string
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshall json response: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code) // setting response status code
	w.Write(data)
}
