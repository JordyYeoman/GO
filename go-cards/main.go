package main

import "fmt"

func main() {
	fmt.Println("Starting Card Deck Sim.....")

	cards := newDeck()
	playerHand, remainingDeck := cards.deal(5)

	playerHand.logAllCards()
	remainingDeck.logAllCards()

	fmt.Println("Cards: ", cards)

	// Call our receiver function on the 'deck' type.
	cards.logAllCards()
}
