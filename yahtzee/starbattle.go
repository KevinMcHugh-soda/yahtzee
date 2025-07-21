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
	Blocked    = "🟫"
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

	// Create a copy to test the move
	testPuzzle := p.DeepCopy()
	testPuzzle.Cells[column][row].State = Starred

	// Check if this would create too many stars in any segment, row, or column
	for color, cells := range testPuzzle.Segments() {
		stars := testPuzzle.countStars(cells)
		if stars > testPuzzle.CorrectStarsPerArea {
			return p, fmt.Errorf("too many stars in segment %s", color)
		}
		empty := testPuzzle.countEmpty(cells)
		if empty < testPuzzle.CorrectStarsPerArea-stars {
			return p, fmt.Errorf("not enough empty cells in segment %s", color)
		}
	}

	for idx, row := range testPuzzle.Rows() {
		stars := testPuzzle.countStars(row)
		if stars > testPuzzle.CorrectStarsPerArea {
			return p, fmt.Errorf("too many stars in row %d", idx)
		}
		empty := testPuzzle.countEmpty(row)
		if empty < testPuzzle.CorrectStarsPerArea-stars {
			return p, fmt.Errorf("not enough empty cells in row %d", idx)
		}
	}

	for col := range testPuzzle.Columns() {
		stars := testPuzzle.countStars(testPuzzle.Cells[col])
		if stars > testPuzzle.CorrectStarsPerArea {
			return p, fmt.Errorf("too many stars in column %s", col)
		}
		empty := testPuzzle.countEmpty(testPuzzle.Cells[col])
		if empty < testPuzzle.CorrectStarsPerArea-stars {
			return p, fmt.Errorf("not enough empty cells in column %s", col)
		}
	}

	// Check if any cell has stars in adjacent cells
	colIndex := letterIndex[column]
	for _, offset := range [][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}} {
		newRow := row + offset[0]
		newCol := colIndex + offset[1]
		if newRow >= 0 && newRow < testPuzzle.Height && newCol >= 0 && newCol < testPuzzle.Width {
			if testPuzzle.Cells[letters[newCol]][newRow].State == Starred {
				return p, fmt.Errorf("star would be adjacent to another star")
			}
		}
	}

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
			fmt.Printf("incorrect number of stars in %s. found %d, need %d\n", color, p.StarsPerSegment(color), p.CorrectStarsPerArea)
			return false
		}
	}

	// Check rows
	for idx := range p.Rows() {
		if p.StarsPerRow(idx) != p.CorrectStarsPerArea {
			fmt.Printf("incorrect number of stars in row %d. found %d, need %d\n", idx, p.StarsPerRow(idx), p.CorrectStarsPerArea)
			return false
		}
	}

	// Check columns
	for letter := range p.Columns() {
		if p.StarsPerColumn(letter) != p.CorrectStarsPerArea {
			fmt.Printf("incorrect number of stars in column %s. found %d, need %d\n", letter, p.StarsPerColumn(letter), p.CorrectStarsPerArea)
			return false
		}
	}

	return true
}

