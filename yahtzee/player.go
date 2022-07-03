package yahtzee

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

func (p *HumanPlayer) GetName() string {
	return "Mr. Human"
}

func (p *HumanPlayer) GetScorecard() Scorecard {
	return p.Scorecard
}

// TODO would be good to indicate roll no./rolls remaining
func (p *HumanPlayer) AssessRoll(hand Hand) RollDecision {
	fmt.Printf("Roll: %d, %d, %d, %d, %d, \n", hand[0], hand[1], hand[2], hand[3], hand[4])
	// TODO present the entire current scoreboard
	fmt.Println("Type y to keep, space to reroll:")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	var bools [5]bool
	for idx, c := range text {
		bools[idx] = c == rune('y')
	}
	return RollDecision(bools)
}

func (p *HumanPlayer) PickScorable(hand Hand) Scoreable {
	fmt.Printf("Hand: %d, %d, %d, %d, %d, \n", hand[0], hand[1], hand[2], hand[3], hand[4])
	scorableNames := []string{"ones", "twos", "threes"}
	fmt.Printf("Choose a row to score this roll")
	// TODO present the entire current scoreboard, with current scores
	for idx, name := range scorableNames {
		fmt.Printf("[%d] to score %s", idx, name)
	}
	fmt.Println("")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
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
