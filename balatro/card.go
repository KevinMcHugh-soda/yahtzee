package balatro

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

type Suit int

const (
	Hearts Suit = iota
	Diamonds
	Clubs
	Spades
)

func (s Suit) String() string {
	switch s {
	case Hearts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	case Spades:
		return "♠"
	default:
		return "?"
	}
}

type Rank int

const (
	Two Rank = iota + 2
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

func (r Rank) String() string {
	switch r {
	case Ace:
		return "A"
	case King:
		return "K"
	case Queen:
		return "Q"
	case Jack:
		return "J"
	case Ten:
		return "10"
	default:
		return fmt.Sprintf("%d", int(r))
	}
}

// GetValue returns the scoring value for a card (face cards = 10, aces = 11)
func (r Rank) GetValue() int {
	if r >= Jack && r <= King {
		return 10
	}
	if r == Ace {
		return 11
	}
	return int(r)
}

type Card struct {
	Suit Suit
	Rank Rank
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}

func (c Card) GetValue() int {
	return c.Rank.GetValue()
}

type Deck struct {
	Cards []Card
}

func NewDeck() *Deck {
	deck := &Deck{Cards: make([]Card, 0, 52)}
	
	for suit := Hearts; suit <= Spades; suit++ {
		for rank := Two; rank <= Ace; rank++ {
			deck.Cards = append(deck.Cards, Card{Suit: suit, Rank: rank})
		}
	}
	
	return deck
}

func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

func (d *Deck) Draw(n int) []Card {
	if n > len(d.Cards) {
		n = len(d.Cards)
	}
	
	drawn := make([]Card, n)
	copy(drawn, d.Cards[:n])
	d.Cards = d.Cards[n:]
	
	return drawn
}

type Hand []Card

func (h Hand) String() string {
	var cards []string
	for _, card := range h {
		cards = append(cards, card.String())
	}
	return strings.Join(cards, " ")
}

// Sort hand by rank for easier evaluation
func (h Hand) Sort() {
	sort.Slice(h, func(i, j int) bool {
		return h[i].Rank < h[j].Rank
	})
}

// GetTotalValue returns the sum of all card values in the hand
func (h Hand) GetTotalValue() int {
	total := 0
	for _, card := range h {
		total += card.GetValue()
	}
	return total
}