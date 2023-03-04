package nws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var GetZonesURL string = "https://api.weather.gov/zones?area=%s"

// GetZones calls the NWS api to get all the zones for the specified area. The area is the code name for the specified area.
// For example, to get the zones for Illinois use IL as the area.
func GetZones(area string) (*AreaZones, error) {
	resp, err := http.Get(fmt.Sprintf(GetZonesURL, area))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response code: StatusCode: %d", resp.StatusCode)
	}

	var expiresAt time.Time

	expiresHeader, ok := resp.Header["Expires"]
	// if for any reason the header is not set by the NWS api, default to 1 day
	if !ok {
		expiresAt = NextDay()
	} else {
		// NWS api formats expires header in a RFC1123 layout: Mon, 03 Apr 2023 19:08:22 GMT
		rfc1123, err := time.Parse(time.RFC1123, expiresHeader[0])
		if err != nil {
			log.Printf("nws api expires header could not be parsed as RFC1123: defaulting to next day: %v\n", err)
			expiresAt = NextDay()
		} else {
			expiresAt = rfc1123.UTC()
		}
	}

	data := zonesResponse{Area: area, ExpiresAt: expiresAt}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.AsAreaZones(), nil
}

// NextDay returns the UTC time exactly 24 hours in the future from the moment the func is called.
func NextDay() time.Time {
	return time.Now().Add(time.Hour * 24).UTC()
}

var GetZoneURL = "https://api.weather.gov/zones/%s/%s"

// GetZoneGeo calls the NWS api to get the geometric data for the specified zone. The zone type (county, forecast, fire, etc.)
// must also be specified.
func GetZoneGeo(zoneID string, zoneType string) (*Geo, error) {
	resp, err := http.Get(fmt.Sprintf(GetZoneURL, zoneType, zoneID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data zoneGeoRespone
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	geo := Geo{geoType: data.Geo.Type}

	var gErr error
	switch geo.geoType {
	case "Polygon":
		gErr = unmarshalGeo[PolygonList](data.Geo.Coordinates, &geo)
	case "MultiPolygon":
		gErr = unmarshalGeo[MultiPolygon](data.Geo.Coordinates, &geo)
	}

	if gErr != nil {
		return nil, gErr
	}

	return &geo, nil
}

// unmarshalGeo will unmarshal json to a supported type and store it into a *Geo.Coordinates.
func unmarshalGeo[T PolygonList | MultiPolygon](data *json.RawMessage, dest *Geo) error {
	var t T
	err := json.Unmarshal(*data, &t)
	if err != nil {
		return err
	}

	dest.coordinates = t
	return nil
}
