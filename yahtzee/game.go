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
	Players []Player
	Winner  []Player
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
			g.playTurn(plyr)
		}
	}
}

func (g *Game) playTurn(p Player) {
	var hand Hand
	rd := RollDecision(make([]bool, 5))
	for i := 0; i < 2; i++ {
		hand = g.getRoll(hand, rd)
		rd = p.AssessRoll(hand)
		if rd.WillKeepAll() {
			break
		}
	}
	scorable := p.PickScorable(hand)
	scorecard := p.GetScorecard()
	scorecard.Score(&hand, scorable)
	// TODO: entering a new score erases your previous - sounds like I'm copying by value instead of by reference?

	fmt.Println(p.GetName(), scorecard.Total())
}