// IsUnsolvable returns true if the puzzle is in an unsolvable state
func (p *Puzzle) IsUnsolvable() bool {
	// Check if any row or column has too many stars
	for _, row := range p.Rows() {
		stars := 0
		for _, cell := range row {
			if cell.State == Starred {
				stars++
			}
		}
		if stars > p.CorrectStarsPerArea {
			return true
		}
	}

	for _, col := range p.Columns() {
		stars := 0
		for _, cell := range col {
			if cell.State == Starred {
				stars++
			}
		}
		if stars > p.CorrectStarsPerArea {
			return true
		}
	}

	// Check if any segment has too many stars
	for color := range p.Segments() {
		stars := p.StarsPerSegment(color)
		if stars > p.CorrectStarsPerArea {
			return true
		}
	}

	// Check if any row, column, or segment can't get enough stars
	for _, row := range p.Rows() {
		availableCells := 0
		for _, cell := range row {
			if cell.State != Blocked {
				availableCells++
			}
		}
		if availableCells < p.CorrectStarsPerArea {
			return true
		}
	}

	for _, col := range p.Columns() {
		availableCells := 0
		for _, cell := range col {
			if cell.State != Blocked {
				availableCells++
			}
		}
		if availableCells < p.CorrectStarsPerArea {
			return true
		}
	}

	for _, cells := range p.Segments() {
		availableCells := 0
		for _, cell := range cells {
			if p.Cells[cell.Column][cell.Row].State != Blocked {
				availableCells++
			}
		}
		if availableCells < p.CorrectStarsPerArea {
			return true
		}
	}

	// Check if any adjacent cells both have stars
	for i := 0; i < p.Height; i++ {
		for j := 0; j < p.Width; j++ {
			if p.Cells[letters[j]][i].State == Starred {
				// Check adjacent cells
				for _, offset := range [][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}} {
					newRow := i + offset[0]
					newCol := j + offset[1]
					if newRow >= 0 && newRow < p.Height && newCol >= 0 && newCol < p.Width {
						if p.Cells[letters[newCol]][newRow].State == Starred {
							return true
						}
					}
				}
			}
		}
	}

	// Check if any row, column, or segment has too many blocked cells
	for _, row := range p.Rows() {
		blocked := 0
		for _, cell := range row {
			if cell.State == Blocked {
				blocked++
			}
		}
		if p.Width-blocked < p.CorrectStarsPerArea {
			return true
		}
	}

	for _, col := range p.Columns() {
		blocked := 0
		for _, cell := range col {
			if cell.State == Blocked {
				blocked++
			}
		}
		if p.Height-blocked < p.CorrectStarsPerArea {
			return true
		}
	}

	for _, cells := range p.Segments() {
		blocked := 0
		for _, cell := range cells {
			if p.Cells[cell.Column][cell.Row].State == Blocked {
				blocked++
			}
		}
		if len(cells)-blocked < p.CorrectStarsPerArea {
			return true
		}
	}

	return false
}

var globalSolveCounter int

// Solve attempts to solve the puzzle
func Solve(puzzle Puzzle) (Puzzle, bool) {
	globalSolveCounter++
	if globalSolveCounter > 50000 {
		fmt.Println("bailing!")
		return puzzle, false
	}

	if puzzle.Solved() {
		return puzzle, true
	}

	if puzzle.IsUnsolvable() {
		return puzzle, false
	}

	// Try to deduce first
	deduced, err := puzzle.Deduce()
	if err != nil {
		return puzzle, false
	}
	puzzle = *deduced

	// Try to place stars systematically by segment
	segments := getSortedSegments(puzzle)
	sort.SliceStable(segments, func(i, j int) bool {
		starsI := puzzle.StarsPerSegment(segments[i].color)
		starsJ := puzzle.StarsPerSegment(segments[j].color)
		if starsI != starsJ {
			return starsI < starsJ
		}
		// Prioritize segments with fewer empty cells
		emptyI := puzzle.countEmpty(segments[i].cells)
		emptyJ := puzzle.countEmpty(segments[j].cells)
		if emptyI != emptyJ {
			return emptyI < emptyJ
		}
		return len(segments[i].cells) < len(segments[j].cells)
	})

	// Try each segment in order
	for _, segment := range segments {
		stars := puzzle.StarsPerSegment(segment.color)
		if stars >= puzzle.CorrectStarsPerArea {
			continue
		}

		// Get empty cells in this segment
		emptyCells := make([]Cell, 0)
		for _, cell := range segment.cells {
			if puzzle.Cells[cell.Column][cell.Row].State == Empty {
				emptyCells = append(emptyCells, cell)
			}
		}

		// Sort empty cells by number of available neighbors
		sort.SliceStable(emptyCells, func(i, j int) bool {
			availableI := puzzle.countAvailableNeighbors(emptyCells[i])
			availableJ := puzzle.countAvailableNeighbors(emptyCells[j])
			return availableI < availableJ
		})

		// Try each empty cell
		for _, cell := range emptyCells {
			// Try placing a star
			newPuzzle := puzzle.DeepCopy()
			newPuzzle.Cells[cell.Column][cell.Row].State = Starred

			// Recursively solve
			solved, success := Solve(*newPuzzle)
			if success {
				return solved, true
			}

			// If placing a star didn't work, try eliminating the cell
			newPuzzle = puzzle.DeepCopy()
			newPuzzle.Cells[cell.Column][cell.Row].State = Eliminated

			// Recursively solve
			solved, success = Solve(*newPuzzle)
			if success {
				return solved, true
			}
		}
	}

	return puzzle, false
}

