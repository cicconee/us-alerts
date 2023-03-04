package nws

import (
	"encoding/json"
)

// zoneGeoResponse holds the data returned from the https://api.weather.gov/zones/{type}/{id} endpoint.
// When calling GetZoneGeo the response is decoded into zoneGeoResponse.
type zoneGeoRespone struct {
	// Geo holds the geometric data returned by the NWS api.
	Geo geometry `json:"geometry"`
}

// geometry holds the geometric data for a zone that is returned by the NWS api.
type geometry struct {
	// Type is the geometric type. The types supported are Polygon and MultiPolygon. It is important to know that
	// a Type of Polygon does not mean a single polygon makes up the zone. There are some cases where a Polygon is
	// in fact made up of multiple polygons.
	//
	// If Type is Polygon, each index in the json array "coordinates" is a single polygon. However, the size of
	// coordinates may be larger than 1. This means multiple polygons may exist.
	//
	// If Type is MultiPolygon, each index in the json array "coordinates" is a multi polygon. Each index of a
	// multi polygon is made up of one or more polygons.
	Type string `json:"type"`

	// Coordinates holds data for either a Polygon or Multipolygon. It cannot be determined before calling NWS api
	// so a type *json.RawMessage is used when retrieving the data. Once the data is retrieved it is later casted to the
	// appropriate type.
	Coordinates *json.RawMessage `json:"coordinates"`
}

// MultiPolygon represents a list of PolygonList that make up a geographical area such as a county or forecast zone.
// Each zone is made up of one or more polygons.
type MultiPolygon []PolygonList

// PolygonList represents a list of polygons that may represent an entire, or partial zone. It is important
// to understand that some zones are made up of multiple PolygonList.
//
// If a zone is made up of a single polygon it will be represented as a PolygonList with a size of 1.
type PolygonList [][]Coordinate

// Coordinate holds a coordinate for a geometric shape. Index 0 refers to the longitude, index 1 refers to the
// latitude. Coordinate does not follow the traditional lat-long format due to the NWS api design.
type Coordinate []float64

// Long returns the longitude.
func (c Coordinate) Long() float64 {
	return c[0]
}

// Lat returns the latitude.
func (c Coordinate) Lat() float64 {
	return c[1]
}

// Geo holds the geographical data for a zone. It is important to understand that if Geo has a Type of Polygon, that
// does not mean the zone is made up of a single Polygon. It is unclear why the NWS api chose this design.
//
// The Polys func handles this confussion for you and will return a list of all the polygons making up the zone.
type Geo struct {
	// geoType is the geometric type. The types supported are Polygon and MultiPolygon.
	geoType string

	// coordinates holds the data to be casted to the appropriate type. In the Geo.Polys func, once data is casted
	// it is constructed into a PolygonList regardless of the type.
	coordinates interface{}
}

func (g *Geo) IsPoly() bool {
	return g.geoType == "Polygon"
}

func (g *Geo) IsMultiPoly() bool {
	return g.geoType == "MultiPolygon"
}

func (g *Geo) Type() string {
	return g.geoType
}

// Polys returns a PolygonList that make up the entire zone regardless if Geo.Type is Polygon or MultiPolygon.
//
// If Geo.Type is Polygon, Poly simply returns a PolygonList that holds all the polygons that make up the zone.
//
// If Geo.Type is MultiPolygon, Poly will flatten the 2d slice of polygons into a single list. The resulting
// PolygonList will hold all the polygons that make up the zone.
func (g *Geo) Polys() PolygonList {
	switch g.geoType {
	case "Polygon":
		return g.coordinates.(PolygonList)
	case "MultiPolygon":
		polys := [][]Coordinate{}

		multi := g.coordinates.(MultiPolygon)
		for _, mp := range multi {
			for _, p := range mp {
				polys = append(polys, p)
			}
		}

		return polys
	default:
		return nil
	}
}
