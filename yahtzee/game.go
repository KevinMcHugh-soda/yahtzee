package yahtzee

import (
	"math/rand"
	"time"
)

type RollDecision [5]bool

func (rd *RollDecision) WillKeepAll() bool {
	for _, b := range rd {
		if !b {
			return false
		}
	}
	return true
}

func (rd *RollDecision) ToReroll() []int {
	torr := make([]int, 0)
	for idx, b := range rd {
		if !b {
			torr = append(torr, idx)
		}
	}
	return torr
}

type Game struct {
	Players []Player
	Winner  []Player
}

func (g *Game) getRoll(hand Hand) Hand {
	retVal := Hand{}
	for idx, val := range hand {
		if val == 0 {
			retVal[idx] = rand.Intn(5) + 1
		}
	}
	return retVal
}

func (g *Game) Play() {
	rand.Seed(time.Now().Unix())
	for idx := 0; idx < ScoreableCount; idx++ {
		for _, plyr := range g.Players {
			g.playTurn(plyr)
		}
	}
}

func (g *Game) playTurn(p Player) {
	var hand Hand
	for i := 0; i < 2; i++ {
		hand = g.getRoll(Hand{})
		rd := p.AssessRoll(hand)
		if rd.WillKeepAll() {
			break
		}
	}
	scorable := p.PickScorable(hand)
	scorecard := p.GetScorecard()
	scorecard.Score(&hand, scorable)
}
