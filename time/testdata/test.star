load('time.star', 'time')
load('assert.star', 'assert')

assert.eq(time.time("2011-04-22T13:33:48Z"), time.time("2011-04-22T13:33:48Z"))
assert.eq(time.zero, time.time("0001-01-01T00:00:00Z"))
assert.true(time.time("2010-04-22T13:33:48Z") < time.time("2011-04-22T13:33:48Z"))
assert.true(time.time("2011-04-22T13:33:48Z") == time.time("2011-04-22T13:33:48Z"))
assert.true(time.time("2012-04-22T13:33:48Z") > time.time("2011-04-22T13:33:48Z"))

t = time.time("2000-01-02T03:04:05Z")
# TODO- make this a field, not a method
assert.eq(t.year(), 2000)

assert.eq(t - t, time.duration("0s"))

d = time.duration("1s")
assert.eq(d + d, time.duration("2s"))
assert.eq(d * 5, time.duration("5s"))
assert.eq(time.duration("0s") + time.duration("3m35s"), time.duration("3m35s"))