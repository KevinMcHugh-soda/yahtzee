package yahtzee

type ScorableName string

const (
	OnesName   = "ones"
	TwosName   = "twos"
	ThreesName = "threes"
	FoursName  = "fours"
	FivesName  = "fives"
	SixesName  = "sixes"

	ThreeOfAKindName  = "3ofakind"
	FourOfAKindName   = "4ofakind"
	FullHouseName     = "fullhouse"
	SmallStraightName = "small straight"
	LargeStraightName = "Large Straight"
	ChanceName        = "Chance"
	YahtzeeName       = "Yahtzee"
	YahtzeeBonusName  = "YahtzeeBonus"

	ErrorName = "error"
)

var ScorableNames = []ScorableName{
	OnesName, TwosName, ThreesName, FoursName, FivesName, SixesName,
	ThreeOfAKindName, FourOfAKindName, FullHouseName, SmallStraightName, LargeStraightName, ChanceName, YahtzeeName,
	YahtzeeBonusName,
	ErrorName,
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
	}

	return nameToPtr[name]
}

func (s *Scorecard) Score(hand *Hand, scoreable Scoreable) int {
	sc := scoreable.Score(*hand)
	score := &sc
	s.scoreYahtzeeBonus(*hand)
	m := *s
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
			val := *m[YahtzeeBonusName] + 100
			m[YahtzeeBonusName] = &val
		}
	}
	return *m[YahtzeeBonusName]
}

func (s *Scorecard) Print() {

}
