package tripoli

import (
	"math/rand"
	"time"
)

// A Scoreboard maps Players to their scores.
type Scoreboard map[Player]int

type Game struct {
	players []Player
	score   Scoreboard
	stake   []int
	kitty   int
}

func NewGame(players []Player) *Game {
	return &Game{
		players: players,
		score:   make(Scoreboard),
		stake:   make([]int, 52),
	}
}

func (g *Game) Play(nhands int) {
	rand.Seed(time.Now().UnixNano())
	for hand := 0; hand < nhands; hand++ {
		g.playHand()
		// rotate the deal
		g.players = append(g.players[:len(g.players)-1], g.players[0])
	}
}

func (g *Game) playHand() {
	var (
		nplayers = len(g.players)
		holder   = make([]int, 52) // holder[c] = the seat of the player who holds Card c, or nplayers for the extra hand/discard
		ncards   = make([]int, nplayers+1)
		hands    = make([][]Card, nplayers+1)
	)

	// Ante
	for _, p := range g.players {
		for c := Hearts.Ten(); c <= Hearts.Ace(); c++ {
			g.score[p]--
			g.stake[c]++
		}
		g.score[p]--
		g.kitty++
	}

	// Shuffle and deal
	holder = rand.Perm(52)
	for c := range holder {
		holder[c] %= nplayers + 1
		hands[holder[c]] = append(hands[holder[c]], Card(c))
	}
	for i, p := range g.players {
		ncards[i] = len(hands[i]) // == (51-i)/(nplayers+1) + 1
		p.Init(
			nplayers,
			i,
			// Do not give Players access to the underlying arrays
			append([]Card{}, hands[i]...),
			append([]int{}, g.stake[Hearts.Ten():]...),
			g.kitty,
		)
	}

	lead := Clubs.Two()
	for holder[lead] == nplayers {
		lead++
	}

	for controller := holder[lead]; ; { // the seat of the Player who last played a card
		for c := lead; c <= lead.Suit().Ace() && holder[c] != nplayers && ncards[controller] > 0; c++ {
			// Play the card
			controller, holder[c] = holder[c], nplayers
			ncards[controller]--

			// Collect points for playing a counter
			g.score[g.players[controller]] += g.stake[c]
			g.stake[c] = 0

			for _, p := range g.players {
				p.Note(controller, c)
			}
		}

		if ncards[controller] == 0 {
			// Winner
			for i, p := range g.players {
				g.score[p] -= ncards[i]
				g.kitty += ncards[i]
			}
			g.score[g.players[controller]] += g.kitty
			g.kitty = 0
			return
		}

		// lowest returns the lowest Card held by a player in the given Suit,
		// and a boolean value that is true if the player holds any Cards in the Suit and false otherwise.
		lowest := func(i int, s Suit) (c Card, ok bool) {
			for c = s.Two(); c <= s.Ace(); c++ {
				if holder[c] == i {
					return c, true
				}
			}
			return 0, false
		}

		// hasColor returns whether a player holds any cards of a given Color.
		hasColor := func(i int, c Color) bool {
			if _, ok := lowest(i, c.Major()); ok {
				return true
			}
			_, ok := lowest(i, c.Minor())
			return ok
		}

		newColor := lead.Color().Other()
		for old := controller; ; controller = (controller + 1) % nplayers {
			if hasColor(controller, newColor) {
				break
			}
			g.score[g.players[controller]]--
			g.kitty++
			if (controller+1)%nplayers == old {
				// Hand is finished, no winner
				for i, p := range g.players {
					g.score[p] -= ncards[i]
					g.kitty += ncards[i]
				}
				return
			}
		}

		minorLead, hasMinor := lowest(controller, newColor.Minor())
		majorLead, hasMajor := lowest(controller, newColor.Major())
		switch {
		case hasMinor && hasMajor:
			if g.players[controller].PlayMajor(newColor) {
				lead = majorLead
			} else {
				lead = minorLead
			}
		case hasMajor:
			lead = majorLead
		case hasMinor:
			lead = minorLead
		}
	}
}
