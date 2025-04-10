package yahtzee

import (
	"fmt"
)

type Color string

const (
	Yellow = "🟨"
	Blue   = "🟦"
	Green  = "🟩"
)

type Segment struct {
	Color Color
	Cells []Cell
}

type State string

const (
	Starred    = "⭐️"
	Empty      = ""
	Eliminated = "❌"
)

type Cell struct {
	Segment Segment
	State   State
	X       int
	Y       int
}

type Puzzle struct {
	Cells               [][]Cell
	Width               int
	Height              int
	CorrectStarsPerArea int
}

func MakeTrivialPuzzle() Puzzle {
	row := []Cell{
		{
			Segment: Segment{
				Color: Yellow,
			},
			State: Empty,
			X:     0,
			Y:     0,
		},
	}
	grid := [][]Cell{row}
	// TODO: I shouldn't have to set these explicitly, but
	p := Puzzle{Cells: grid, CorrectStarsPerArea: 1, Width: 1, Height: 1}
	return p
}

func MakeEasyPuzzle() Puzzle {
	row1 := []Cell{
		{
			Segment: Segment{
				Color: Yellow,
			},
			State: Empty,
			X:     0,
			Y:     0,
		},
		{
			Segment: Segment{
				Color: Blue,
			},
			State: Empty,
			X:     1,
			Y:     0,
		},
		{
			Segment: Segment{
				Color: Blue,
			},
			State: Empty,
			X:     2,
			Y:     0,
		},
	}
	row2 := []Cell{
		{
			Segment: Segment{
				Color: Blue,
			},
			State: Empty,
			X:     0,
			Y:     1,
		},
		{
			Segment: Segment{
				Color: Blue,
			},
			State: Empty,
			X:     1,
			Y:     1,
		},
		{
			Segment: Segment{
				Color: Blue,
			},
			State: Empty,
			X:     2,
			Y:     1,
		},
	}
	row3 := []Cell{
		{
			Segment: Segment{
				Color: Green,
			},
			State: Empty,
			X:     0,
			Y:     2,
		},
		{
			Segment: Segment{
				Color: Green,
			},
			State: Empty,
			X:     1,
			Y:     2,
		},
		{
			Segment: Segment{
				Color: Green,
			},
			State: Empty,
			X:     1,
			Y:     2,
		},
	}
	grid := [][]Cell{row1, row2, row3}
	// TODO: I shouldn't have to set these explicitly, but
	p := Puzzle{Cells: grid, CorrectStarsPerArea: 1, Width: 3, Height: 3}
	return p
}

func (p *Puzzle) Rows() [][]Cell {
	return p.Cells
}

// TODO memoize these 2, maybe
func (p *Puzzle) Columns() [][]Cell {
	ret := make([][]Cell, p.Height)

	for idx := 0; idx < p.Height; idx += 1 {
		ret[idx] = make([]Cell, p.Width)
	}

	// fmt.Println(p.Cells)
	for idx, row := range p.Cells {
		for jdx, cell := range row {
			// fmt.Println("jdx, idx, row-len, ret-len", jdx, idx, len(row), len(ret))
			ret[jdx][idx] = cell
		}
	}

	return ret
}

func (p *Puzzle) Segments() map[Color][]Cell {
	m := make(map[Color][]Cell)
	for _, row := range p.Cells {
		for _, cell := range row {
			m[cell.Segment.Color] = append(m[cell.Segment.Color], cell)
		}
	}

	return m
}

func (p *Puzzle) Print() {
	for _, row := range p.Cells {
		str := ""
		for _, c := range row {
			if c.State == Empty {
				str += string(c.Segment.Color)
			} else {
				str += string(c.State)
			}
		}
		fmt.Println(str)
	}
}

func coord(x int, y int) map[string]int {
	return map[string]int{"x": x, "y": y}
}

