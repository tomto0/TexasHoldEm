package main

import "fmt"

func main() {
	hands := []string{
		// royal flush
		"CT CJ CQ CK CA",
		// straight flush
		"D8 DQ DJ DT D9",
		"H8 HT HJ H7 H9",
		// poker
		"HT SQ ST DT CT",
		"HT SK ST DT CT",
		"H8 SQ S8 D8 C8",
		"H7 SK S7 D7 C7",
		// full house
		"H2 SQ C2 D2 CQ",
		"H2 SJ C2 D2 CJ",
		// flush
		"HK HQ H2 H4 H5",
		"D5 D4 D2 DQ DK",
		// straight
		"H3 S7 H5 D6 H4",
		"C9 CT SJ D7 H8",
		"H4 S5 HA D3 H2",
		// three of a kind
		"H2 SQ S2 D2 CK",
		"H2 S7 S2 D2 C9",
		"H2 S8 S2 D2 C9",
		// two pairs
		"H5 SQ C5 DT CT",
		"H9 SQ C9 DT CT",
		// one pair
		"H3 S8 H5 D8 CA",
		"S4 DA H3 CA HT",
		// high card
		"H3 S8 H5 DK CA",
		"H3 S8 H5 DK CT",
		"H3 S8 H5 DK C2",
	}

	for _, handStr := range hands {
		hand := parseHand(handStr)
		result := hand.evaluateHand()
		fmt.Printf("Hand: %s -> %s\n", handStr, result)
	}
}
