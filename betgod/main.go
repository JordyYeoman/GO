package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
)

// Notes
// 1.For GoColly, you need to set up the listeners first then run your visit() command.
// 2.

func main() {
	fmt.Println("System Online and Ready Sir")

	//baseUrl := "https://afltables.com/afl/seas/"
	// afl links to each season
	//aflSeasonsList := getPageLinks(baseUrl)

	// Test navigation to one page + scrape
	//testUrl := "https://afltables.com/afl/seas/2023.html"
	//getPageStats(testUrl)

	// Test string mapping / interpolation
	handleStringMagic()

	fmt.Println("Scraping finished")
	//fmt.Println(aflSeasonsList)
}

func (s *GameStats) setQuarterScore(quarter, score int) {
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

func extractGameStats(line, team string) GameStats {
	var stats GameStats
	stats.TeamName = team

	fmt.Println(line)

	parts := strings.Fields(line) // Split the line by spaces
	fmt.Println(parts)

	// Final Score
	stats.FinalScore = GetFinalScore(parts[5])

	// Match data
	stats.MatchData =
	// Quarters
	for i := 1; i < 5; i++ {
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
		quarter := i // Quarter 1 corresponds to index 1, Quarter 2 to index 2, and so on
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
	//
	//fmt.Println("Matches: ")
	//fmt.Println(matches)
	//
	//if len(matches) > 0 {
	//	for i, match := range matches[0][1:] {
	//		score, err := strconv.Atoi(match)
	//		if err != nil {
	//			fmt.Println("Error converting score to int:", err)
	//			continue
	//		}
	//		quarter := (i / 2) + 1
	//		if i%2 == 0 {
	//			stats.setQuarterScore(quarter, score)
	//		}
	//	}
	//}

	return stats
}

func handleStringMagic() {
	testStr := "Richmond  1.4   2.4   7.8  8.10  58Thu 16-Mar-2023 7:20 PM (6:20 PM) Att: 88,084 Venue: M.C.G.\nCarlton  3.1   4.6   6.9  8.10  58Match drawn [Match stats]\n"

	lines := strings.Split(testStr, "\n")

	// Extract team name:
	// 1. Separate into two variables by newline

	// 2. Extract team names + create struct to store data
	for _, line := range lines {
		for team := range TeamNames {
			if strings.Contains(line, team) {
				fmt.Println("Found team:", team)
				stats := extractGameStats(line, team)
				fmt.Println("Game Stats:", stats)
				break
			}
		}
	}
	// 3. Extract Quarter by quarter scores
	// 4. Extract Match details
	// 5. Extract result of match
}

func getPageStats(url string) {
	fmt.Println("Scraping: ")
	fmt.Println(url)
	count := 1

	c := colly.NewCollector()
	c.OnHTML("table", func(e *colly.HTMLElement) {
		if count == 3 {
			return
		}

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

func getPageLinks(rootURL string) []string {
	// Scrape AFL season data
	c := colly.NewCollector()
	var aflSeasonsList []string
	// We only want the last 25 seasons (for now)
	// Because the season data starts back in 1897 through to 2023
	linkCount := 0
	totalGames := 126
	totalGamesToScrap := 25

	c.OnHTML("a", func(e *colly.HTMLElement) {
		if linkCount < (totalGames - totalGamesToScrap) {
			linkCount++
			return
		}
		// Returning all <a> tag links on page
		aflSeasonsList = append(aflSeasonsList, fmt.Sprintf("%s%s", rootURL, e.Attr("href")))
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
	})

	// Root site, used to find URL addresses for all seasons
	err := c.Visit(fmt.Sprintf("%s%s", rootURL, "season_idx.html"))
	if err != nil {
		log.Printf("Error occured bra: %+v", err)
		log.Fatal(err)
	}

	return aflSeasonsList
}
