/*Package hmac defines hash primitives for starlark.

  outline: hmac
    hmac defines hmac primitives for starlark.
    path: hmac
    functions:
      md5(key, string) string
        returns an md5 hash for a string with a key of "secret"
        examples:
          basic
            calculate an md5 checksum for "hello world" with a key of "secret"
            code:
              load("hmac.star", "hmac")
              sum = hmac.md5("secret", "hello world!")
              print(sum)
              # Output: 0a0461e10e89506d7c31a145663bed93
      sha1(key, string) string
        returns a SHA1 hash for a string
        examples:
          basic
            calculate an SHA1 checksum for "hello world" with a key of "secret"
            code:
              load("hmac.star", "hmac")
              sum = hmac.sha1("secret", "hello world!")
              print(sum)
              # Output: a4df5f9d237ab0ca3241f042bcf6059a4ef491c4
      sha256(key, string) string
        returns an SHA2-256 hash for a string
        examples:
          basic
            calculate an SHA2-256 checksum for "hello world" with a key of "secret"
            code:
              load("hmac.star", "hmac")
              sum = hmac.sha256("secret", "hello world!")
              print(sum)
              # Output: 72069731bf291b463aecb218bc227abce3d403d76da67faef2d48d3cb43b2f54
*/
package hmac
