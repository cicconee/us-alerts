package nws

import "time"

type ZonesResponse struct {
	Zones     []ZoneFeature `json:"features"`
	Area      string
	ExpiresAt time.Time
}

type ZoneFeature struct {
	URI        string                `json:"id"`
	Properties ZoneFeatureProperties `json:"properties"`
}

type ZoneFeatureProperties struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"`
	Name           string    `json:"name"`
	EffectiveDate  time.Time `json:"effectiveDate"`
	ExpirationDate time.Time `json:"expirationDate"`
	State          string    `json:"state"`
}
