package nws

import "time"

// zonesResponse holds the data returned from the https://api.weather.gov/zones?area={area} endpoint.
// When calling GetZones the response is decoded into zonesResponse.
type zonesResponse struct {
	// Zones corresponds to the json array features. Each feature in the array represents a zone.
	Zones []zoneFeature `json:"features"`

	// Area is the area code the zones belong to. This is typically a state abbrevation i.e. IL, AK, etc.
	Area string

	// ExpiresAt holds the date in which the data is no longer valid. This is retrieved from the Expires header.
	// The NWS api formats this date in a RFC 1123 layout.
	ExpiresAt time.Time
}

// zoneFeature holds the data for a specific zone that is returned by the NWS api.
type zoneFeature struct {
	// URI is a url to the zone. The format is https://api.weather.gov/zones/{type}/{id}.
	URI string `json:"id"`

	// Properties holds the details about the zone.
	Properties zoneFeatureProperties `json:"properties"`
}

// zoneFeatureProperties holds the details of a zone that is returned by the NWS api.
type zoneFeatureProperties struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"`
	Name           string    `json:"name"`
	EffectiveDate  time.Time `json:"effectiveDate"`
	ExpirationDate time.Time `json:"expirationDate"`
	State          string    `json:"state"`
}

// AsAreaZones converts a zonesReponse into a AreaZones.
func (z *zonesResponse) AsAreaZones() *AreaZones {
	zones := []Zone{}
	for _, z := range z.Zones {
		zones = append(zones, Zone{
			ID:             z.Properties.ID,
			Type:           z.Properties.Type,
			Name:           z.Properties.Name,
			EffectiveDate:  z.Properties.EffectiveDate,
			ExpirationDate: z.Properties.ExpirationDate,
			State:          z.Properties.State,
		})
	}

	return &AreaZones{
		Zones:     zones,
		Area:      z.Area,
		ExpiresAt: z.ExpiresAt,
	}
}

// AreaZones holds all of the zones for an area. AreaZones is created from zonesResponse and is returned when
// calling GetZones.
type AreaZones struct {
	Zones     []Zone
	Area      string
	ExpiresAt time.Time
}

// Zone holds the data for a zone. Zone is used when returning AreaZones from GetZones.
type Zone struct {
	ID             string
	Type           string
	Name           string
	EffectiveDate  time.Time
	ExpirationDate time.Time
	State          string
}
