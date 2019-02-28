package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/qri-io/dataset/dsio/replacecr"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// ModuleName defines the expected name for this Module when used
// in starlark's load() function, eg: load('csv.star', 'csv')
const ModuleName = "encoding/csv.star"

var (
	once      sync.Once
	csvModule starlark.StringDict
)

// LoadModule loads the base64 module.
// It is concurrency-safe and idempotent.
func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		csvModule = starlark.StringDict{
			"csv": starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
				"read_all": starlark.NewBuiltin("read_all", ReadAll),
			}),
		}
	})
	return csvModule, nil
}

// ReadAll gets all values from a csv source
func ReadAll(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		r io.Reader

		source                       starlark.Value
		lazyQuotes, trimLeadingSpace starlark.Bool
		skip                         = starlark.MakeInt(0)
		fieldsPerRecord              = starlark.MakeInt(0)
		_comma, _comment             starlark.String
	)
	err := starlark.UnpackArgs("read_all", args, kwargs,
		"source", &source,
		"comma?", &_comma,
		"comment", &_comment,
		"lazy_quotes", &lazyQuotes,
		"trim_leading_space", &trimLeadingSpace,
		"fields_per_record", &fieldsPerRecord,
		"skip", &skip)

	if err != nil {
		return nil, err
	}

	switch source.Type() {
	case "string":
		str := string(source.(starlark.String))
		r = strings.NewReader(str)
	}
	csvr := csv.NewReader(replacecr.Reader(r))
	csvr.LazyQuotes = bool(lazyQuotes)
	csvr.TrimLeadingSpace = bool(trimLeadingSpace)

	comma := string(_comma)
	if comma == "" {
		comma = ","
	} else if len(comma) != 1 {
		return starlark.None, fmt.Errorf("expected comma param to be a single-character string")
	}
	csvr.Comma = []rune(comma)[0]

	comment := string(_comment)
	if comment != "" && len(comment) != 1 {
		return starlark.None, fmt.Errorf("expected comment param to be a single-character string")
	} else if comment != "" {
		csvr.Comment = []rune(comment)[0]
	}

	if fpr, ok := fieldsPerRecord.Int64(); ok && fpr != 0 {
		csvr.FieldsPerRecord = int(fpr)
	}

	if s, ok := skip.Int64(); ok {
		for i := 0; i < int(s); i++ {
			csvr.Read()
		}
	}

	strs, err := csvr.ReadAll()
	if err != nil {
		return starlark.None, err
	}

	vals := make([]starlark.Value, len(strs))
	for i, rowStr := range strs {
		row := make([]starlark.Value, len(rowStr))
		for j, cell := range rowStr {
			row[j] = starlark.String(cell)
		}
		vals[i] = starlark.NewList(row)
	}
	return starlark.NewList(vals), nil
}
