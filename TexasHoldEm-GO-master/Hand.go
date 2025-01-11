package main

import (
	"sort"
	"strings"
)

type Hand struct {
	Cards []Card
	Score int // New field to store the score of the hand
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
	parts := strings.Fields(input)
	cards := []Card{}
	for _, part := range parts {
		if part == "-" {
			continue
		}
		if len(part) != 2 {
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

func (h Hand) isStraight() bool {
	counts := h.countValues()
	consecutive := 0

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

	// Special case: Ace-low Straight (A, 2, 3, 4, 5)
	if counts[4][12] > 0 && counts[4][0] > 0 && counts[4][1] > 0 &&
		counts[4][2] > 0 && counts[4][3] > 0 {
		return true
	}

	return false
}

func (h Hand) isFourOfAKind() bool {
	counts := h.countValues()
	for col := 0; col < 13; col++ {
		if counts[4][col] == 4 {
			return true
		}
	}
	return false
}

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

func (h Hand) isThreeOfAKind() bool {
	counts := h.countValues()
	for col := 0; col < 13; col++ {
		if counts[4][col] == 3 {
			return true
		}
	}
	return false
}

func (h Hand) isTwoPairs() bool {
	counts := h.countValues()
	pairCount := 0

	for col := 0; col < 13; col++ {
		if counts[4][col] == 2 {
			pairCount++
		}
	}

	return pairCount >= 2
}

func (h Hand) isOnePair() bool {
	counts := h.countValues()
	for col := 0; col < 13; col++ {
		if counts[4][col] == 2 {
			return true
		}
	}
	return false
}

func (h Hand) containsHighestStraight() bool {
	counts := h.countValues()
	return counts[4][8] >= 1 && counts[4][9] >= 1 && counts[4][10] >= 1 &&
		counts[4][11] >= 1 && counts[4][12] >= 1
}

func (h Hand) isStraightFlush() bool {
	return h.isFlush() && h.isStraight()
}

func (h Hand) evaluateHand() (string, int) {
	if h.isStraightFlush() {
		return "Straight Flush", 8
	}
	if h.isFourOfAKind() {
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
	if h.isThreeOfAKind() {
		return "Three of a Kind", 3
	}
	if h.isTwoPairs() {
		return "Two Pairs", 2
	}
	if h.isOnePair() {
		return "One Pair", 1
	}
	return "High Card", 0
}

func compareKickers(h1, h2 Hand, community Hand) int {
	combined1 := append(h1.Cards, community.Cards...)
	combined2 := append(h2.Cards, community.Cards...)

	counts1 := h1.countValues()
	counts2 := h2.countValues()

	// Handle tie-breaking for Two Pairs
	if h1.isTwoPairs() && h2.isTwoPairs() {
		pairs1 := findPairs(counts1)
		pairs2 := findPairs(counts2)

		// Compare the higher pair
		if pairs1[0] > pairs2[0] {
			return 1
		} else if pairs1[0] < pairs2[0] {
			return -1
		}

		// Compare the second pair
		if pairs1[1] > pairs2[1] {
			return 1
		} else if pairs1[1] < pairs2[1] {
			return -1
		}

		// Compare the kicker
		return compareHighestRemaining(combined1, combined2)
	}

	// Default tie-breaking for other hands
	return compareHighestRemaining(combined1, combined2)
}

func findPairs(counts [5][13]int) []int {
	pairs := []int{}
	for col := 12; col >= 0; col-- {
		if counts[4][col] == 2 {
			pairs = append(pairs, col+2)
		}
	}
	return pairs
}

func compareHighestRemaining(combined1, combined2 []Card) int {
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

func extractCardValues(h Hand) []int {
	values := make([]int, len(h.Cards))
	for i, card := range h.Cards {
		values[i] = card.Value
	}
	sort.Sort(sort.Reverse(sort.IntSlice(values)))
	return values
}

func (h Hand) CompareTo(other Hand, community Hand) int {
	_, hand1Score := h.evaluateHand()
	_, hand2Score := other.evaluateHand()

	// Compare hand scores
	if hand1Score > hand2Score {
		return 1
	} else if hand1Score < hand2Score {
		return -1
	}

	// Tie-breaking logic
	return compareKickers(h, other, community)
}
