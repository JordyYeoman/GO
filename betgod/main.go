package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/google/uuid"
	"log"
	"strconv"
	"strings"
)

// Notes
// 1.For GoColly, you need to set up the listeners first then run your visit() command.
// 2.

func main() {
	fmt.Println("System Online and Ready Sir")

	var aflSeasonList []string
	totalSeasons := 30
	lastSeason := 2023 // Season we want to start counting back from
	for i := 0; i < totalSeasons; i++ {
		// Convert lastSeason - i to string
		seasonStr := strconv.Itoa(lastSeason - i)

		// Concatenate the URL parts into a slice of strings
		urlParts := []string{"https://afltables.com/afl/seas/", seasonStr, ".html"}

		// Join the URL parts with an empty separator
		url := strings.Join(urlParts, "")

		// Append the URL to aflSeasonList
		aflSeasonList = append(aflSeasonList, url)
	}

	// Test navigation to one page + scrape
	//testUrl := "https://afltables.com/afl/seas/2023.html"
	//getPageStats(testUrl)

	// Test string mapping / interpolation
	//handleStringMagic()

	fmt.Println("Scraping finished")
	//fmt.Println(aflSeasonsList)
}

func ExtractMatchStats(gameURL string) {
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

	fmt.Println(MatchResult)
}

func getPageStats(url string) {
	fmt.Println("Scraping: ")
	fmt.Println(url)
	count := 1
	endOfRelevantPage := false // Exiting before finals to ease scraping, can come back and add into data.

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
		}
		count++
	})

	// Extract data
	// 1. Get team name
	// 2. Get quarter by quarter team scores
	// 3. Get quarter by quarter leader
	// 4. Game Result
	// 5. Game details (string)

	err := c.Visit(url)
	if err != nil {
		log.Printf("Error occured bra: %+v", err)
		log.Fatal(err)
	}
}
