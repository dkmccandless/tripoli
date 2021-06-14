package game

import (
	"reflect"
	"testing"

	"github.com/dkmccandless/tripoli/card"
)

// minor and major are Players with deterministic behavior for testing purposes.

// minor plays a minor suit whenever possible.
type minor struct{ n int }

func (m *minor) Init(int, int, []card.Card, map[card.Card]int, int) {}

func (m *minor) Note(int, card.Card) {}

func (m *minor) PlayMajor(card.Color) bool { return false }

var pa, pb = &minor{0}, &minor{1}

// major plays a major suit whenever possible.
type major struct{ n int }

func (m *major) Init(int, int, []card.Card, map[card.Card]int, int) {}

func (m *major) Note(int, card.Card) {}

func (m *major) PlayMajor(card.Color) bool { return true }

var pc, pd = &major{0}, &major{1}

func TestNew(t *testing.T) {
	for _, test := range []struct {
		players []Player
		g       *Game
	}{
		{
			[]Player{pa, pb},
			&Game{
				score: map[Player]int{pa: 0, pb: 0},
				stake: counters(0, 0, 0, 0, 0),
			},
		},
		{
			[]Player{pa, pb, pc},
			&Game{
				score: map[Player]int{pa: 0, pb: 0, pc: 0},
				stake: counters(0, 0, 0, 0, 0),
			},
		},
		{
			[]Player{pa, pb, pc, pd},
			&Game{
				score: map[Player]int{pa: 0, pb: 0, pc: 0, pd: 0},
				stake: counters(0, 0, 0, 0, 0),
			},
		},
	} {
		if g := New(test.players); !reflect.DeepEqual(g, test.g) {
			t.Errorf("New(%v): Game is %+v, expected %+v",
				test.players, g, test.g,
			)
		}
	}
}

func TestInit(t *testing.T) {
	for _, test := range []struct {
		players []Player
		n       []int
	}{
		{
			[]Player{pa, pb},
			[]int{18, 17},
		},
		{
			[]Player{pa, pb, pc},
			[]int{13, 13, 13},
		},
		{
			[]Player{pa, pb, pc, pd},
			[]int{11, 11, 10, 10},
		},
	} {
		g := New(test.players)
		r := g.init()
		if g != r.g {
			t.Fatalf("round %+v wraps Game %+v, expected %+v", r, r.g, g)
		}
		if len(r.p) != len(test.players) {
			t.Fatalf("len(p) is %v, expected %v", len(r.p), len(test.players))
		}
		for i := range test.players {
			var ok bool
			for j := range r.p {
				if test.players[i] == r.p[j] {
					ok = true
					break
				}
			}
			if !ok {
				t.Errorf("player %v is not in p", i)
			}
		}
		if len(r.deck) != 52 {
			t.Errorf("len(deck) is %v, expected 52", len(r.deck))
		}
		if !reflect.DeepEqual(r.n, test.n) {
			t.Errorf("n is %v, expected %v", r.n, test.n)
		}
	}
}

func TestFirstLead(t *testing.T) {
	for _, test := range []struct {
		r *round
		c card.Card
	}{
		{r: &round{deck: make([]int, 52)}, c: 0},
		{r: &round{deck: append([]int{-1}, make([]int, 51)...)}, c: 1},
		{r: &round{deck: append([]int{0, -1}, make([]int, 50)...)}, c: 0},
		{r: &round{deck: append([]int{-1, -1}, make([]int, 50)...)}, c: 2},
	} {
		if c := test.r.firstLead(); c != test.c {
			t.Errorf("firstLead(%+v): got %v, expected %v", test.r, c, test.c)
		}
	}
}

