/*
Package zipfile defines zipfileimatical functions, it's intented to be a drop-in
subset of python's zipfile module for starlark:
https://docs.python.org/3/library/zipfile.html
*/
package zipfile

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/google/skylark"
	"github.com/google/skylark/skylarkstruct"
)

// ModuleName defines the expected name for this Module when used
// in skylark's load() function, eg: load('zipfile.sky', 'zipfile')
const ModuleName = "zipfile.sky"

var (
	once          sync.Once
	zipfileModule skylark.StringDict
)

// LoadModule loads the zipfile module.
// It is concurrency-safe and idempotent.
func LoadModule() (skylark.StringDict, error) {
	once.Do(func() {
		zipfileModule = skylark.StringDict{
			"ZipFile": skylark.NewBuiltin("ZipFile", newZipFile),
		}
	})
	return zipfileModule, nil
}

// newZipfile opens a zip archive ZipFile(file, mode='r', compression=ZIP_STORED, allowZip64=True, compresslevel=None)
func newZipFile(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var file skylark.String
	if err := skylark.UnpackArgs("ZipFile", args, kwargs, "file", &file); err != nil {
		return nil, err
	}

	rdr := strings.NewReader(string(file))
	zr, err := zip.NewReader(rdr, int64(len(file)))
	if err != nil {
		return skylark.None, err
	}

	return ZipFile{zr}.Struct(), nil
}

type ZipFile struct {
	*zip.Reader
}

// Struct turns zipFile into a skylark struct value
func (zf ZipFile) Struct() *skylarkstruct.Struct {
	return skylarkstruct.FromStringDict(skylark.String("ZipFile"), skylark.StringDict{
		"namelist": skylark.NewBuiltin("namelist", zf.namelist),
		"open":     skylark.NewBuiltin("open", zf.open),
	})
}

func (zf ZipFile) namelist(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var names []skylark.Value
	for _, f := range zf.File {
		names = append(names, skylark.String(f.Name))
	}
	return skylark.NewList(names), nil
}

func (zf ZipFile) open(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var name skylark.String
	if err := skylark.UnpackArgs("open", args, kwargs, "name", &name); err != nil {
		return nil, err
	}
	n := string(name)
	for _, f := range zf.File {
		if n == f.Name {
			return ZipInfo{f}.Struct(), nil
		}
	}

	return skylark.None, fmt.Errorf("not found")
}

type ZipInfo struct {
	*zip.File
}

// Struct turns zipInfo into a skylark struct value
func (zf ZipInfo) Struct() *skylarkstruct.Struct {
	return skylarkstruct.FromStringDict(skylark.String("ZipFile"), skylark.StringDict{
		"read": skylark.NewBuiltin("read", zf.read),
	})
}

func (f ZipInfo) read(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	rc, err := f.File.Open()
	if err != nil {
		return skylark.None, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return skylark.None, err
	}

	return skylark.String(string(data)), nil
}

// func (f ZipInfo) String() string        { return time.Duration(d).String() }
// func (f ZipInfo) Type() string          { return "ZipInfo" }
// func (f ZipInfo) Freeze()               {} // TODO - ???
// func (f ZipInfo) Hash() (uint32, error) { return hashString(d.String()), nil }
// func (f ZipInfo) Truth() skylark.Bool   { return f == nil }
// func (f ZipInfo) Attr(name string) (skylark.Value, error) {
// 	return builtinAttr(d, name, durationMethods)
// }
// func (d Duration) AttrNames() []string { return builtinAttrNames(durationMethods) }

// type File struct {
// }

// type Directory struct {
// }
