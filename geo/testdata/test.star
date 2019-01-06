load('geo.star', 'geo')
load('assert.star','assert')

p1 = geo.Point(-44.34, 33)

assert.eq(p1.x, -44.34)
assert.eq(p1.lat, -44.34)
assert.eq(p1.y, 33)
assert.eq(p1.lng, 33)

p2 = geo.Point(-44, 33)
assert.eq(p1.distanceGeodesic(p2), 31742.52939277697)

planar_d = geo.Point(1,1).distance(geo.Point(2,1))
assert.eq(planar_d, 1)

line = geo.Line([[1,2], geo.Point(2,2)])

assert.eq(line.length(), 1)
assert.eq(line.lengthGeodesic(), 111251.67796723428)

p = geo.Polygon([
        # Outer boundary
        [
          [-93.515625, 54.16243396806779],
          [-99.49218749999999,42.5530802889558],
          [-72.0703125, 32.24997445586331],
          [-72.0703125, 43.83452678223682],
          [-72.0703125, 54.36775852406841],
          [-80.85937499999999, 57.326521225217064],
          [-93.515625, 54.16243396806779],
        ],

        # Hole in Polygon
        [
          [-87.1875, 49.61070993807422],
          [-87.890625, 44.59046718130883],
          [-81.2109375, 43.83452678223682],
          [-80.5078125, 48.22467264956519],
          [-87.1875, 49.61070993807422],
        ],
])


p3 = geo.Point(-65.390625,48.45835188280866)
p4 = geo.Point(-93.5,53.1)

poly1 = geo.Polygon([
      [
        [-93.515625, 54.16243396806779],
        [-99.49218749999999, 42.5530802889558],
        [-72.0703125, 32.24997445586331],
        [-72.0703125, 43.83452678223682],
        [-72.0703125, 54.36775852406841],
        [-80.85937499999999, 57.326521225217064],
        [-93.515625, 54.16243396806779],
      ]
 ])

assert.eq(geo.within(p3, poly1), False)
assert.eq(geo.within(p4, poly1), True)

combine_poly = geo.MultiPolygon([poly1, p])


geoJSONString = '''{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "properties": {
        "ufo_sightings": 20,
        "population": 40
      },
      "geometry": {
        "type": "Point",
        "coordinates": [
          -103.0078125,
          38.272688535980976
        ]
      }
    },
    {
      "type": "Feature",
      "properties": {
        "ufo_sightings": 100,
        "population": 2
      },
      "geometry": {
        "type": "Point",
        "coordinates": [
          -87.1875,
          47.989921667414194
        ]
      }
    }
  ]
}'''

geoms, properties = geo.parseGeoJSON(geoJSONString)

assert.eq(geoms[0].lat,-103.0078125)
assert.eq(geoms[0].lng, 38.272688535980976)
assert.eq(properties[1]["ufo_sightings"], 100)