/*Package hash defines hash primitives for starlark.

  outline: hash
    hash defines hash primitives for starlark.
    path: hash
    functions:
      md5(string) string
        returns an md5 hash for a string
        examples:
          basic
            calculate an md5 checksum for "hello world"
            code:
              load("hash.star", "hash")
              sum = hash.md5("hello world!")
              print(sum)
              # Output: fc3ff98e8c6a0d3087d515c0473f8677
      sha1(string) string
        returns a SHA1 hash for a string
        examples:
          basic
            calculate an SHA1 checksum for "hello world"
            code:
              load("hash.star", "hash")
              sum = hash.sha1("hello world!")
              print(sum)
              # Output: 430ce34d020724ed75a196dfc2ad67c77772d169
      sha256(string) string
        returns an SHA2-256 hash for a string
        examples:
          basic
            calculate an SHA2-256 checksum for "hello world"
            code:
              load("hash.star", "hash")
              sum = hash.sha256("hello world!")
              print(sum)
              # Output: 7509e5bda0c762d2bac7f90d758b5b2263fa01ccbc542ab5e3df163be08e6ca9
*/
package hash
