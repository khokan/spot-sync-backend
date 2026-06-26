package dto

type CreateZoneResponse struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	TotalCapacity int     `json:"total_capacity"`
	PricePerHour  float64 `json:"price_per_hour"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type ZoneResponse struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	TotalCapacity  int     `json:"total_capacity"`
	PricePerHour   float64 `json:"price_per_hour"`
	AvailableSpots int     `json:"available_spots"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at,omitempty"`
}
