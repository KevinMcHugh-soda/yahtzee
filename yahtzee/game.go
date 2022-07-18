package yahtzee

import (
	"fmt"
	"math/rand"
	"sort"
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
	Winner  []Player
}

func (g *Game) getRoll(hand Hand, rd RollDecision) Hand {
	retSlice := make([]int, 5)
	for idx, keep := range rd {
		if keep {
			retSlice[idx] = hand[idx]
		} else {
			retSlice[idx] = rollDie()
		}
	}
	sort.Ints(retSlice)
	return Hand{retSlice[0], retSlice[1], retSlice[2], retSlice[3], retSlice[4]}
}

func rollDie() int {
	return rand.Intn(6) + 1
}

func (g *Game) Play() {
	rand.Seed(time.Now().Unix())
	for idx := 0; idx < ScoreableCount; idx++ {
		for _, plyr := range g.Players {
			g.playTurn(*plyr)
		}
	}
	topScore := 0
	for _, plyr := range g.Players {
		total := (*plyr).GetScorecard().Total()
		if topScore > total {
			topScore = total
			g.Winner = []Player{*plyr}
		} else if topScore == total {
			g.Winner = append(g.Winner, *plyr)
		}
	}
}

func (g *Game) playTurn(p Player) {
	hSlice := []int{rollDie(), rollDie(), rollDie(), rollDie(), rollDie()}
	sort.Ints(hSlice)
	hand1 := Hand{hSlice[0], hSlice[1], hSlice[2], hSlice[3], hSlice[4]}
	// hand1 := Hand{6, 6, 6, 6, 6}
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
	fmt.Println(p.GetName())
	fmt.Println(p.GetScorecard().Print())
}
