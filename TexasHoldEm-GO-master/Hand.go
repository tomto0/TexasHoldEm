package main

import (
	"sort"
	"strings"
)

type Hand struct {
	Cards    []Card // Enthält die Karten der Hand
	Score    int    // Speichert die Punktzahl der Hand
	HandType string // Speichert den Handtyp (z.B. "Straight Flush")
}

// Zuordnung von Zeichen zu Kartenwerten
var valueMap = map[byte]int{
	'2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'T': 10,
	'J': 11, 'Q': 12, 'K': 13, 'A': 14,
}

// Zuordnung von Zeichen zu Farben (Pik, Herz, Karo, Kreuz)
var suitMap = map[rune]int{
	'S': 0, // Spades
	'H': 1, // Hearts
	'D': 2, // Diamonds
	'C': 3, // Clubs
}

// Konvertiert eine Eingabe in ein Hand-Objekt
func parseHand(input string) Hand {
	parts := strings.Fields(input)
	var cards []Card
	for _, part := range parts {
		if part == "-" {
			continue
		}
		if len(part) != 2 {
			panic("Falsches Format der Karte: " + part)
		}
		value, ok := valueMap[part[1]]
		if !ok {
			panic("Falscher Wert  der Karte: " + string(part[1]))
		}
		cards = append(cards, Card{
			Suit:  rune(part[0]),
			Value: value,
		})
	}
	return Hand{Cards: cards}
}

// Erstellt eine 2D-Array-Darstellung der Kartenhand
func (h Hand) countValues() [5][13]int {
	var handArray [5][13]int
	for _, card := range h.Cards {
		suitIndex := suitMap[card.Suit]
		valueIndex := card.Value - 2
		handArray[suitIndex][valueIndex]++
		handArray[4][valueIndex]++ // Add to totals row
	}
	return handArray
}

// Zählt die Anzahl der Paare in der Hand
func (h Hand) countPairs() int {
	counts := h.countValues()
	pairCount := 0
	for col := 0; col < 13; col++ {
		if counts[4][col] == 2 {
			pairCount++
		}
	}
	return pairCount
}

// Prüft, ob die Hand eine bestimmte Anzahl gleicher Karten enthält
func (h Hand) hasNOfAKind(n int) bool {
	counts := h.countValues()
	for col := 0; col < 13; col++ {
		if counts[4][col] == n {
			return true
		}
	}
	return false
}

// Prüft, ob die Hand ein Flush ist (5 Karten gleicher Farbe)
func (h Hand) isFlush() bool {
	counts := h.countValues()
	for row := 0; row < 4; row++ {
		rowSum := 0
		for col := 0; col < 13; col++ {
			rowSum += counts[row][col]
		}
		if rowSum >= 5 {
			return true
		}
	}
	return false
}

// Prüft, ob die Hand eine Straße enthält (5 aufeinanderfolgende Werte)
func (h Hand) isStraight() bool {
	counts := h.countValues()
	consecutive := 0

	// Prüfen auf normale Straße
	for col := 0; col < 13; col++ {
		if counts[4][col] > 0 {
			consecutive++
			if consecutive == 5 {
				return true
			}
		} else {
			consecutive = 0
		}
	}

	// Prüfen auf spezielle Straße (Ass bis 5)
	return counts[4][12] > 0 && counts[4][0] > 0 && counts[4][1] > 0 &&
		counts[4][2] > 0 && counts[4][3] > 0
}

// Prüft, ob die Hand eine Straße und einen Flush enthält (Straight Flush)
func (h Hand) isStraightFlush() bool {
	return h.isFlush() && h.isStraight()
}

// Prüft, ob die Hand ein Full House ist
func (h Hand) isFullHouse() bool {
	counts := h.countValues()
	hasThree := false
	hasPair := false

	for col := 0; col < 13; col++ {
		if counts[4][col] == 3 {
			hasThree = true
		} else if counts[4][col] == 2 {
			hasPair = true
		}
	}

	return hasThree && hasPair
}

func (h Hand) isRoyalFlush() bool {
	// Ein Royal Flush ist ein Straight Flush mit den höchsten Karten: 10, J, Q, K, A
	counts := h.countValues()
	if counts[4][8] >= 1 && counts[4][9] >= 1 && counts[4][10] >= 1 && counts[4][11] >= 1 && counts[4][12] >= 1 {
		// Prüfen, ob die Karten dieselbe Farbe haben (Flush)
		for suit := 0; suit < 4; suit++ {
			if counts[suit][8] >= 1 && counts[suit][9] >= 1 && counts[suit][10] >= 1 && counts[suit][11] >= 1 && counts[suit][12] >= 1 {
				return true
			}
		}
	}
	return false
}

