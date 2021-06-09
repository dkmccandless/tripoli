// Package card defines a standard deck of 52 playing cards.
package card

// A Card is a standard playing card. Cards are ordered by suit first, then rank
// (e.g. two of clubs = 0, three of clubs = 1, ace of hearts = 51).
type Card int

// Color returns a Card's Color.
func (c Card) Color() Color { return c.Suit().Color() }

// Suit returns a Card's Suit.
func (c Card) Suit() Suit { return Suit(c / 13) }

// Rank returns a Card's Rank.
func (c Card) Rank() Rank { return Rank(c % 13) }

// A Rank is a playing card rank. Aces are high.
type Rank int

//go:generate stringer -type=Rank
const (
	Two Rank = iota
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

// Suit returns the Card of the given Suit.
func (r Rank) Suit(s Suit) Card { return Card(13*s) + Card(r) }

// A Suit is a standard playing card suit.
// Clubs and Diamonds are the "minor suits"
// and Spades and Hearts the "major suits".
type Suit int

//go:generate stringer -type=Suit
const (
	Clubs Suit = iota
	Diamonds
	Spades
	Hearts
)

// Rank returns the Card of the given Rank.
func (s Suit) Rank(r Rank) Card { return Card(13*s) + Card(r) }

// Color returns the color of a Suit.
func (s Suit) Color() Color { return Color(s % 2) }

// A Color is a suit's color.
type Color int

//go:generate stringer -type=Color
const (
	Black Color = iota
	Red
)

// Minor returns a Color's minor Suit (Clubs or Diamonds).
func (c Color) Minor() Suit { return Suit(c) }

// Major returns a Color's major Suit (Spades or Hearts).
func (c Color) Major() Suit { return Suit(c) + 2 }

// Opp returns the opposite Color.
func (c Color) Opp() Color { return 1 - c }
