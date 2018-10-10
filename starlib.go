package starlib

import (
	"fmt"

	starlark "github.com/google/skylark"
	"github.com/qri-io/starlib/html"
	"github.com/qri-io/starlib/http"
	"github.com/qri-io/starlib/re"
	"github.com/qri-io/starlib/time"
	"github.com/qri-io/starlib/xlsx"
	"github.com/qri-io/starlib/zipfile"
)

// Loader presents the starlib library as a loader
func Loader(thread *starlark.Thread, module string) (dict starlark.StringDict, err error) {
	switch module {
	case time.ModuleName:
		return time.LoadModule()
	case http.ModuleName:
		return http.LoadModule()
	case xlsx.ModuleName:
		return xlsx.LoadModule()
	case html.ModuleName:
		return html.LoadModule()
	case zipfile.ModuleName:
		return zipfile.LoadModule()
	case re.ModuleName:
		return re.LoadModule()
	}

	return nil, fmt.Errorf("invalid module")
}
