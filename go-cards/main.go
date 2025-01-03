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
	//fmt.Println(cards2)

	err = cards2.shuffle()
	if err != nil {
		log.Fatal("unable to shuffle deck")
	}

	//playerHand, remainingDeck := cards.deal(5)

	//playerHand.logAllCards()
	//remainingDeck.logAllCards()

	// Call our receiver function on the 'deck' type.
	cards2.logAllCards()

	// Test method
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	IsEvenOrOdd(numbers)

	// Basic person struct manipulation
	DoStuff()
}
