package yahtzee

import (
	"fmt"
	"sort"
	"strings"
)

// Color represents a segment color in the puzzle
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

// State represents the state of a cell
type State string

const (
	Starred    = "⭐️"
	Empty      = ""
	Eliminated = "❌"
)

// Cell represents a single cell in the puzzle
type Cell struct {
	Segment Color
	State   State
	Row     int
	Column  string
}

func (c *Cell) Coords() string {
	return fmt.Sprintf("%s%d", c.Column, c.Row)
}

// Puzzle represents a Star Battle puzzle
type Puzzle struct {
	Cells               map[string][]Cell
	Width               int
	Height              int
	CorrectStarsPerArea int
}

// Helper functions for common operations
func (p *Puzzle) countStars(cells []Cell) int {
	stars := 0
	for _, cell := range cells {
		if cell.State == Starred {
			stars++
		}
	}
	return stars
}

func (p *Puzzle) countEmpty(cells []Cell) int {
	empty := 0
	for _, cell := range cells {
		if cell.State == Empty {
			empty++
		}
	}
	return empty
}

// getFreshCells takes a slice of cells and returns a new slice containing the current state
// of those cells from the puzzle's Cells map. This is needed because the input cells may be
// outdated copies, while the puzzle's Cells map contains the current state.
func (p *Puzzle) getFreshCells(cells []Cell) []Cell {
	var fresh []Cell
	for _, cell := range cells {
		fresh = append(fresh, p.Cells[cell.Column][cell.Row])
	}
	return fresh
}

// MakeEasyPuzzle creates a simple 5x5 puzzle
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

// ParsePuzzle creates a new puzzle from a list of row strings
func ParsePuzzle(rowStrs []string, starsPerArea int) (*Puzzle, error) {
	if len(rowStrs) == 0 {
		return nil, fmt.Errorf("empty puzzle")
	}

	cols := make(map[string][]Cell)
	split := strings.Split(rowStrs[0], "")
	width := len(split)

	// Initialize columns
	for idx := 0; idx < width; idx++ {
		letter := letters[idx]
		cols[letter] = make([]Cell, width)
	}

	// Parse rows
	for rowIdx, rowStr := range rowStrs {
		split := strings.Split(rowStr, "")
		if len(split) != width {
			return nil, fmt.Errorf("lines are not equal length! First line was %d, line #%d is %d", width, rowIdx, len(split))
		}
		for jdx := 0; jdx < width; jdx++ {
			letter := letters[jdx]
			cols[letter][rowIdx] = Cell{
				Segment: Color(split[jdx]),
				State:   Empty,
				Row:     rowIdx,
				Column:  letter,
			}
		}
	}
	return &Puzzle{Cells: cols, CorrectStarsPerArea: starsPerArea, Width: width, Height: len(rowStrs)}, nil
}

// Rows returns the puzzle rows as a 2D slice
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

// Columns returns the puzzle columns as a map
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

// Star places a star at the given position and applies constraints
func (p *Puzzle) Star(row int, column string) (*Puzzle, error) {
	// Validate input
	if _, ok := p.Cells[column]; !ok {
		return p, fmt.Errorf("invalid column: %s", column)
	}
	if row < 0 || row >= p.Height {
		return p, fmt.Errorf("invalid row: %d", row)
	}

	cell := p.Cells[column][row]
	if cell.State != Empty {
		return p, fmt.Errorf("cell already (%s) has state %s", cell.Coords(), cell.State)
	}

	// Check constraints
	if p.StarsPerSegment(cell.Segment) >= p.CorrectStarsPerArea {
		return p, fmt.Errorf("too many stars in this segment")
	}
	if p.StarsPerRow(row) >= p.CorrectStarsPerArea {
		return p, fmt.Errorf("too many stars in this row")
	}
	if p.StarsPerColumn(column) >= p.CorrectStarsPerArea {
		return p, fmt.Errorf("too many stars in this column")
	}

	// Create a copy to test the move
	testPuzzle := p.DeepCopy()
	testPuzzle.Cells[column][row].State = Starred

	// Get cells to eliminate
	elimCells := testPuzzle.getCellsToEliminate(row, column, cell)

	// Eliminate cells
	for _, coord := range elimCells {
		if coord.col == "Z" || coord.row < 0 || letterIndex[coord.col] >= testPuzzle.Width {
			continue
		}
		other := testPuzzle.Cells[coord.col][coord.row]
		if cell.Row == coord.row && cell.Column == coord.col || other.State == Eliminated {
			continue
		}
		if other.State == Starred {
			return testPuzzle, fmt.Errorf("attempting to eliminate a starred cell! %s", other.Coords())
		}
		testPuzzle.Cells[coord.col][coord.row].State = Eliminated
	}

	// Apply the changes
	*p = *testPuzzle
	return p.Deduce()
}

