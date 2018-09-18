// Package xlsx implements excel file readers in skylark
package xlsx

import (
	"net/http"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/google/skylark"
	"github.com/google/skylark/skylarkstruct"
)

// ModuleName defines the expected name for this Module when used
// in skylark's load() function, eg: load('xlsx.sky', 'xlsx')
const ModuleName = "xlsx.sky"

// LoadModule creates an xlsx Module
func LoadModule() (skylark.StringDict, error) {
	m := &Module{}
	ns := skylark.StringDict{
		"xlsx": m.Struct(),
	}
	return ns, nil
}

// Module joins http tools to a dataset, allowing dataset
// to follow along with http requests
type Module struct {
}

// Struct returns this module's methods as a skylark Struct
func (m *Module) Struct() *skylarkstruct.Struct {
	return skylarkstruct.FromStringDict(skylarkstruct.Default, m.StringDict())
}

// StringDict returns all module methods in a skylark.StringDict
func (m *Module) StringDict() skylark.StringDict {
	return skylark.StringDict{
		"get_url": skylark.NewBuiltin("get_url", m.GetURL),
	}
}

// GetURL gets a file for a given URL
func (m *Module) GetURL(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var url string
	if err := skylark.UnpackPositionalArgs("get_url", args, kwargs, 1, &url); err != nil {
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

// Struct turns a file into a *skylark.Struct
func (f *File) Struct() *skylarkstruct.Struct {
	return skylarkstruct.FromStringDict(skylarkstruct.Default, skylark.StringDict{
		"get_sheets": skylark.NewBuiltin("get_sheets", f.GetSheets),
		"get_rows":   skylark.NewBuiltin("get_rows", f.GetRows),
	})
}

// GetSheets returns a map of ints to sheet names, sheet numbers are 1-based
func (f *File) GetSheets(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	sheets := &skylark.Dict{}
	for idx, name := range f.xlsx.GetSheetMap() {
		sheets.Set(skylark.MakeInt(idx), skylark.String(name))
	}
	return sheets, nil
}

// GetRows grabs rows for a given sheet
func (f *File) GetRows(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var sheet string
	if err := skylark.UnpackPositionalArgs("get_rows", args, kwargs, 1, &sheet); err != nil {
		return nil, err
	}

	xRows := f.xlsx.GetRows(sheet)
	rows := make([]skylark.Value, len(xRows))
	for i, xRow := range xRows {
		col := make([]skylark.Value, len(xRow))
		for j, xCell := range xRow {
			col[j] = skylark.String(xCell)
		}
		rows[i] = skylark.NewList(col)
	}

	return skylark.NewList(rows), nil
}
