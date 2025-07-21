package balatro

import (
	"sort"
)

type HandType int

const (
	HighCard HandType = iota
	Pair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

func (ht HandType) String() string {
	switch ht {
	case HighCard:
		return "High Card"
	case Pair:
		return "Pair"
	case TwoPair:
		return "Two Pair"
	case ThreeOfAKind:
		return "Three of a Kind"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case FullHouse:
		return "Full House"
	case FourOfAKind:
		return "Four of a Kind"
	case StraightFlush:
		return "Straight Flush"
	case RoyalFlush:
		return "Royal Flush"
	default:
		return "Unknown"
	}
}

// GetBaseMultiplier returns the base scoring multiplier for each hand type
func (ht HandType) GetBaseMultiplier() int {
	switch ht {
	case HighCard:
		return 1
	case Pair:
		return 2
	case TwoPair:
		return 3
	case ThreeOfAKind:
		return 4
	case Straight:
		return 5
	case Flush:
		return 6
	case FullHouse:
		return 8
	case FourOfAKind:
		return 10
	case StraightFlush:
		return 15
	case RoyalFlush:
		return 25
	default:
		return 1
	}
}

type HandEvaluation struct {
	Type       HandType
	Multiplier int
	CardValue  int
	TotalScore int
}

func EvaluateHand(hand Hand) HandEvaluation {
	if len(hand) == 0 {
		return HandEvaluation{
			Type:       HighCard,
			Multiplier: 1,
			CardValue:  0,
			TotalScore: 0,
		}
	}

	handType := determineHandType(hand)
	multiplier := handType.GetBaseMultiplier()
	cardValue := hand.GetTotalValue()
	totalScore := cardValue * multiplier

	return HandEvaluation{
		Type:       handType,
		Multiplier: multiplier,
		CardValue:  cardValue,
		TotalScore: totalScore,
	}
}

func determineHandType(hand Hand) HandType {
	if len(hand) < 2 {
		return HighCard
	}

	// Sort hand for easier evaluation
	sortedHand := make(Hand, len(hand))
	copy(sortedHand, hand)
	sortedHand.Sort()

	isFlush := checkFlush(sortedHand)
	isStraight := checkStraight(sortedHand)
	
	if isFlush && isStraight {
		if isRoyalFlush(sortedHand) {
			return RoyalFlush
		}
		return StraightFlush
	}

	rankCounts := getRankCounts(sortedHand)
	counts := make([]int, 0)
	for _, count := range rankCounts {
		counts = append(counts, count)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))

	if len(counts) > 0 && counts[0] == 4 {
		return FourOfAKind
	}
	if len(counts) >= 2 && counts[0] == 3 && counts[1] == 2 {
		return FullHouse
	}
	if isFlush {
		return Flush
	}
	if isStraight {
		return Straight
	}
	if len(counts) > 0 && counts[0] == 3 {
		return ThreeOfAKind
	}
	if len(counts) >= 2 && counts[0] == 2 && counts[1] == 2 {
		return TwoPair
	}
	if len(counts) > 0 && counts[0] == 2 {
		return Pair
	}

	return HighCard
}

func checkFlush(hand Hand) bool {
	if len(hand) < 5 {
		return false
	}
	suit := hand[0].Suit
	for _, card := range hand {
		if card.Suit != suit {
			return false
		}
	}
	return true
}

func checkStraight(hand Hand) bool {
	if len(hand) < 5 {
		return false
	}
	
	// Check for regular straight
	for i := 1; i < len(hand); i++ {
		if int(hand[i].Rank) != int(hand[i-1].Rank)+1 {
			// Check for ace-low straight (A-2-3-4-5)
			if i == len(hand)-1 && hand[0].Rank == Two && hand[len(hand)-1].Rank == Ace {
				continue
			}
			return false
		}
	}
	return true
}

func isRoyalFlush(hand Hand) bool {
	if len(hand) != 5 {
		return false
	}
	return hand[0].Rank == Ten && hand[4].Rank == Ace
}

func getRankCounts(hand Hand) map[Rank]int {
	counts := make(map[Rank]int)
	for _, card := range hand {
		counts[card.Rank]++
	}
	return counts
}