package main

import (
	"sort"
	"strings"
)

type Hand struct {
	Cards []Card
}

var valueMap = map[byte]int{
	'2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'T': 10,
	'J': 11, 'Q': 12, 'K': 13, 'A': 14,
}

func parseHand(input string) Hand {
	parts := strings.Split(input, " ")
	cards := make([]Card, len(parts))
	for i, part := range parts {
		cards[i] = Card{
			Suit:  rune(part[0]),
			Value: valueMap[part[1]],
		}
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
	values := make([]int, len(h.Cards))
	for i, card := range h.Cards {
		values[i] = card.Value
	}
	sort.Ints(values)

	// Check regular straight
	isRegularStraight := true
	for i := 1; i < len(values); i++ {
		if values[i] != values[i-1]+1 {
			isRegularStraight = false
			break
		}
	}

	// Check Ace-low straight (A, 2, 3, 4, 5)
	isAceLowStraight := values[0] == 2 && values[1] == 3 && values[2] == 4 && values[3] == 5 && values[4] == 14

	return isRegularStraight || isAceLowStraight
}

func (h Hand) countValues() map[int]int {
	counts := make(map[int]int)
	for _, card := range h.Cards {
		counts[card.Value]++
	}
	return counts
}

func (h Hand) evaluateHand() string {
	isFlush := h.isFlush()
	isStraight := h.isStraight()
	counts := h.countValues()

	switch {
	case isFlush && isStraight && h.Cards[0].Value == 10:
		return "Royal Flush"
	case isFlush && isStraight:
		return "Straight Flush"
	case containsCount(counts, 4):
		return "Poker"
	case containsCount(counts, 3) && containsCount(counts, 2):
		return "Full House"
	case isFlush:
		return "Flush"
	case isStraight:
		return "Straight"
	case containsCount(counts, 3):
		return "Three of a Kind"
	case countPairs(counts) == 2:
		return "Two Pairs"
	case containsCount(counts, 2):
		return "One Pair"
	default:
		return "High Card"
	}
}

func containsCount(counts map[int]int, n int) bool {
	for _, count := range counts {
		if count == n {
			return true
		}
	}
	return false
}

func countPairs(counts map[int]int) int {
	pairs := 0
	for _, count := range counts {
		if count == 2 {
			pairs++
		}
	}
	return pairs
}
