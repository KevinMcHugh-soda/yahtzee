package yahtzee

import (
	"fmt"
	"strings"
)

type Color string

const (
	Yellow = "🟨"
	Blue   = "🟦"
	Green  = "🟩"
	Red    = "🟥"
	Orange = "🟧"
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

// Not actually a valid/solvable puzzle
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

func MakeRealPuzzle() Puzzle {
	rows := []string{fiveXfive1, fiveXfive2, fiveXfive3, fiveXfive4, fiveXfive5}
	return ParsePuzzle(rows)
}

var fiveXfive1 = "🟨🟨🟩🟩🟩"
var fiveXfive2 = "🟦🟨🟩🟩🟩"
var fiveXfive3 = "🟦🟥🟧🟧🟩"
var fiveXfive4 = "🟥🟥🟧🟧🟩"
var fiveXfive5 = "🟥🟥🟥🟥🟥"

func ParsePuzzle(rowStrs []string) Puzzle {
	rows := make([][]Cell, len(rowStrs))
	for idx, rowStr := range rowStrs {
		rows[idx] = make([]Cell, len(strings.Split(rowStr, "")))
	}
	for idx, rowStr := range rowStrs {
		for jdx, str := range strings.Split(rowStr, "") {
			rows[idx][jdx] = Cell{
				Segment: Segment{Color: Color(str)},
				State:   Empty, // this shouldn't need to be specified
				X:       idx,
				Y:       jdx,
			}
		}
	}
	// TODO: I shouldn't have to set these explicitly, but
	p := Puzzle{Cells: rows, CorrectStarsPerArea: 1, Width: 5, Height: 5}
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

func (p *Puzzle) Print(msg string) {
	fmt.Println(msg)

	// fmt.Println(" 0️⃣||1️⃣||2️⃣||3️⃣4️⃣")
	for idx, row := range p.Cells {
		str := fmt.Sprintf("%d|", idx)
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

type Coordinate struct {
	x int
	y int
}

func coord(x int, y int) Coordinate {
	return Coordinate{x: x, y: y}
}

func (p *Puzzle) Star(x int, y int) (*Puzzle, error) {
	p.Print(fmt.Sprintf("state before placing star at (%d,%d)", x, y))
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

	p.Cells[y][x].State = Starred

	// Eliminate nearby cells
	// First construct a list of all possible potentialCoordinates
	potentialCoordinates := []Coordinate{
		coord(x-1, y-1),
		coord(x-1, y),
		coord(x-1, y+1),
		coord(x, y-1),
		coord(x, y+1),
		coord(x+1, y-1),
		coord(x+1, y),
		coord(x+1, y+1),
	}
	finalCoordinates := make([]Coordinate, 0)
	// identify any coordinates that are illegal
	for _, coordinate := range potentialCoordinates {
		legalX := true
		legalY := true

		if coordinate.x < 0 || coordinate.x >= p.Width {
			legalX = false
		}

		if coordinate.y < 0 || coordinate.y >= p.Height {
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
			finalCoordinates = append(finalCoordinates, coord(other.X, other.Y))
		}
	}

	// Was that the last star in this row?
	if starsInRow+1 == p.CorrectStarsPerArea {
		// Eliminate everything else in this row
		for _, other := range p.Rows()[y] {
			finalCoordinates = append(finalCoordinates, coord(other.X, other.Y))
		}
	}

	// Was that the last star in this segment?
	if starsInColumn+1 == p.CorrectStarsPerArea {
		// Eliminate everything else in this column
		for _, other := range p.Columns()[x] {
			finalCoordinates = append(finalCoordinates, coord(other.X, other.Y))
		}
	}

	// fmt.Printf("eliminating (%s)\n", final)

	eliminationStr := "Eliminating: "
	// then we'll remove any illegal coordinates
	for _, coordinate := range finalCoordinates {
		// same
		if cell.X == coordinate.x && cell.Y == coordinate.y {
			continue
		}
		other := p.Cells[coordinate.y][coordinate.x]
		if other.State == Starred {
			return nil, fmt.Errorf("attempting to eliminate a starred cell! %d, %d", other.X, other.Y)
		}
		p.Cells[coordinate.y][coordinate.x].State = Eliminated

		eliminationStr = fmt.Sprintf("%s;(%d,%d)", eliminationStr, coordinate.y, coordinate.x)
	}

	// Did this make the puzzle unsolvable?
	// Check all segments to see if they're unsolvable:
	for color, segment := range p.Segments() {
		emptyCells, starredCells := 0, 0
		for _, cell := range segment {
			if cell.State == Empty {
				emptyCells += 1
			} else if cell.State == Starred {
				starredCells += 1
			}
		}

		if emptyCells < p.CorrectStarsPerArea-starredCells {
			return nil, fmt.Errorf("not enough cells left to solve %s", color)
		}
	}

	for idx, segment := range p.Rows() {
		emptyCells, starredCells := 0, 0
		for _, cell := range segment {
			if cell.State == Empty {
				emptyCells += 1
			} else if cell.State == Starred {
				starredCells += 1
			}
		}

		if emptyCells < p.CorrectStarsPerArea-starredCells {
			return nil, fmt.Errorf("not enough cells left to solve row %d", idx)
		}
	}

	for idx, segment := range p.Columns() {
		emptyCells, starredCells := 0, 0
		for _, cell := range segment {
			if cell.State == Empty {
				emptyCells += 1
			} else if cell.State == Starred {
				starredCells += 1
			}
		}

		if emptyCells < p.CorrectStarsPerArea-starredCells {
			return nil, fmt.Errorf("not enough cells left to solve column %d", idx)
		}
	}

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
OUTER:
	for !puzzle.Solved() {
		globalSolveCounter += 1
		fmt.Println("calls to Solve:", globalSolveCounter)

		if globalSolveCounter > 5 {
			fmt.Println("clearly busted, bailing")
			break OUTER
		}

		puzzle.Print("going to find a spot")

		for _, row := range puzzle.Cells {
			for _, cell := range row {
				if cell.State != Empty {
					fmt.Printf("skipping cell (%d,%d) which already has a state of %s\n", cell.X, cell.Y, cell.State)
					continue
				}

				// 🔁 Deep copy before attempting a star
				clone := puzzle.DeepCopy()
				newPuzzle, err := clone.Star(cell.X, cell.Y)
				if err != nil {
					// ❌ Don't touch the original puzzle
					puzzle.Print(fmt.Sprintf("uh oh, erorr! we'll have to backtrack. state with star placed: %s\n", err))
					continue
				}

				// ✅ Either solved, or keep searching
				if newPuzzle.Solved() {
					return *newPuzzle, true
				}

				result, solved := Solve(*newPuzzle)
				if solved {
					return result, true
				}
			}
		}
	}

	return puzzle, puzzle.Solved()
}

//	func clone(original [][]Cell) [][]Cell {
//		cloned := make([][]Cell, len(original))
//		for i := range original {
//			cloned[i] = make([]Cell, len(original[i]))
//			copy(cloned[i], original[i])
//		}
//		return cloned
//	}
func (p *Puzzle) DeepCopy() *Puzzle {
	newCells := make([][]Cell, len(p.Cells))
	for y := range p.Cells {
		newCells[y] = make([]Cell, len(p.Cells[y]))
		for x, cell := range p.Cells[y] {
			newCells[y][x] = Cell{
				X:       cell.X,
				Y:       cell.Y,
				State:   cell.State,   // this is what really matters
				Segment: cell.Segment, // safe to share if immutable
			}
		}
	}

	return &Puzzle{
		Cells:               newCells,
		Width:               p.Width,
		Height:              p.Height,
		CorrectStarsPerArea: p.CorrectStarsPerArea,
	}
}
