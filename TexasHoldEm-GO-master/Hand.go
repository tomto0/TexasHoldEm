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

func compareKickers(h1, h2 Hand, community Hand) int {
	// Kombiniere die Handkarten mit den Gemeinschaftskarten
	combined1 := append(h1.Cards, community.Cards...)
	combined2 := append(h2.Cards, community.Cards...)

	// Extrahiere und sortiere die Kartenwerte
	values1 := extractCardValues(Hand{Cards: combined1})
	values2 := extractCardValues(Hand{Cards: combined2})

	// Vergleiche die Kartenwerte (höchste zuerst)
	for i := 0; i < len(values1) && i < len(values2); i++ {
		if values1[i] > values2[i] {
			return 1
		} else if values1[i] < values2[i] {
			return -1
		}
	}

	// Wenn beide Hände gleich sind
	return 0
}

// Bewertet die Hand und gibt den Handtyp und die Punktzahl zurück
func (h Hand) evaluateHand() (string, int) {
	if h.isRoyalFlush() {
		return "Royal Flush", 9
	}
	if h.isStraightFlush() {
		return "Straight Flush", 8
	}
	if h.hasNOfAKind(4) {
		return "Four of a Kind", 7
	}
	if h.isFullHouse() {
		return "Full House", 6
	}
	if h.isFlush() {
		return "Flush", 5
	}
	if h.isStraight() {
		return "Straight", 4
	}
	if h.hasNOfAKind(3) {
		return "Three of a Kind", 3
	}
	if h.countPairs() >= 2 {
		return "Two Pairs", 2
	}
	if h.hasNOfAKind(2) {
		return "One Pair", 1
	}
	return "High Card", 0
}

// Vergleicht zwei Hände anhand der Punktzahl und Kickern
func compareHighestCards(h1, h2 Hand, community Hand) int {
	combined1 := append(h1.Cards, community.Cards...)
	combined2 := append(h2.Cards, community.Cards...)

	values1 := extractCardValues(Hand{Cards: combined1})
	values2 := extractCardValues(Hand{Cards: combined2})

	sort.Sort(sort.Reverse(sort.IntSlice(values1)))
	sort.Sort(sort.Reverse(sort.IntSlice(values2)))

	for i := 0; i < len(values1) && i < len(values2); i++ {
		if values1[i] > values2[i] {
			return 1
		} else if values1[i] < values2[i] {
			return -1
		}
	}
	return 0
}

// Extrahiert die Kindergartener aus einer Hand
func extractCardValues(h Hand) []int {
	values := make([]int, len(h.Cards))
	for i, card := range h.Cards {
		values[i] = card.Value
	}
	sort.Sort(sort.Reverse(sort.IntSlice(values))) // Absteigend sortieren
	return values
}

// CompareTo Vergleicht zwei Hände
func (h Hand) CompareTo(other Hand, community Hand) int {
	_, hand1Score := h.evaluateHand()
	_, hand2Score := other.evaluateHand()

	// Standardmäßiger Vergleich nach Punktzahl
	if hand1Score > hand2Score {
		return 1
	} else if hand1Score < hand2Score {
		return -1
	}

	// Fallback für gleiche Punktzahl: Kartenwert vergleichen
	return compareKickers(h, other, community)
}

func (h Hand) String() string {
	cardStrings := []string{}
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
		// Tiebreaker logic if scores are equal
		return hl[i].HandType < hl[j].HandType
	}
	return hl[i].Score > hl[j].Score
}
func (hl HandList) Swap(i, j int) {
	hl[i], hl[j] = hl[j], hl[i]
}
