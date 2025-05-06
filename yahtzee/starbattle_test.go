package yahtzee

import (
	"testing"
)

func TestParsePuzzle(t *testing.T) {
	tests := []struct {
		name         string
		rows         []string
		starsPerArea int
		wantErr      bool
	}{
		{
			name:         "valid 5x5 puzzle",
			rows:         []string{fiveXfive1, fiveXfive2, fiveXfive3, fiveXfive4, fiveXfive5},
			starsPerArea: 1,
			wantErr:      false,
		},
		{
			name:         "unequal row lengths",
			rows:         []string{"🟨🟨🟩", "🟦🟨🟩🟩"},
			starsPerArea: 1,
			wantErr:      true,
		},
		{
			name:         "empty puzzle",
			rows:         []string{},
			starsPerArea: 1,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePuzzle(tt.rows, tt.starsPerArea)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePuzzle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("ParsePuzzle() returned nil puzzle when no error expected")
			}
		})
	}
}

func TestCellCoords(t *testing.T) {
	tests := []struct {
		name     string
		cell     Cell
		expected string
	}{
		{
			name:     "A0",
			cell:     Cell{Row: 0, Column: "A"},
			expected: "A0",
		},
		{
			name:     "B1",
			cell:     Cell{Row: 1, Column: "B"},
			expected: "B1",
		},
		{
			name:     "C2",
			cell:     Cell{Row: 2, Column: "C"},
			expected: "C2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cell.Coords(); got != tt.expected {
				t.Errorf("Cell.Coords() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPuzzleStar(t *testing.T) {
	puzzle, _ := ParsePuzzle([]string{fiveXfive1, fiveXfive2, fiveXfive3, fiveXfive4, fiveXfive5}, 1)

	tests := []struct {
		name    string
		row     int
		column  string
		wantErr bool
	}{
		{
			name:    "valid star placement",
			row:     0,
			column:  "A",
			wantErr: false,
		},
		{
			name:    "invalid column",
			row:     0,
			column:  "Z",
			wantErr: true,
		},
		{
			name:    "invalid row",
			row:     10,
			column:  "A",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := puzzle.Star(tt.row, tt.column)
			if (err != nil) != tt.wantErr {
				t.Errorf("Puzzle.Star() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPuzzleSolved(t *testing.T) {
	// Create a solved puzzle
	rows := []string{
		"🟨🟨🟩🟩🟩",
		"🟦🟨🟩🟩🟩",
		"🟦🟥🟧🟧🟩",
		"🟥🟥🟧🟧🟩",
		"🟥🟥🟥🟥🟥",
	}
	puzzle, _ := ParsePuzzle(rows, 1)

	// Place stars in a valid configuration - one per segment, row, and column
	puzzle.Cells["B"][0].State = Starred // Yellow segment
	puzzle.Cells["D"][1].State = Starred // Green segment
	puzzle.Cells["A"][2].State = Starred // Blue segment
	puzzle.Cells["C"][3].State = Starred // Orange segment
	puzzle.Cells["E"][4].State = Starred // Red segment

	// Print the puzzle state for debugging
	puzzle.Print("Puzzle state after placing stars:")

	// Print segment information for debugging
	for color, cells := range puzzle.Segments() {
		stars := 0
		for _, cell := range cells {
			if cell.State == Starred {
				stars++
			}
		}
		t.Logf("Segment %s has %d stars", color, stars)
	}

	// Print row information for debugging
	for i := 0; i < puzzle.Height; i++ {
		stars := 0
		for _, cell := range puzzle.Rows()[i] {
			if cell.State == Starred {
				stars++
			}
		}
		t.Logf("Row %d has %d stars", i, stars)
	}

	// Print column information for debugging
	for _, col := range puzzle.ColumnNames() {
		stars := 0
		for _, cell := range puzzle.Columns()[col] {
			if cell.State == Starred {
				stars++
			}
		}
		t.Logf("Column %s has %d stars", col, stars)
	}

	tests := []struct {
		name     string
		puzzle   *Puzzle
		expected bool
	}{
		{
			name:     "solved puzzle",
			puzzle:   puzzle,
			expected: true,
		},
		{
			name:     "unsolved puzzle",
			puzzle:   &Puzzle{CorrectStarsPerArea: 1, Width: 5, Height: 5},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.puzzle.Solved(); got != tt.expected {
				t.Errorf("Puzzle.Solved() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPuzzleDeduce(t *testing.T) {
	rows := []string{
		"🟨🟨🟩🟩🟩",
		"🟦🟨🟩🟩🟩",
		"🟦🟥🟧🟧🟩",
		"🟥🟥🟧🟧🟩",
		"🟥🟥🟥🟥🟥",
	}
	puzzle, _ := ParsePuzzle(rows, 1)

	tests := []struct {
		name    string
		puzzle  *Puzzle
		wantErr bool
	}{
		{
			name:    "valid deduction",
			puzzle:  puzzle,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.puzzle.Deduce()
			if (err != nil) != tt.wantErr {
				t.Errorf("Puzzle.Deduce() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPuzzleDeepCopy(t *testing.T) {
	rows := []string{
		"🟨🟨🟩🟩🟩",
		"🟦🟨🟩🟩🟩",
		"🟦🟥🟧🟧🟩",
		"🟥🟥🟧🟧🟩",
		"🟥🟥🟥🟥🟥",
	}
	original, _ := ParsePuzzle(rows, 1)

	// Place a star in the original
	original.Cells["A"][0].State = Starred

	// Make a copy
	copy := original.DeepCopy()

	// Modify the copy
	copy.Cells["B"][1].State = Starred

	// Check that original wasn't modified
	if original.Cells["B"][1].State != Empty {
		t.Error("DeepCopy() modified original puzzle")
	}

	// Check that copy has both stars
	if copy.Cells["A"][0].State != Starred || copy.Cells["B"][1].State != Starred {
		t.Error("DeepCopy() didn't preserve all state")
	}
}

func TestSegmentConstraint(t *testing.T) {
	rows := []string{
		"🟨🟨🟩🟩🟩",
		"🟦🟨🟩🟩🟩",
		"🟦🟥🟧🟧🟩",
		"🟥🟥🟧🟧🟩",
		"🟥🟥🟥🟥🟥",
	}
	puzzle, _ := ParsePuzzle(rows, 1)

	// Test segment constraint
	constraint := SegmentConstraint{Segment: Yellow}

	// Place a star in the yellow segment
	puzzle.Cells["A"][0].State = Starred

	// Apply constraint
	changed := constraint.Apply(puzzle)

	// Verify that other cells in yellow segment are eliminated
	if puzzle.Cells["B"][0].State != Eliminated {
		t.Error("SegmentConstraint did not eliminate other cells in segment")
	}

	if !changed {
		t.Error("SegmentConstraint did not report changes")
	}
}

func TestRowConstraint(t *testing.T) {
	rows := []string{
		"🟨🟨🟩🟩🟩",
		"🟦🟨🟩🟩🟩",
		"🟦🟥🟧🟧🟩",
		"🟥🟥🟧🟧🟩",
		"🟥🟥🟥🟥🟥",
	}
	puzzle, _ := ParsePuzzle(rows, 1)

	// Test row constraint
	constraint := RowConstraint{Row: 0}

	// Place a star in the row
	puzzle.Cells["A"][0].State = Starred

	// Apply constraint
	changed := constraint.Apply(puzzle)

	// Verify that other cells in row are eliminated
	if puzzle.Cells["B"][0].State != Eliminated {
		t.Error("RowConstraint did not eliminate other cells in row")
	}

	if !changed {
		t.Error("RowConstraint did not report changes")
	}
}

func TestColumnConstraint(t *testing.T) {
	rows := []string{
		"🟨🟨🟩🟩🟩",
		"🟦🟨🟩🟩🟩",
		"🟦🟥🟧🟧🟩",
		"🟥🟥🟧🟧🟩",
		"🟥🟥🟥🟥🟥",
	}
	puzzle, _ := ParsePuzzle(rows, 1)

	// Test column constraint
	constraint := ColumnConstraint{Column: "A"}

	// Place a star in the column
	puzzle.Cells["A"][0].State = Starred

	// Apply constraint
	changed := constraint.Apply(puzzle)

	// Verify that other cells in column are eliminated
	if puzzle.Cells["A"][1].State != Eliminated {
		t.Error("ColumnConstraint did not eliminate other cells in column")
	}

	if !changed {
		t.Error("ColumnConstraint did not report changes")
	}
}

func TestRowSegmentConstraint(t *testing.T) {
	rows := []string{
		"🟨🟨🟩🟩🟩",
		"🟦🟨🟩🟩🟩",
		"🟦🟥🟧🟧🟩",
		"🟥🟥🟧🟧🟩",
		"🟥🟥🟥🟥🟥",
	}
	puzzle, _ := ParsePuzzle(rows, 1)

	// Test row segment constraint
	constraint := RowSegmentConstraint{Row: 0, Segment: Yellow}

	// Place a star from another segment in the row
	puzzle.Cells["C"][0].State = Starred

	// Apply constraint
	changed := constraint.Apply(puzzle)

	// Verify that cells in yellow segment in row 0 are eliminated
	if puzzle.Cells["A"][0].State != Eliminated || puzzle.Cells["B"][0].State != Eliminated {
		t.Error("RowSegmentConstraint did not eliminate cells in segment")
	}

	if !changed {
		t.Error("RowSegmentConstraint did not report changes")
	}
}

func TestColumnSegmentConstraint(t *testing.T) {
	rows := []string{
		"🟨🟨🟩🟩🟩",
		"🟦🟨🟩🟩🟩",
		"🟦🟥🟧🟧🟩",
		"🟥🟥🟧🟧🟩",
		"🟥🟥🟥🟥🟥",
	}
	puzzle, _ := ParsePuzzle(rows, 1)

	// Test column segment constraint
	constraint := ColumnSegmentConstraint{Column: "A", Segment: Yellow}

	// Place a star from another segment in the column
	puzzle.Cells["A"][2].State = Starred

	// Apply constraint
	changed := constraint.Apply(puzzle)

	// Verify that cells in yellow segment in column A are eliminated
	if puzzle.Cells["A"][0].State != Eliminated {
		t.Error("ColumnSegmentConstraint did not eliminate cells in segment")
	}

	if !changed {
		t.Error("ColumnSegmentConstraint did not report changes")
	}
}

func TestAllConstraints(t *testing.T) {
	testRows := []string{
		"🟨🟨🟩🟩🟩",
		"🟦🟨🟩🟩🟩",
		"🟦🟥🟧🟧🟩",
		"🟥🟥🟧🟧🟩",
		"🟥🟥🟥🟥🟥",
	}
	puzzle, _ := ParsePuzzle(testRows, 1)

	constraints := puzzle.AllConstraints()

	// Print all constraints for debugging
	t.Logf("Found %d constraints:", len(constraints))
	for _, c := range constraints {
		t.Logf("  %s", c.String())
	}

	// Count each type of constraint
	var segments, numRows, cols, rowSegs, colSegs int

	for _, c := range constraints {
		switch c.(type) {
		case SegmentConstraint:
			segments++
		case RowConstraint:
			numRows++
		case ColumnConstraint:
			cols++
		case RowSegmentConstraint:
			rowSegs++
		case ColumnSegmentConstraint:
			colSegs++
		}
	}

	t.Logf("Segments: %d, Rows: %d, Cols: %d, RowSegs: %d, ColSegs: %d",
		segments, numRows, cols, rowSegs, colSegs)

	// Verify we have the expected number of constraints
	expectedConstraints := len(puzzle.Segments()) + // Segment constraints
		puzzle.Height + // Row constraints
		puzzle.Width + // Column constraints
		puzzle.Height + // Row segment constraints (one per row)
		puzzle.Width // Column segment constraints (one per column)

	if len(constraints) != expectedConstraints {
		t.Errorf("AllConstraints() returned %d constraints, want %d", len(constraints), expectedConstraints)
	}
}
