package starlib

import (
	"fmt"

	"github.com/google/skylark"
	"github.com/qri-io/starlib/html"
	"github.com/qri-io/starlib/http"
	"github.com/qri-io/starlib/time"
	"github.com/qri-io/starlib/xlsx"
)

// Loader presents the starlib library as a loader
func Loader(thread *skylark.Thread, module string) (dict skylark.StringDict, err error) {
	switch module {
	case time.ModuleName:
		return time.LoadModule()
	case http.ModuleName:
		return http.LoadModule()
	case xlsx.ModuleName:
		return xlsx.LoadModule()
	case html.ModuleName:
		return html.LoadModule()
	}

	return nil, fmt.Errorf("invalid module")
}
