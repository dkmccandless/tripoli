package game

import "github.com/dkmccandless/tripoli/card"

// A Player can participate in a game of Tripoli.
type Player interface {
	// Init informs the Player of the number of Players in the game,
	// the Player's position in the deal (counting from the dealer's left),
	// the Player's cards, and the values of the counter stakes and Kitty.
	// It is called once per hand, after the ante and before play begins.
	Init(n, pos int, hand []card.Card, stake map[card.Card]int, kitty int)

	// Note informs the Player whenever any Player plays a card.
	Note(pos int, card card.Card)

	// PlayMajor reports whether the Player decides to play from a color's
	// "major" suit (spades or hearts) or "minor" suit (clubs or diamonds).
	// It is only called when the Player must decide which suit to play.
	PlayMajor(color card.Color) bool
}
