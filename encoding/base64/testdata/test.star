load('encoding/base64.star', 'base64')
load('assert.star', 'assert')

assert.eq(base64.encode("hello"), "aGVsbG8=")
assert.eq(base64.encode("hello", encoding="standard_raw"), "aGVsbG8")
assert.eq(base64.encode("hello friend!", encoding="url"), "aGVsbG8gZnJpZW5kIQ==")
assert.eq(base64.encode("hello friend!", encoding="url_raw"), "aGVsbG8gZnJpZW5kIQ")

assert.eq(base64.decode("aGVsbG8="),"hello")
assert.eq(base64.decode("aGVsbG8", encoding="standard_raw"),"hello")
assert.eq(base64.decode("aGVsbG8gZnJpZW5kIQ==", encoding="url"),"hello friend!")
assert.eq(base64.decode("aGVsbG8gZnJpZW5kIQ", encoding="url_raw"),"hello friend!")