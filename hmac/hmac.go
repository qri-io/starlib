package hmac

import (
	"crypto/hmac"
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
// in starlark's load() function, eg: load('hmac.star', 'hmac')
const ModuleName = "hmac.star"

var (
	once       sync.Once
	hmacModule starlark.StringDict
	hmacError  error
)

// LoadModule loads the time module.
// It is concurrency-safe and idempotent
func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		hmacModule = starlark.StringDict{
			"hmac": &starlarkstruct.Module{
				Name: "hmac",
				Members: starlark.StringDict{
					"md5":    starlark.NewBuiltin("hmac.md5", fnHmac(md5.New)),
					"sha1":   starlark.NewBuiltin("hmac.sha1", fnHmac(sha1.New)),
					"sha256": starlark.NewBuiltin("hmac.sha256", fnHmac(sha256.New)),
				},
			},
		}

	})
	return hmacModule, hmacError

}

func fnHmac(hashFunc func() hash.Hash) func(*starlark.Thread, *starlark.Builtin, starlark.Tuple, []starlark.Tuple) (starlark.Value, error) {
	return func(t *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var (
			key starlark.String
			s 	starlark.String
		)
		if err := starlark.UnpackPositionalArgs(fn.Name(), args, kwargs, 1, &key, &s); err != nil {
			return nil, err
		}

		h :=  hmac.New(hashFunc, []byte(string(key)))

		if _, err := h.Write([]byte(string(s))); err != nil {
			return starlark.None, err
		}
		return starlark.String(fmt.Sprintf("%x", h.Sum(nil))), nil
	}
}
