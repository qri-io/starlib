package dataframe

import (
	"encoding/json"
	"fmt"
	"strconv"

	"go.starlark.net/starlark"
)

// rowIter iterates a DataFrame's body one row at a time, with an optional ordering
type rowIter struct {
	idx   int
	df    *DataFrame
	limit int
	// optionally, a rowIter can have an order to change how the rows are traversed
	order []int
	// only used if order exists
	step int
}

func newRowIter(df *DataFrame) *rowIter {
	return &rowIter{idx: 0, df: df, limit: df.NumRows()}
}

// newRowIterWithOrder returns a rowIter that iterates the rows in a specific order
func newRowIterWithOrder(df *DataFrame, order []int) *rowIter {
	if order == nil {
		return newRowIter(df)
	}
	idx := order[0]
	return &rowIter{idx: idx, df: df, limit: df.NumRows(), order: order, step: 0}
}

func (r *rowIter) Done() bool {
	return r.idx >= r.limit
}

func (r *rowIter) Next() {
	if r.order != nil {
		r.step++
		if r.step < len(r.order) {
			r.idx = r.order[r.step]
		} else {
			r.idx = r.limit
		}
		return
	}
	r.idx++
}

func (r *rowIter) Marshal(pos int) string {
	if pos == -1 {
		data, err := json.Marshal(r.GetRow().data)
		if err != nil {
			return "?"
		}
		return string(data)
	}
	return r.GetStr(pos)
}

func (r *rowIter) GetRow() *rowTuple {
	items := make([]interface{}, len(r.df.body))
	for k := 0; k < len(r.df.body); k++ {
		items[k] = r.df.body[k].At(r.idx)
	}
	return &rowTuple{index: r.df.columns, data: items}
}

func (r *rowIter) GetStr(pos int) string {
	return r.df.body[pos].StrAt(r.idx)
}

func (r *rowIter) Index() int {
	return r.idx
}

func (r *rowIter) RowSize() int {
	return len(r.GetRow().data)
}

func (r *rowIter) MergeWith(q *rowIter, leftKey, rightKey, ignore int) *rowTuple {
	leftElem := r.GetStr(leftKey)
	rightElem := q.GetStr(rightKey)
	if leftElem == rightElem {
		mine := r.GetRow()
		theirs := q.GetRow().removeFromStringList(ignore)
		return &rowTuple{index: r.df.columns, data: mine.concat(theirs)}
	}
	return nil
}

type rowTuple struct {
	index *Index
	data  []interface{}
}

func (rt *rowTuple) removeFromStringList(i int) *rowTuple {
	if i == -1 {
		return rt
	}
	ls := rt.data
	a := make([]interface{}, len(ls))
	copy(a, ls)
	copy(a[i:], a[i+1:])
	a[len(a)-1] = nil
	return &rowTuple{data: a[:len(a)-1], index: rt.index}
}

func (rt *rowTuple) concat(other *rowTuple) []interface{} {
	return append(rt.data, other.data...)
}

func (rt *rowTuple) padToSize(num int) *rowTuple {
	if len(rt.data) < num {
		pad := make([]interface{}, num-len(rt.data))
		rt.data = append(rt.data, pad...)
	}
	return rt
}

func (rt *rowTuple) toTuple() starlark.Tuple {
	var items []interface{}
	for _, r := range rt.data {
		v := fmt.Sprintf("%v", r)
		items = append(items, v)
	}

	rowSeries := &Series{which: typeObj, valObjs: items, index: rt.index}
	arguments := starlark.Tuple{rowSeries}
	return arguments
}

func (rt *rowTuple) StrAt(pos int) string {
	it := rt.data[pos]
	if num, ok := it.(int); ok {
		return strconv.Itoa(num)
	}
	return fmt.Sprintf("%s", it)
}
