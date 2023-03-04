package nws

import "time"

type ZonesResponse struct {
	Zones     []Zone `json:"features"`
	ExpiresAt time.Time
}

type Zone struct {
	URI        string         `json:"id"`
	Properties ZoneProperties `json:"properties"`
}

type ZoneProperties struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"`
	Name           string    `json:"name"`
	EffectiveDate  time.Time `json:"effectiveDate"`
	ExpirationDate time.Time `json:"expirationDate"`
	State          string    `json:"state"`
}
