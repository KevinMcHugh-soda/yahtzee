package yahtzee

import (
	"fmt"
	"math/rand"
	"time"
)

// TODO it's nice to have size in here ?
type RollDecision []bool

func (rd *RollDecision) WillKeepAll() bool {
	for _, b := range *rd {
		if !b {
			return false
		}
	}
	return true
}

type Game struct {
	Players []*Player
	Winner  []*Player
}

func (g *Game) getRoll(hand Hand, rd RollDecision) Hand {
	retVal := Hand{}
	for idx, keep := range rd {
		if keep {
			retVal[idx] = hand[idx]
		} else {
			retVal[idx] = rand.Intn(5) + 1
		}
	}
	return retVal
}

func (g *Game) Play() {
	rand.Seed(time.Now().Unix())
	for idx := 0; idx < ScoreableCount; idx++ {
		for _, plyr := range g.Players {
			g.playTurn(*plyr)
		}
	}
}

func (g *Game) playTurn(p Player) {
	hand1 := Hand{rand.Intn(5) + 1, rand.Intn(5) + 1, rand.Intn(5) + 1, rand.Intn(5) + 1, rand.Intn(5) + 1}
	rd1 := p.AssessRoll(hand1)
	if rd1.WillKeepAll() {
		score(p, hand1)
	} else {
		hand2 := g.getRoll(hand1, rd1)
		rd2 := p.AssessRoll(hand2)

		if rd2.WillKeepAll() {
			score(p, hand2)
		} else {
			hand3 := g.getRoll(hand2, rd2)
			score(p, hand3)
		}
	}
}

func score(p Player, hand Hand) {
	scorable := p.PickScorable(hand)
	scorecard := p.GetScorecard()

	scorecard.Score(&hand, scorable)
	fmt.Println(p.GetName(), scorecard.Total())
}