// getCellsToEliminate returns all cells that should be eliminated after placing a star
func (p *Puzzle) getCellsToEliminate(row int, column string, cell Cell) []Coordinate {
	colIndex := letterIndex[column]
	coords := []Coordinate{
		p.coordInt(row-1, colIndex-1),
		p.coordInt(row-1, colIndex),
		p.coordInt(row-1, colIndex+1),
		p.coordInt(row, colIndex-1),
		p.coordInt(row, colIndex+1),
		p.coordInt(row+1, colIndex-1),
		p.coordInt(row+1, colIndex),
		p.coordInt(row+1, colIndex+1),
	}

	// Add cells from same segment if needed
	if p.StarsPerSegment(cell.Segment)+1 == p.CorrectStarsPerArea {
		for _, other := range p.Segments()[cell.Segment] {
			coords = append(coords, coord(other.Row, other.Column))
		}
	}

	// Add cells from same row if needed
	if p.StarsPerRow(row)+1 == p.CorrectStarsPerArea {
		for _, col := range p.ColumnNames() {
			coords = append(coords, coord(row, col))
		}
	}

	// Add cells from same column if needed
	if p.StarsPerColumn(column)+1 == p.CorrectStarsPerArea {
		for idx := 0; idx < p.Height; idx++ {
			coords = append(coords, coord(idx, column))
		}
	}

	return coords
}

// StarsPerSegment returns the number of stars in a segment
func (p *Puzzle) StarsPerSegment(segment Color) int {
	return p.countStars(p.Segments()[segment])
}

// StarsPerRow returns the number of stars in a row
func (p *Puzzle) StarsPerRow(row int) int {
	return p.countStars(p.Rows()[row])
}

// StarsPerColumn returns the number of stars in a column
func (p *Puzzle) StarsPerColumn(column string) int {
	return p.countStars(p.Columns()[column])
}

// Solved checks if the puzzle is solved
func (p *Puzzle) Solved() bool {
	// Check segments
	for color := range p.Segments() {
		if p.StarsPerSegment(color) != p.CorrectStarsPerArea {
			fmt.Printf("not enough stars in %s. found %d, need %d\n", color, p.StarsPerSegment(color), p.CorrectStarsPerArea)
			return false
		}
	}

	// Check rows
	for idx := range p.Rows() {
		if p.StarsPerRow(idx) != p.CorrectStarsPerArea {
			fmt.Printf("not enough stars in row %d. found %d, need %d\n", idx, p.StarsPerRow(idx), p.CorrectStarsPerArea)
			return false
		}
	}

	// Check columns
	for letter := range p.Columns() {
		if p.StarsPerColumn(letter) != p.CorrectStarsPerArea {
			fmt.Printf("not enough stars in column %s. found %d, need %d\n", letter, p.StarsPerColumn(letter), p.CorrectStarsPerArea)
			return false
		}
	}

	return true
}

// IsUnsolvable checks if the puzzle is unsolvable
func (p *Puzzle) IsUnsolvable() bool {
	// Check segments
	for _, segment := range p.Segments() {
		empty := p.countEmpty(segment)
		stars := p.countStars(segment)
		if empty < p.CorrectStarsPerArea-stars {
			return true
		}
	}

	// Check rows
	for _, row := range p.Rows() {
		empty := p.countEmpty(row)
		stars := p.countStars(row)
		if empty < p.CorrectStarsPerArea-stars {
			return true
		}
	}

	// Check columns
	for _, colName := range p.ColumnNames() {
		empty := p.countEmpty(p.Cells[colName])
		stars := p.countStars(p.Cells[colName])
		if empty < p.CorrectStarsPerArea-stars {
			return true
		}
	}

	return false
}

