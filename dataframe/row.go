package dataframe

import (
	"encoding/json"
	"fmt"

	"go.starlark.net/starlark"
)

type rowIter struct {
	idx   int
	df    *DataFrame
	limit int
	order []int
	step  int // only used if order exists
}

func newRowIter(df *DataFrame) *rowIter {
	return &rowIter{idx: 0, df: df, limit: df.numRows()}
}

func newRowIterWithOrder(df *DataFrame, order []int) *rowIter {
	if order == nil {
		return newRowIter(df)
	}
	idx := order[0]
	return &rowIter{idx: idx, df: df, limit: df.numRows(), order: order, step: 0}
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
		data, err := json.Marshal(r.GetRow().items)
		if err != nil {
			return "?"
		}
		return string(data)
	}
	return r.GetStr(pos)
}

func (r *rowIter) GetRow() *rowTuple {
	var items []interface{}
	items = make([]interface{}, len(r.df.body))
	for k := 0; k < len(r.df.body); k++ {
		items[k] = r.df.body[k].at(r.idx)
	}
	return &rowTuple{df: r.df, idx: r.idx, items: items}
}

func (r *rowIter) GetStr(pos int) string {
	return r.df.body[pos].strAt(r.idx)
}

func (r *rowIter) Index() int {
	return r.idx
}

func (r *rowIter) MergeWith(q *rowIter, leftKey, rightKey, ignore int) *rowTuple {
	leftElem := r.GetStr(leftKey)
	rightElem := q.GetStr(rightKey)
	if leftElem == rightElem {
		items := r.GetRow()
		other := q.GetRow().removeFromStringList(ignore)
		return &rowTuple{idx: r.idx, df: r.df, items: items.concat(other)}
	}
	return nil
}

type rowTuple struct {
	idx   int
	df    *DataFrame
	items []interface{}
}

func (rt *rowTuple) removeFromStringList(i int) *rowTuple {
	if i == -1 {
		return rt
	}
	ls := rt.items
	a := make([]interface{}, len(ls))
	copy(a, ls)
	copy(a[i:], a[i+1:])
	a[len(a)-1] = nil
	return &rowTuple{items: a[:len(a)-1], idx: rt.idx, df: rt.df}
}

func (rt *rowTuple) concat(other *rowTuple) []interface{} {
	return append(rt.items, other.items...)
}

func (rt *rowTuple) toTuple() starlark.Tuple {
	var items []string
	for _, r := range rt.items {
		v := fmt.Sprintf("%v", r)
		items = append(items, v)
	}

	rowSeries := &Series{which: typeObj, valObjs: items, index: rt.df.columnNames}
	arguments := starlark.Tuple{rowSeries}
	return arguments
}

func (rt *rowTuple) strAt(pos int) string {
	it := rt.items[pos]
	return fmt.Sprintf("%s", it)
}

type rowCollect struct {
	df      *DataFrame
	collect []Series
}

func newRowCollect(df *DataFrame) *rowCollect {
	collect := make([]Series, len(df.body))
	return &rowCollect{df: df, collect: collect}
}

func newRowCollectOfSize(df *DataFrame, size int) *rowCollect {
	collect := make([]Series, size)
	return &rowCollect{df: df, collect: collect}
}

func (rc *rowCollect) Push(rt *rowTuple) {
	items := rt.items
	for k := 0; k < len(items); k++ {
		rc.collect[k].push(items[k])
	}
}

func (rc *rowCollect) Body() []Series {
	return rc.collect
}
