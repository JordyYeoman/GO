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

	poo := getPageStats("https://afltables.com/afl/seas/2023.html")
	fmt.Println(poo)

	fmt.Println("Scraping finished")
}

func ExtractMatchStats(gameURL string) MatchStats {
	// Struct to contain full match data
	var MatchResult = MatchStats{}
	teamOneSet := false

	lines := strings.Split(gameURL, "\n")

	for _, line := range lines {
		for team := range TeamNames {
			// Temp string so we don't include the final match summary and accidentally include a team twice.
			tempLine := strings.Fields(line)
			var tempString string
			// Check if there are at least 5 elements in tempLine to avoid out of bounds error
			if len(tempLine) < 5 {
				// Handle error if required AND/OR Continue to the next iteration
				continue
			}
			// Join the first 5 elements of tempLine with a space separator
			tempString = strings.Join(tempLine[:5], " ")

			if strings.Contains(tempString, team) {
				fmt.Println()
				fmt.Printf("Found Team: %+v", team)
				fmt.Println()
				// Slice team name from string
				adjustedLine := RemoveTeamName(line, team)
				stats := ExtractTeamStats(adjustedLine, team)
				// First team is away
				if !teamOneSet {
					MatchResult.TeamOne = stats
					teamOneSet = true
				}
				// Second team is home
				MatchResult.TeamTwo = stats

				break
			}
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

func getPageStats(url string) []MatchStats {
	fmt.Println("Scraping: ")
	fmt.Println(url)
	count := 1
	endOfRelevantPage := false // Exiting before finals to ease scraping, can come back and add into data.

	var sliceOMatchStats []MatchStats

	c := colly.NewCollector()
	c.OnHTML("table", func(e *colly.HTMLElement) {
		if endOfRelevantPage {
			return
		}

		//if strings.Contains(e.Text, "Ladder") {
		//	endOfRelevantPage = true
		//}
		// Every 2nd table on the page has the data we require
		// Ignore round number + we start at round 1.
		if count%2 == 0 {
			fmt.Println(e.Text)
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
