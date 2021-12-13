package starlib

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	outline "github.com/b5/outline/lib"
	"github.com/qri-io/starlib/dataframe"
	"go.starlark.net/starlark"
)

func glob(root string, fn func(string) bool) []string {
	var matches []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if fn(path) {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return matches
}

var skipExamplesFiles = map[string]struct{}{
	"compress/gzip/doc.go": struct{}{},
	"zipfile/doc.go":       struct{}{},
}

func TestDocExamples(t *testing.T) {
	docFiles := glob(".", func(path string) bool {
		return strings.Contains(path, "/doc.go")
	})

	for _, path := range docFiles {
		if _, ok := skipExamplesFiles[path]; ok {
			continue
		}
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		doc, err := outline.ParseFirst(f)
		if err != nil {
			t.Fatal(err)
		}

		for _, eg := range doc.Examples() {
			t.Run(fmt.Sprintf("%s_%s", doc.Name, eg.Name), func(t *testing.T) {
				printBuf := &bytes.Buffer{}
				thread := &starlark.Thread{Load: Loader, Print: func(thread *starlark.Thread, msg string) { printBuf.WriteString(msg) }}
				dataframe.SetOutputSize(thread, 0, 0)
				g := starlark.StringDict{}
				_, err := starlark.ExecFile(thread, eg.Name, eg.Code, g)
				if err != nil {
					t.Fatal(err)
				}

				// Parse comments of the form:
				// # Output:
				lines := strings.Split(eg.Code, "\n")
				foundOutput := false
				expectBuf := &bytes.Buffer{}
				for _, line := range lines {
					if !foundOutput {
						if strings.Contains(line, "# Output: ") {
							foundOutput = true
							outSuffix := strings.TrimPrefix(line, "# Output: ")
							if len(outSuffix) > 0 {
								expectBuf.WriteString(outSuffix + "\n")
							}
						}
					}
				}

				got := strings.TrimSpace(printBuf.String())
				expect := strings.TrimSpace(expectBuf.String())
				if len(got) > 0 && len(expect) == 0 {
					t.Error("example printed output, no '# Output:' comment present")
				} else if expect != got {
					t.Errorf("output mistmatch\nwant: %s\ngot:  %s", expect, got)
				}
			})
		}
	}
}