func TestAnte(t *testing.T) {
	for name, test := range map[string]struct{ r, want *round }{
		"2 players": {
			&round{
				g: &Game{
					score: map[Player]int{pa: 0, pb: 0},
					stake: counters(0, 0, 0, 0, 0),
				},
			},
			&round{
				g: &Game{
					score: map[Player]int{pa: -5, pb: 0},
					stake: counters(1, 1, 1, 1, 1),
				},
			},
		},
		"3 players": {
			&round{
				g: &Game{
					score: map[Player]int{pa: 3, pb: -4, pc: -5},
					stake: counters(0, 0, 0, 3, 3),
				},
			},
			&round{
				g: &Game{
					score: map[Player]int{pa: -2, pb: -4, pc: -5},
					stake: counters(1, 1, 1, 4, 4),
				},
			},
		},
		"4 players": {
			&round{
				g: &Game{
					score: map[Player]int{pa: -10, pb: 3, pc: -4, pd: 1},
					stake: counters(0, 4, 0, 0, 0),
					kitty: 6,
				},
			},
			&round{
				g: &Game{
					score: map[Player]int{pa: -15, pb: 3, pc: -4, pd: 1},
					stake: counters(1, 5, 1, 1, 1),
					kitty: 6,
				},
			},
		},
	} {
		if test.r.ante(pa); !reflect.DeepEqual(test.r, test.want) {
			t.Errorf("ante(%q): round is %+v, expected %+v",
				name, test.r, test.want,
			)
		}
	}
}

func TestCollect(t *testing.T) {
	for name, test := range map[string]struct {
		r    *round
		c    card.Card
		want *round
	}{
		"no counter": {
			r: &round{
				g: &Game{
					score: map[Player]int{pa: 0, pb: 0, pc: 0},
					stake: counters(3, 3, 3, 3, 3),
				},
			},
			c: 0,
			want: &round{
				g: &Game{
					score: map[Player]int{pa: 0, pb: 0, pc: 0},
					stake: counters(3, 3, 3, 3, 3),
				},
			},
		},
		"counter": {
			r: &round{
				g: &Game{
					score: map[Player]int{pa: 0, pb: 0, pc: 0},
					stake: counters(3, 3, 3, 3, 3),
				},
			},
			c: 50,
			want: &round{
				g: &Game{
					score: map[Player]int{pa: 3, pb: 0, pc: 0},
					stake: counters(3, 3, 3, 0, 3),
				},
			},
		},
	} {
		if test.r.collect(pa, test.c); !reflect.DeepEqual(test.r, test.want) {
			t.Errorf("collect(%q): round is %+v, expected %+v",
				name, test.r, test.want,
			)
		}
	}
}

func TestPayKitty(t *testing.T) {
	for name, test := range map[string]struct {
		r    *round
		n    int
		want *round
	}{
		"2 players": {
			&round{
				g: &Game{
					score: map[Player]int{pa: 0, pb: 0},
				},
			},
			1,
			&round{
				g: &Game{
					score: map[Player]int{pa: -1, pb: 0},
					kitty: 1,
				},
			},
		},
		"3 players": {
			&round{
				g: &Game{
					score: map[Player]int{pa: 3, pb: -4, pc: -5},
					kitty: 1,
				},
			},
			1,
			&round{
				g: &Game{
					score: map[Player]int{pa: 2, pb: -4, pc: -5},
					kitty: 2,
				},
			},
		},
		"4 players": {
			&round{
				g: &Game{
					score: map[Player]int{pa: -10, pb: 3, pc: -4, pd: 1},
					kitty: 6,
				},
			},
			4,
			&round{
				g: &Game{
					score: map[Player]int{pa: -14, pb: 3, pc: -4, pd: 1},
					kitty: 10,
				},
			},
		},
	} {
		if test.r.payKitty(pa, test.n); !reflect.DeepEqual(test.r, test.want) {
			t.Errorf("payKitty(%q): round is %+v, expected %+v",
				name, test.r, test.want,
			)
		}
	}
}

func TestCollectKitty(t *testing.T) {
	for name, test := range map[string]struct{ r, want *round }{
		"2 players": {
			&round{
				g: &Game{
					score: map[Player]int{pa: 0, pb: 0},
					kitty: 1,
				},
			},
			&round{
				g: &Game{
					score: map[Player]int{pa: 1, pb: 0},
					kitty: 0,
				},
			},
		},
		"3 players": {
			&round{
				g: &Game{
					score: map[Player]int{pa: 3, pb: -4, pc: -5},
					kitty: 1,
				},
			},
			&round{
				g: &Game{
					score: map[Player]int{pa: 4, pb: -4, pc: -5},
					kitty: 0,
				},
			},
		},
		"4 players": {
			&round{
				g: &Game{
					score: map[Player]int{pa: -10, pb: 3, pc: -4, pd: 1},
					kitty: 6,
				},
			},
			&round{
				g: &Game{
					score: map[Player]int{pa: -4, pb: 3, pc: -4, pd: 1},
					kitty: 0,
				},
			},
		},
	} {
		if test.r.collectKitty(pa); !reflect.DeepEqual(test.r, test.want) {
			t.Errorf("collectKitty(%q): round is %+v, expected %+v",
				name, test.r, test.want,
			)
		}
	}
}

