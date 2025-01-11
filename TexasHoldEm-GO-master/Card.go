package main

import "fmt"

type Card struct {
	Suit  rune // 'S', 'H', 'D', 'C' (Kartenfarben: Pik, Herz, Karo, Kreuz)
	Value int  // 2-14 (Ass ist 14)
}

// Gibt die Karte als lesbaren String zurück
func (c Card) String() string {
	var suit string
	switch c.Suit {
	case 'S':
		suit = "Pik"
	case 'H':
		suit = "Herz"
	case 'D':
		suit = "Karo"
	case 'C':
		suit = "Kreuz"
	}
	return fmt.Sprintf("%d von %s", c.Value, suit)
}