func (p *Puzzle) Star(x int, y int) (*Puzzle, error) {
	fmt.Println("state before placing star")
	p.Print()
	fmt.Printf("attempting to place a star at (%d,%d)\n", x, y)
	cell := p.Cells[y][x]
	if cell.State != Empty {
		return nil, fmt.Errorf("cell already (%d,%d) has state %s", x, y, cell.State)
	}
	starsInSegment := p.StarsPerSegment(cell.Segment.Color)
	if starsInSegment >= p.CorrectStarsPerArea {
		return nil, fmt.Errorf("too many stars in this segment")
	}
	starsInRow := p.StarsPerRow(y)
	if starsInRow >= p.CorrectStarsPerArea {
		return nil, fmt.Errorf("too many stars in this row")
	}
	starsInColumn := p.StarsPerColumn(x)
	if starsInColumn >= p.CorrectStarsPerArea {
		return nil, fmt.Errorf("too many stars in this column")
	}

	cell.State = Starred
	p.Cells[y][x].State = Starred

	// Eliminate nearby cells
	// First construct a list of all possible potentialCoordinates
	potentialCoordinates := []map[string]int{
		coord(x-1, y-1),
		coord(x-1, y),
		coord(x-1, y+1),
		coord(x, y-1),
		coord(x, y+1),
		coord(x+1, y-1),
		coord(x+1, y),
		coord(x+1, y+1),
	}
	finalCoordinates := []map[string]int{}
	// identify any coordinates that are illegal
	for _, coordinate := range potentialCoordinates {
		legalX := true
		legalY := true

		if coordinate["x"] < 0 || coordinate["x"] >= p.Width {
			legalX = false
		}

		if coordinate["y"] < 0 || coordinate["y"] >= p.Height {
			legalY = false
		}

		if legalX && legalY {
			finalCoordinates = append(finalCoordinates, coordinate)
		}
	}

	// Was that the last star in this segment?
	if starsInSegment+1 == p.CorrectStarsPerArea {
		// Eliminate everything else in this segment
		for _, other := range p.Segments()[cell.Segment.Color] {
			if cell.X == other.X && cell.Y == other.Y {
				continue
			}
			if other.State == Starred {
				return nil, fmt.Errorf("attempting to eliminate a starred cell! %d, %d", other.X, other.Y)
			}
			finalCoordinates = append(finalCoordinates, coord(other.X, other.Y))
		}
	}

	// Was that the last star in this row?
	if starsInRow+1 == p.CorrectStarsPerArea {
		// Eliminate everything else in this row
		for _, other := range p.Rows()[x] {
			if cell.X == other.X && cell.Y == other.Y {
				continue
			}
			if other.State == Starred {
				return nil, fmt.Errorf("attempting to eliminate a starred cell! %d, %d", other.X, other.Y)
			}
			finalCoordinates = append(finalCoordinates, coord(other.X, other.Y))
		}
	}

	// Was that the last star in this segment?
	if starsInColumn+1 == p.CorrectStarsPerArea {
		// Eliminate everything else in this column
		for _, other := range p.Columns()[y] {
			finalCoordinates = append(finalCoordinates, coord(other.X, other.Y))
		}
	}

	// then we'll eliminate any illegal coordinates
	for _, coordinate := range finalCoordinates {
		if cell.X == coordinate["x"] && cell.Y == coordinate["y"] {
			continue
		}
		fmt.Printf("eliminating (%d,%d)\n", coordinate["x"], coordinate["y"])
		other := p.Cells[coordinate["y"]][coordinate["x"]]
		if other.State == Starred {
			return nil, fmt.Errorf("attempting to eliminate a starred cell! %d, %d", other.X, other.Y)
		}
		p.Cells[coordinate["y"]][coordinate["x"]].State = Eliminated
		fmt.Printf("state of cell (%d, %d) after elimination: %s\n", other.X, other.Y, p.Cells[other.X][other.Y].State)
	}
	fmt.Println("printing after starring and eliminating")
	// The cells here don't feature any eliminated spaces, so what the hell
	p.Print()
	return p, nil
}

func (p *Puzzle) StarsPerSegment(segment Color) int {
	stars := 0
	for _, cell := range p.Segments()[segment] {
		if cell.State == Starred {
			stars += 1
		}
	}

	return stars
}

func (p *Puzzle) StarsPerRow(row int) int {
	stars := 0
	for _, cell := range p.Rows()[row] {
		if cell.State == Starred {
			stars += 1
		}
	}

	return stars
}

func (p *Puzzle) StarsPerColumn(column int) int {
	stars := 0
	for _, cell := range p.Columns()[column] {
		if cell.State == Starred {
			stars += 1
		}
	}

	return stars
}

func (p *Puzzle) Solved() bool {
	for color := range p.Segments() {
		stars := p.StarsPerSegment(color)
		if stars != p.CorrectStarsPerArea {
			fmt.Printf("not enough stars in %s. found %d, need %d\n", color, stars, p.CorrectStarsPerArea)
			return false
		}
	}

	for idx := range p.Rows() {
		stars := p.StarsPerRow(idx)
		if stars != p.CorrectStarsPerArea {
			return false
		}
	}

	for idx := range p.Columns() {
		stars := p.StarsPerColumn(idx)
		if stars != p.CorrectStarsPerArea {
			return false
		}
	}

	return true
}

var globalSolveCounter int

// We're building a solver for the "star battle" puzzle.
// It's going to use constraint propagation.
func Solve(puzzle Puzzle) (Puzzle, bool) {
	// iterate through the cells
	// identify cells that cannot possibly be checked
	puzzle.Print()

OUTER:
	for !puzzle.Solved() {
		globalSolveCounter += 1
		fmt.Println("calls to Solve:", globalSolveCounter)

		if globalSolveCounter > 10 {
			fmt.Println("clearly busted, bailing")
			break OUTER
		}
		puzzle.Print()
		for _, row := range puzzle.Cells {
			for _, cell := range row {
				if cell.State != Empty {
					fmt.Printf("skipping cell (%d,%d) which already has a state of %s\n", cell.X, cell.Y, cell.State)
					continue
				}
				puzzle.Print()
				p, err := puzzle.Star(cell.X, cell.Y)
				if err != nil {
					fmt.Println("uh oh, error: ", err)
					puzzle.Print()
					continue
				} else if puzzle.Solved() {
					break OUTER
				} else {
					// recurse
					_, solved := Solve(*p)
					if solved {
						return *p, true
					}
				}
			}
		}
	}

	return puzzle, true
}
