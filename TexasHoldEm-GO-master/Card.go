package main

import "fmt"

type Card struct {
	Suit  rune // 'S', 'H', 'D', 'C'
	Value int  // 2-14 (Ace is 14)
}

func (c Card) String() string {
	var suit string
	switch c.Suit {
	case 'S':
		suit = "Spades"
	case 'H':
		suit = "Hearts"
	case 'D':
		suit = "Diamonds"
	case 'C':
		suit = "Clubs"
	}
	return fmt.Sprintf("%d of %s", c.Value, suit)
}