func TestPlayCard(t *testing.T) {
	for name, test := range map[string]struct {
		r    *round
		c    card.Card
		want *round
	}{
		"no counter": {
			r: &round{
				g: &Game{
					score: map[Player]int{pa: 0, pb: 0, pc: 0},
					stake: counters(3, 3, 3, 3, 3),
				},
				p: []Player{pb, pc, pa},
				n: []int{13, 13, 13},
				deck: []int{
					2, 0, 0, -1, -1, 0, -1, -1, 0, 0, 0, 0, 1,
					-1, 1, 2, -1, 2, 0, 2, 1, 2, -1, 2, 0, 0,
					0, 0, -1, 1, -1, 1, 1, 1, 1, 2, 2, -1, 1,
					2, 2, 1, -1, 2, 2, 2, 1, -1, -1, 1, 0, 1,
				},
			},
			c: 0,
			want: &round{
				g: &Game{
					score: map[Player]int{pa: 0, pb: 0, pc: 0},
					stake: counters(3, 3, 3, 3, 3),
				},
				p: []Player{pb, pc, pa},
				n: []int{13, 13, 12},
				deck: []int{
					-1, 0, 0, -1, -1, 0, -1, -1, 0, 0, 0, 0, 1,
					-1, 1, 2, -1, 2, 0, 2, 1, 2, -1, 2, 0, 0,
					0, 0, -1, 1, -1, 1, 1, 1, 1, 2, 2, -1, 1,
					2, 2, 1, -1, 2, 2, 2, 1, -1, -1, 1, 0, 1,
				},
			},
		},
		"counter": {
			r: &round{
				g: &Game{
					score: map[Player]int{pa: 0, pb: 0, pc: 0},
					stake: counters(3, 3, 3, 3, 3),
				},
				p: []Player{pb, pc, pa},
				n: []int{11, 13, 12},
				deck: []int{
					-1, -1, -1, -1, -1, 0, -1, -1, 0, 0, 0, 0, 1,
					-1, 1, 2, -1, 2, 0, 2, 1, 2, -1, 2, 0, 0,
					0, 0, -1, 1, -1, 1, 1, 1, 1, 2, 2, -1, 1,
					2, 2, 1, -1, 2, 2, 2, 1, -1, -1, 1, 0, 1,
				},
			},
			c: 50,
			want: &round{
				g: &Game{
					score: map[Player]int{pa: 0, pb: 3, pc: 0},
					stake: counters(3, 3, 3, 0, 3),
				},
				p: []Player{pb, pc, pa},
				n: []int{10, 13, 12},
				deck: []int{
					-1, -1, -1, -1, -1, 0, -1, -1, 0, 0, 0, 0, 1,
					-1, 1, 2, -1, 2, 0, 2, 1, 2, -1, 2, 0, 0,
					0, 0, -1, 1, -1, 1, 1, 1, 1, 2, 2, -1, 1,
					2, 2, 1, -1, 2, 2, 2, 1, -1, -1, 1, -1, 1,
				},
			},
		},
		"out": {
			r: &round{
				g: &Game{
					score: map[Player]int{pa: 0, pb: 3, pc: 3},
					stake: counters(3, 3, 3, 0, 0),
				},
				p: []Player{pb, pc, pa},
				n: []int{8, 1, 3},
				deck: []int{
					-1, -1, -1, -1, -1, 0, -1, -1, 0, 0, 0, 0, -1,
					-1, -1, -1, -1, 2, 0, 2, -1, -1, -1, 2, 0, 0,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 1, -1, -1,
				},
			},
			c: 49,
			want: &round{
				g: &Game{
					score: map[Player]int{pa: 0, pb: 3, pc: 6},
					stake: counters(3, 3, 0, 0, 0),
				},
				p: []Player{pb, pc, pa},
				n: []int{8, 0, 3},
				deck: []int{
					-1, -1, -1, -1, -1, 0, -1, -1, 0, 0, 0, 0, -1,
					-1, -1, -1, -1, 2, 0, 2, -1, -1, -1, 2, 0, 0,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
				},
			},
		},
	} {
		if test.r.playCard(test.c); !reflect.DeepEqual(test.r, test.want) {
			t.Errorf("playCard(%q): round is %+v, expected %+v",
				name, test.r, test.want,
			)
		}
	}
}

