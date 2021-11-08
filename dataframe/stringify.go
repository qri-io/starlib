package dataframe

import (
	"fmt"
	"strings"
)

type stringConfig struct {
	terminalWidth  int
	terminalHeight int
}

func (df *DataFrame) stringify() string {
	// TODO(dustmop): Detect terminal size, instead of using hard-coded values
	df.stringConf.terminalWidth = 80
	df.stringConf.terminalHeight = 24

	stopRow, renewRow := df.determineRowsToShow(df.stringConf)
	labelWidth, cellWidths := df.determineCellWidths(stopRow, renewRow)
	stopCol, renewCol := df.determineColsToShow(df.stringConf, cellWidths)

	text0 := df.stringifyColumns(stopCol, renewCol, labelWidth, cellWidths)
	text1 := df.stringifyRows(stopRow, renewRow, stopCol, renewCol, labelWidth, cellWidths)
	return text0 + text1
}

func (df *DataFrame) determineRowsToShow(conf stringConfig) (int, int) {
	// TODO(dustmop): Move magic constants into config
	stopLine := conf.terminalHeight - 7
	renewLine := df.NumRows() - 3
	if stopLine >= renewLine {
		return -1, -1
	}
	return stopLine, renewLine
}

func (df *DataFrame) determineColsToShow(conf stringConfig, cellWidths []int) (int, int) {
	if len(cellWidths) < 2 {
		return -1, -1
	}

	endWidths := 0
	endWidths += cellWidths[len(cellWidths)-2]
	endWidths += cellWidths[len(cellWidths)-1]

	// TODO(dustmop): Move magic constants into config
	limitWidth := conf.terminalWidth - endWidths - 11

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

	renewIndex := len(cellWidths) - 2

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
		for _, str := range df.index.texts {
			k := len(str)
			if k > labelWidth {
				labelWidth = k
			}
		}
	}
	if renewRow >= 0 && labelWidth < 3 {
		labelWidth = 3
	}

	// Create array of max widths, starting at 0
	cellWidths := make([]int, df.NumCols())
	colTexts := []string{}
	if df.columns != nil {
		colTexts = df.columns.texts
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
		colTexts = df.columns.texts
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
		if i > stopRow && i < renewRow {
			continue
		}

		render := []string{""}
		// Render the index number or label to start the line
		if df.index == nil {
			render[0] = padString(i, labelWidth)
		} else {
			render[0] = padString(df.index.texts[i], labelWidth)
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
