package starlib

import (
	"fmt"

	"github.com/qri-io/starlib/bsoup"
	"github.com/qri-io/starlib/encoding/base64"
	"github.com/qri-io/starlib/encoding/csv"
	"github.com/qri-io/starlib/encoding/json"
	"github.com/qri-io/starlib/encoding/yaml"
	"github.com/qri-io/starlib/geo"
	"github.com/qri-io/starlib/html"
	"github.com/qri-io/starlib/http"
	"github.com/qri-io/starlib/math"
	"github.com/qri-io/starlib/re"
	"github.com/qri-io/starlib/time"
	"github.com/qri-io/starlib/xlsx"
	"github.com/qri-io/starlib/zipfile"
	"go.starlark.net/starlark"
)

// Version is the current semver for the entire starlib library
const Version = "0.4.1"

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
	case bsoup.ModuleName:
		return bsoup.LoadModule()
	case zipfile.ModuleName:
		return zipfile.LoadModule()
	case re.ModuleName:
		return re.LoadModule()
	case base64.ModuleName:
		return base64.LoadModule()
	case csv.ModuleName:
		return csv.LoadModule()
	case json.ModuleName:
		return json.LoadModule()
	case yaml.ModuleName:
		return yaml.LoadModule()
	case geo.ModuleName:
		return geo.LoadModule()
	case math.ModuleName:
		return math.LoadModule()
	}

	return nil, fmt.Errorf("invalid module '%s'", module)
}
