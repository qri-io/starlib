# geo
geo defines geographic operations in two-dimensional space

## Functions

#### `Line(points) Line`
Line constructor. Takes either an array of coordinate pairs or an array of point objects and returns the line that connects them. Points do not need to be collinear, providing a single point returns a line with a length of 0

**parameters:**

| name | type | description |
|------|------|-------------|
| `points` | `[[]float|Point]` | list of points on the line |


#### `MultiPolygon(polygons) MultiPolygon`
MultiPolygon constructor. MultiPolygon groups a list of polygons to behave like a single polygon

**parameters:**

| name | type | description |
|------|------|-------------|
| `polygons` | `[Polygon]` |  |


#### `Point(x,y) Point`
Point constructor, takes an x(longitude) and y(latitude) value and returns a Point object

**parameters:**

| name | type | description |
|------|------|-------------|
| `x` | `float` | x-dimension value (longitude if using geodesic space) |
| `y` | `float` | y-dimension value (latitude if using geodesic space) |


#### `Polygon(rings) Polygon`
Polygon constructor. Takes a list of lists of coordinate pairs (or point objects) that define the outer boundary and any holes / inner boundaries that represent a polygon. In GIS tradition, lists of coordinates that wind clockwise are filled regions and  anti-clockwise represent holes.

**parameters:**

| name | type | description |
|------|------|-------------|
| `rings` | `[Line]` | list of closed lines that constitute the polygon |


#### `parseGeoJSON(data) (geoms, properties)`
Parses string data in IETF-7946 format (https://tools.ietf.org/html/rfc7946) returning a list of geometries and equal-length list of properties for each geometry

**parameters:**

| name | type | description |
|------|------|-------------|
| `data` | `string` | string of GeoJSON data |


#### `within(geom,polygon) bool`
Returns True if geom is entirely contained by polygon

**parameters:**

| name | type | description |
|------|------|-------------|
| `geom` | `[point,line,polygon]` | maybe-inner geometry |
| `polygon` | `[Polygon,MultiPolygon]` | maybe-outer polygon |



## Types

### `Line`
an ordered list of points that define a line
**Methods**
#### `length() float`
Euclidean Length

#### `geodesicLength() float`
Line length on the surface of a sphere with the same radius as Earth

### `MultiPolygon`
MultiPolygon groups a list of polygons to behave like a single polygon### `Point`
a two-dimensional point in space
**Methods**
#### `distance(p2) float`
Euclidean Distance to the other point

**parameters:**

| name | type | description |
|------|------|-------------|
| `p2` | `` | point to measure distance to |


#### `distanceGeodesic(p2) float`
Distance on the surface of a sphere with the same radius as Earth

**parameters:**

| name | type | description |
|------|------|-------------|
| `p2` | `point` | point to measure distance to |


### `Polygon`
an ordered list of closed lines (rings) that define a shape. lists of coordinates that wind clockwise are filled regions and  anti-clockwise represent holes.