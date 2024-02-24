package main

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
)

// TODO: Add weighting from betfair model: https://www.betfair.com.au/hub/sports/afl/afl-predictions-model/

func getMatchesWhereTeamWonFirstQuarterAndWon(matches []MatchStats, teamName string) []MatchStats {
	var filteredMatches []MatchStats

	for _, match := range matches {
		if match.TeamOne.TeamName == teamName && match.TeamOne.QuarterOneResult == "WIN" && match.TeamOne.MatchResult == "WIN" || match.TeamTwo.TeamName == teamName && match.TeamTwo.QuarterOneResult == "WIN" && match.TeamOne.MatchResult == "WIN" {
			filteredMatches = append(filteredMatches, match)
		}
	}

	return filteredMatches
}

func getQuarterResult(team TeamStatsWithMatchId, quarter int) string {
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

func getMatchesWhereTeamOneQuarterXAndLost(matches []MatchStats, teamName string, quarter int) []MatchStats {
	var filteredMatches []MatchStats

	for _, match := range matches {
		team := match.TeamOne
		if match.TeamTwo.TeamName == teamName {
			team = match.TeamTwo
		}

		if team.TeamName == teamName && getQuarterResult(team, quarter) == "WIN" && team.MatchResult == "LOSS" {
			filteredMatches = append(filteredMatches, match)
		}
	}

	return filteredMatches
}

func getMatchesWhereTeamWonFirstQuarterAndLost(matches []MatchStats, teamName string) []MatchStats {
	var filteredMatches []MatchStats

	for _, match := range matches {
		if match.TeamOne.TeamName == teamName && match.TeamOne.QuarterOneResult == "WIN" && match.TeamOne.MatchResult == "LOSS" || match.TeamTwo.TeamName == teamName && match.TeamTwo.QuarterOneResult == "WIN" && match.TeamOne.MatchResult == "LOSS" {
			filteredMatches = append(filteredMatches, match)
		}
	}

	return filteredMatches
}

func getMatchesWhereTeamOneWonFirstQuarter(matches []MatchStats, teamName string) []MatchStats {
	var filteredMatches []MatchStats

	for _, match := range matches {
		// Check if teamOne won the first quarter
		if match.TeamOne.TeamName == teamName && match.TeamOne.QuarterOneResult == "WIN" || match.TeamTwo.TeamName == teamName && match.TeamTwo.QuarterOneResult == "WIN" {
			filteredMatches = append(filteredMatches, match)
		}
	}

	return filteredMatches
}

// FYI - All queries run over 30 seasons of scraped data.
func main() {
	fmt.Println("Connect to DB and analyse data")

	// 1. Create queries to question teams quarter performance
	// 2. Create queries to question team vs team quarter performance
	// 3. Create query to get half / fulltime result average of teams
	// 4. Create query to question half / fulltime result of 2 specific teams

	// Connect to DB
	db := connectToDB()
	teamOne := "Collingwood"
	teamTwo := "Carlton"

	// allTimeTeamStats := getAllTeamStats(db, teamOne)
	allTimeTeamVsTeamStats, allTimeTeamErr := getTeamVsTeamStats(db, teamOne, teamTwo)
	if allTimeTeamErr != nil {
		log.Fatal(allTimeTeamErr)
	}

	filteredMatches := getMatchesWhereTeamOneWonFirstQuarter(allTimeTeamVsTeamStats, teamOne)
	filteredMatchesTwo := getMatchesWhereTeamOneWonFirstQuarter(allTimeTeamVsTeamStats, teamTwo)
	filteredMatchesThree := getMatchesWhereTeamWonFirstQuarterAndWon(allTimeTeamVsTeamStats, teamOne)
	filteredMatchesFour := getMatchesWhereTeamWonFirstQuarterAndLost(allTimeTeamVsTeamStats, teamOne)
	filteredMatchesFive := getMatchesWhereTeamOneQuarterXAndLost(allTimeTeamVsTeamStats, teamOne, 1)
	filteredMatchesSix := getMatchesWhereTeamOneQuarterXAndLost(allTimeTeamVsTeamStats, teamOne, 2)

	fmt.Println()
	fmt.Printf("Number of times collingwood won first quarter against carlton in last 30 years: %+v, out of total games: %+v", len(filteredMatches), len(allTimeTeamVsTeamStats))
	fmt.Println()
	fmt.Printf("Number of times carlton won first quarter against collingwood in last 30 years: %+v, out of total games: %+v", len(filteredMatchesTwo), len(allTimeTeamVsTeamStats))
	fmt.Println()
	fmt.Printf("Number of times collingwood won first quarter against carlton and WON in last 30 years: %+v, out of total games: %+v", len(filteredMatchesThree), len(allTimeTeamVsTeamStats))
	fmt.Println()
	fmt.Printf("Number of times collingwood won first quarter against carlton and LOST in last 30 years: %+v, out of total games: %+v", len(filteredMatchesFour), len(allTimeTeamVsTeamStats))
	fmt.Println()
	fmt.Printf("Number of times collingwood won first quarter against carlton and LOST in last 30 years: %+v, out of total games: %+v", len(filteredMatchesFive), len(allTimeTeamVsTeamStats))
	fmt.Println()
	fmt.Printf("Number of times collingwood won second quarter against carlton and LOST in last 30 years: %+v, out of total games: %+v", len(filteredMatchesSix), len(allTimeTeamVsTeamStats))
	fmt.Println()
	//
	//var quarterOneWins []TeamStatsWithMatchId
	//var quarterOneLosses []TeamStatsWithMatchId
	//var quarterOneDraws []TeamStatsWithMatchId
	//
	//for _, team := range allTimeTeamStats {
	//	if team.QuarterOneResult == "LOSS" {
	//		quarterOneLosses = append(quarterOneLosses, team)
	//	} else if team.QuarterOneResult == "WIN" {
	//		quarterOneWins = append(quarterOneWins, team)
	//	} else {
	//		quarterOneDraws = append(quarterOneDraws, team)
	//	}
	//}
	//
	//fmt.Println()
	//fmt.Printf("%+v Quarter One Wins %+v", teamName, len(quarterOneWins))
	//fmt.Println()
	//fmt.Printf("%+v Quarter One Losses %+v", teamName, len(quarterOneLosses))
	//fmt.Println()
	//fmt.Printf("%+v Quarter One Draws %+v", teamName, len(quarterOneDraws))
	//fmt.Println()

	// Disconnect DB
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db) // Defer means run this when the wrapping function terminates
}

