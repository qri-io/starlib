package starlib

import (
	"fmt"

	"github.com/google/skylark"
	"github.com/qri-io/starlib/time"
)

// Loader presents the starlib library as a loader
func Loader(thread *skylark.Thread, module string) (dict skylark.StringDict, err error) {
	switch module {
	case time.ModuleName:
		return time.LoadTimeModule()
	}

	return nil, fmt.Errorf("invalid module")
}
