load('geo.star', 'geo')
load('assert.star','assert')

p1 = geo.point(-44.34, 33)

assert.eq(p1.x, -44.34)
assert.eq(p1.lat, -44.34)
assert.eq(p1.y, 33)
assert.eq(p1.lng, 33)

p2 = geo.point(-44, 33)
assert.eq(p1.distanceGeodesic(p2), 31742.52939277697)

planar_d = geo.point(1,1).distance(geo.point(2,1))
assert.eq(planar_d, 1)