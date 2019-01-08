package geo

import (
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/qri-io/starlib/util"
	"go.starlark.net/starlark"
)

func parseGeoJSON(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (v starlark.Value, err error) {
	var (
		data      starlark.Value
		dataBytes []byte
		fc        *geojson.FeatureCollection
	)
	v = starlark.None

	if err = starlark.UnpackArgs("parseGeoJSON", args, kwargs, "data", &data); err != nil {
		return
	}

	switch val := data.(type) {
	case starlark.String:
		dataBytes = []byte(val)
	default:
		err = fmt.Errorf("parseGeoJSON: invalid argument type, expected string")
		return
	}

	if fc, err = geojson.UnmarshalFeatureCollection(dataBytes); err != nil {
		return
	}

	geoms := make([]starlark.Value, len(fc.Features))
	properties := make([]starlark.Value, len(fc.Features))
	for i, feat := range fc.Features {
		if geoms[i], err = geomFromOrbGeom(feat.Geometry); err != nil {
			return
		}
		if properties[i], err = dictFromGeoJSONProperties(feat.Properties); err != nil {
			return
		}
	}

	return starlark.Tuple([]starlark.Value{
		starlark.NewList(geoms),
		starlark.NewList(properties),
	}), nil
}

func geomFromOrbGeom(orbGeom orb.Geometry) (starlark.Value, error) {
	switch geom := orbGeom.(type) {
	case orb.Point:
		return Point(geom), nil
	case orb.LineString:
		line := make(Line, len(geom))
		for i, pt := range geom {
			line[i] = Point(pt)
		}
		return line, nil
	case orb.Ring:
		line := make(Line, len(geom))
		for i, pt := range geom {
			line[i] = Point(pt)
		}
		return line, nil
	case orb.Polygon:
		poly := make(Polygon, len(geom))
		for i, r := range geom {
			line := make(Line, len(r))
			for j, pt := range r {
				line[j] = Point(pt)
			}
			poly[i] = line
		}
		return poly, nil
	default:
		return starlark.None, fmt.Errorf("unrecognized geoJSON type: %s", orbGeom.GeoJSONType())
	}
}

func dictFromGeoJSONProperties(props geojson.Properties) (starlark.Value, error) {
	return util.Marshal(map[string]interface{}(props))
}
