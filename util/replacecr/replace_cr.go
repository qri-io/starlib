// Package replacecr defines a wrapper for replacing solo carriage return characters (\r)
// with carriage-return + line feed (\r\n)
package replacecr

import (
	"bufio"
	"io"
)

// Reader wraps an io.Reader. on every call of Read. it looks for
// for instances of lonely \r replacing them with \r\n before returning to the end consumer
// lots of files in the wild will come without "proper" line breaks, which irritates go's
// standard csv package. This'll fix by wrapping the reader passed to csv.NewReader:
// 		rdr, err := csv.NewReader(replacecr.Reader(r))
// because Reader adds '\n' characters, the number of bytes reported from the underlying
// reader can/will differ from what the underlyng reader would return
// if read from directly. This can cause issues with checksums and byte counts.
// Use with caution.
func Reader(data io.Reader) io.Reader {
	return crlfReplaceReader{
		rdr: bufio.NewReader(data),
	}
}

// crlfReplaceReader wraps a reader
type crlfReplaceReader struct {
	rdr *bufio.Reader
}

// Read implements io.Reader for crlfReplaceReader
func (c crlfReplaceReader) Read(p []byte) (n int, err error) {
	lenP := len(p)
	if lenP == 0 {
		return
	}

	for {
		if n == lenP {
			return
		}

		p[n], err = c.rdr.ReadByte()
		if err != nil {
			return
		}

		// any time we encounter \r & still have space, check to see if \n follows
		// ff next char is not \n, add it in manually
		if p[n] == '\r' && n < lenP-1 {
			if pk, err := c.rdr.Peek(1); (err == nil && pk[0] != '\n') || (err != nil && err.Error() == "EOF") {
				n++
				p[n] = '\n'
			}
		}

		n++
	}
}
