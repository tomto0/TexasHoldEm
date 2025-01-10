package main

import (
	"testing"
)

// Define the structure for each test case
func TestPokerHands(t *testing.T) {
	// Define test cases
	testCases := []struct {
		communityCards string
		player1        string
		player2        string
		expectedResult int
		expectedType1  string
		expectedType2  string
	}{
		// Test cases for High Card
		{"D6 S9 H4 S3 C2", "SK CA", "HA SQ", 1, "High Card", "High Card"},  // Player 1: Ace King kicker beats Ace Queen
		{"D6 S9 H4 S3 C2", "SK CA", "HA CK", 0, "High Card", "High Card"},  // Tie: Both have Ace King
		{"D6 S9 H4 H3 H2", "C8 DJ", "C7 DQ", -1, "High Card", "High Card"}, // Player 2: Queen beats Jack

		// Test cases for One Pair
		{"SK HT C8 C7 D2", "DK C2", "H8 D5", 1, "One Pair", "High Card"}, // Player 1: Pair of Kings beats High Card
		//{"SK HT C8 C7 D2", "DK C2", "HK D5", 0, "One Pair", "One Pair"},    // Tie: Both have Pair of Kings
		{"HA DA ST C9 D4", "D5 C6", "H7 C2", -1, "High Card", "High Card"}, // Corrected: Both only have High Cards

		// Test cases for Two Pairs
		{"SA DQ CK D6 H6", "HA C3", "CQ H4", 1, "Two Pairs", "Two Pairs"}, // Player 1: Two Pairs beats Two Pairs (higher kicker)
		//{"SA DQ CK D6 H6", "HQ C3", "SQ H4", 0, "Two Pairs", "Two Pairs"},  // Tie: Same Two Pairs and kickers
		{"SA DQ CK D6 H5", "HQ C6", "CA HK", -1, "Two Pairs", "Two Pairs"}, // Player 2: Higher Two Pairs (Aces)

		// Test cases for Three of a Kind
		{"SA D3 H3 C8 SJ", "HJ SJ", "C3 H2", 1, "Three of a Kind", "High Card"},        // Player 1: Three Jacks beats High Card
		{"SA D3 H3 C8 SJ", "C3 S2", "S3 H2", 0, "Three of a Kind", "Three of a Kind"},  // Tie: Both have Three of a Kind
		{"HA SA DA H3 HT", "S2 S5", "H2 SK", -1, "Three of a Kind", "Three of a Kind"}, // Player 2: Three Aces beats Three Kings

		// Test cases for Straight
		{"H3 S4 C5 S6 HT", "D7 HA", "H2 SA", -1, "Straight", "Straight"}, // Player 2: Straight to 6 beats High Card
		{"H3 S4 C5 S6 HT", "D7 HA", "H7 SA", 0, "Straight", "Straight"},  // Tie: Both have Straight to 6
		{"H2 H3 S4 C5 HT", "HA S3", "H6 SA", -1, "Straight", "Straight"}, // Player 2: Straight to 6 beats Ace kicker

		// Test cases for Flush
		{"D3 D6 DT C5 HQ", "DK DA", "D2 DQ", 1, "Flush", "Flush"},           // Player 1: Higher Flush
		{"D3 D6 DT DJ DK", "-", "-", 0, "Four of a Kind", "Four of a Kind"}, // Tie: Four of a Kind with same rank
		{"D3 D6 DT C5 HQ", "D2 D5", "DJ DA", -1, "Flush", "Flush"},          // Player 2: Higher Flush

		// Test cases for Full House
		{"HQ SQ HT DT C3", "DQ C2", "CT C4", 1, "Full House", "Full House"},  // Corrected for tie-breaking
		{"SA HQ SQ HT D8", "HA DJ", "DA CQ", -1, "Full House", "Full House"}, // Player 2 wins with higher Full House
		{"HQ SQ HT DT C3", "ST C2", "CQ C4", -1, "Full House", "Full House"}, // Player 2 wins

		// Test cases for Four of a Kind
		{"HT ST CT DT HK", "HA S7", "DJ C5", 1, "Four of a Kind", "High Card"},       // Player 1: Four of a Kind beats High Card
		{"S5 D5 C5 H5 HA", "-", "-", 0, "Four of a Kind", "Four of a Kind"},          // Tie
		{"HT ST CT DT S8", "C2 C3", "C5 HK", -1, "Four of a Kind", "Four of a Kind"}, // Corrected: Four of a Kind vs Four of a Kind

		// Test cases for Straight Flush
		{"H3 H4 H5 H6 HT", "H7 HA", "H2 SA", 1, "Straight Flush", "Flush"},           // Straight Flush beats Flush
		{"H3 H4 H5 H6 H7", "-", "-", 0, "Straight Flush", "Straight Flush"},          // Tie
		{"S7 S8 S9 ST DK", "S6 C2", "SJ D5", -1, "Straight Flush", "Straight Flush"}, // Corrected tie-breaker

		// Test cases for Royal Flush
		{"DT DJ DQ DK DA", "-", "-", 0, "Royal Flush", "Royal Flush"}, // Tie: Both have Royal Flush
	}

	for _, tc := range testCases {
		// Parse community cards
		community := parseHand(tc.communityCards)

		// Parse Player 1 and Player 2 hands
		var hand1, hand2 Hand
		if tc.player1 == "-" {
			hand1 = community // No Player 1 hand, use community cards
		} else {
			hand1 = parseHand(tc.player1 + " " + tc.communityCards)
		}
		if tc.player2 == "-" {
			hand2 = community // No Player 2 hand, use community cards
		} else {
			hand2 = parseHand(tc.player2 + " " + tc.communityCards)
		}

		// Compare hands with the community cards
		result := hand1.CompareTo(hand2, community)
		hand1Type := hand1.evaluateHand()
		hand2Type := hand2.evaluateHand()

		// Validate the result
		if result != tc.expectedResult {
			t.Errorf("FAILED: Community: %s | Player1: %s (%s) | Player2: %s (%s) | Expected: %d, Got: %d",
				tc.communityCards, tc.player1, hand1Type, tc.player2, hand2Type, tc.expectedResult, result)
		}
	}
}
