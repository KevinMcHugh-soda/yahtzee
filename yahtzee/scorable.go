package yahtzee

// If we add more categories, increment this? yuck
const ScoreableCount = 13

type Hand [5]int

type Scoreable interface {
	Score(hand Hand, hadYahtzee bool) int
	MaxPossible() int
	ProbabilityToHit(hand Hand, rollsRemaining int) float64
}

type Ones struct{}
type Twos struct{}
type Threes struct{}
type Fours struct{}
type Fives struct{}
type Sixes struct{}
type ThreeOfAKind struct{}
type FourOfAKind struct{}
type FullHouse struct{}
type SmallStraight struct{}
type LargeStraight struct{}
type Chance struct{}
type Yahtzee struct{}
type ErrorScore struct{}

func (s ErrorScore) Score(hand Hand, hadYahtzee bool) int {
	return 0
}

func (s Ones) Score(hand Hand, hadYahtzee bool) int {
	return scoreFaceValues(hand, 1)
}
func (s Twos) Score(hand Hand, hadYahtzee bool) int {
	return scoreFaceValues(hand, 2)
}
func (s Threes) Score(hand Hand, hadYahtzee bool) int {
	return scoreFaceValues(hand, 3)
}
func (s Fours) Score(hand Hand, hadYahtzee bool) int {
	return scoreFaceValues(hand, 4)
}
func (s Fives) Score(hand Hand, hadYahtzee bool) int {
	return scoreFaceValues(hand, 5)
}
func (s Sixes) Score(hand Hand, hadYahtzee bool) int {
	return scoreFaceValues(hand, 6)
}
func (s ThreeOfAKind) Score(hand Hand, hadYahtzee bool) int {
	scoring := false
	for _, count := range valueCounts(hand) {
		if count >= 3 {
			scoring = true
		}
	}
	if scoring {
		score := 0
		for _, value := range hand {
			score += value
		}
		return score
	}
	return 0
}
func (s FourOfAKind) Score(hand Hand, hadYahtzee bool) int {
	scoring := false
	for _, count := range valueCounts(hand) {
		if count >= 4 {
			scoring = true
		}
	}
	if scoring {
		score := 0
		for _, value := range hand {
			score += value
		}
		return score
	}
	return 0
}
func (s FullHouse) Score(hand Hand, hadYahtzee bool) int {
	if isJoker(hand, hadYahtzee) {
		return 25
	}

	hasTwo, hasThree := false, false
	for _, count := range valueCounts(hand) {
		if count == 2 {
			hasTwo = true
		} else if count == 3 {
			hasThree = true
		} else if count == 5 {
			hasTwo = true
			hasThree = true
		}
	}

	if hasTwo && hasThree {
		return 25
	}
	return 0
}

func (s SmallStraight) Score(hand Hand, hadYahtzee bool) int {
	if isJoker(hand, hadYahtzee) {
		return 30
	}
	valueCounts := valueCounts(hand)

	scoring := (valueCounts[3] >= 1 && valueCounts[4] >= 1) &&
		((valueCounts[1] >= 1 && valueCounts[2] >= 1) ||
			(valueCounts[2] >= 1 && valueCounts[5] >= 1) ||
			(valueCounts[5] >= 1 && valueCounts[6] >= 1))
	if scoring {
		return 30
	}
	return 0
}
func (s LargeStraight) Score(hand Hand, hadYahtzee bool) int {
	if isJoker(hand, hadYahtzee) {
		return 40
	}
	valueCounts := valueCounts(hand)

	scoring := (valueCounts[2] == 1 && valueCounts[3] == 1 && valueCounts[4] == 1 && valueCounts[5] == 1) && ((valueCounts[1] == 1) || (valueCounts[6] == 1))
	if scoring {
		return 40
	}
	return 0
}
func (s Chance) Score(hand Hand, hadYahtzee bool) int {
	score := 0
	for _, value := range hand {
		score = score + value
	}
	return score
}

func (s Yahtzee) Score(hand Hand, hadYahtzee bool) int {
	if isYahtzee(hand) {
		return 50
	}
	return 0
}

func valueCounts(hand Hand) map[int]int {
	valueCounts := make(map[int]int, 0)

	for _, value := range hand {
		valueCounts[value] = valueCounts[value] + 1
	}
	return valueCounts
}

func scoreFaceValues(hand Hand, value int) int {
	foundValue := 0
	for _, die := range hand {
		if die == value {
			foundValue += die
		}
	}
	return foundValue
}

func isYahtzee(hand Hand) bool {
	for _, count := range valueCounts(hand) {
		if count == 5 {
			return true
		}
	}

	return false
}

func isJoker(hand Hand, hadYahztee bool) bool {
	return isYahtzee(hand) && hadYahztee
}
