package dataframe

import (
	"fmt"
	"math"
)

type tableBuilder struct {
	numCols  int
	builders []*typedSliceBuilder
	tableErr error
	nameMap  map[string]int
	names    []string
}

func newTableBuilder(numCols, rowCapacity int) *tableBuilder {
	builders := make([]*typedSliceBuilder, numCols)
	for x := 0; x < numCols; x++ {
		builders[x] = newTypedSliceBuilder(rowCapacity)
	}
	return &tableBuilder{numCols: numCols, builders: builders}
}

func (b *tableBuilder) setDtypes(body []Series) {
	for k, series := range body {
		b.builders[k].dType = series.dtype
	}
}

func (b *tableBuilder) pushNamedRow(row *namedRow) {
	if err := row.Validate(); err != nil {
		b.tableErr = err
		return
	}
	if b.nameMap == nil {
		b.nameMap = make(map[string]int)
	}

	// how many rows are in the table being built
	currSize := 0
	if len(b.builders) > 0 {
		currSize = b.builders[0].Len()
	}

	// track which existing column builders are added to
	satisfy := make([]bool, len(b.builders))
	for i, name := range row.names {
		idx, ok := b.nameMap[name]
		if !ok {
			// no column with this name, add a new one, fill
			// with null values
			idx = len(b.builders)
			b.nameMap[name] = idx
			newColumn := newTypedSliceBuilderNaNFilled(currSize)
			b.builders = append(b.builders, newColumn)
			b.names = append(b.names, name)
			b.numCols = len(b.builders)
		}
		if idx < len(satisfy) {
			satisfy[idx] = true
		}
		b.builders[idx].push(row.values[i])
	}

	// if any columns were not added to, give them NaN
	for j, sat := range satisfy {
		if !sat {
			b.builders[j].push(math.NaN())
		}
	}
}

func (b *tableBuilder) pushRow(row []interface{}) {
	if len(row) != b.numCols {
		b.tableErr = fmt.Errorf("size of row %d does not match size of body: %d", len(row), b.numCols)
		return
	}
	for x := 0; x < b.numCols; x++ {
		b.builders[x].push(row[x])
	}
}

func (b *tableBuilder) pushTextRow(row []string) {
	if len(row) != b.numCols {
		b.tableErr = fmt.Errorf("size of row %d does not match size of body: %d", len(row), b.numCols)
		return
	}
	for x := 0; x < b.numCols; x++ {
		b.builders[x].parsePush(row[x])
	}
}

func (b *tableBuilder) body() ([]Series, error) {
	if b == nil {
		return []Series{}, nil
	}
	if b.tableErr != nil {
		return nil, b.tableErr
	}
	for x := 0; x < b.numCols; x++ {
		if b.builders[x].buildError != nil {
			return nil, b.builders[x].buildError
		}
	}
	result := make([]Series, b.numCols)
	for x := 0; x < b.numCols; x++ {
		result[x] = b.builders[x].toSeries(nil, "")
	}
	return result, nil
}

func (b *tableBuilder) colNames() []string {
	return b.names
}

type namedRow struct {
	names  []string
	values []interface{}
}

func newNamedRow(names []string, values []interface{}) *namedRow {
	return &namedRow{names: names, values: values}
}

func (r *namedRow) Validate() error {
	if r.names == nil {
		return fmt.Errorf("namedRow has nil names")
	}
	if r.values == nil {
		return fmt.Errorf("namedRow has nil values")
	}
	if len(r.names) != len(r.values) {
		return fmt.Errorf("namedRow does not have same length for names and values")
	}
	return nil
}

func (r *namedRow) Len() int {
	return len(r.names)
}