func TestPlayRun(t *testing.T) {
	for name, test := range map[string]struct {
		r    *round
		c    card.Card
		pos  int
		won  bool
		want *round
	}{
		"extra hand": {
			&round{
				g: &Game{
					score: map[Player]int{pa: -1, pb: 2, pc: -4},
					stake: counters(0, 0, 0, 0, 0),
					kitty: 3,
				},
				p: []Player{pa, pb, pc},
				n: []int{8, 9, 6},
				deck: []int{
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					-1, -1, 1, 0, 0, -1, 1, -1, -1, -1, 0, 1, 0,
					-1, 2, 2, 0, 2, 2, 2, 1, 2, 1, 1, 0, 0,
					-1, 1, 1, 1, -1, 0, -1, -1, -1, -1, -1, -1, -1,
				},
			},
			15,
			0,
			false,
			&round{
				g: &Game{
					score: map[Player]int{pa: -1, pb: 2, pc: -4},
					stake: counters(0, 0, 0, 0, 0),
					kitty: 3,
				},
				p: []Player{pa, pb, pc},
				n: []int{6, 8, 6},
				deck: []int{
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					-1, -1, -1, -1, -1, -1, 1, -1, -1, -1, 0, 1, 0,
					-1, 2, 2, 0, 2, 2, 2, 1, 2, 1, 1, 0, 0,
					-1, 1, 1, 1, -1, 0, -1, -1, -1, -1, -1, -1, -1,
				},
			},
		},
		"ace": {
			&round{
				g: &Game{
					score: map[Player]int{pa: -1, pb: 2, pc: -4},
					stake: counters(0, 0, 0, 0, 0),
					kitty: 3,
				},
				p: []Player{pa, pb, pc},
				n: []int{8, 9, 6},
				deck: []int{
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					-1, -1, 1, 0, 0, -1, 1, -1, -1, -1, 0, 1, 0,
					-1, 2, 2, 0, 2, 2, 2, 1, 2, 1, 1, 0, 0,
					-1, 1, 1, 1, -1, 0, -1, -1, -1, -1, -1, -1, -1,
				},
			},
			29,
			0,
			false,
			&round{
				g: &Game{
					score: map[Player]int{pa: -1, pb: 2, pc: -4},
					stake: counters(0, 0, 0, 0, 0),
					kitty: 3,
				},
				p: []Player{pa, pb, pc},
				n: []int{5, 6, 2},
				deck: []int{
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					-1, -1, 1, 0, 0, -1, 1, -1, -1, -1, 0, 1, 0,
					-1, 2, 2, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					-1, 1, 1, 1, -1, 0, -1, -1, -1, -1, -1, -1, -1,
				},
			},
		},
		"out": {
			&round{
				g: &Game{
					score: map[Player]int{pa: -1, pb: 2, pc: -4},
					stake: counters(0, 0, 0, 0, 0),
					kitty: 3,
				},
				p: []Player{pa, pb, pc},
				n: []int{8, 9, 6},
				deck: []int{
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					-1, -1, 1, 0, 0, -1, 1, -1, -1, -1, 0, 1, 0,
					-1, 2, 2, 0, 2, 2, 2, 1, 2, 1, 1, 0, 0,
					-1, 1, 1, 1, -1, 0, -1, -1, -1, -1, -1, -1, -1,
				},
			},
			27,
			2,
			true,
			&round{
				g: &Game{
					score: map[Player]int{pa: -1, pb: 2, pc: -4},
					stake: counters(0, 0, 0, 0, 0),
					kitty: 3,
				},
				p: []Player{pa, pb, pc},
				n: []int{7, 8, 0},
				deck: []int{
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					-1, -1, 1, 0, 0, -1, 1, -1, -1, -1, 0, 1, 0,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, 1, 1, 0, 0,
					-1, 1, 1, 1, -1, 0, -1, -1, -1, -1, -1, -1, -1,
				},
			},
		},
	} {
		pos, won := test.r.playRun(test.c)
		if pos != test.pos || won != test.won {
			t.Errorf("playRun(%q): got %v, %v; expected %v, %v",
				name, pos, won, test.pos, test.won,
			)
		}
		if !reflect.DeepEqual(test.r, test.want) {
			t.Errorf("playRun(%q): round is %+v, expected %+v",
				name, test.r, test.want,
			)
		}
	}
}

