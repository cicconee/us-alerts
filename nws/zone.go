package nws

import "time"

type zonesResponse struct {
	Zones     []zoneFeature `json:"features"`
	Area      string
	ExpiresAt time.Time
}

type zoneFeature struct {
	URI        string                `json:"id"`
	Properties zoneFeatureProperties `json:"properties"`
}

type zoneFeatureProperties struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"`
	Name           string    `json:"name"`
	EffectiveDate  time.Time `json:"effectiveDate"`
	ExpirationDate time.Time `json:"expirationDate"`
	State          string    `json:"state"`
}

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

type AreaZones struct {
	Zones     []Zone
	Area      string
	ExpiresAt time.Time
}

type Zone struct {
	ID             string
	Type           string
	Name           string
	EffectiveDate  time.Time
	ExpirationDate time.Time
	State          string
}