// compareKickers vergleicht die Kicker der beiden Hände unter Berücksichtigung der Gemeinschaftskarten.
// Dabei werden die Karten der Spielerhände mit den Gemeinschaftskarten kombiniert und nach Wert sortiert.
// Anschließend werden die höchsten Karten der kombinierten Hände paarweise verglichen, um den Gewinner zu bestimmen.
// Rückgabewerte:
// - 1: Die erste Hand (hand1) hat den besseren Kicker.
// - -1: Die zweite Hand (hand2) hat den besseren Kicker.
// - 0: Beide Hände sind gleichwertig.
func compareKickers(hand1 Hand, hand2 Hand, community Hand) int {
	// Zuerst nach Handtypen vergleichen
	if hand1.Score != hand2.Score {
		if hand1.Score > hand2.Score {
			return 1
		}
		return -1
	}

	// Wenn Handtypen gleich sind, Kicker vergleichen
	allCards1 := append(hand1.Cards, community.Cards...)
	allCards2 := append(hand2.Cards, community.Cards...)
	sort.Slice(allCards1, func(i, j int) bool { return allCards1[i].Value > allCards1[j].Value })
	sort.Slice(allCards2, func(i, j int) bool { return allCards2[i].Value > allCards2[j].Value })

	for i := 0; i < len(allCards1) && i < len(allCards2); i++ {
		if allCards1[i].Value != allCards2[i].Value {
			if allCards1[i].Value > allCards2[i].Value {
				return 1
			}
			return -1
		}
	}
	return 0
}

// evaluateHand bewertet die Stärke einer Hand basierend auf den gegebenen Karten.
// Sie kombiniert die Karten der Hand mit den Gemeinschaftskarten und prüft auf die stärkste
func (h Hand) evaluateHand() (string, int) {
	handType := 0 // Standardmäßig "High Card"

	switch {
	case h.isRoyalFlush():
		handType = 9
	case h.isStraightFlush():
		handType = 8
	case h.hasNOfAKind(4):
		handType = 7
	case h.isFullHouse():
		handType = 6
	case h.isFlush():
		handType = 5
	case h.isStraight():
		handType = 4
	case h.hasNOfAKind(3):
		handType = 3
	case h.countPairs() == 2:
		handType = 2
	case h.countPairs() == 1:
		handType = 1
	}
	return handTypeToString(handType), handType
}

func handTypeToString(handType int) string {
	types := []string{"High Card", "One Pair", "Two Pairs", "Three of a Kind", "Straight", "Flush", "Full House", "Four of a Kind", "Straight Flush", "Royal Flush"}
	if handType >= 0 && handType < len(types) {
		return types[handType]
	}
	return "Unbekannt"
}

func (h Hand) CompareTo(other Hand, community Hand) int {
	if h.Score != other.Score {
		if h.Score > other.Score {
			return 1
		}
		return -1
	}
	return compareKickers(h, other, community)
}

func (h Hand) String() string {
	// Sortiert die Karten absteigend
	sort.Slice(h.Cards, func(i, j int) bool {
		return h.Cards[i].Value > h.Cards[j].Value
	})

	// Generiert einen String aus den Kartenzahlen
	var cardStrings []string
	for _, card := range h.Cards {
		cardStrings = append(cardStrings, card.String())
	}
	return strings.Join(cardStrings, ", ")
}

// HandList ist ein benutzerdefinierter Typ für eine Liste von 	Händen
type HandList []Hand

func (hl HandList) Len() int {
	return len(hl)
}

func (hl HandList) Less(i, j int) bool {
	if hl[i].Score == hl[j].Score {
		community := Hand{}                                // Hier die tatsächlichen Community-Karten übergeben
		return compareKickers(hl[i], hl[j], community) > 0 // Beachte: > 0 für bessere Hand
	}
	return hl[i].Score > hl[j].Score // Höherer Handtyp kommt zuerst
}

func (hl HandList) Swap(i, j int) {
	hl[i], hl[j] = hl[j], hl[i]
}
