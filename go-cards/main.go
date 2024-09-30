package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Starting Card Deck Sim.....")

	cards := newDeck()
	err := cards.saveToFile("my_cards")
	if err != nil {
		log.Fatal("Error saving cards")
	}

	cards2, _ := getSavedDeckFromFile("my_cards")
	fmt.Println(cards2)

	//playerHand, remainingDeck := cards.deal(5)

	//playerHand.logAllCards()
	//remainingDeck.logAllCards()

	// Call our receiver function on the 'deck' type.
	//cards.logAllCards()
}
