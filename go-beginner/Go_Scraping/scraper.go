package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
    // siteUrl := "https://www.betfair.com.au/exchange/plus/"
    siteUrl := "https://scrapeme.live/shop/"

    // initializing the slice of structs that will contain the scraped data 
    var pokemonProducts []PokemonProduct

	c := colly.NewCollector()
    c.UserAgent = "Go program"
    
    c.OnHTML("#base > div > div.centeredContainer_f19gskyt > div > div.contentAndFooterContainer_fx60ys1.layoutContainersWhenNextRacesRibbonShown_f1k13kpz > div > div > div.childContentWithRibbonShown_fphioe2 > div > div > div > div > div.content_f12s7fql > div:nth-child(4) > div > div.carouselDesktop_f1f1ygnm > div:nth-child(1) > div > div > div > div > div.priceButtons_f1hf8rho > div:nth-child(1) > div.priceButtonContainer_fse8zei > div > div > div > button > div > div > div > span > div", func(e *colly.HTMLElement) { 
        // initializing a new PokemonProduct instance 
        pokemonProduct := PokemonProduct{} 
    
        // scraping the data of interest 
        pokemonProduct.url = e.ChildAttr("a", "href") 
        pokemonProduct.image = e.ChildAttr("img", "src") 
        pokemonProduct.name = e.ChildText("h2") 
        pokemonProduct.price = e.ChildText(".price") 
    
        // adding the product instance with scraped data to the list of products 
        pokemonProducts = append(pokemonProducts, pokemonProduct) 
        
    })

    c.Visit(siteUrl)
}

func loggingSiteDetails(c *colly.Collector) {
    c.OnRequest(func(r *colly.Request) {

        for key, value := range *r.Headers {
            fmt.Printf("%s: %s\n", key, value)
        }

        fmt.Println(r.Method)
    })

    c.OnHTML("title", func(e *colly.HTMLElement) {

        fmt.Println("-----------------------------")

        fmt.Println(e.Text)
    })

    c.OnResponse(func(r *colly.Response) {

        fmt.Println("-----------------------------")

        fmt.Println(r.StatusCode)

        for key, value := range *r.Headers {
            fmt.Printf("%s: %s\n", key, value)
        }
    })
}


// Types

// defining a data structure to store the scraped data 
type PokemonProduct struct { 
	url, image, name, price string 
}