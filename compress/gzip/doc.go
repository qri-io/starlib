/*Package gzip defines gzip encoding & decoding functions

  outline: gzip
    decompress files like the GNU programs gzip and gunzip would.
    path: compress/gzip
    functions:
      decompress(data) bytes
        Return a bytes object containing the uncompressed data.
				params:
					data string, bytes
						data can be a string or bytes of compressed gzip data
				examples:
					basic
						download a gzip file & decompress the contents
						code:
							load("compress/gzip.star", "gzip")
							load("http.star", "http")
							url = "http://www.mygziprepo.com/sample.gz"
							raw = http.get(url).body()
							data = gzip.decompress(raw)
*/
package gzip