// Function to fetch team stats for a given match_id and team
func getTeamStats(db *sql.DB, matchID string, teamName string) (TeamStatsWithMatchId, error) {
	var teamStats TeamStatsWithMatchId

	// Query team stats
	//err := db.QueryRow("SELECT * FROM team_stats WHERE match_id = ? AND team_name = ?", matchID, teamName).Scan(&teamStats.MatchID, &teamStats.TeamName, &teamStats.QuarterOneScore, &teamStats.QuarterOneResult, &teamStats.QuarterOneData, &teamStats.QuarterTwoScore, &teamStats.QuarterTwoResult, &teamStats.QuarterTwoData, &teamStats.QuarterThreeScore, &teamStats.QuarterThreeResult, &teamStats.QuarterThreeData, &teamStats.QuarterFourScore, &teamStats.QuarterFourResult, &teamStats.QuarterFourData, &teamStats.MatchResult, &teamStats.MatchData, &teamStats.FinalScore)
	err := db.QueryRow("SELECT match_id, team_name, quarter_one_score, quarter_one_result, quarter_one_data, quarter_two_score, quarter_two_result, quarter_two_data, quarter_three_score, quarter_three_result, quarter_three_data, quarter_four_score, quarter_four_result, quarter_four_data, match_result, match_data, final_score FROM team_stats WHERE match_id = ? AND team_name = ?", matchID, teamName).Scan(&teamStats.MatchID, &teamStats.TeamName, &teamStats.QuarterOneScore, &teamStats.QuarterOneResult, &teamStats.QuarterOneData, &teamStats.QuarterTwoScore, &teamStats.QuarterTwoResult, &teamStats.QuarterTwoData, &teamStats.QuarterThreeScore, &teamStats.QuarterThreeResult, &teamStats.QuarterThreeData, &teamStats.QuarterFourScore, &teamStats.QuarterFourResult, &teamStats.QuarterFourData, &teamStats.MatchResult, &teamStats.MatchData, &teamStats.FinalScore)
	if err != nil {
		return teamStats, err
	}

	return teamStats, nil
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
		teamOneStats, err := getTeamStats(db, match_id, teamOne)
		if err != nil {
			return nil, err
		}
		teamTwoStats, err := getTeamStats(db, match_id, teamTwo)
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

func getAllTeamStats(db *sql.DB, teamName string) []TeamStatsWithMatchId {
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

	rows, err := db.Query("SELECT match_id, team_name, quarter_one_score, quarter_one_result, quarter_one_data, quarter_two_score, quarter_two_result, quarter_two_data, quarter_three_score, quarter_three_data, quarter_three_result, quarter_four_score, quarter_four_data, quarter_four_result, match_result, match_data, final_score from team_stats WHERE team_name = ?", teamName)
	if err != nil {
		log.WithError(err).Warn("Error querying db")
	}

	defer rows.Close()

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
