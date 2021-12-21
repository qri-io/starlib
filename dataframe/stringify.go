package dataframe

import (
	"fmt"
	"strings"
)

// OutputConfig specifies the size of the output area that a stringified
// DataFrame will be shown on
type OutputConfig struct {
	Width  int
	Height int
	// how many rows to show at the bottom when some rows are skipped
	RowsAtBottom int
	// how much blank space to provide in the output area, in rows
	BlankRows int
	// how many columns to show at the right when some columns are skipped
	ColsAtRight int
	// when skipping horizontally, how much remaining space to require
	HorizontalAllowance int
	// when skipping vertically, how much remaining space to require
	VerticalAllowance int
}

func (df *DataFrame) stringify() string {
	outconf := df.outconf
	// If no OutputConfig is assigned, construct one with default values
	if outconf == nil {
		outconf = &OutputConfig{}
	}
	if outconf.Width == 0 {
		outconf.Width = 80
	}
	if outconf.Height == 0 {
		outconf.Height = 24
	}

	// Set magic numbers that lead to a nice, default display
	outconf.RowsAtBottom = 3
	outconf.BlankRows = 4
	outconf.HorizontalAllowance = 12
	// TODO(dustmop): Determine this as max(some fraction of num columns, 2)
	outconf.ColsAtRight = 2
	// 2 is for the header row plus the skipping row (dots)
	outconf.VerticalAllowance = outconf.RowsAtBottom + outconf.BlankRows + 2

	stopRow, renewRow := df.determineRowsToShow(outconf)
	labelWidth, cellWidths := df.determineCellWidths(stopRow, renewRow)
	stopCol, renewCol := df.determineColsToShow(outconf, cellWidths)

	text0 := df.stringifyColumns(stopCol, renewCol, labelWidth, cellWidths)
	text1 := df.stringifyRows(stopRow, renewRow, stopCol, renewCol, labelWidth, cellWidths)
	return text0 + text1
}

func (df *DataFrame) determineRowsToShow(outconf *OutputConfig) (int, int) {
	stopLine := outconf.Height - outconf.VerticalAllowance
	renewLine := df.NumRows() - outconf.RowsAtBottom
	if stopLine >= renewLine {
		return -1, -1
	}
	return stopLine, renewLine
}

func (df *DataFrame) determineColsToShow(outconf *OutputConfig, cellWidths []int) (int, int) {
	if len(cellWidths) < 2 {
		return -1, -1
	}

	endWidths := 0
	for i := 0; i < outconf.ColsAtRight; i++ {
		k := len(cellWidths) - i - 1
		if k < 0 {
			break
		}
		endWidths += cellWidths[k]
	}

	limitWidth := outconf.Width - endWidths - outconf.HorizontalAllowance

	runningSum := 4
	stopIndex := -1
	for i := range cellWidths {
		sum := runningSum + cellWidths[i] + 2
		if sum > limitWidth {
			stopIndex = i
			break
		}
		runningSum = sum
	}

	renewIndex := len(cellWidths) - outconf.ColsAtRight

	if stopIndex == -1 {
		return -1, -1
	}
	return stopIndex, renewIndex
}

func (df *DataFrame) determineCellWidths(stopRow, renewRow int) (int, []int) {
	// Get width of the left-hand label
	labelWidth := 0
	if df.index == nil {
		bodyHeight := df.NumRows()
		k := toWidth(bodyHeight)
		if k > labelWidth {
			labelWidth = k
		}
	} else {
		for _, str := range df.index.Columns() {
			k := len(str)
			if k > labelWidth {
				labelWidth = k
			}
		}
	}
	if stopRow >= 0 && labelWidth < 3 {
		labelWidth = 3
	}

	// Create array of max widths, starting at 0
	cellWidths := make([]int, df.NumCols())
	colTexts := []string{}
	if df.columns != nil {
		colTexts = df.columns.Columns()
	}
	for i, name := range colTexts {
		w := len(name)
		if w > cellWidths[i] {
			cellWidths[i] = w
		}
	}
	for i := 0; i < df.NumRows(); i++ {
		for j, col := range df.body {
			elem := col.StrAt(i)
			w := len(elem)
			if w > cellWidths[j] {
				cellWidths[j] = w
			}
		}
	}

	return labelWidth, cellWidths
}

func (df *DataFrame) stringifyColumns(stopIndex, renewIndex, labelWidth int, cellWidths []int) string {
	colTexts := []string{}
	if df.columns != nil {
		colTexts = df.columns.Columns()
	}

	// Render columns
	header := make([]string, 0, len(colTexts))
	if len(colTexts) > 0 {
		// Render the column names
		for i, name := range colTexts {
			if stopIndex == -1 || i < stopIndex || i >= renewIndex {
				header = append(header, padString(name, cellWidths[i]))
			} else if i == stopIndex {
				header = append(header, "...")
			} else {
				continue
			}
		}
	} else {
		// Render the column indicies
		for i := range df.body {
			header = append(header, padString(i, cellWidths[i]))
		}
	}
	// TODO(dustmop): Move padding sizes into config
	padding := strings.Repeat(" ", labelWidth)
	answer := fmt.Sprintf("%s    %s\n", padding, strings.Join(header, "  "))
	return answer
}

func (df *DataFrame) stringifyRows(stopRow, renewRow, stopCol, renewCol, labelWidth int, cellWidths []int) string {
	collect := []string{}
	// Render each row
	for i := 0; i < df.NumRows(); i++ {
		// For the seam, add "..." for each visible column
		if i == stopRow {
			render := []string{padString("...", labelWidth) + "  "}
			for j := range df.body {
				if j < stopCol || j >= renewCol {
					render = append(render, padString("...", cellWidths[j]))
				} else if j == stopCol {
					render = append(render, "...")
				} else {
					continue
				}
			}
			collect = append(collect, strings.Join(render, "  "))
			continue
		}

		// If this is a row being omitted, skip it
		if stopRow > 0 && i > stopRow && i < renewRow {
			continue
		}

		render := []string{""}
		// Render the index number or label to start the line
		if df.index == nil {
			render[0] = padString(i, labelWidth)
		} else {
			render[0] = padString(df.index.Columns()[i], labelWidth)
		}
		// 2 extra spaces after the lhs label
		render[0] += "  "
		// Render each element of the row
		for j, col := range df.body {
			if j < stopCol || j >= renewCol {
				render = append(render, padString(col.StrAt(i), cellWidths[j]))
			} else if j == stopCol {
				render = append(render, "...")
			} else {
				continue
			}
		}
		collect = append(collect, strings.Join(render, "  "))
	}
	return strings.Join(collect, "\n")
}

func toWidth(n int) int {
	return len(fmt.Sprintf("%d", n))
}

func padString(obj interface{}, n int) string {
	return fmt.Sprintf("%*v", n, obj)
}
