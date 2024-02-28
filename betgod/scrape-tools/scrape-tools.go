package main

//
//import (
//	"database/sql"
//	"fmt"
//	"github.com/gocolly/colly"
//	"github.com/google/uuid"
//	"strconv"
//	"strings"
//)
//
//func scrapePageData() {
//	fmt.Println("System Online and Ready Sir")
//
//	// Generate season data
//	var aflSeasonList []AFLSeasonList
//	totalSeasons := 30 // Total amount of seasons to record
//	lastSeason := 2023 // Season we want to start counting back from
//
//	for i := 0; i < totalSeasons; i++ {
//		var season AFLSeasonList
//		// Convert lastSeason - i to string
//		seasonYear := strconv.Itoa(lastSeason - i)
//
//		// Concatenate the URL parts into a slice of strings
//		urlParts := []string{"https://afltables.com/afl/seas/", seasonYear, ".html"}
//
//		// Join the URL parts with an empty separator
//		url := strings.Join(urlParts, "")
//
//		season.seasonLink = url
//		season.seasonYear = seasonYear
//		// Append the URL to aflSeasonList
//		aflSeasonList = append(aflSeasonList, season)
//	}
//
//	// Create large slice of slices of matches
//	var pageData [][]MatchStats
//	//Loop over each page link and create dataset
//	for _, season := range aflSeasonList {
//		p, err := getPageStats(season.seasonLink, season.seasonYear)
//		// P1
//		// Subroutines should only 'handle' the error if it can recover from it,
//		// Return the error ^up if it can.
//		// Bubble this bad boi
//
//		// P2 - Graceful error handling.
//		// Never log AND return the error.
//		// What else can you do with the err??
//		if err != nil {
//			log.Fatal(err)
//			//log.WithError(err).Warn("Getting page stats is DEAD")
//		}
//
//		pageData = append(pageData, p)
//	}
//
//	fmt.Println(pageData)
//
//	db := connectToDB()
//
//	// Create team_stats table (If required)
//	createTeamStatsDBTables(db)
//	// Create match_stats table (If required)
//	createMatchStatsDBTables(db)
//
//	insertTeamStatsWithMatchId(db, pageData)
//
//	// Disconnect DB
//	defer func(db *sql.DB) {
//		err := db.Close()
//		if err != nil {
//			log.Fatal(err)
//		}
//	}(db) // Defer means run this when the wrapping function terminates
//}
//
//// TODO:
//// Return this out somewhere??!?!
//// Every () should return an error.
//func ExtractMatchStats(gameURL string, season string) (MatchStats, error) {
//	// Struct to contain full match data
//	var MatchResult = MatchStats{
//		Season:  season,
//		MatchID: uuid.New().String(),
//	}
//	teamOneSet := false
//
//	lines := strings.Split(gameURL, "\n")
//
//	for _, line := range lines {
//		// Extract team name and find the actual team name in map
//		tempLine := strings.Fields(line)
//
//		if len(tempLine) < 5 {
//			continue
//		}
//
//		tempStrC := strings.Join(tempLine[:5], " ")
//		// Find which team to use
//		teamToUse := FindCorrectTeamName(tempStrC)
//
//		if teamToUse != "" {
//			//fmt.Println("Found team:", teamToUse)
//
//			// Slice team name from string
//			adjustedLine := RemoveTeamName(line, teamToUse)
//			stats := ExtractTeamStats(adjustedLine, teamToUse)
//
//			if !teamOneSet {
//				MatchResult.TeamOne = TeamStatsWithMatchId{
//					MatchID:            MatchResult.MatchID,
//					TeamName:           stats.TeamName,
//					QuarterOneScore:    stats.QuarterOneScore,
//					QuarterOneResult:   stats.QuarterOneResult,
//					QuarterOneData:     stats.QuarterOneData,
//					QuarterTwoScore:    stats.QuarterTwoScore,
//					QuarterTwoResult:   stats.QuarterTwoResult,
//					QuarterTwoData:     stats.QuarterTwoData,
//					QuarterThreeScore:  stats.QuarterThreeScore,
//					QuarterThreeResult: stats.QuarterThreeResult,
//					QuarterThreeData:   stats.QuarterThreeData,
//					QuarterFourScore:   stats.QuarterFourScore,
//					QuarterFourResult:  stats.QuarterFourResult,
//					QuarterFourData:    stats.QuarterFourData,
//					MatchResult:        stats.MatchResult,
//					MatchData:          stats.MatchData,
//					FinalScore:         stats.FinalScore,
//				}
//				teamOneSet = true
//				continue
//			}
//
//			MatchResult.TeamTwo = TeamStatsWithMatchId{
//				MatchID:            MatchResult.MatchID,
//				TeamName:           stats.TeamName,
//				QuarterOneScore:    stats.QuarterOneScore,
//				QuarterOneResult:   stats.QuarterOneResult,
//				QuarterOneData:     stats.QuarterOneData,
//				QuarterTwoScore:    stats.QuarterTwoScore,
//				QuarterTwoResult:   stats.QuarterTwoResult,
//				QuarterTwoData:     stats.QuarterTwoData,
//				QuarterThreeScore:  stats.QuarterThreeScore,
//				QuarterThreeResult: stats.QuarterThreeResult,
//				QuarterThreeData:   stats.QuarterThreeData,
//				QuarterFourScore:   stats.QuarterFourScore,
//				QuarterFourResult:  stats.QuarterFourResult,
//				QuarterFourData:    stats.QuarterFourData,
//				MatchResult:        stats.MatchResult,
//				MatchData:          stats.MatchData,
//				FinalScore:         stats.FinalScore,
//			}
//			break
//		}
//	}
//
//	// Find match winner
//	tempTeamOneOutcome := MatchResult.TeamOne.FinalScore
//	tempTeamTwoOutcome := MatchResult.TeamTwo.FinalScore
//	// Set Quarter Results for each team (Needed when doing large single team analysis)
//	if MatchResult.TeamOne.QuarterOneScore > MatchResult.TeamTwo.QuarterOneScore {
//		MatchResult.TeamOne.QuarterOneResult = "WIN"
//		MatchResult.TeamTwo.QuarterOneResult = "LOSS"
//	} else if MatchResult.TeamOne.QuarterOneScore < MatchResult.TeamTwo.QuarterOneScore {
//		MatchResult.TeamOne.QuarterOneResult = "LOSS"
//		MatchResult.TeamTwo.QuarterOneResult = "WIN"
//	} else {
//		MatchResult.TeamOne.QuarterOneResult = "DRAW"
//		MatchResult.TeamTwo.QuarterOneResult = "DRAW"
//	}
//	updateQuarterResult(&MatchResult.TeamOne, &MatchResult.TeamTwo)
//
//	if tempTeamOneOutcome > tempTeamTwoOutcome {
//		MatchResult.TeamOne.MatchResult = "WIN"
//		MatchResult.TeamTwo.MatchResult = "LOSS"
//		MatchResult.WinningTeam = MatchResult.TeamOne.TeamName
//	} else if tempTeamOneOutcome < tempTeamTwoOutcome {
//		MatchResult.TeamOne.MatchResult = "LOSS"
//		MatchResult.TeamTwo.MatchResult = "WIN"
//		MatchResult.WinningTeam = MatchResult.TeamTwo.TeamName
//	} else {
//		// It's a draw
//		MatchResult.TeamOne.MatchResult = "DRAW"
//		MatchResult.TeamTwo.MatchResult = "DRAW"
//	}
//
//	// If no team, return nothing
//	if MatchResult.TeamTwo.TeamName == "" {
//		return MatchStats{}, nil
//	}
//
//	return MatchResult, nil
//}
//
//func getPageStats(url string, year string) ([]MatchStats, error) {
//	fmt.Println("Scraping: ")
//	fmt.Println(url)
//	endOfRelevantPage := false // Exiting before finals to ease scraping, can come back and add into data.
//
//	var sliceOMatchStats []MatchStats
//
//	c := colly.NewCollector()
//	var err error
//	c.OnHTML("table", func(e *colly.HTMLElement) {
//		var matchStats MatchStats
//		if endOfRelevantPage { // When we reach the final ladder 'year + season'
//			return
//		}
//
//		// TODO: Error checking here for tables that aren't match stats
//
//		if strings.Contains(e.Text, "Ladder") {
//			return
//		}
//
//		if strings.Contains(e.Text, year+" Ladder") {
//			endOfRelevantPage = true
//		}
//
//		// Every 2nd table on the page has the data we require
//		// Ignore round number + we start at round 1.
//		//fmt.Println(e.Text)
//		matchStats, err = ExtractMatchStats(e.Text, year)
//		if err != nil {
//			// Handle error
//			return
//		}
//		// handle err for above method.
//
//		// Only add match stats if team names exist
//		if matchStats.TeamOne.TeamName != "" {
//			sliceOMatchStats = append(sliceOMatchStats, matchStats)
//		}
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	if err := c.Visit(url); err != nil {
//		return nil, err
//	}
//
//	return sliceOMatchStats, nil
//}
//
//func insertTeamStatsWithMatchId(db *sql.DB, seasons [][]MatchStats) {
//	// For each match, create the match stats and two teams stats entries
//	// Remember to add season number also...
//	for _, matches := range seasons {
//		for _, match := range matches {
//			season, err := strconv.Atoi(match.Season)
//			if err != nil {
//				log.Fatal(err)
//			}
//			// Create a team entry for each team
//			var teamOne = TeamStatsWithMatchId{
//				MatchID:            match.MatchID,
//				TeamName:           match.TeamOne.TeamName,
//				QuarterOneScore:    match.TeamOne.QuarterOneScore,
//				QuarterOneData:     match.TeamOne.QuarterOneData,
//				QuarterOneResult:   match.TeamOne.QuarterOneResult,
//				QuarterTwoScore:    match.TeamOne.QuarterTwoScore,
//				QuarterTwoData:     match.TeamOne.QuarterTwoData,
//				QuarterTwoResult:   match.TeamOne.QuarterTwoResult,
//				QuarterThreeScore:  match.TeamOne.QuarterThreeScore,
//				QuarterThreeData:   match.TeamOne.QuarterThreeData,
//				QuarterThreeResult: match.TeamOne.QuarterThreeResult,
//				QuarterFourScore:   match.TeamOne.QuarterFourScore,
//				QuarterFourData:    match.TeamOne.QuarterFourData,
//				QuarterFourResult:  match.TeamOne.QuarterFourResult,
//				MatchResult:        match.TeamOne.MatchResult,
//				MatchData:          match.TeamOne.MatchData,
//				FinalScore:         match.TeamOne.FinalScore,
//				Season:             season,
//			}
//
//			var teamTwo = TeamStatsWithMatchId{
//				MatchID:            match.MatchID,
//				TeamName:           match.TeamTwo.TeamName,
//				QuarterOneScore:    match.TeamTwo.QuarterOneScore,
//				QuarterOneData:     match.TeamTwo.QuarterOneData,
//				QuarterOneResult:   match.TeamTwo.QuarterOneResult,
//				QuarterTwoScore:    match.TeamTwo.QuarterTwoScore,
//				QuarterTwoData:     match.TeamTwo.QuarterTwoData,
//				QuarterTwoResult:   match.TeamTwo.QuarterTwoResult,
//				QuarterThreeScore:  match.TeamTwo.QuarterThreeScore,
//				QuarterThreeData:   match.TeamTwo.QuarterThreeData,
//				QuarterThreeResult: match.TeamTwo.QuarterThreeResult,
//				QuarterFourScore:   match.TeamTwo.QuarterFourScore,
//				QuarterFourData:    match.TeamTwo.QuarterFourData,
//				QuarterFourResult:  match.TeamTwo.QuarterFourResult,
//				MatchResult:        match.TeamTwo.MatchResult,
//				MatchData:          match.TeamTwo.MatchData,
//				FinalScore:         match.TeamTwo.FinalScore,
//				Season:             season,
//			}
//			// Create separate team entries
//			insertTeamStats(db, teamOne)
//			insertTeamStats(db, teamTwo)
//			// Create match entry
//			insertMatchStats(db, match)
//		}
//	}
//}
//
//func createTeamStatsDBTables(db *sql.DB) {
//	query := `CREATE TABLE IF NOT EXISTS team_stats (
//		id INT AUTO_INCREMENT PRIMARY KEY,
//		match_id TEXT,
// 	team_name TEXT,
//		quarter_one_score BIGINT,
//		quarter_one_result TEXT,
//		quarter_one_data TEXT,
//		quarter_two_score BIGINT,
//		quarter_two_result TEXT,
//		quarter_two_data TEXT,
//		quarter_three_score BIGINT,
//		quarter_three_data TEXT,
//		quarter_three_result TEXT,
//		quarter_four_score BIGINT,
//		quarter_four_data TEXT,
//		quarter_four_result TEXT,
//		match_result TEXT,
//		match_data TEXT,
//		final_score BIGINT,
//		season BIGINT
//	);`
//
//	_, err := db.Exec(query) // Execute query against DB without returning any rows
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
//func createMatchStatsDBTables(db *sql.DB) {
//	query := `CREATE TABLE IF NOT EXISTS match_stats (
//		id INT AUTO_INCREMENT PRIMARY KEY,
//  	match_id TEXT,
//		team_one VARCHAR(255),
//		team_two VARCHAR(255),
//		winning_team TEXT,
//		season TEXT
//	)`
//
//	_, err := db.Exec(query) // Execute query against DB without returning any rows
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
//func insertMatchStats(db *sql.DB, matchStats MatchStats) int {
//	query := "INSERT INTO match_stats (match_id, team_one, team_two, winning_team, season) VALUES (?, ?, ?, ?, ?);"
//	result, err := db.Exec(query, matchStats.MatchID, matchStats.TeamOne.TeamName, matchStats.TeamTwo.TeamName, matchStats.WinningTeam, matchStats.Season)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Retrieve the last inserted ID
//	pk, err := result.LastInsertId()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return int(pk)
//}
