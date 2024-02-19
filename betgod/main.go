package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/google/uuid"
	"log"
	"strconv"
	"strings"
)

type AFLSeasonList struct {
	seasonLink string
	seasonYear string
}

func main() {
	fmt.Println("System Online and Ready Sir")

	var aflSeasonList []AFLSeasonList
	totalSeasons := 30
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

	//var pageData []MatchStats
	// Loop over each page link and create dataset
	//for data := range aflSeasonList {
	//	pageData = append(pageData, getPageStats(data))
	//}

	// getPageStats("https://afltables.com/afl/seas/2023.html", "2023")

	testStr := "Port Adelaide  4.1   5.6   8.7  9.16  70Sat 16-Sep-2023 7:10 PM (7:40 PM) Att: 45,520 Venue: Adelaide Oval\nGreater Western Sydney  4.4  9.11 11.15 13.15  93Greater Western Sydney won by 23 pts [Match stats]\n"
	r := ExtractMatchStats(testStr)
	fmt.Println()
	fmt.Printf("Match Stats: %+v\n", r)
	fmt.Println()

	fmt.Println("Scraping finished")
}

func ExtractMatchStats(gameURL string) MatchStats {
	// Struct to contain full match data
	var MatchResult = MatchStats{}
	teamOneSet := false

	lines := strings.Split(gameURL, "\n")

	for _, line := range lines {
		// Extract team name and find the actual team name in map
		tempLine := strings.Fields(line)

		if len(tempLine) < 5 {
			continue
		}

		tempStrC := strings.Join(tempLine[:5], " ")

		//fmt.Println()
		//fmt.Printf("Temp String: %+v", tempStrC)
		//fmt.Println()

		// Find which team to use
		teamToUse := FindCorrectTeamName(tempStrC)

		if teamToUse != "" {
			fmt.Println("Found team:", teamToUse)

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

	return MatchResult
}

func getPageStats(url string, year string) []MatchStats {
	fmt.Println("Scraping: ")
	fmt.Println(url)
	count := 1
	endOfRelevantPage := false // Exiting before finals to ease scraping, can come back and add into data.

	var sliceOMatchStats []MatchStats

	c := colly.NewCollector()
	c.OnHTML("table", func(e *colly.HTMLElement) {
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
		if count%2 == 0 {
			//fmt.Println(e.Text)
			sliceOMatchStats = append(sliceOMatchStats, ExtractMatchStats(e.Text))
		}
		count++
	})

	err := c.Visit(url)
	if err != nil {
		log.Printf("Error occured bra: %+v", err)
		log.Fatal(err)
	}

	return sliceOMatchStats
}
