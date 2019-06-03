package xlsx

import (
	"net/http"

	"github.com/360EntSecGroup-Skylar/excelize"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// ModuleName defines the expected name for this Module when used
// in starlark's load() function, eg: load('xlsx.star', 'xlsx')
const ModuleName = "xlsx.star"

// LoadModule creates an xlsx Module
func LoadModule() (starlark.StringDict, error) {
	m := &Module{}
	ns := starlark.StringDict{
		"xlsx": m.Struct(),
	}
	return ns, nil
}

// Module joins http tools to a dataset, allowing dataset
// to follow along with http requests
type Module struct {
}

// Struct returns this module's methods as a starlark Struct
func (m *Module) Struct() *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(starlarkstruct.Default, m.StringDict())
}

// StringDict returns all module methods in a starlark.StringDict
func (m *Module) StringDict() starlark.StringDict {
	return starlark.StringDict{
		"get_url": starlark.NewBuiltin("get_url", m.GetURL),
	}
}

// GetURL gets a file for a given URL
func (m *Module) GetURL(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var url string
	if err := starlark.UnpackPositionalArgs("get_url", args, kwargs, 1, &url); err != nil {
		return nil, err
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	xlsx, err := excelize.OpenReader(res.Body)
	if err != nil {
		return nil, err
	}
	if err = res.Body.Close(); err != nil {
		return nil, err
	}

	f := &File{
		xlsx: xlsx,
	}
	return f.Struct(), nil
}

// File is an xlsx-format excel file
type File struct {
	xlsx *excelize.File
}

// Struct turns a file into a *starlark.Struct
func (f *File) Struct() *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(starlarkstruct.Default, starlark.StringDict{
		"get_sheets": starlark.NewBuiltin("get_sheets", f.GetSheets),
		"get_rows":   starlark.NewBuiltin("get_rows", f.GetRows),
	})
}

// GetSheets returns a map of ints to sheet names, sheet numbers are 1-based
func (f *File) GetSheets(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	sheets := &starlark.Dict{}
	for idx, name := range f.xlsx.GetSheetMap() {
		sheets.SetKey(starlark.MakeInt(idx), starlark.String(name))
	}
	return sheets, nil
}

// GetRows grabs rows for a given sheet
func (f *File) GetRows(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var sheet string
	if err := starlark.UnpackPositionalArgs("get_rows", args, kwargs, 1, &sheet); err != nil {
		return nil, err
	}

	xRows := f.xlsx.GetRows(sheet)

	rows := make([]starlark.Value, len(xRows))
	for i, xRow := range xRows {
		col := make([]starlark.Value, len(xRow))
		for j, xCell := range xRow {
			col[j] = starlark.String(xCell)
		}
		rows[i] = starlark.NewList(col)
	}

	return starlark.NewList(rows), nil
}
