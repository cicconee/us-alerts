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
func GetZones(area string) (*ZonesResponse, error) {
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

	data := ZonesResponse{Area: area, ExpiresAt: expiresAt}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

// NextDay returns the UTC time exactly 24 hours in the future from the moment the func is called.
func NextDay() time.Time {
	return time.Now().Add(time.Hour * 24).UTC()
}
