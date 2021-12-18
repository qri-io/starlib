package dataframe

import (
	"bufio"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCleanReader(t *testing.T) {
	prevBufferSize := cleanReaderBufferSize
	defer func() { cleanReaderBufferSize = prevBufferSize }()
	cleanReaderBufferSize = 10

	text := `a,b,c              
d,e,f                      
g,h,i        
j,k,l     
`
	cr := newCleanReader(strings.NewReader(text))

	lines := make([]string, 0)
	scanner := bufio.NewScanner(cr)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}

	expect := []string{
		"a,b,c",
		"d,e,f",
		"g,h,i",
		"j,k,l",
	}
	if diff := cmp.Diff(expect, lines); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}
}

func TestCleanReaderNoEndline(t *testing.T) {
	text := `id,animal,sound
1,cat,meow
2,dog,bark
3,eel,zap
1,cat,meow
20,dog,barks
3,eel,zap
4,frog,ribbit`

	cr := newCleanReader(strings.NewReader(text))

	lines := make([]string, 0)
	scanner := bufio.NewScanner(cr)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}

	expect := []string{
		"id,animal,sound",
		"1,cat,meow",
		"2,dog,bark",
		"3,eel,zap",
		"1,cat,meow",
		"20,dog,barks",
		"3,eel,zap",
		"4,frog,ribbit",
	}
	if diff := cmp.Diff(expect, lines); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}
}

func TestCleanReaderLongLine(t *testing.T) {
	prevBufferSize := cleanReaderBufferSize
	defer func() { cleanReaderBufferSize = prevBufferSize }()
	cleanReaderBufferSize = 4

	text := `a,b,c,d,e,f,g,h,i
j,k,l,m,n,o,p,q,r`
	cr := newCleanReader(strings.NewReader(text))

	lines := make([]string, 0)
	scanner := bufio.NewScanner(cr)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}

	expect := []string{
		"a,b,c,d,e,f,g,h,i",
		"j,k,l,m,n,o,p,q,r",
	}
	if diff := cmp.Diff(expect, lines); diff != "" {
		t.Errorf("mismatch (-want +got):%s\n", diff)
	}
}
