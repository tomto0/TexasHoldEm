package main

import (
	"fmt"
	"sort"
	"testing"
)

// Definiert die Testfälle
func TestPokerHands(t *testing.T) {
	// Testfälle definieren
	testCases := []struct {
		communityCards string
		player1        string
		player2        string
		expectedResult int
		expectedType1  string
		expectedType2  string
	}{

		//Testfälle für High Card
		{"D6 S9 H4 S3 C2", "SK CA", "HA SQ", 1, "High Card", "High Card"},  // Spieler 1: Ass-König-Kicker schlägt Ass-Dame
		{"D6 S9 H4 S3 C2", "SK CA", "HA CK", 0, "High Card", "High Card"},  // Unentschieden: Beide haben Ass-König
		{"D6 S9 H4 H3 H2", "C8 DJ", "C7 DQ", -1, "High Card", "High Card"}, // Spieler 2: Dame schlägt Bube

		// Testfälle für One Pair
		{"SK HT C8 C7 D2", "DK C2", "H8 D5", 1, "Two Pairs", "High Card"},  // Spieler 1: Paar Könige schlägt High Card
		{"SK HT C8 C7 D2", "DK C5", "HK D5", 0, "Two Pairs", "One Pair"},   // Spieler 1 gewinnt: Zwei Paare (Könige und Zweien) gegen ein Paar
		{"HA DA ST C9 D4", "D5 C6", "H7 C2", -1, "High Card", "High Card"}, // Spieler 2: Höhere High Card gewinnt

		// Testfälle für Two Pairs
		{"SA DQ CK D6 H6", "HA C3", "CQ H4", 1, "Two Pairs", "Two Pairs"},  // Spieler 1: Zwei Paare mit höherem Kicker schlagen Spieler 2
		{"SA DQ CK D6 H6", "HQ C3", "SQ H3", 0, "Two Pairs", "Two Pairs"},  // Unentschieden: Gleiche Paare und Kicker
		{"SA DQ CK D6 H5", "HQ C6", "CA HK", -1, "Two Pairs", "Two Pairs"}, // Spieler 2: Höhere Zwei Paare (Ass)

		// Testfälle für Three of a Kind
		{"SA D3 H3 C8 SJ", "HJ SJ", "C3 H2", 1, "Three of a Kind", "High Card"},        // Spieler 1: Drilling Buben schlägt High Card
		{"SA D3 H3 C8 SJ", "C3 S2", "S3 H2", 0, "Three of a Kind", "Three of a Kind"},  // Unentschieden: Beide haben Drilling
		{"HA SA DA H3 HT", "S2 S5", "H2 SK", -1, "Three of a Kind", "Three of a Kind"}, // Spieler 2: Drilling Asse schlägt Drilling Könige

		// Testfälle für Straight
		{"H3 S4 C5 S6 HT", "D7 HA", "H2 SA", 1, "Straight", "Straight"},  // Spieler 1: Straße bis 6 schlägt Straße bis 5
		{"H3 S4 C5 S6 HT", "D7 HA", "H7 SA", 0, "Straight", "Straight"},  // Unentschieden: Beide haben Straße bis 6
		{"H2 H3 S4 C5 HT", "HA S3", "H6 SA", -1, "Straight", "Straight"}, // Spieler 2: Höhere Straße (6) schlägt niedrigere (5)

		// Testfälle für Flush
		{"D3 D6 DT C5 HQ", "DK DA", "D2 DQ", 1, "Flush", "Flush"},           // Spieler 1: Höherer Flush gewinnt
		{"D3 D6 DT DJ DK", "-", "-", 0, "Four of a Kind", "Four of a Kind"}, // Unentschieden: Beide haben Vierling mit gleichem Rang
		{"D3 D6 DT C5 HQ", "D2 D5", "DJ DA", -1, "Flush", "Flush"},          // Spieler 2: Höherer Flush gewinnt

		// Testfälle für Full House
		{"HQ SQ HT DT C3", "DQ C2", "CT C4", 1, "Full House", "Full House"},  // Spieler 1 gewinnt: Tie-Break mit höherem Drilling
		{"SA HQ SQ HT D8", "HA DJ", "DA CQ", -1, "Full House", "Full House"}, // Spieler 2 gewinnt mit höherem Full House
		{"HQ SQ HT DT C3", "ST C2", "CQ C4", -1, "Full House", "Full House"}, // Spieler 2 gewinnt mit höherem Drilling

		// Testfälle für Four of a Kind
		{"HT ST CT DT HK", "HA S7", "DJ C5", 1, "Four of a Kind", "High Card"},       // Spieler 1: Vierling schlägt High Card
		{"S5 D5 C5 H5 HA", "-", "-", 0, "Four of a Kind", "Four of a Kind"},          // Unentschieden: Beide haben gleichen Vierling
		{"HT ST CT DT S8", "C2 C3", "C5 HK", -1, "Four of a Kind", "Four of a Kind"}, // Spieler 2 gewinnt: Vierling mit höherem Kicker

		// Testfälle für Straight Flush
		{"H3 H4 H5 H6 HT", "H7 HA", "H2 SA", 1, "Straight Flush", "Flush"},           // Straight Flush schlägt Flush
		{"H3 H4 H5 H6 H7", "-", "-", 0, "Straight Flush", "Straight Flush"},          // Unentschieden: Beide haben gleichen Straight Flush
		{"S7 S8 S9 ST DK", "S6 C2", "SJ D5", -1, "Straight Flush", "Straight Flush"}, // Spieler 2 gewinnt mit höherem Straight Flush

		// Testfälle für Royal Flush
		{"DT DJ DQ DK DA", "-", "-", 0, "Royal Flush", "Royal Flush"}, // Unentschieden: Beide haben Royal Flush
	}

	// Testet, ob das Programm für die Testfälle richtige Ergebnisse liefert
	for _, tc := range testCases {
		// Parse der Gemeinschaftskarten (Community Cards)
		community := parseHand(tc.communityCards)

		// Parse der Hände von Spieler 1 und Spieler 2
		var hand1, hand2 Hand
		if tc.player1 == "-" {
			hand1 = community // Keine Hand für Spieler 1, verwende nur die Gemeinschaftskarten
		} else {
			hand1 = parseHand(tc.player1 + " " + tc.communityCards) // Kombiniere Spieler 1 Karten mit den Gemeinschaftskarten
		}
		if tc.player2 == "-" {
			hand2 = community // Keine Hand für Spieler 2, verwende nur die Gemeinschaftskarten
		} else {
			hand2 = parseHand(tc.player2 + " " + tc.communityCards) // Kombiniere Spieler 2 Karten mit den Gemeinschaftskarten
		}

		// Vergleich der Hände unter Einbeziehung der Gemeinschaftskarten
		result := hand1.CompareTo(hand2, community)
		hand1Type, hand1Score := hand1.evaluateHand() // Handtyp und Punktzahl für Spieler 1
		hand2Type, hand2Score := hand2.evaluateHand() // Handtyp und Punktzahl für Spieler 2

		// Validierung des Ergebnisses
		if result != tc.expectedResult {
			t.Errorf("FEHLGESCHLAGEN: Gemeinschaftskarten: %s | Spieler 1: %s (%s, %d) | Spieler 2: %s (%s, %d) | Erwartet: %d, Erhalten: %d",
				tc.communityCards, tc.player1, hand1Type, hand1Score, tc.player2, hand2Type, hand2Score, tc.expectedResult, result)
		}
	}

	// Erstellt die Rangliste
	for _, tc := range testCases {
		// Parse der Gemeinschaftskarten (Community Cards)
		community := parseHand(tc.communityCards)

		// Parse der Hände von Spieler 1 und Spieler 2
		var hand1, hand2 Hand
		if tc.player1 == "-" {
			hand1 = community // Keine Hand für Spieler 1, verwende nur die Gemeinschaftskarten
		} else {
			hand1 = parseHand(tc.player1 + " " + tc.communityCards) // Kombiniere Spieler 1 Karten mit den Gemeinschaftskarten
		}
		if tc.player2 == "-" {
			hand2 = community // Keine Hand für Spieler 2, verwende nur die Gemeinschaftskarten
		} else {
			hand2 = parseHand(tc.player2 + " " + tc.communityCards) // Kombiniere Spieler 2 Karten mit den Gemeinschaftskarten
		}

		// Bewertung der Hände
		hand1Type, hand1Score := hand1.evaluateHand() // Extrahiere Typ und Punktzahl für Hand 1
		hand2Type, hand2Score := hand2.evaluateHand() // Extrahiere Typ und Punktzahl für Hand 2

		// Erstelle eine Liste der Hände mit Typ und Punktzahl
		hands := HandList{
			Hand{Cards: hand1.Cards, Score: hand1Score, HandType: hand1Type},
			Hand{Cards: hand2.Cards, Score: hand2Score, HandType: hand2Type},
		}

		// Sortiere die Hände innerhalb des Spiels
		sort.Sort(hands)

		// Ausgabe der Rangliste für das aktuelle Spiel
		fmt.Printf("Rangliste für Spiel (Community: %s):\n", tc.communityCards)
		rank := 1 // Initial rank
		for i := 0; i < len(hands); i++ {
			if i > 0 && hands[i].Score == hands[i-1].Score && compareKickers(hands[i], hands[i-1], community) == 0 {
				// If tied, keep the same rank
				fmt.Printf("Platz %d: Hand: [%s], Typ: %s, Punktzahl: %d\n",
					rank, hands[i].String(), hands[i].HandType, hands[i].Score)
			} else {
				// Assign new rank
				rank = i + 1
				fmt.Printf("Platz %d: Hand: [%s], Typ: %s, Punktzahl: %d\n",
					rank, hands[i].String(), hands[i].HandType, hands[i].Score)
			}
		}

		// Überprüfung der Sortierung für das Spiel
		for i := 0; i < len(hands)-1; i++ {
			if hands[i].Score < hands[i+1].Score {
				t.Errorf("FEHLGESCHLAGEN: Hand %v mit Score %d kommt vor Hand %v mit Score %d",
					hands[i].Cards, hands[i].Score, hands[i+1].Cards, hands[i+1].Score)
			}
		}
	}

}