// countAvailableNeighbors counts the number of empty or starred neighbors a cell has
func (p *Puzzle) countAvailableNeighbors(cell Cell) int {
	count := 0
	for _, offset := range [][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}} {
		row := letterIndex[cell.Column] + offset[0]
		col := cell.Row + offset[1]
		if row >= 0 && row < p.Width && col >= 0 && col < p.Height {
			state := p.Cells[letters[row]][col].State
			if state == Empty || state == Starred {
				count++
			}
		}
	}
	return count
}

// countSegmentConstraints returns the number of constraints affecting a segment
func countSegmentConstraints(p Puzzle, cells []Cell) int {
	constraints := 0
	for _, cell := range cells {
		if p.Cells[cell.Column][cell.Row].State != Empty {
			constraints++
		}
		// Count adjacent stars and eliminated cells
		for _, offset := range [][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}} {
			row := cell.Row + offset[0]
			col := letterIndex[cell.Column] + offset[1]
			if row >= 0 && row < p.Height && col >= 0 && col < p.Width {
				state := p.Cells[letters[col]][row].State
				if state == Starred || state == Eliminated {
					constraints++
				}
			}
		}
	}
	return constraints
}

// countCellConstraints returns the number of constraints affecting a cell
func countCellConstraints(p Puzzle, cell Cell) int {
	constraints := 0
	// Count adjacent stars and eliminated cells
	for _, offset := range [][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}} {
		row := cell.Row + offset[0]
		col := letterIndex[cell.Column] + offset[1]
		if row >= 0 && row < p.Height && col >= 0 && col < p.Width {
			state := p.Cells[letters[col]][row].State
			if state == Starred || state == Eliminated {
				constraints++
			}
		}
	}
	return constraints
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// generateCombinations generates all possible combinations of k cells from a slice of cells
func generateCombinations(cells []Cell, k int) [][]Cell {
	if k == 0 {
		return [][]Cell{{}}
	}
	if len(cells) == 0 {
		return nil
	}

	first := cells[0]
	rest := cells[1:]

	// Combinations that include the first element
	withFirst := generateCombinations(rest, k-1)
	for i := range withFirst {
		withFirst[i] = append([]Cell{first}, withFirst[i]...)
	}

	// Combinations that exclude the first element
	withoutFirst := generateCombinations(rest, k)

	return append(withFirst, withoutFirst...)
}

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

