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
	Black  = "⬛️"
	Purple = "🟪"
	White  = "⬜️"
	Brown  = "🟫"
	Beer   = "🍺"
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

func (c *Cell) Coords() string {
	return fmt.Sprintf("%s%d", c.Column, c.Row)
}

type Puzzle struct {
	Cells               map[string][]Cell
	Width               int
	Height              int
	CorrectStarsPerArea int
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

var letters = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
var letterIndex = map[string]int{"A": 0, "B": 1, "C": 2, "D": 3, "E": 4, "F": 5, "G": 6, "H": 7, "I": 8, "J": 9}

func ParsePuzzle(rowStrs []string) Puzzle {
	cols := make(map[string][]Cell)
	for idx, letter := range letters {
		rowStr := rowStrs[idx]
		cols[letter] = make([]Cell, len(strings.Split(rowStr, "")))
	}
	for rowIdx, rowStr := range rowStrs {
		split := strings.Split(rowStr, "")
		for jdx, letter := range letters {
			cols[letter][rowIdx] = Cell{
				Segment: Segment{Color: Color(split[jdx])},
				State:   Empty, // this shouldn't need to be specified
				Row:     rowIdx,
				Column:  letter,
			}
		}
	}
	// TODO: I shouldn't have to set these explicitly, but
	p := Puzzle{Cells: cols, CorrectStarsPerArea: 1, Width: 5, Height: 5}
	return p
}

func (p *Puzzle) Rows() [][]Cell {
	ret := make([][]Cell, p.Height)

	for idx := range ret {
		ret[idx] = make([]Cell, p.Width)
	}

	for cdx, letter := range letters {
		column := p.Cells[letter]
		for rdx, cell := range column {
			ret[rdx][cdx] = cell
		}
	}
	return ret
}

// TODO memoize these 2, maybe
func (p *Puzzle) Columns() map[string][]Cell {
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

	fmt.Println(" |A B C D E F G H I J")
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

func (p *Puzzle) coordInt(row int, col int) Coordinate {
	r := -8
	c := "Q"
	if row < 0 || row >= p.Width {
		r = -9
	} else {
		r = row
	}

	if col < 0 || col >= p.Height {
		c = "Z"
	} else {
		c = letters[col]
	}

	return Coordinate{row: r, col: c}
}

func (p *Puzzle) Star(row int, column string) (*Puzzle, error) {
	cell := p.Cells[column][row]
	// p.Print(fmt.Sprintf("state before placing star at (%s)(%s)", cell.Coords(), cell.Segment.Color))
	if cell.State != Empty {
		return p, fmt.Errorf("cell already (%s) has state %s", cell.Coords(), cell.State)
	}
	starsInSegment := p.StarsPerSegment(cell.Segment.Color)
	if starsInSegment >= p.CorrectStarsPerArea {
		return p, fmt.Errorf("too many stars in this segment")
	}
	starsInRow := p.StarsPerRow(row)
	if starsInRow >= p.CorrectStarsPerArea {
		return p, fmt.Errorf("too many stars in this row")
	}
	starsInColumn := p.StarsPerColumn(column)
	if starsInColumn >= p.CorrectStarsPerArea {
		return p, fmt.Errorf("too many stars in this column")
	}

	p.Cells[column][row].State = Starred

	// Eliminate nearby cells
	// First construct a list of all possible potentialCoordinates
	colIndex := letterIndex[column]
	finalCoordinates := []Coordinate{
		p.coordInt(row-1, colIndex-1),
		p.coordInt(row-1, colIndex),
		p.coordInt(row-1, colIndex+1),
		p.coordInt(row, colIndex-1),
		p.coordInt(row, colIndex+1),
		p.coordInt(row+1, colIndex-1),
		p.coordInt(row+1, colIndex),
		p.coordInt(row+1, colIndex+1),
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
		for _, col := range letters {
			finalCoordinates = append(finalCoordinates, coord(row, col))
		}
	}

	// Was that the last star in this column?
	if starsInColumn+1 == p.CorrectStarsPerArea {
		// Eliminate everything else in this column
		for idx := 0; idx < p.Height; idx++ {
			finalCoordinates = append(finalCoordinates, coord(idx, column))
		}
	}

	// then we'll remove any illegal coordinates
	for _, coordinate := range finalCoordinates {
		if coordinate.col == "Z" || coordinate.row < 0 {
			continue
		}
		other := p.Cells[coordinate.col][coordinate.row]

		if cell.Row == coordinate.row && cell.Column == coordinate.col || other.State == Eliminated {
			continue
		}
		if other.State == Starred {
			return p, fmt.Errorf("attempting to eliminate a starred cell! %s", other.Coords())
		}
		p.Cells[coordinate.col][coordinate.row].State = Eliminated
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
			p.Print("state with errror")
			return p, fmt.Errorf("not enough cells left to solve %s", color)
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
			return p, fmt.Errorf("not enough cells left to solve row %d", idx)
		}
	}

	for _, letter := range letters {
		segment := p.Columns()[letter]
		emptyCells, starredCells := 0, 0
		for _, cell := range segment {
			if cell.State == Empty {
				emptyCells += 1
			} else if cell.State == Starred {
				starredCells += 1
			}
		}

		if emptyCells < p.CorrectStarsPerArea-starredCells {
			return p, fmt.Errorf("not enough cells left to solve column %s", letter)
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
			fmt.Printf("not enough stars in row %d. found %d, need %d\n", idx, stars, p.CorrectStarsPerArea)
			return false
		}
	}

	for letter := range p.Columns() {
		stars := p.StarsPerColumn(letter)
		if stars != p.CorrectStarsPerArea {
			fmt.Printf("not enough stars in column %s. found %d, need %d\n", letter, stars, p.CorrectStarsPerArea)
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

		for _, letter := range letters {
			column := puzzle.Columns()[letter]
			for _, cell := range column {
				if cell.State != Empty {
					continue
				}

				// Deep copy before attempting a star
				clone := puzzle.DeepCopy()
				newPuzzle, err := clone.Star(cell.Row, cell.Column)
				if err != nil {
					puzzle.Print(fmt.Sprintf("uh oh, erorr! we'll have to backtrack. state with star placed: %s\n", err))
					// since we couldn't star this cell, we can eliminate it.
					puzzle.Cells[cell.Column][cell.Row].State = Eliminated
					continue
				}

				// Either solved, or keep searching
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