// Solve attempts to solve the puzzle
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

	// Try solving using smallest segments first
	segments := getSortedSegments(puzzle)
	for _, seg := range segments {
		if result, solved := trySolveSegment(puzzle, seg); solved {
			return result, true
		}
	}

	return puzzle, false
}

var globalSolveCounter int

func getSortedSegments(puzzle Puzzle) []segmentInfo {
	var segments []segmentInfo
	for color, cells := range puzzle.Segments() {
		segments = append(segments, segmentInfo{color, cells})
	}
	sort.Slice(segments, func(i, j int) bool {
		return len(segments[i].cells) < len(segments[j].cells)
	})
	return segments
}

type segmentInfo struct {
	color Color
	cells []Cell
}

func trySolveSegment(puzzle Puzzle, seg segmentInfo) (Puzzle, bool) {
	for _, cell := range seg.cells {
		if cell.State != Empty {
			continue
		}

		// Try placing a star
		starClone := puzzle.DeepCopy()
		newPuzzle, err := starClone.Star(cell.Row, cell.Column)
		if err == nil {
			if newPuzzle.Solved() {
				return *newPuzzle, true
			}
			if result, solved := Solve(*newPuzzle); solved {
				return result, true
			}
		}

		// Try eliminating the cell
		elimClone := puzzle.DeepCopy()
		elimClone.Cells[cell.Column][cell.Row].State = Eliminated
		if elimClone.Solved() {
			return *elimClone, true
		}
		if result, solved := Solve(*elimClone); solved {
			return result, true
		}
	}
	return puzzle, false
}

// Deduce applies all constraints to the puzzle
func (p *Puzzle) Deduce() (*Puzzle, error) {
	changed := true
	for changed {
		changed = false
		for _, constraint := range p.AllConstraints() {
			if constraint.Apply(p) {
				fmt.Printf("🧠 Deduction applied: %s\n", constraint.String())
				changed = true
			}
		}
	}
	return p, nil
}

// DeepCopy creates a deep copy of the puzzle
func (p *Puzzle) DeepCopy() *Puzzle {
	newCells := make(map[string][]Cell, len(p.Cells))
	for y := range p.Cells {
		newCells[y] = make([]Cell, len(p.Cells[y]))
		for x, cell := range p.Cells[y] {
			newCells[y][x] = Cell{
				Row:     cell.Row,
				Column:  cell.Column,
				State:   cell.State,
				Segment: cell.Segment,
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

// Constraint interface for puzzle constraints
type Constraint interface {
	Apply(p *Puzzle) bool
	String() string
}

// Base constraint implementation
type baseConstraint struct {
	applyFunc  func(p *Puzzle) bool
	stringFunc func() string
}

func (bc baseConstraint) Apply(p *Puzzle) bool {
	return bc.applyFunc(p)
}

func (bc baseConstraint) String() string {
	return bc.stringFunc()
}

// SegmentConstraint ensures correct number of stars in a segment
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

// RowConstraint ensures correct number of stars in a row
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

// ColumnConstraint ensures correct number of stars in a column
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

// RowSegmentConstraint ensures correct number of stars in a row-segment intersection
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

// ColumnSegmentConstraint ensures correct number of stars in a column-segment intersection
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

// AllConstraints returns all constraints for the puzzle
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
		segment := p.Rows()[row][0].Segment
		constraints = append(constraints, RowSegmentConstraint{Row: row, Segment: segment})
	}

	// Add one column-segment constraint per column
	for _, col := range p.ColumnNames() {
		segment := p.Cells[col][0].Segment
		constraints = append(constraints, ColumnSegmentConstraint{Column: col, Segment: segment})
	}

	return constraints
}