// Deduce applies all constraints to the puzzle
func (p *Puzzle) Deduce() (*Puzzle, error) {
	changed := true
	for changed {
		changed = false

		// Apply row constraints
		for i := 0; i < p.Height; i++ {
			row := p.Rows()[i]
			stars := p.countStars(row)
			if stars == p.CorrectStarsPerArea {
				// Mark all remaining empty cells as eliminated
				for _, cell := range row {
					if cell.State == Empty {
						cell.State = Eliminated
						changed = true
					}
				}
			}
		}

		// Apply column constraints
		for _, col := range p.ColumnNames() {
			stars := p.StarsPerColumn(col)
			if stars == p.CorrectStarsPerArea {
				// Mark all remaining empty cells as eliminated
				for _, cell := range p.Cells[col] {
					if cell.State == Empty {
						cell.State = Eliminated
						changed = true
					}
				}
			}
		}

		// Apply segment constraints
		for color, cells := range p.Segments() {
			stars := p.StarsPerSegment(color)
			if stars == p.CorrectStarsPerArea {
				// Mark all remaining empty cells as eliminated
				for _, cell := range cells {
					if p.Cells[cell.Column][cell.Row].State == Empty {
						p.Cells[cell.Column][cell.Row].State = Eliminated
						changed = true
					}
				}
			}
		}

		// Apply adjacency constraints
		for i := 0; i < p.Height; i++ {
			for j := 0; j < p.Width; j++ {
				if p.Cells[letters[j]][i].State == Starred {
					// Eliminate adjacent cells
					for _, offset := range [][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}} {
						newRow := i + offset[0]
						newCol := j + offset[1]
						if newRow >= 0 && newRow < p.Height && newCol >= 0 && newCol < p.Width {
							if p.Cells[letters[newCol]][newRow].State == Empty {
								p.Cells[letters[newCol]][newRow].State = Eliminated
								changed = true
							}
						}
					}
				}
			}
		}

		// Apply forced star placement
		for i := 0; i < p.Height; i++ {
			for j := 0; j < p.Width; j++ {
				if p.Cells[letters[j]][i].State == Empty {
					// Check if this is the only possible position for a star in its row, column, or segment
					row := p.Rows()[i]
					col := p.Cells[letters[j]]
					segment := p.Segments()[p.GetCellColor(letters[j], i)]

					rowStars := p.countStars(row)
					colStars := p.StarsPerColumn(letters[j])
					segStars := p.StarsPerSegment(p.GetCellColor(letters[j], i))

					rowEmpty := p.countEmpty(row)
					colEmpty := p.countEmpty(col)
					segEmpty := p.countEmpty(segment)

					if (rowStars < p.CorrectStarsPerArea && rowEmpty == 1) ||
						(colStars < p.CorrectStarsPerArea && colEmpty == 1) ||
						(segStars < p.CorrectStarsPerArea && segEmpty == 1) {
						p.Cells[letters[j]][i].State = Starred
						changed = true
					}
				}
			}
		}

		// Apply forced elimination
		for i := 0; i < p.Height; i++ {
			for j := 0; j < p.Width; j++ {
				if p.Cells[letters[j]][i].State == Empty {
					// Check if placing a star here would force too many stars in a row, column, or segment
					rowStars := p.StarsPerRow(i)
					colStars := p.StarsPerColumn(letters[j])
					segStars := p.StarsPerSegment(p.GetCellColor(letters[j], i))

					if rowStars >= p.CorrectStarsPerArea ||
						colStars >= p.CorrectStarsPerArea ||
						segStars >= p.CorrectStarsPerArea {
						p.Cells[letters[j]][i].State = Eliminated
						changed = true
					}
				}
			}
		}

		// Apply pattern-based deductions
		for i := 0; i < p.Height; i++ {
			for j := 0; j < p.Width; j++ {
				if p.Cells[letters[j]][i].State == Empty {
					// Check for patterns that force star placement or elimination
					if p.isForceStarPattern(i, j) {
						p.Cells[letters[j]][i].State = Starred
						changed = true
					} else if p.isForceEliminationPattern(i, j) {
						p.Cells[letters[j]][i].State = Eliminated
						changed = true
					}
				}
			}
		}
	}

	return p, nil
}

