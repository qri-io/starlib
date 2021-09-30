/*Package geo defines geographic operations

outline: geo
  geo defines geographic operations in two-dimensional space
  path: geo
  functions:
    Point(x,y) Point
      Point constructor, takes an x(longitude) and y(latitude) value and
      returns a Point object
      params:
        x float
          x-dimension value (longitude if using geodesic space)
        y float
          y-dimension value (latitude if using geodesic space)
      examples:
        stonehenge
          create a point at the Stonehenge prehistoric monument in the United Kingdom
          code:
            load("geo.star", "geo")
            # create a point at 51.1789° N, 1.8262° W, use negative y (latitude) value for
            # west quadrant
            stonehenge = geo.Point(51.1789, -1.8262)
            print(stonehenge)
            # Output: (51.178900,-1.826200)
    Line(points) Line
      Line constructor. Takes either an array of coordinate pairs or an array
      of point objects and returns the line that connects them. Points do not
      need to be collinear, providing a single point returns a line with a
      length of 0
      params:
        points [[]float|Point]
          list of points on the line
    Polygon(rings) Polygon
      Polygon constructor. Takes a list of lists of coordinate pairs (or point
      objects) that define the outer boundary and any holes / inner boundaries
      that represent a polygon. In GIS tradition, lists of coordinates that
      wind clockwise are filled regions and  anti-clockwise represent holes.
      params:
        rings [Line]
          list of closed lines that constitute the polygon
    MultiPolygon(polygons) MultiPolygon
      MultiPolygon constructor. MultiPolygon groups a list of polygons to
      behave like a single polygon
      params:
        polygons [Polygon]
    within(geom,polygon) bool
      Returns True if geom is entirely contained by polygon
      params:
        geom [point,line,polygon]
          maybe-inner geometry
        polygon [Polygon,MultiPolygon]
          maybe-outer polygon
    parseGeoJSON(data) (geoms, properties)
      Parses string data in IETF-7946 (GeoJSON) format (https://tools.ietf.org/html/rfc7946)
      returning a list of geometries and equal-length list of properties for each geometry
      params:
        data string
          string of GeoJSON data
      examples:
        FeatureCollection
          parse example
          code:
            load("geo.star", "geo")
            geo_json_string = """
            {
              "type": "FeatureCollection",
              "features": [{
                "type": "Feature",
                "properties": {
                  "name": "Coors Field"
                },
                "geometry": {
                  "type": "Point",
                  "coordinates": [-104.99404, 39.75621]
                }
            }, {
                "type": "Feature",
                "properties": {
                  "name": "Busch Field"
                },
                "geometry": {
                  "type": "Point",
                  "coordinates": [-104.98404, 39.74621]
                }
              }]
            }
            """
            (geoms, props) = geo.parseGeoJSON(geo_json_string)
            print(props)
            # Output: [{"name": "Coors Field"}, {"name": "Busch Field"}]

  types:
    Point
      a two-dimensional point in space
      methods:
        distance(p2) float
          Euclidean Distance to the other point
          params:
            p2  point
              point to measure distance to
        distanceGeodesic(p2) float
          Distance on the surface of a sphere with the same radius as Earth
          params:
            p2 point
              point to measure distance to
    Line
      an ordered list of points that define a line
      methods:
        length() float
          Euclidean Length
        lengthGeodesic() float
          Line length on the surface of a sphere with the same radius as Earth
    Polygon
      an ordered list of closed lines (rings) that define a shape. lists of
      coordinates that wind clockwise are filled regions and  anti-clockwise
      represent holes.
    MultiPolygon
      MultiPolygon groups a list of polygons to behave like a single polygon

*/
package geo
