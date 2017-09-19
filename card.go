package tripoli

const ranks = 13

// A Color represents the color of a Suit.
type Color int

const (
	Black Color = iota
	Red
)

// Minor returns a Color's minor Suit (Clubs or Diamonds).
func (c Color) Minor() Suit { return Suit(c) }

// Major returns a Color's major Suit (Spades or Hearts).
func (c Color) Major() Suit { return Suit(c) + 2 }

// Other returns the opposite of the given Color.
func (c Color) Other() Color { return 1 &^ c }

// A Suit represents a standard playing card suit in the order specified in the constant declaration.
// Clubs and Diamonds are the "minor suits" and Spades and Hearts the "major suits".
// Behavior is undefined for values outside the range [0, 4).
type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Spades
	Hearts
)

// Card returns the Card of the given Rank and Suit.
func (s Suit) Card(r Rank) Card { return Card(ranks*s) + Card(r) }

// Color returns the color of a Suit.
func (s Suit) Color() Color { return Color(s & 1) }

// A Rank represents a standard playing card rank. Aces are high.
// Behavior is undefined for values outside the range [0, 13).
type Rank int

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

// Card returns the Card of the given Rank and Suit.
func (r Rank) Card(s Suit) Card { return Card(ranks*s) + Card(r) }

// A Card represents a card in a standard deck in suit-rank order
// (two of clubs = 0, three of clubs = 1, ace of hearts = 51).
// Behavior is undefined for values outside the range [0, 52).
type Card int

// Color returns the Color of a Card.
func (c Card) Color() Color { return c.Suit().Color() }

// Suit returns the Suit of a Card.
func (c Card) Suit() Suit { return Suit(c / ranks) }

// Rank returns the Rank of a Card.
func (c Card) Rank() Rank { return Rank(c % ranks) }
