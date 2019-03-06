/*
Package zipfile defines zipfileimatical functions, it's intended to be a drop-in
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

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// ModuleName defines the expected name for this Module when used
// in starlark's load() function, eg: load('zipfile.star', 'zipfile')
const ModuleName = "zipfile.star"

var (
	once          sync.Once
	zipfileModule starlark.StringDict
)

// LoadModule loads the zipfile module.
// It is concurrency-safe and idempotent.
func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		zipfileModule = starlark.StringDict{
			"ZipFile": starlark.NewBuiltin("ZipFile", newZipFile),
		}
	})
	return zipfileModule, nil
}

// newZipfile opens a zip archive ZipFile(file, mode='r', compression=ZIP_STORED, allowZip64=True, compresslevel=None)
func newZipFile(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var file starlark.String
	if err := starlark.UnpackArgs("ZipFile", args, kwargs, "file", &file); err != nil {
		return nil, err
	}

	rdr := strings.NewReader(string(file))
	zr, err := zip.NewReader(rdr, int64(len(file)))
	if err != nil {
		return starlark.None, err
	}

	return ZipFile{zr}.Struct(), nil
}

// ZipFile is a starlark zip file
type ZipFile struct {
	*zip.Reader
}

// Struct turns zipFile into a starlark struct value
func (zf ZipFile) Struct() *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(starlark.String("ZipFile"), starlark.StringDict{
		"namelist": starlark.NewBuiltin("namelist", zf.namelist),
		"open":     starlark.NewBuiltin("open", zf.open),
	})
}

func (zf ZipFile) namelist(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var names []starlark.Value
	for _, f := range zf.File {
		names = append(names, starlark.String(f.Name))
	}
	return starlark.NewList(names), nil
}

func (zf ZipFile) open(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name starlark.String
	if err := starlark.UnpackArgs("open", args, kwargs, "name", &name); err != nil {
		return nil, err
	}
	n := string(name)
	for _, f := range zf.File {
		if n == f.Name {
			return ZipInfo{f}.Struct(), nil
		}
	}

	return starlark.None, fmt.Errorf("not found")
}

// ZipInfo is a starlark information object for a Zip archive component
type ZipInfo struct {
	*zip.File
}

// Struct turns zipInfo into a starlark struct value
func (zi ZipInfo) Struct() *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(starlark.String("ZipFile"), starlark.StringDict{
		"read": starlark.NewBuiltin("read", zi.read),
	})
}

func (zi ZipInfo) read(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	rc, err := zi.File.Open()
	if err != nil {
		return starlark.None, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return starlark.None, err
	}

	return starlark.String(string(data)), nil
}

// func (f ZipInfo) String() string        { return time.Duration(d).String() }
// func (f ZipInfo) Type() string          { return "ZipInfo" }
// func (f ZipInfo) Freeze()               {} // TODO - ???
// func (f ZipInfo) Hash() (uint32, error) { return hashString(d.String()), nil }
// func (f ZipInfo) Truth() starlark.Bool   { return f == nil }
// func (f ZipInfo) Attr(name string) (starlark.Value, error) {
// 	return builtinAttr(d, name, durationMethods)
// }
// func (d Duration) AttrNames() []string { return builtinAttrNames(durationMethods) }

// type File struct {
// }

// type Directory struct {
// }
