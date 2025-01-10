package main

import (
	"sort"
	"strings"
)

type Hand struct {
	Cards []Card
	Score int
}

var valueMap = map[byte]int{
	'2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'T': 10,
	'J': 11, 'Q': 12, 'K': 13, 'A': 14,
}

var suitMap = map[rune]int{
	'S': 0, // Spades
	'H': 1, // Hearts
	'D': 2, // Diamonds
	'C': 3, // Clubs
}

func parseHand(input string) Hand {
	parts := strings.Fields(input) // Split input by spaces
	cards := []Card{}
	for _, part := range parts {
		if part == "-" { // Skip placeholders
			continue
		}
		if len(part) != 2 { // Check if each card is exactly 2 characters long
			panic("Invalid card format: " + part)
		}
		value, ok := valueMap[part[1]]
		if !ok {
			panic("Invalid card value: " + string(part[1]))
		}
		cards = append(cards, Card{
			Suit:  rune(part[0]),
			Value: value,
		})
	}
	return Hand{Cards: cards}
}

func (h Hand) isFlush() bool {
	suit := h.Cards[0].Suit
	for _, card := range h.Cards {
		if card.Suit != suit {
			return false
		}
	}
	return true
}

func (h Hand) isStraight() bool {
	values := extractCardValues(h)
	sort.Ints(values)

	// Regular straight check
	isRegularStraight := true
	for i := 1; i < len(values); i++ {
		if values[i] != values[i-1]+1 {
			isRegularStraight = false
			break
		}
	}

	// Ace-low straight check
	isAceLowStraight := values[0] == 2 && values[1] == 3 && values[2] == 4 && values[3] == 5 && values[len(values)-1] == 14

	return isRegularStraight || isAceLowStraight
}

func (h Hand) countValues() [4][13]int {
	var handArray [4][13]int
	for _, card := range h.Cards {
		suitIndex := suitMap[card.Suit]
		valueIndex := card.Value - 2
		handArray[suitIndex][valueIndex]++
	}
	return handArray
}

func (h Hand) evaluateHand() string {
	isFlush := h.isFlush()
	isStraight := h.isStraight()
	counts := h.countValues()

	// Check for Straight Flush
	if isFlush && isStraight {
		return "Straight Flush"
	}

	// Check other hand types
	if containsCount(counts, 4) {
		return "Four of a Kind"
	}
	if containsCount(counts, 3) && containsCount(counts, 2) {
		return "Full House"
	}
	if isFlush {
		return "Flush"
	}
	if isStraight {
		return "Straight"
	}
	if containsCount(counts, 3) {
		return "Three of a Kind"
	}
	if countPairs(counts) == 2 {
		return "Two Pairs"
	}
	if containsCount(counts, 2) {
		return "One Pair"
	}
	return "High Card"
}

func containsCount(counts [4][13]int, n int) bool {
	for _, row := range counts {
		for _, count := range row {
			if count == n {
				return true
			}
		}
	}
	return false
}

func countPairs(counts [4][13]int) int {
	pairs := 0
	for _, row := range counts {
		for _, count := range row {
			if count == 2 {
				pairs++
			}
		}
	}
	return pairs
}

func (h Hand) CompareTo(other Hand, community Hand) int {
	hand1Type := h.evaluateHand()
	hand2Type := other.evaluateHand()

	handRanks := map[string]int{
		"Royal Flush":     10,
		"Straight Flush":  9,
		"Four of a Kind":  8,
		"Full House":      7,
		"Flush":           6,
		"Straight":        5,
		"Three of a Kind": 4,
		"Two Pairs":       3,
		"One Pair":        2,
		"High Card":       1,
	}

	// Compare hand types based on ranking
	if handRanks[hand1Type] > handRanks[hand2Type] {
		return 1
	} else if handRanks[hand1Type] < handRanks[hand2Type] {
		return -1
	}

	// If hand types are equal, compare kickers
	return compareKickers(h, other, community)
}

func compareKickers(h1, h2 Hand, community Hand) int {
	// Combine player hands with community cards
	combined1 := append(h1.Cards, community.Cards...)
	combined2 := append(h2.Cards, community.Cards...)

	// Extract and sort card values
	values1 := extractCardValues(Hand{Cards: combined1})
	values2 := extractCardValues(Hand{Cards: combined2})
	sort.Sort(sort.Reverse(sort.IntSlice(values1)))
	sort.Sort(sort.Reverse(sort.IntSlice(values2)))

	// Compare sorted values card-by-card
	for i := 0; i < len(values1) && i < len(values2); i++ {
		if values1[i] > values2[i] {
			return 1
		} else if values1[i] < values2[i] {
			return -1
		}
	}

	// If all values are equal, return 0 for a tie
	return 0
}

// Helper to compare card values
func compareCardValues(h1, h2 Hand) int {
	// Extract the card values and sort them
	h1Values := extractCardValues(h1)
	h2Values := extractCardValues(h2)

	// Compare highest values, then second highest, etc.
	for i := 0; i < len(h1Values); i++ {
		if h1Values[i] > h2Values[i] {
			return 1
		} else if h1Values[i] < h2Values[i] {
			return -1
		}
	}
	return 0
}

// Helper to extract and sort card values for comparison
func extractCardValues(h Hand) []int {
	values := make([]int, len(h.Cards))
	for i, card := range h.Cards {
		values[i] = card.Value
	}
	sort.Sort(sort.Reverse(sort.IntSlice(values))) // Highest to lowest
	return values
}

type Hands []Hand

func (h Hands) Len() int      { return len(h) }
func (h Hands) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h Hands) Less(i, j int) bool {
	community := Hand{} // Placeholder. Use actual community cards as needed.
	return h[i].CompareTo(h[j], community) > 0
}
