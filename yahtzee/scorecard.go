package yahtzee

type Scorecard struct {
	ones   *int
	twos   *int
	threes *int
	fours  *int
	fives  *int
	sixes  *int

	threeOfAKind   *int
	fourOfAKind    *int
	fullHouse      *int
	smallStraight  *int
	largeStraight  *int
	chance         *int
	yahtzee        *int
	yahtzeeBonuses []int
}

func (s *Scorecard) Score(hand *Hand, scoreable Scoreable) int {
	sc := scoreable.Score(*hand)
	score := &sc
	s.scoreYahtzeeBonus(*hand)
	switch scoreable.(type) {
	case Ones:
		s.ones = score
	case Twos:
		s.twos = score
	case Threes:
		s.threes = score
	case Fours:
		s.fours = score
	case Fives:
		s.fives = score
	case Sixes:
		s.sixes = score
	case ThreeOfAKind:
		s.threeOfAKind = score
	case FourOfAKind:
		s.fourOfAKind = score
	case FullHouse:
		s.fullHouse = score
	case SmallStraight:
		s.smallStraight = score
	case LargeStraight:
		s.largeStraight = score
	case Chance:
		s.chance = score
	case Yahtzee:
		s.yahtzee = score
	}
	return *score
}

func (s *Scorecard) Subtotal() int {
	return *s.ones + *s.twos + *s.threes + *s.fours + *s.fives + *s.sixes
}
func (s *Scorecard) Total() int {
	sub := s.Subtotal()
	total := sub
	if sub > 63 {
		total = total + 25
	}
	bonusPoints := 0
	for _, bonus := range s.yahtzeeBonuses {
		bonusPoints = bonusPoints + bonus
	}
	return total + *s.threeOfAKind + *s.fourOfAKind + *s.fullHouse + *s.smallStraight + *s.largeStraight + *s.chance + *s.yahtzee + bonusPoints
}

func (s *Scorecard) scoreYahtzeeBonus(hand Hand) int {
	if *s.yahtzee == 0 {
		return 0
	}
	for _, count := range valueCounts(hand) {
		if count == 5 {
			s.yahtzeeBonuses = append(s.yahtzeeBonuses, 100)
		}
	}
	return 100
}
