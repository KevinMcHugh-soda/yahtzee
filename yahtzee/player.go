package yahtzee

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Player interface {
	GetName() string
	GetScorecard() Scorecard
	AssessRoll(hand Hand) RollDecision
	PickScorable(hand Hand) Scoreable
}

type HumanPlayer struct {
	Scorecard Scorecard
}

func (p HumanPlayer) GetName() string {
	return "Mr. Human"
}

func (p HumanPlayer) GetScorecard() Scorecard {
	return p.Scorecard
}

// TODO would be good to indicate roll no./rolls remaining
func (p HumanPlayer) AssessRoll(hand Hand) RollDecision {
	fmt.Printf("Roll: %d, %d, %d, %d, %d, \n", hand[0], hand[1], hand[2], hand[3], hand[4])
	// TODO present the entire current scoreboard
	fmt.Println("Type y to keep, space to reroll:")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	bools := make([]bool, 5)
	for idx, c := range text {
		if idx < len(bools) {
			bools[idx] = c == rune('y')
		}
	}
	return RollDecision(bools)
}

func (p HumanPlayer) PickScorable(hand Hand) Scoreable {
	fmt.Printf("Hand: %d, %d, %d, %d, %d, \n", hand[0], hand[1], hand[2], hand[3], hand[4])
	scorableNames := []string{"ones", "twos", "threes"}
	fmt.Println("Choose a row to score this roll")
	// TODO present the entire current scoreboard, with current scores
	// TODO don't prompt someone to use a row twice
	for idx, name := range scorableNames {
		fmt.Printf("[%d] to score %s; ", idx, name)
	}
	fmt.Println("")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\n")
	// TODO Handle invalid selection - build a function which can prompt until a valid selection is made.
	choice, err := strconv.Atoi(text)
	if err != nil {
		panic(err)
	}

	switch scorableNames[choice] {
	case "ones":
		return Ones{}
	case "twos":
		return Twos{}
	case "threes":
		return Threes{}
	}

	return ErrorScore{}
}

func NewHumanPlayer() HumanPlayer {
	return HumanPlayer{
		Scorecard: Scorecard{
			// TODO what why
			// oh it's because they're nil in subtotal, fuck me right
			ones:           new(int),
			twos:           new(int),
			threes:         new(int),
			fours:          new(int),
			fives:          new(int),
			sixes:          new(int),
			threeOfAKind:   new(int),
			fourOfAKind:    new(int),
			fullHouse:      new(int),
			smallStraight:  new(int),
			largeStraight:  new(int),
			chance:         new(int),
			yahtzee:        new(int),
			yahtzeeBonuses: []int{},
		},
	}
}
