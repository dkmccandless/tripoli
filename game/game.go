// Package game implements a simplified version of the Michigan portion of a
// game of Tripoli.
package game

import (
	"math/rand"

	"github.com/dkmccandless/tripoli/card"
)

// A Game administers a game of Tripoli.
type Game struct {
	score map[Player]int
	stake map[card.Card]int
	kitty int
}

// New initializes a new Game.
func New(players []Player) *Game {
	score := make(map[Player]int)
	for _, p := range players {
		score[p] = 0
	}
	return &Game{
		score: score,
		stake: map[card.Card]int{
			card.Hearts.Rank(card.Ten):   0,
			card.Hearts.Rank(card.Jack):  0,
			card.Hearts.Rank(card.Queen): 0,
			card.Hearts.Rank(card.King):  0,
			card.Hearts.Rank(card.Ace):   0,
		},
	}
}

// Score returns the players' scores.
func (g *Game) Score() map[Player]int {
	score := make(map[Player]int)
	for p, n := range g.score {
		score[p] = n
	}
	return score
}

// Stake returns the values of the stake pots.
func (g *Game) Stake() map[card.Card]int {
	stake := make(map[card.Card]int)
	for c, n := range g.stake {
		stake[c] = n
	}
	return stake
}

// Kitty returns the value of the Kitty.
func (g *Game) Kitty() int { return g.kitty }

// a round administers a round of Tripoli.
type round struct {
	g *Game

	// p records each player's position in the deal.
	p []Player

	// n records the number of cards in each player's hand.
	n []int

	// deck records the location of each card.
	// A value of -1 indicates the extra hand/discard.
	deck []int
}

// init initializes a round with a shuffled deck, seats the players in a random
// order, and antes for each player. init calls each Player's Init method.
func (g *Game) init() *round {
	r := &round{g: g}
	n := len(g.score)

	for p := range g.score {
		r.p = append(r.p, p)
		r.ante(p)
	}
	rand.Shuffle(n, func(i, j int) { r.p[i], r.p[j] = r.p[j], r.p[i] })

	r.deck = rand.Perm(52)
	hand := make([][]card.Card, n)
	for i, v := range r.deck {
		if v %= n + 1; v <= n {
			r.deck[i] = v
			hand[v] = append(hand[v], card.Card(i))
		} else {
			r.deck[i] = -1
		}
	}

	for i, p := range r.p {
		r.n[i] = len(hand[i])
		stake := make(map[card.Card]int)
		for c, v := range g.stake {
			stake[c] = v
		}
		p.Init(n, i, hand[i], stake, g.kitty)
	}

	return r
}

// Play plays a round of Tripoli.
// The players' positions in the deal are randomly assigned.
func (g *Game) Play() {
	r := g.init()
	var pos int
	var won bool
	for lead := r.firstLead(); ; {
		if pos, won = r.playRun(lead); won {
			break
		}
		var ok bool
		if lead, ok = r.nextLead(pos, lead.Color().Opp()); !ok {
			break
		}
	}
	for i := range r.p {
		r.payKitty(r.p[i], r.n[i])
	}
	if won {
		r.collectKitty(r.p[pos])
	}
}

// firstLead returns the lowest club held by any player.
func (r *round) firstLead() card.Card {
	c := card.Clubs.Rank(card.Two)
	for r.deck[c] == -1 {
		c++
	}
	return c
}

// playRun plays a sequence of consecutive cards until a player is out of cards
// or no player holds the next card. It returns the position of the player who
// played the last card and a boolean value reporting whether they have won the
// round (by playing the last card in their hand).
func (r *round) playRun(lead card.Card) (pos int, won bool) {
	for c := lead; c.Suit() == lead.Suit() && r.deck[c] != -1; c++ {
		pos = r.deck[c]
		r.playCard(c)
		if r.n[pos] == 0 {
			return pos, true
		}
	}
	return pos, false
}

// playCard plays a card and calls each Player's Note method.
func (r *round) playCard(c card.Card) {
	pos := r.deck[c]
	r.deck[c] = -1
	r.n[pos]--
	r.collect(r.p[pos], c)
	for _, p := range r.p {
		p.Note(pos, c)
	}
}

// nextLead returns the lead card that begins the next run and a boolean value
// reporting whether the game continues.
// Control of the lead begins at the indicated position and passes as necessary
// to the first player in order with a card in the correct color. Passed players
// pay one point to the Kitty.
// If the Player to lead holds cards in both suits, nextLead calls that Player's
// PlayMajor method.
func (r *round) nextLead(pos int, color card.Color) (lead card.Card, ok bool) {
	// next returns the next position in order.
	next := func(n int) int { return (n + 1) % len(r.p) }
	for old := pos; ; pos = next(pos) {
		minor, hasMinor := r.lowest(pos, color.Minor())
		major, hasMajor := r.lowest(pos, color.Major())
		switch p := r.p[pos]; {
		case hasMinor && hasMajor:
			if p.PlayMajor(color) {
				return major, true
			} else {
				return minor, true
			}
		case hasMajor:
			return major, true
		case hasMinor:
			return minor, true
		default:
			r.payKitty(p, 1)
			if next(pos) == old {
				return 0, false
			}
		}
	}
}

// lowest returns the lowest card held by a player in a suit,
// and a boolean value reporting whether the player holds any cards in the suit.
func (r *round) lowest(pos int, s card.Suit) (card.Card, bool) {
	for c := s.Rank(card.Two); c.Suit() == s; c++ {
		if r.deck[c] == pos {
			return c, true
		}
	}
	return 0, false
}

// ante transfers one point from a player's score to each stake pot.
func (r *round) ante(p Player) {
	for c := range r.g.stake {
		r.g.score[p]--
		r.g.stake[c]++
	}
}

// collect transfers a card's stake to a player's score.
func (r *round) collect(p Player, c card.Card) {
	r.g.score[p] += r.g.stake[c]
	r.g.stake[c] = 0
}

// payKitty transfers an amount from a player's score to the kitty.
func (r *round) payKitty(p Player, n int) {
	r.g.score[p] -= n
	r.g.kitty += n
}

// collectKitty transfers the kitty to a player's score.
func (r *round) collectKitty(p Player) {
	r.g.score[p] += r.g.kitty
	r.g.kitty = 0
}
