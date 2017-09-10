package tripoli

type Player interface {
	// Init informs the Player of the number of Players in the game, the Player's position in the deal (counting from the dealer's left),
	// the cards the Player is dealt, and the values of the counter stakes (in the order Ten, Jack, Queen, King, Ace of Hearts) and kitty.
	// It is called once per hand, after the ante and before play begins.
	Init(nplayers int, seat int, hand []Card, stake []int, kitty int)

	// Note informs the Player when a player plays a card.
	Note(seat int, card Card)

	// PlayMajor reports whether the Player decides to play from the "major" suit (spades or hearts) of the specified color or the "minor" suit (clubs or diamonds).
	// PlayMajor will only be called if the Player holds at least one card in both suits.
	PlayMajor(color Color) bool
}
