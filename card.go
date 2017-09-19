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

func (s Suit) Two() Card   { return Card(ranks*s) + Card(Two) }
func (s Suit) Three() Card { return Card(ranks*s) + Card(Three) }
func (s Suit) Four() Card  { return Card(ranks*s) + Card(Four) }
func (s Suit) Five() Card  { return Card(ranks*s) + Card(Five) }
func (s Suit) Six() Card   { return Card(ranks*s) + Card(Six) }
func (s Suit) Seven() Card { return Card(ranks*s) + Card(Seven) }
func (s Suit) Eight() Card { return Card(ranks*s) + Card(Eight) }
func (s Suit) Nine() Card  { return Card(ranks*s) + Card(Nine) }
func (s Suit) Ten() Card   { return Card(ranks*s) + Card(Ten) }
func (s Suit) Jack() Card  { return Card(ranks*s) + Card(Jack) }
func (s Suit) Queen() Card { return Card(ranks*s) + Card(Queen) }
func (s Suit) King() Card  { return Card(ranks*s) + Card(King) }
func (s Suit) Ace() Card   { return Card(ranks*s) + Card(Ace) }

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

func (r Rank) Clubs() Card    { return Card(ranks*Clubs) + Card(r) }
func (r Rank) Diamonds() Card { return Card(ranks*Diamonds) + Card(r) }
func (r Rank) Spades() Card   { return Card(ranks*Spades) + Card(r) }
func (r Rank) Hearts() Card   { return Card(ranks*Hearts) + Card(r) }

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
