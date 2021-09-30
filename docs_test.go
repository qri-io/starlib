package starlib

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	outline "github.com/b5/outline/lib"
	"go.starlark.net/starlark"
)

func TestDocExamples(t *testing.T) {
	docFiles, err := filepath.Glob("./*/*doc.go")
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range docFiles {
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
