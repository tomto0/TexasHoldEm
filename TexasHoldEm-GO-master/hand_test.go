package main

import (
	"testing"
)

// Testfälle für die Bewertung verschiedener Pokerhände.
func TestEvaluateHands(t *testing.T) {
	tests := []struct {
		handStr  string
		expected string
	}{
		// Royal Flush
		{"CT CJ CQ CK CA", "Royal Flush"},
		// Straight Flush
		{"D8 DQ DJ DT D9", "Straight Flush"},
		{"H8 HT HJ H7 H9", "Straight Flush"},
		// Poker
		{"HT SQ ST DT CT", "Poker"},
		{"HT SK ST DT CT", "Poker"},
		{"H8 SQ S8 D8 C8", "Poker"},
		{"H7 SK S7 D7 C7", "Poker"},
		// Full House
		{"H2 SQ C2 D2 CQ", "Full House"},
		{"H2 SJ C2 D2 CJ", "Full House"},
		// Flush
		{"HK HQ H2 H4 H5", "Flush"},
		{"D5 D4 D2 DQ DK", "Flush"},
		// Straight
		{"H3 S7 H5 D6 H4", "Straight"},
		{"C9 CT SJ D7 H8", "Straight"},
		{"H4 S5 HA D3 H2", "Straight"},
		// Three of a Kind
		{"H2 SQ S2 D2 CK", "Three of a Kind"},
		{"H2 S7 S2 D2 C9", "Three of a Kind"},
		{"H2 S8 S2 D2 C9", "Three of a Kind"},
		// Two Pairs
		{"H5 SQ C5 DT CT", "Two Pairs"},
		{"H9 SQ C9 DT CT", "Two Pairs"},
		// One Pair
		{"H3 S8 H5 D8 CA", "One Pair"},
		{"S4 DA H3 CA HT", "One Pair"},
		// High Card
		{"H3 S8 H5 DK CA", "High Card"},
		{"H3 S8 H5 DK CT", "High Card"},
		{"H3 S8 H5 DK C2", "High Card"},
	}

	for _, test := range tests {
		hand := parseHand(test.handStr)
		result := Hand.evaluateHand(hand)
		if result != test.expected {
			t.Errorf("Hand: %s, Erwartet: %s, Erhalten: %s", test.handStr, test.expected, result)
		}
	}
}

// Benchmark für die Bewertung der oben genannten Hände.
func BenchmarkEvaluateHands(b *testing.B) {
	hands := []string{
		"CT CJ CQ CK CA",
		"D8 DQ DJ DT D9",
		"H8 HT HJ H7 H9",
		"HT SQ ST DT CT",
		"H3 S8 H5 DK CA",
	}
	for i := 0; i < b.N; i++ {
		for _, handStr := range hands {
			hand := parseHand(handStr)
			Hand.evaluateHand(hand)
		}
	}
}
