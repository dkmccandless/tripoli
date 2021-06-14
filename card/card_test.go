package card

import "testing"

func TestCard(t *testing.T) {
	for _, test := range []struct {
		card  Card
		color Color
		suit  Suit
		rank  Rank
	}{
		{0, Black, Clubs, Two},
		{1, Black, Clubs, Three},
		{2, Black, Clubs, Four},
		{3, Black, Clubs, Five},
		{4, Black, Clubs, Six},
		{5, Black, Clubs, Seven},
		{6, Black, Clubs, Eight},
		{7, Black, Clubs, Nine},
		{8, Black, Clubs, Ten},
		{9, Black, Clubs, Jack},
		{10, Black, Clubs, Queen},
		{11, Black, Clubs, King},
		{12, Black, Clubs, Ace},
		{13, Red, Diamonds, Two},
		{14, Red, Diamonds, Three},
		{15, Red, Diamonds, Four},
		{16, Red, Diamonds, Five},
		{17, Red, Diamonds, Six},
		{18, Red, Diamonds, Seven},
		{19, Red, Diamonds, Eight},
		{20, Red, Diamonds, Nine},
		{21, Red, Diamonds, Ten},
		{22, Red, Diamonds, Jack},
		{23, Red, Diamonds, Queen},
		{24, Red, Diamonds, King},
		{25, Red, Diamonds, Ace},
		{26, Black, Spades, Two},
		{27, Black, Spades, Three},
		{28, Black, Spades, Four},
		{29, Black, Spades, Five},
		{30, Black, Spades, Six},
		{31, Black, Spades, Seven},
		{32, Black, Spades, Eight},
		{33, Black, Spades, Nine},
		{34, Black, Spades, Ten},
		{35, Black, Spades, Jack},
		{36, Black, Spades, Queen},
		{37, Black, Spades, King},
		{38, Black, Spades, Ace},
		{39, Red, Hearts, Two},
		{40, Red, Hearts, Three},
		{41, Red, Hearts, Four},
		{42, Red, Hearts, Five},
		{43, Red, Hearts, Six},
		{44, Red, Hearts, Seven},
		{45, Red, Hearts, Eight},
		{46, Red, Hearts, Nine},
		{47, Red, Hearts, Ten},
		{48, Red, Hearts, Jack},
		{49, Red, Hearts, Queen},
		{50, Red, Hearts, King},
		{51, Red, Hearts, Ace},
	} {
		if color := test.card.Color(); color != test.color {
			t.Errorf("Card(%v).Color: got %v, expected %v",
				test.card, color, test.color,
			)
		}
		if suit := test.card.Suit(); suit != test.suit {
			t.Errorf("Card(%v).Suit: got %v, expected %v",
				test.card, suit, test.suit,
			)
		}
		if rank := test.card.Rank(); rank != test.rank {
			t.Errorf("Card(%v).Rank: got %v, expected %v",
				test.card, rank, test.rank,
			)
		}
		if card := test.rank.Suit(test.suit); card != test.card {
			t.Errorf("%v.Suit(%v): got %v, expected %v",
				test.rank, test.suit, card, test.card,
			)
		}
		if card := test.suit.Rank(test.rank); card != test.card {
			t.Errorf("%v.Rank(%v): got %v, expected %v",
				test.suit, test.rank, card, test.card,
			)
		}
	}
}

func TestColor(t *testing.T) {
	for _, test := range []struct {
		color, opp   Color
		minor, major Suit
	}{
		{Black, Red, Clubs, Spades},
		{Red, Black, Diamonds, Hearts},
	} {
		if opp := test.color.Opp(); opp != test.opp {
			t.Errorf("%v.Opp: got %v, expected %v",
				test.color, opp, test.opp,
			)
		}
		if minor := test.color.Minor(); minor != test.minor {
			t.Errorf("%v.Minor: got %v, expected %v",
				test.color, minor, test.minor,
			)
		}
		if major := test.color.Major(); major != test.major {
			t.Errorf("%v.Major: got %v, expected %v",
				test.color, major, test.major,
			)
		}
	}
}
