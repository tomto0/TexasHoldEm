package main

import "fmt"

type Card struct {
	Suit  rune // 'S', 'H', 'D', 'C' (Kartenfarben: Pik, Herz, Karo, Kreuz)
	Value int  // 2-14 (Ass ist 14)
}

// Gibt die Karte als lesbaren String zurück
func (c Card) String() string {
	return fmt.Sprintf("%d von %s", c.Value, suitToString(c.Suit))
}

func suitToString(suit rune) string {
	switch suit {
	case 'S':
		return "Pik"
	case 'H':
		return "Herz"
	case 'D':
		return "Karo"
	case 'C':
		return "Kreuz"
	default:
		return "Unbekannt"
	}
}
