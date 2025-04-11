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
	Row     int
	Column  string
}

type Puzzle struct {
	Cells               map[string][]Cell
	Width               int
	Height              int
	CorrectStarsPerArea int
}

func MakeTrivialPuzzle() Puzzle {
	row := map[string][]Cell{
		"A": {{
			Segment: Segment{
				Color: Yellow,
			},
			State:  Empty,
			Row:    0,
			Column: "A",
		}},
	}
	// TODO: I shouldn't have to set these explicitly, but
	p := Puzzle{Cells: row, CorrectStarsPerArea: 1, Width: 1, Height: 1}
	return p
}

// Not actually a valid/solvable puzzle
func MakeEasyPuzzle() Puzzle {
	row1 := []Cell{
		{
			Segment: Segment{
				Color: Yellow,
			},
			State:  Empty,
			Row:    0,
			Column: "A",
		},
		{
			Segment: Segment{
				Color: Blue,
			},
			State:  Empty,
			Row:    1,
			Column: "A",
		},
		{
			Segment: Segment{
				Color: Blue,
			},
			State:  Empty,
			Row:    2,
			Column: "A",
		},
	}
	row2 := []Cell{
		{
			Segment: Segment{
				Color: Blue,
			},
			State:  Empty,
			Row:    0,
			Column: "B",
		},
		{
			Segment: Segment{
				Color: Blue,
			},
			State:  Empty,
			Row:    1,
			Column: "B",
		},
		{
			Segment: Segment{
				Color: Blue,
			},
			State:  Empty,
			Row:    2,
			Column: "B",
		},
	}
	row3 := []Cell{
		{
			Segment: Segment{
				Color: Green,
			},
			State:  Empty,
			Row:    0,
			Column: "C",
		},
		{
			Segment: Segment{
				Color: Green,
			},
			State:  Empty,
			Row:    1,
			Column: "C",
		},
		{
			Segment: Segment{
				Color: Green,
			},
			State:  Empty,
			Row:    1,
			Column: "C",
		},
	}
	grid := map[string][]Cell{"A": row1, "B": row2, "C": row3}
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

var letters = []string{"A", "B", "C", "D", "E"}
var letterIndex = map[string]int{"A": 1, "B": 2, "C": 3, "D": 4, "E": 5}

func ParsePuzzle(rowStrs []string) Puzzle {
	rows := make(map[string][]Cell)
	for idx, letter := range letters {
		rowStr := rowStrs[idx]
		rows[letter] = make([]Cell, len(strings.Split(rowStr, "")))
	}
	for idx, letter := range letters {
		rowStr := rowStrs[idx]
		for jdx, str := range strings.Split(rowStr, "") {
			rows[letter][jdx] = Cell{
				Segment: Segment{Color: Color(str)},
				State:   Empty, // this shouldn't need to be specified
				Row:     idx,
				Column:  letter,
			}
		}
	}
	// TODO: I shouldn't have to set these explicitly, but
	p := Puzzle{Cells: rows, CorrectStarsPerArea: 1, Width: 5, Height: 5}
	return p
}

func (p *Puzzle) Rows() [][]Cell {
	ret := make([][]Cell, p.Height)

	for idx := 0; idx < p.Height; idx += 1 {
		ret[idx] = make([]Cell, p.Width)
	}

	for idx, letter := range letters {
		row := p.Cells[letter]
		for jdx, cell := range row {
			ret[jdx][idx] = cell
		}
	}
	return ret
}

// TODO memoize these 2, maybe
func (p *Puzzle) Columns() map[string][]Cell {
	// ret := make([][]Cell, p.Height)

	// for idx := 0; idx < p.Height; idx += 1 {
	// 	ret[idx] = make([]Cell, p.Width)
	// }

	// // fmt.Println(p.Cells)
	// for idx, letter := range letters {
	// 	row := p.Cells[letter]
	// 	for jdx, cell := range row {
	// 		ret[jdx][idx] = cell
	// 	}
	// }

	return p.Cells
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

	fmt.Println(" |A B C D E")
	for idx, row := range p.Rows() {
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
	row int
	col string
}

func (c *Coordinate) colIndex() int {
	return letterIndex[c.col]
}

func coord(row int, col string) Coordinate {
	return Coordinate{row: row, col: col}
}

func coordInt(x int, y int) Coordinate {
	if y < 0 || y > len(letters) {
		return Coordinate{row: -1}
	}
	return Coordinate{row: x, col: letters[y]}
}

func (p *Puzzle) Star(row int, column string) (*Puzzle, error) {
	cell := p.Cells[column][row]
	p.Print(fmt.Sprintf("state before placing star at (%d,%s)(%s)", row, column, cell.Segment.Color))
	if cell.State != Empty {
		return nil, fmt.Errorf("cell already (%d,%s) has state %s", row, column, cell.State)
	}
	starsInSegment := p.StarsPerSegment(cell.Segment.Color)
	if starsInSegment >= p.CorrectStarsPerArea {
		return nil, fmt.Errorf("too many stars in this segment")
	}
	starsInRow := p.StarsPerRow(row)
	if starsInRow >= p.CorrectStarsPerArea {
		return nil, fmt.Errorf("too many stars in this row")
	}
	starsInColumn := p.StarsPerColumn(column)
	if starsInColumn >= p.CorrectStarsPerArea {
		return nil, fmt.Errorf("too many stars in this column")
	}

	p.Cells[column][row].State = Starred

	// Eliminate nearby cells
	// First construct a list of all possible potentialCoordinates
	colIndex := letterIndex[column]
	potentialCoordinates := []Coordinate{
		coordInt(row-1, colIndex-1),
		coordInt(row-1, colIndex),
		coordInt(row-1, colIndex+1),
		coordInt(row, colIndex-1),
		coordInt(row, colIndex+1),
		coordInt(row+1, colIndex-1),
		coordInt(row+1, colIndex),
		coordInt(row+1, colIndex+1),
	}

	fmt.Println("current coords:", row, column)
	fmt.Println("potential coords: ", potentialCoordinates)
	finalCoordinates := make([]Coordinate, 0)
	// identify any coordinates that are illegal
	for _, coordinate := range potentialCoordinates {
		legalX := true
		legalY := true

		if coordinate.row < 0 || coordinate.row >= p.Width {
			legalX = false
		}

		if coordinate.colIndex() < 0 || coordinate.colIndex() >= p.Height {
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
			finalCoordinates = append(finalCoordinates, coord(other.Row, other.Column))
		}
	}

	// Was that the last star in this row?
	if starsInRow+1 == p.CorrectStarsPerArea {
		// Eliminate everything else in this row
		for _, other := range p.Rows()[row] {
			finalCoordinates = append(finalCoordinates, coord(other.Row, other.Column))
		}
	}

	// Was that the last star in this segment?
	if starsInColumn+1 == p.CorrectStarsPerArea {
		// Eliminate everything else in this column
		for _, other := range p.Columns()[column] {
			finalCoordinates = append(finalCoordinates, coord(other.Row, other.Column))
		}
	}

	// fmt.Printf("eliminating (%s)\n", final)

	eliminationStr := "Eliminating: "
	// then we'll remove any illegal coordinates
	for _, coordinate := range finalCoordinates {
		other := p.Cells[coordinate.col][coordinate.row]

		if cell.Row == coordinate.row && cell.Column == coordinate.col || other.State == Eliminated {
			continue
		}
		if other.State == Starred {
			return nil, fmt.Errorf("attempting to eliminate a starred cell! %d, %s", other.Row, other.Column)
		}
		p.Cells[coordinate.col][coordinate.row].State = Eliminated

		eliminationStr = fmt.Sprintf("%s;(%d,%s)", eliminationStr, coordinate.col, coordinate.row)
	}

	fmt.Println(eliminationStr)

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

func (p *Puzzle) StarsPerColumn(column string) int {
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

		puzzle.Print("state while looking for a cell to attempt eliminating")

		for _, row := range puzzle.Cells {
			for _, cell := range row {
				if cell.State != Empty {
					fmt.Printf("skipping cell (%d,%s) which already has a state of %s\n", cell.Row, cell.Column, cell.State)
					continue
				}

				// 🔁 Deep copy before attempting a star
				clone := puzzle.DeepCopy()
				newPuzzle, err := clone.Star(cell.Row, cell.Column)
				if err != nil {
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

func (p *Puzzle) DeepCopy() *Puzzle {
	newCells := make(map[string][]Cell, len(p.Cells))
	for y := range p.Cells {
		newCells[y] = make([]Cell, len(p.Cells[y]))
		for x, cell := range p.Cells[y] {
			newCells[y][x] = Cell{
				Row:     cell.Row,
				Column:  cell.Column,
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
