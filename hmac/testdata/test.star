"""
Test data for hmac module
"""
load('hmac.star', 'hmac')
load('assert.star', 'assert')

assert.eq(hmac.md5("secret", "helloworld"), "8bd4df4530c3c2cafabf6986740e44bd")
assert.eq(hmac.sha1("secret", "helloworld"), "e92eb69939a8b8c9843a75296714af611c73fb53")
assert.eq(hmac.sha256("secret", "helloworld"), "7a7c2bf41973489be3b318ad2f16c75fc875c340deecb12a3f79b28bb7135c97")