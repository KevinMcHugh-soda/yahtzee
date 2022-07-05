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
	GetScorecard() *Scorecard
	AssessRoll(hand Hand) RollDecision
	PickScorable(hand Hand) Scoreable
}

type HumanPlayer struct {
	Scorecard *Scorecard
}

func (p HumanPlayer) GetName() string {
	return "Mr. Human!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!"
}

func (p HumanPlayer) GetScorecard() *Scorecard {
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

func (p HumanPlayer) EnsureValidResponse(prompt string, isValid func(string) bool) string {
	var input string
	for {
		fmt.Println(prompt)
		reader := bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		input = strings.Trim(input, "\n")

		if isValid(input) {
			break
		}
	}

	return input
}

func (p HumanPlayer) PickScorable(hand Hand) Scoreable {
	fmt.Printf("Hand: %d, %d, %d, %d, %d, \n", hand[0], hand[1], hand[2], hand[3], hand[4])
	prompt := "Choose a row to score this roll\n"
	// TODO present the entire current scoreboard, with current scores

	// TODO 1 index this at some point
	validScorableSelections := make(map[int]bool, len(ScorableNames))
	for idx, name := range ScorableNames {
		if p.Scorecard.NameToScorePtr(name) == nil {
			prompt += fmt.Sprintf("[%d] to score %s; ", idx, name)
			validScorableSelections[idx] = true
		}
	}
	input := p.EnsureValidResponse(prompt, func(input string) bool {
		val, err := strconv.Atoi(input)

		return err == nil && val >= 0 && val < len(ScorableNames) && validScorableSelections[val]
	})
	choice, _ := strconv.Atoi(input)

	return ScoreableByName(ScorableNames[choice])
}

func NewHumanPlayer() *Player {
	scoreCard := Scorecard{}
	hp := HumanPlayer{
		Scorecard: &scoreCard,
	}

	p := Player(hp)

	return &p
}