// isForceStarPattern checks if the cell must contain a star based on surrounding patterns
func (p *Puzzle) isForceStarPattern(row, col int) bool {
	// Check if this is the only possible position for a star in a 2x2 area
	for _, area := range [][4][2]int{
		{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
		{{-1, 0}, {-1, 1}, {0, 0}, {0, 1}},
		{{0, -1}, {0, 0}, {1, -1}, {1, 0}},
		{{-1, -1}, {-1, 0}, {0, -1}, {0, 0}},
	} {
		emptyCount := 0
		validArea := true
		for _, pos := range area {
			newRow := row + pos[0]
			newCol := col + pos[1]
			if newRow < 0 || newRow >= p.Height || newCol < 0 || newCol >= p.Width {
				validArea = false
				break
			}
			if p.Cells[letters[newCol]][newRow].State == Empty {
				emptyCount++
			} else if p.Cells[letters[newCol]][newRow].State == Starred {
				validArea = false
				break
			}
		}
		if validArea && emptyCount == 1 {
			return true
		}
	}
	return false
}

// isForceEliminationPattern checks if the cell cannot contain a star based on surrounding patterns
func (p *Puzzle) isForceEliminationPattern(row, col int) bool {
	// Check if placing a star here would create an impossible pattern
	for _, pattern := range [][3][2]int{
		{{0, 0}, {0, 1}, {0, 2}},
		{{0, 0}, {1, 0}, {2, 0}},
		{{0, 0}, {1, 1}, {2, 2}},
		{{0, 2}, {1, 1}, {2, 0}},
	} {
		starCount := 0
		validPattern := true
		for _, pos := range pattern {
			newRow := row + pos[0]
			newCol := col + pos[1]
			if newRow < 0 || newRow >= p.Height || newCol < 0 || newCol >= p.Width {
				validPattern = false
				break
			}
			if p.Cells[letters[newCol]][newRow].State == Starred {
				starCount++
			}
		}
		if validPattern && starCount >= 2 {
			return true
		}
	}
	return false
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
		// If we have enough stars, eliminate all other cells
		changed := false
		for _, cell := range empties {
			p.Cells[cell.Column][cell.Row].State = Eliminated
			changed = true
		}
		return changed
	} else if needed == len(empties) {
		// If we need exactly the number of empty cells, they must all be stars
		changed := false
		for _, cell := range empties {
			p.Cells[cell.Column][cell.Row].State = Starred
			changed = true
		}
		return changed
	} else if needed < 0 {
		// If we have too many stars, the puzzle is unsolvable
		return false
	} else if len(empties) < needed {
		// If we don't have enough empty cells, the puzzle is unsolvable
		return false
	}

	// Check if any empty cells are adjacent to stars
	changed := false
	for _, cell := range empties {
		colIndex := letterIndex[cell.Column]
		adjacentToStar := false
		for _, offset := range [][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}} {
			row := cell.Row + offset[0]
			col := colIndex + offset[1]
			if row >= 0 && row < p.Height && col >= 0 && col < p.Width {
				adjCell := p.Cells[letters[col]][row]
				if adjCell.State == Starred {
					adjacentToStar = true
					break
				}
			}
		}
		if adjacentToStar {
			p.Cells[cell.Column][cell.Row].State = Eliminated
			changed = true
		}
	}
	return changed
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
			} else if cell.State == Starred {
				starsFromOtherSegments++
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
	} else if needed == len(segCellsInRow) {
		changed := false
		for _, cell := range segCellsInRow {
			if p.Cells[cell.Column][cell.Row].State == Empty {
				p.Cells[cell.Column][cell.Row].State = Starred
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
			} else if cell.State == Starred {
				starsFromOtherSegments++
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
	} else if needed == len(segCellsInCol) {
		changed := false
		for _, cell := range segCellsInCol {
			if p.Cells[cell.Column][cell.Row].State == Empty {
				p.Cells[cell.Column][cell.Row].State = Starred
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

	// Add row-segment constraints for each segment in each row
	for row := 0; row < p.Height; row++ {
		segmentsInRow := make(map[Color]bool)
		for _, cell := range p.Rows()[row] {
			if !segmentsInRow[cell.Segment] {
				constraints = append(constraints, RowSegmentConstraint{Row: row, Segment: cell.Segment})
				segmentsInRow[cell.Segment] = true
			}
		}
	}

	// Add column-segment constraints for each segment in each column
	for _, col := range p.ColumnNames() {
		segmentsInCol := make(map[Color]bool)
		for _, cell := range p.Cells[col] {
			if !segmentsInCol[cell.Segment] {
				constraints = append(constraints, ColumnSegmentConstraint{Column: col, Segment: cell.Segment})
				segmentsInCol[cell.Segment] = true
			}
		}
	}

	return constraints
}
