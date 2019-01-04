/*Package geo defines geographic operations

  outline: geo
    geo defines geographic operations
    functions:
      point(lat,lng)
        Point constructor takes an x(longitude) and y(latitude) value and returns a Point object
        params
          lat float
          lng float
      within(geomA,geomB)
        Returns True if geometry A is entirely contained by geometry B
        params:
          a [point,line,polygon]
            maybe-inner geometry
          b [point,line,polygon]
            maybe-outer geometery
      intersects(geomA,geomB)
        Similar to within but part of geometry B can lie outside of geometry A and it will still return True
    types:
      point
        methods:
          buffer(x int)
            Generates a buffered region of x units around a point
          distance(p2 point)
            Euclidian Distance
          distanceGeodesic(p2 point)
            Distance on the surface of a sphere with the same radius as Earth
          KNN()
            Given a target point T and an array of other points, return the K nearest points to T
          greatCircle(p2 point)
            Returns the great circle line segment to point 2
      line
        methods:
          buffer()
          length()
          geodesicLength()
      polygon
*/
package geo
