package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type AFLSeasonList struct {
	seasonLink string
	seasonYear string
}

func main() {
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
			log.WithError(err).Warn("Getting page stats is DEAD")
		}

		pageData = append(pageData, p)
	}

	fmt.Println(pageData)
	// Connect to DB
	handleDBConnection(pageData)
}

// TODO:
// Return this out somewhere??!?!
// Every () should return an error.
func ExtractMatchStats(gameURL string, season string) (MatchStats, error) {
	// Struct to contain full match data
	var MatchResult = MatchStats{
		Season: season,
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
				MatchResult.TeamOne = stats
				teamOneSet = true
				continue
			}

			MatchResult.TeamTwo = stats
			break
		}
	}

	MatchResult.MatchID = uuid.New().String()
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
