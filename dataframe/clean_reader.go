package dataframe

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
)

// cleanReader reads and strips trailing whitespace from each line, and replaces
// crlf with just linefeed (\n)
type cleanReader struct {
	rdr   *bufio.Reader
	chunk []byte
	pend  []byte
	nPend int
}

var (
	trailingws = regexp.MustCompile(` +\r?\n`)
	// can be overridden for tests
	cleanReaderBufferSize = 2048
)

func newCleanReader(r io.Reader) *cleanReader {
	rdr := bufio.NewReaderSize(r, cleanReaderBufferSize)
	return &cleanReader{
		rdr:   rdr,
		chunk: make([]byte, cleanReaderBufferSize),
		pend:  make([]byte, cleanReaderBufferSize),
	}
}

func (cr *cleanReader) Read(bs []byte) (int, error) {
	buff := make([]byte, cr.nPend)
	copy(buff, cr.pend[:cr.nPend])

	var endErr error
	for {
		n, err := cr.rdr.Read(cr.chunk)
		if err == io.EOF {
			// if the reader is done, handle any final pending bytes
			if cr.nPend == 0 {
				return n, err
			}
			cr.nPend = 0
			endErr = err
			break
		} else if err != nil {
			return n, err
		}

		// find the last endline char, or read more if none found
		pos := bytes.LastIndex(cr.chunk[:n], []byte{10})
		if pos == -1 {
			buff = append(buff, cr.chunk[:n]...)
			continue
		}

		// save pending bytes after the last endline
		buff = append(buff, cr.chunk[:pos+1]...)
		copy(cr.pend, cr.chunk[pos+1:n])
		cr.nPend = n - pos - 1
		break
	}

	// remove trailing whitespace from each line
	lines := string(buff)
	trimmed := []byte(trailingws.ReplaceAllString(lines, "\n"))
	copy(bs, trimmed)
	// if the reader is done, return the proper error
	return len(trimmed), endErr
}
