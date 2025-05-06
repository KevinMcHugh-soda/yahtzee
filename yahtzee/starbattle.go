package yahtzee

import (
	"fmt"
	"sort"
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

type State string

const (
	Starred    = "⭐️"
	Empty      = ""
	Eliminated = "❌"
)

type Cell struct {
	Segment Color
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

func MakeEasyPuzzle() Puzzle {
	rows := []string{fiveXfive1, fiveXfive2, fiveXfive3, fiveXfive4, fiveXfive5}
	p, _ := ParsePuzzle(rows, 1)

	return *p
}

var fiveXfive1 = "🟨🟨🟩🟩🟩"
var fiveXfive2 = "🟦🟨🟩🟩🟩"
var fiveXfive3 = "🟦🟥🟧🟧🟩"
var fiveXfive4 = "🟥🟥🟧🟧🟩"
var fiveXfive5 = "🟥🟥🟥🟥🟥"

var letters = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
var letterIndex = map[string]int{"A": 0, "B": 1, "C": 2, "D": 3, "E": 4, "F": 5, "G": 6, "H": 7, "I": 8, "J": 9}

func ParsePuzzle(rowStrs []string, starsPerArea int) (*Puzzle, error) {
	if len(rowStrs) == 0 {
		return nil, fmt.Errorf("empty puzzle")
	}

	cols := make(map[string][]Cell)
	split := strings.Split(rowStrs[0], "")
	width := len(split)

	for idx := 0; idx < width; idx += 1 {
		letter := letters[idx]
		cols[letter] = make([]Cell, width)
	}

	for rowIdx, rowStr := range rowStrs {
		split := strings.Split(rowStr, "")
		if len(split) != width {
			return nil, fmt.Errorf("lines are not equal length! First line was %d, line #%d is %d", width, rowIdx, len(split))
		}
		for jdx := 0; jdx < width; jdx += 1 {
			letter := letters[jdx]

			cols[letter][rowIdx] = Cell{
				Segment: Color(split[jdx]),
				State:   Empty, // this shouldn't need to be specified
				Row:     rowIdx,
				Column:  letter,
			}
		}
	}
	p := Puzzle{Cells: cols, CorrectStarsPerArea: starsPerArea, Width: width, Height: len(rowStrs)}
	return &p, nil
}

func (p *Puzzle) Rows() [][]Cell {
	ret := make([][]Cell, p.Height)

	for idx := range ret {
		ret[idx] = make([]Cell, p.Width)
	}

	for cdx, letter := range p.ColumnNames() {
		column := p.Cells[letter]
		for rdx, cell := range column {
			ret[rdx][cdx] = cell
		}
	}
	return ret
}

func (p *Puzzle) Columns() map[string][]Cell {
	return p.Cells
}

func (p *Puzzle) Segments() map[Color][]Cell {
	m := make(map[Color][]Cell)
	for _, row := range p.Cells {
		for _, cell := range row {
			m[cell.Segment] = append(m[cell.Segment], cell)
		}
	}

	return m
}
func (p *Puzzle) ColumnNames() []string {
	return letters[0:p.Width]
}

func (p *Puzzle) Print(msg string) {
	fmt.Println(msg)

	fmt.Printf(" |%s\n", strings.Join(p.ColumnNames(), " "))
	for idx, row := range p.Rows() {
		str := fmt.Sprintf("%d|", idx)
		for _, c := range row {
			if c.State == Empty {
				str += string(c.Segment)
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
	// Validate column exists
	if _, ok := p.Cells[column]; !ok {
		return p, fmt.Errorf("invalid column: %s", column)
	}

	// Validate row is in bounds
	if row < 0 || row >= p.Height {
		return p, fmt.Errorf("invalid row: %d", row)
	}

	cell := p.Cells[column][row]
	if cell.State != Empty {
		return p, fmt.Errorf("cell already (%s) has state %s", cell.Coords(), cell.State)
	}
	starsInSegment := p.StarsPerSegment(cell.Segment)
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

	// Create a copy to test the move
	testPuzzle := p.DeepCopy()
	testPuzzle.Cells[column][row].State = Starred

	// Eliminate nearby cells
	colIndex := letterIndex[column]
	finalCoordinates := []Coordinate{
		testPuzzle.coordInt(row-1, colIndex-1),
		testPuzzle.coordInt(row-1, colIndex),
		testPuzzle.coordInt(row-1, colIndex+1),
		testPuzzle.coordInt(row, colIndex-1),
		testPuzzle.coordInt(row, colIndex+1),
		testPuzzle.coordInt(row+1, colIndex-1),
		testPuzzle.coordInt(row+1, colIndex),
		testPuzzle.coordInt(row+1, colIndex+1),
	}

	// Was that the last star in this segment?
	if starsInSegment+1 == testPuzzle.CorrectStarsPerArea {
		// Eliminate everything else in this segment
		for _, other := range testPuzzle.Segments()[cell.Segment] {
			finalCoordinates = append(finalCoordinates, coord(other.Row, other.Column))
		}
	}

	// Was that the last star in this row?
	if starsInRow+1 == testPuzzle.CorrectStarsPerArea {
		// Eliminate everything else in this row
		for _, col := range testPuzzle.ColumnNames() {
			finalCoordinates = append(finalCoordinates, coord(row, col))
		}
	}

	// Was that the last star in this column?
	if starsInColumn+1 == testPuzzle.CorrectStarsPerArea {
		// Eliminate everything else in this column
		for idx := 0; idx < testPuzzle.Height; idx++ {
			finalCoordinates = append(finalCoordinates, coord(idx, column))
		}
	}

	// then we'll remove any illegal coordinates
	for _, coordinate := range finalCoordinates {
		if coordinate.col == "Z" || coordinate.row < 0 {
			continue
		}
		if letterIndex[coordinate.col] >= testPuzzle.Width {
			continue
		}
		other := testPuzzle.Cells[coordinate.col][coordinate.row]

		if cell.Row == coordinate.row && cell.Column == coordinate.col || other.State == Eliminated {
			continue
		}
		if other.State == Starred {
			return testPuzzle, fmt.Errorf("attempting to eliminate a starred cell! %s", other.Coords())
		}
		testPuzzle.Cells[coordinate.col][coordinate.row].State = Eliminated
	}

	// If the move is valid, apply it to the original puzzle
	*p = *testPuzzle
	return p.Deduce()
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

func isHorizontalSegment(cells []Cell) bool {
	firstRow := cells[0].Row
	for _, c := range cells {
		if c.Row != firstRow {
			return false
		}
	}
	return true
}

func isVerticalSegment(cells []Cell) bool {
	firstCol := cells[0].Column
	for _, c := range cells {
		if c.Column != firstCol {
			return false
		}
	}
	return true
}

// We're building a solver for the "star battle" puzzle.
// It's going to use constraint propagation.
// func Solve(puzzle Puzzle) (Puzzle, bool) {
// OUTER:
// 	for !puzzle.Solved() {
// 		globalSolveCounter += 1
// 		fmt.Println("calls to Solve:", globalSolveCounter, "cells considered:", totalCellsConsidered)

// 		if globalSolveCounter > 1000 {
// 			fmt.Println("clearly busted, bailing")
// 			break OUTER
// 		}

// 		for _, letter := range puzzle.ColumnNames() {
// 			column := puzzle.Columns()[letter]
// 			for _, cell := range column {
// 				if cell.State != Empty {
// 					continue
// 				}
// 				totalCellsConsidered += 1
// 				puzzle.Print(fmt.Sprintf("state before starring %s", cell.Coords()))

// 				// Deep copy before attempting a star
// 				clone := puzzle.DeepCopy()
// 				newPuzzle, err := clone.Star(cell.Row, cell.Column)
// 				if err != nil {
// 					puzzle.Print(fmt.Sprintf("uh oh, erorr! we'll have to backtrack. state with star placed: %s\n", err))
// 					// since we couldn't star this cell, we can eliminate it.
// 					// puzzle.Cells[cell.Column][cell.Row].State = Eliminated
// 					continue
// 				}

// 				// Either solved, or keep searching
// 				if newPuzzle.Solved() {
// 					return *newPuzzle, true
// 				}

// 				result, solved := Solve(*newPuzzle)
// 				if solved {
// 					return result, true
// 				}
// 			}
// 		}
// 	}

// 	return puzzle, puzzle.Solved()
// }

func (p *Puzzle) IsUnsolvable() bool {
	for _, segment := range p.Segments() {
		empty, stars := 0, 0
		for _, cell := range segment {
			switch cell.State {
			case Empty:
				empty++
			case Starred:
				stars++
			}
		}
		if empty < p.CorrectStarsPerArea-stars {
			return true
		}
	}

	for _, row := range p.Rows() {
		empty, stars := 0, 0
		for _, cell := range row {
			switch cell.State {
			case Empty:
				empty++
			case Starred:
				stars++
			}
		}
		if empty < p.CorrectStarsPerArea-stars {
			return true
		}
	}

	for _, colName := range p.ColumnNames() {
		empty, stars := 0, 0
		for _, cell := range p.Cells[colName] {
			switch cell.State {
			case Empty:
				empty++
			case Starred:
				stars++
			}
		}
		if empty < p.CorrectStarsPerArea-stars {
			return true
		}
	}

	return false
}

func Solve(puzzle Puzzle) (Puzzle, bool) {
	globalSolveCounter++
	fmt.Println("calls to Solve:", globalSolveCounter)
	if globalSolveCounter > 1000 {
		fmt.Println("bailing!")
		return puzzle, false
	}

	if puzzle.Solved() {
		return puzzle, true
	}

	if puzzle.IsUnsolvable() {
		return puzzle, false
	}

	// Step 1: Get segments sorted by size (ascending)
	type segmentInfo struct {
		color Color
		cells []Cell
	}
	var segments []segmentInfo
	for color, cells := range puzzle.Segments() {
		segments = append(segments, segmentInfo{color, cells})
	}
	sort.Slice(segments, func(i, j int) bool {
		return len(segments[i].cells) < len(segments[j].cells)
	})

	// Step 2: Try to solve using the smallest segment first
	for _, seg := range segments {
		for _, cell := range seg.cells {
			if cell.State != Empty {
				continue
			}

			puzzle.Print(fmt.Sprintf("state before starring at %s", cell.Coords()))
			// Try placing a star
			starClone := puzzle.DeepCopy()
			newPuzzle, err := starClone.Star(cell.Row, cell.Column)
			if err == nil {
				if newPuzzle.Solved() {
					return *newPuzzle, true
				}
				result, solved := Solve(*newPuzzle)
				if solved {
					return result, true
				}
			}

			// Try eliminating the cell
			elimClone := puzzle.DeepCopy()
			elimClone.Cells[cell.Column][cell.Row].State = Eliminated

			if elimClone.Solved() {
				return *elimClone, true
			}
			result, solved := Solve(*elimClone)
			if solved {
				return result, true
			}
		}
	}

	// No valid paths
	return puzzle, false
}

func (p *Puzzle) Deduce() (*Puzzle, error) {
	changed := true

	for changed {
		changed = false
		constraints := p.AllConstraints()

		for _, constraint := range constraints {
			if constraint.Apply(p) {
				fmt.Printf("🧠 Deduction applied: %s\n", constraint.String())
				changed = true
			}
		}
	}

	return p, nil
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

type Constraint interface {
	Apply(p *Puzzle) bool
	String() string
}
type SegmentConstraint struct {
	Segment Color
}

func (sc SegmentConstraint) Apply(p *Puzzle) bool {
	empties := []Cell{}
	stars := 0
	for _, cell := range p.Segments()[sc.Segment] {
		live := p.Cells[cell.Column][cell.Row]
		if live.State == Starred {
			stars++
		} else if live.State == Empty {
			empties = append(empties, live)
		}
	}

	needed := p.CorrectStarsPerArea - stars
	if needed == 0 {
		changed := false
		for _, cell := range empties {
			p.Cells[cell.Column][cell.Row].State = Eliminated
			changed = true
		}
		return changed
	} else if needed == len(empties) {
		changed := false
		for _, cell := range empties {
			p.Cells[cell.Column][cell.Row].State = Starred
			changed = true
		}
		return changed
	}

	return false
}

func (sc SegmentConstraint) String() string {
	return fmt.Sprintf("SegmentConstraint(Segment %s)", sc.Segment)
}

type RowConstraint struct {
	Row int
}

func (rc RowConstraint) Apply(p *Puzzle) bool {
	empties := []Cell{}
	stars := 0
	for _, cell := range p.Rows()[rc.Row] {
		live := p.Cells[cell.Column][cell.Row]
		if live.State == Starred {
			stars++
		} else if live.State == Empty {
			empties = append(empties, live)
		}
	}

	needed := p.CorrectStarsPerArea - stars
	if needed == 0 {
		changed := false
		for _, cell := range empties {
			p.Cells[cell.Column][cell.Row].State = Eliminated
			changed = true
		}
		return changed
	} else if needed == len(empties) {
		changed := false
		for _, cell := range empties {
			p.Cells[cell.Column][cell.Row].State = Starred
			changed = true
		}
		return changed
	}

	return false
}

func (rc RowConstraint) String() string {
	return fmt.Sprintf("RowConstraint(Row %d)", rc.Row)
}

type ColumnConstraint struct {
	Column string
}

func (cc ColumnConstraint) Apply(p *Puzzle) bool {
	empties := []Cell{}
	stars := 0
	for _, cell := range p.Columns()[cc.Column] {
		live := p.Cells[cell.Column][cell.Row]
		if live.State == Starred {
			stars++
		} else if live.State == Empty {
			empties = append(empties, live)
		}
	}

	needed := p.CorrectStarsPerArea - stars
	if needed == 0 {
		changed := false
		for _, cell := range empties {
			p.Cells[cell.Column][cell.Row].State = Eliminated
			changed = true
		}
		return changed
	} else if needed == len(empties) {
		changed := false
		for _, cell := range empties {
			p.Cells[cell.Column][cell.Row].State = Starred
			changed = true
		}
		return changed
	}

	return false
}

func (cc ColumnConstraint) String() string {
	return fmt.Sprintf("ColumnConstraint(Column %s)", cc.Column)
}

type RowSegmentConstraint struct {
	Row     int
	Segment Color
}

func (rsc RowSegmentConstraint) Apply(p *Puzzle) bool {
	allRowCells := p.Rows()[rsc.Row]
	var segCellsInRow []Cell
	var starsFromOtherSegments int

	for _, c := range allRowCells {
		cell := p.Cells[c.Column][c.Row]
		if cell.Segment == rsc.Segment {
			if cell.State == Empty {
				segCellsInRow = append(segCellsInRow, cell)
			}
		} else if cell.State == Starred {
			starsFromOtherSegments++
		}
	}

	needed := p.CorrectStarsPerArea - starsFromOtherSegments
	if needed <= 0 {
		// This segment can't contribute any more stars to this row
		changed := false
		for _, cell := range segCellsInRow {
			if p.Cells[cell.Column][cell.Row].State == Empty {
				p.Cells[cell.Column][cell.Row].State = Eliminated
				changed = true
			}
		}
		return changed
	}

	return false
}

func (rsc RowSegmentConstraint) String() string {
	return fmt.Sprintf("RowSegmentConstraint(Row %d, Segment %s)", rsc.Row, rsc.Segment)
}

type ColumnSegmentConstraint struct {
	Column  string
	Segment Color
}

func (csc ColumnSegmentConstraint) Apply(p *Puzzle) bool {
	allColCells := p.Columns()[csc.Column]
	var segCellsInCol []Cell
	var starsFromOtherSegments int

	for _, c := range allColCells {
		cell := p.Cells[c.Column][c.Row]
		if cell.Segment == csc.Segment {
			if cell.State == Empty {
				segCellsInCol = append(segCellsInCol, cell)
			}
		} else if cell.State == Starred {
			starsFromOtherSegments++
		}
	}

	needed := p.CorrectStarsPerArea - starsFromOtherSegments
	if needed <= 0 {
		changed := false
		for _, cell := range segCellsInCol {
			if p.Cells[cell.Column][cell.Row].State == Empty {
				p.Cells[cell.Column][cell.Row].State = Eliminated
				changed = true
			}
		}
		return changed
	}

	return false
}

func (csc ColumnSegmentConstraint) String() string {
	return fmt.Sprintf("ColumnSegmentConstraint(Column %s, Segment %s)", csc.Column, csc.Segment)
}

func (p *Puzzle) AllConstraints() []Constraint {
	var constraints []Constraint

	// Add segment constraints
	for color := range p.Segments() {
		constraints = append(constraints, SegmentConstraint{color})
	}

	// Add row constraints
	for row := 0; row < p.Height; row++ {
		constraints = append(constraints, RowConstraint{row})
	}

	// Add column constraints
	for _, col := range p.ColumnNames() {
		constraints = append(constraints, ColumnConstraint{col})
	}

	// Add one row-segment constraint per row
	for row := 0; row < p.Height; row++ {
		// Just use the first segment in the row
		segment := p.Rows()[row][0].Segment
		constraints = append(constraints, RowSegmentConstraint{Row: row, Segment: segment})
	}

	// Add one column-segment constraint per column
	for _, col := range p.ColumnNames() {
		// Just use the first segment in the column
		segment := p.Cells[col][0].Segment
		constraints = append(constraints, ColumnSegmentConstraint{Column: col, Segment: segment})
	}

	return constraints
}
