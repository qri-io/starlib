package dataframe

import (
	"fmt"
)

type tableBuilder struct {
	numCols  int
	builders []*typedSliceBuilder
}

func newTableBuilder(numCols, rowCapacity int) *tableBuilder {
	builders := make([]*typedSliceBuilder, numCols)
	for x := 0; x < numCols; x++ {
		builders[x] = newTypedSliceBuilder(rowCapacity)
	}
	return &tableBuilder{numCols: numCols, builders: builders}
}

func (b *tableBuilder) pushRow(row []interface{}) {
	if len(row) != b.numCols {
		panic(fmt.Errorf("size of row %d does not match size of body: %d", len(row), b.numCols))
	}
	for x := 0; x < b.numCols; x++ {
		b.builders[x].push(row[x])
	}
}

func (b *tableBuilder) pushTextRow(row []string) {
	if len(row) != b.numCols {
		panic(fmt.Errorf("size of row %d does not match size of body: %d", len(row), b.numCols))
	}
	for x := 0; x < b.numCols; x++ {
		b.builders[x].parsePush(row[x])
	}
}

func (b *tableBuilder) body() ([]Series, error) {
	if b == nil {
		return []Series{}, nil
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
