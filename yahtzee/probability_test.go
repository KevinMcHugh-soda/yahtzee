package yahtzee_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"kevinmchugh.me/yahtzee/m/v2/yahtzee"
)

func TestFives_ProbabilityToHit(t *testing.T) {
	sut := yahtzee.Fives{}

	oneSixth := 1.0 / 6.0
	testCases := []struct {
		name        string
		hand        yahtzee.Hand
		probability float64
	}{
		{
			name:        "all fives",
			hand:        [5]int{5, 5, 5, 5, 5},
			probability: 1.0,
		}, {
			name:        "no fives",
			hand:        [5]int{1, 1, 1, 1, 1},
			probability: 1 - (oneSixth * 5),
		}, {
			name:        "four fives",
			hand:        [5]int{5, 5, 5, 5, 1},
			probability: 1 - (oneSixth * 1),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.InDelta(t, testCase.probability, sut.ProbabilityToHit(testCase.hand, 1), 0.001)
		})
	}
}

func TestLargeStraight_ProbabilityToHit(t *testing.T) {
	sut := yahtzee.LargeStraight{}

	oneSixth := 1.0 / 6.0
	testCases := []struct {
		name        string
		hand        yahtzee.Hand
		probability float64
	}{
		{
			name:        "completed low straight",
			hand:        [5]int{1, 2, 3, 4, 5},
			probability: 1.0,
		}, {
			name:        "completed high straight",
			hand:        [5]int{2, 3, 4, 5, 6},
			probability: 1.0,
		}, {
			name:        "one roll missing for low",
			hand:        [5]int{1, 2, 3, 4, 1},
			probability: 1 - oneSixth,
		}, {
			name:        "one roll missing for high",
			hand:        [5]int{2, 3, 4, 5, 2},
			probability: 1 - oneSixth,
		}, {
			name:        "four rolls missing for low",
			hand:        [5]int{1, 1, 1, 1, 1},
			probability: 1 - (oneSixth * 4),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.InDelta(t, testCase.probability, sut.ProbabilityToHit(testCase.hand, 1), 0.001)
		})
	}
}
