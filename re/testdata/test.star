load("re.star", "re")
load("assert.star", "assert")

pattern = "(\w*)\s*(ADD|REM|DEL|EXT|TRF)\s*(.*)\s*(NAT|INT)\s*(.*)\s*(\(\w{2}\))\s*(.*)"
test = "EDM ADD FROM INJURED NAT Jordan BEAULIEU (DB) Western University"

assert.eq(re.match(pattern,test), [(test, "EDM", "ADD", "FROM INJURED ", "NAT", "Jordan BEAULIEU ", "(DB)", "Western University")])
assert.eq(re.sub(pattern, "", test), "")

# assert.eq()
# re.fullmatch("foo")
# re.split("foo")
# re.findall("foo")
# re.finditer("foo")
# re.sub("foo")
# re.subn("foo")
# re.escape("foo")
# re.purge("foo")