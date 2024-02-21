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
		winning_team TEXT,
		season TEXT
	)`

	_, err := db.Exec(query) // Execute query against DB without returning any rows
	if err != nil {
		log.Fatal(err)
	}
}
