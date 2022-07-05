package yahtzee

import (
	"fmt"
	"strconv"
)

type ScorableName string

const (
	OnesName   = "Ones"
	TwosName   = "Twos"
	ThreesName = "Threes"
	FoursName  = "Fours"
	FivesName  = "Fives"
	SixesName  = "Sixes"

	ThreeOfAKindName  = "3 Of A Kind"
	FourOfAKindName   = "4 Of A Kind"
	FullHouseName     = "Full House"
	SmallStraightName = "Small Straight"
	LargeStraightName = "Large Straight"
	ChanceName        = "Chance"
	YahtzeeName       = "Yahtzee"
	YahtzeeBonusName  = "Yahtzee Bonus"

	ErrorName = "error"
)

var ScorableNames = []ScorableName{
	OnesName, TwosName, ThreesName, FoursName, FivesName, SixesName,
	ThreeOfAKindName, FourOfAKindName, FullHouseName, SmallStraightName, LargeStraightName, ChanceName, YahtzeeName,
	// Well, this isn't actually _scorable_, you can't record it, so, hrm.
	YahtzeeBonusName,
}

func ScoreableByName(name ScorableName) Scoreable {
	scorablesByName := map[ScorableName]Scoreable{
		OnesName:          Ones{},
		TwosName:          Twos{},
		ThreesName:        Threes{},
		FoursName:         Fours{},
		FivesName:         Fives{},
		SixesName:         Sixes{},
		ThreeOfAKindName:  ThreeOfAKind{},
		FourOfAKindName:   FourOfAKind{},
		FullHouseName:     FullHouse{},
		SmallStraightName: SmallStraight{},
		LargeStraightName: LargeStraight{},
		ChanceName:        Chance{},
		YahtzeeName:       Yahtzee{},
	}

	return scorablesByName[name]
}

type Scorecard map[ScorableName]*int

func (s *Scorecard) NameToScorePtr(name ScorableName) *int {
	m := *s
	// consider constructing this with a loop
	nameToPtr := map[ScorableName]*int{
		OnesName:   m[OnesName],
		TwosName:   m[TwosName],
		ThreesName: m[ThreesName],
		FoursName:  m[FoursName],
		FivesName:  m[FivesName],
		SixesName:  m[SixesName],

		ThreeOfAKindName:  m[ThreeOfAKindName],
		FourOfAKindName:   m[FourOfAKindName],
		SmallStraightName: m[SmallStraightName],
		LargeStraightName: m[LargeStraightName],
		ChanceName:        m[ChanceName],
		YahtzeeName:       m[YahtzeeName],
		YahtzeeBonusName:  m[YahtzeeBonusName],
	}

	return nameToPtr[name]
}

func (s *Scorecard) Score(hand *Hand, scoreable Scoreable) int {
	sc := scoreable.Score(*hand)
	score := &sc
	s.scoreYahtzeeBonus(*hand)
	m := *s
	// TODO I would love to get rid of the mapping here, somehow
	switch scoreable.(type) {
	case Ones:
		m[OnesName] = score
	case Twos:
		m[TwosName] = score
	case Threes:
		m[ThreesName] = score
	case Fours:
		m[FoursName] = score
	case Fives:
		m[FivesName] = score
	case Sixes:
		m[SixesName] = score
	case ThreeOfAKind:
		m[ThreeOfAKindName] = score
	case FourOfAKind:
		m[FourOfAKindName] = score
	case FullHouse:
		m[FullHouseName] = score
	case SmallStraight:
		m[SmallStraightName] = score
	case LargeStraight:
		m[LargeStraightName] = score
	case Chance:
		m[ChanceName] = score
	case Yahtzee:
		m[YahtzeeName] = score
	}
	return *score
}

func (s *Scorecard) Subtotal() int {
	m := *s
	return ValOrZero(m[OnesName]) + ValOrZero(m[TwosName]) + ValOrZero(m[ThreesName]) +
		ValOrZero(m[FoursName]) + ValOrZero(m[FivesName]) + ValOrZero(m[SixesName])
}

func ValOrZero(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func (s *Scorecard) Total() int {
	sub := s.Subtotal()
	total := sub
	if sub > 63 {
		total = total + 25
	}
	m := *s
	return total + ValOrZero(m[ThreeOfAKindName]) + ValOrZero(m[FourOfAKindName]) + ValOrZero(m[FullHouseName]) +
		ValOrZero(m[SmallStraightName]) + ValOrZero(m[LargeStraightName]) + ValOrZero(m[ChanceName]) + ValOrZero(m[YahtzeeName]) + ValOrZero((m[YahtzeeBonusName]))
}

func (s *Scorecard) scoreYahtzeeBonus(hand Hand) int {
	m := *s
	if m[YahtzeeName] == nil || *m[YahtzeeName] == 0 {
		return 0
	}
	for _, count := range valueCounts(hand) {
		if count == 5 {
			val := 100
			if m[YahtzeeBonusName] != nil {
				val = *m[YahtzeeBonusName] + 100
			}
			m[YahtzeeBonusName] = &val
		}
	}
	return ValOrZero(m[YahtzeeBonusName])
}

func (s *Scorecard) Print() string {
	return s.PrintWithDecorator(func(ScorableName) string { return "" })
}

func (s *Scorecard) PrintWithDecorator(decFn func(ScorableName) string) string {
	str := "-------------------------------------\n"
	str += "| name                         score|\n"
	m := (*s)
	for _, name := range ScorableNames {
		valPtr := m[name]
		val := "-"
		if valPtr != nil {
			val = strconv.Itoa(*valPtr)
		}
		str += fmt.Sprintf("| %-14s                 %3s|", name, val)
		str += decFn(name) + "\n"
	}
	str += fmt.Sprintf("| %-14s                 %3d|\n", "Total", s.Total())
	str += "-------------------------------------\n"
	return str
}