func TestLowest(t *testing.T) {
	r := &round{
		deck: []int{
			-1, -1, -1, -1, 0, -1, 0, -1, 1, 1, -1, 1, 1,
			-1, -1, -1, -1, -1, -1, 2, 0, 2, -1, -1, 2, 2,
			1, 0, -1, 2, -1, 1, -1, 2, 2, 0, -1, -1, 2,
			-1, 1, 0, -1, 0, -1, 1, -1, 0, 0, 1, 1, 2,
		},
	}
	for name, test := range map[string]struct {
		pos int
		s   card.Suit
		c   card.Card
		ok  bool
	}{
		"two":  {1, card.Spades, card.Spades.Rank(card.Two), true},
		"mid":  {0, card.Clubs, card.Clubs.Rank(card.Six), true},
		"ace":  {2, card.Hearts, card.Hearts.Rank(card.Ace), true},
		"void": {1, card.Diamonds, 0, false},
	} {
		c, ok := r.lowest(test.pos, test.s)
		if c != test.c || ok != test.ok {
			t.Errorf("lowest(%q): got %v, %v; expected %v, %v",
				name, c, ok, test.c, test.ok,
			)
		}
	}
}

func TestNextLead(t *testing.T) {
	for name, test := range map[string]struct {
		r     *round
		pos   int
		color card.Color
		lead  card.Card
		ok    bool
		g     *Game
	}{
		"continue, force minor": {
			&round{
				g: &Game{
					score: map[Player]int{pa: 0, pc: 0},
				},
				p: []Player{pa, pc},
				n: []int{1, 1},
				deck: []int{
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
				},
			},
			0,
			card.Red,
			card.Diamonds.Rank(card.Ace),
			true,
			&Game{
				score: map[Player]int{pa: 0, pc: 0},
			},
		},
		"continue, force major": {
			&round{
				g: &Game{
					score: map[Player]int{pa: 0, pc: 0},
				},
				p: []Player{pa, pc},
				n: []int{1, 1},
				deck: []int{
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
				},
			},
			0,
			card.Black,
			card.Spades.Rank(card.Ace),
			true,
			&Game{
				score: map[Player]int{pa: 0, pc: 0},
			},
		},
		"continue, choose minor": {
			&round{
				g: &Game{
					score: map[Player]int{pa: 0, pc: 0},
				},
				p: []Player{pa, pc},
				n: []int{2, 2},
				deck: []int{
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0,
				},
			},
			0,
			card.Red,
			card.Diamonds.Rank(card.Ace),
			true,
			&Game{
				score: map[Player]int{pa: 0, pc: 0},
			},
		},
		"continue, choose major": {
			&round{
				g: &Game{
					score: map[Player]int{pa: 0, pc: 0},
				},
				p: []Player{pa, pc},
				n: []int{2, 2},
				deck: []int{
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0,
				},
			},
			1,
			card.Black,
			card.Spades.Rank(card.Ace),
			true,
			&Game{
				score: map[Player]int{pa: 0, pc: 0},
			},
		},
		"pass, force minor": {
			&round{
				g: &Game{
					score: map[Player]int{pa: 0, pc: 0},
				},
				p: []Player{pa, pc},
				n: []int{1, 1},
				deck: []int{
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
				},
			},
			1,
			card.Red,
			card.Diamonds.Rank(card.Ace),
			true,
			&Game{
				score: map[Player]int{pa: 0, pc: -1},
				kitty: 1,
			},
		},
		"pass, force major": {
			&round{
				g: &Game{
					score: map[Player]int{pa: 0, pc: 0},
				},
				p: []Player{pa, pc},
				n: []int{1, 1},
				deck: []int{
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
				},
			},
			1,
			card.Black,
			card.Spades.Rank(card.Ace),
			true,
			&Game{
				score: map[Player]int{pa: 0, pc: -1},
				kitty: 1,
			},
		},
		"pass, choose minor": {
			&round{
				g: &Game{
					score: map[Player]int{pa: 0, pc: 0},
				},
				p: []Player{pa, pc},
				n: []int{2, 2},
				deck: []int{
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0,
				},
			},
			1,
			card.Red,
			card.Diamonds.Rank(card.Ace),
			true,
			&Game{
				score: map[Player]int{pa: 0, pc: -1},
				kitty: 1,
			},
		},
		"pass, choose major": {
			&round{
				g: &Game{
					score: map[Player]int{pa: 0, pc: 0},
				},
				p: []Player{pa, pc},
				n: []int{2, 2},
				deck: []int{
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0,
				},
			},
			0,
			card.Black,
			card.Spades.Rank(card.Ace),
			true,
			&Game{
				score: map[Player]int{pa: -1, pc: 0},
				kitty: 1,
			},
		},
		"round over": {
			&round{
				g: &Game{
					score: map[Player]int{pa: 0, pc: 0},
				},
				p: []Player{pa, pc},
				n: []int{1, 1},
				deck: []int{
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 0,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 1,
				},
			},
			1,
			card.Black,
			0,
			false,
			&Game{
				score: map[Player]int{pa: -1, pc: -1},
				kitty: 2,
			},
		},
	} {
		lead, ok := test.r.nextLead(test.pos, test.color)
		if lead != test.lead || ok != test.ok {
			t.Errorf("nextLead(%q): got %v, %v; expected %v, %v",
				name, lead, ok, test.lead, test.ok,
			)
		}
		if !reflect.DeepEqual(test.r.g, test.g) {
			t.Errorf("nextLead(%q): Game is %+v, expected %+v",
				name, test.r.g, test.g,
			)
		}
	}
}

