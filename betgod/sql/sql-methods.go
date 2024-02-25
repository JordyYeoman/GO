package sqlMethods

import (
	"database/sql"
	"log"
)

func createTeamStatsDBTables(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS team_stats (
		id INT AUTO_INCREMENT PRIMARY KEY,
		match_id TEXT,
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
		final_score BIGINT,
		season BIGINT
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
		winning_team TEXT,
		season TEXT
	)`

	_, err := db.Exec(query) // Execute query against DB without returning any rows
	if err != nil {
		log.Fatal(err)
	}
}

//func insertTeamStatsWithMatchId(db *sql.DB, seasons [][]MatchStats) {
//	// For each match, create the match stats and two teams stats entries
//	// Remember to add season number also...
//	for _, matches := range seasons {
//		for _, match := range matches {
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
//			}
//			// Create separate team entries
//			insertTeamStats(db, teamOne)
//			insertTeamStats(db, teamTwo)
//			// Create match entry
//			insertMatchStats(db, match)
//		}
//	}
//}
