package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"sync"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// ModuleName defines the expected name for this Module when used
// in starlark's load() function, eg: load('hash.star', 'hash')
const ModuleName = "hash.star"

var (
	once       sync.Once
	hashModule starlark.StringDict
	hashError  error
)

// LoadModule loads the time module.
// It is concurrency-safe and idempotent
func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		hashModule = starlark.StringDict{
			"hash": &starlarkstruct.Module{
				Name: "hash",
				Members: starlark.StringDict{
					"md5":    starlark.NewBuiltin("hash.md5", fnHash(md5.New)),
					"sha1":   starlark.NewBuiltin("hash.sha1", fnHash(sha1.New)),
					"sha256": starlark.NewBuiltin("hash.sha256", fnHash(sha256.New)),
				},
			},
		}

	})
	return hashModule, hashError

}

func fnHash(hash func() hash.Hash) func(*starlark.Thread, *starlark.Builtin, starlark.Tuple, []starlark.Tuple) (starlark.Value, error) {
	return func(t *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var s starlark.String
		if err := starlark.UnpackPositionalArgs(fn.Name(), args, kwargs, 1, &s); err != nil {
			return nil, err
		}

		h := hash()
		h.Write([]byte(string(s)))
		return starlark.String(fmt.Sprintf("%x", h.Sum(nil))), nil
	}
}
