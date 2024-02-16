package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
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
	testUrl := "https://afltables.com/afl/seas/2023.html"
	getPageStats(testUrl)

	fmt.Println("Scraping finished")
	//fmt.Println(aflSeasonsList)
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
		// Ignore round number
		if count%2 == 0 {
			e.ForEach("table", func(t int, z *colly.HTMLElement) {
				fmt.Println("Testing")
				fmt.Println(t)
				fmt.Println(z)
			})
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
