package tripoli

// A Player can participate in a game of Tripoli.
type Player interface {
	// Init informs the Player of the number of Players in the game,
	// the Player's position in the deal (counting from the dealer's left),
	// the Cards the Player is dealt, and the values of the counter stakes
	// (in the order Ten, Jack, Queen, King, Ace of Hearts) and Kitty.
	// It is called once per hand, after the ante and before play begins.
	Init(nplayers int, seat int, hand []Card, stake []int, kitty int)

	// Note informs the Player whenever any Player plays a card.
	Note(seat int, card Card)

	// PlayMajor reports whether the Player decides to play from the specified color's "major" suit
	// (spades or hearts) or the "minor" suit (clubs or diamonds).
	// It is called when the Player must decide which suit to play,
	// and only if the Player holds at least one card in both suits.
	PlayMajor(color Color) bool
}