func TestPlay(t *testing.T) {
	for name, test := range map[string]struct{ r, want *round }{
		"winner": {
			&round{
				g: &Game{
					score: map[Player]int{pa: -5, pb: -5, pc: -5, pd: -5},
					stake: counters(4, 4, 4, 4, 4),
				},
				p: []Player{pa, pb, pc, pd},
				n: []int{11, 11, 10, 10},
				deck: []int{
					-1, 2, 2, 0, 3, 3, -1, 1, -1, -1, 2, 2, 3,
					1, -1, 3, -1, 3, 0, 1, -1, 2, 1, 1, 1, 3,
					-1, 2, 3, 0, 1, 1, -1, 1, 0, 3, 2, 2, 0,
					3, 0, 0, 0, 3, 1, 0, 0, 1, 0, 2, 2, -1,
				},
			},
			&round{
				g: &Game{
					score: map[Player]int{pa: 9, pb: -4, pc: 0, pd: -9},
					stake: counters(0, 0, 0, 0, 4),
				},
				p: []Player{pa, pb, pc, pd},
				n: []int{0, 3, 3, 4},
				deck: []int{
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 2, 2, 3,
					-1, -1, 3, -1, 3, -1, -1, -1, 2, 1, 1, 1, 3,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
				},
			},
		},
		"no winner": {
			&round{
				g: &Game{
					score: map[Player]int{pa: -5, pb: -5, pc: -5, pd: -5},
					stake: counters(4, 4, 4, 4, 4),
				},
				p: []Player{pa, pb, pc, pd},
				n: []int{11, 11, 10, 10},
				deck: []int{
					1, 3, 0, 2, 2, 1, 2, 2, 0, 0, 1, 0, 0,
					2, 3, 1, 0, -1, -1, -1, 3, 0, -1, 0, -1, -1,
					-1, 1, 1, 1, 0, 3, 1, 3, 2, 0, 3, 3, -1,
					1, -1, 2, -1, 2, 3, 3, 3, 1, 1, 2, 0, 2,
				},
			},
			&round{
				g: &Game{
					score: map[Player]int{pa: -5, pb: 1, pc: -2, pd: -9},
					stake: counters(0, 0, 0, 0, 0),
					kitty: 15,
				},
				p: []Player{pa, pb, pc, pd},
				n: []int{2, 1, 3, 2},
				deck: []int{
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					2, 3, -1, -1, -1, -1, -1, 3, 0, -1, 0, -1, -1,
					-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
					1, -1, 2, -1, 2, -1, -1, -1, -1, -1, -1, -1, -1,
				},
			},
		},
	} {
		if test.r.play(); !reflect.DeepEqual(test.r, test.want) {
			t.Errorf("play(%q): round is %+v, expected %+v",
				name, test.r, test.want,
			)
		}
	}
}
