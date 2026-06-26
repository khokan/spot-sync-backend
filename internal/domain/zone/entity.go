package zone

import "time"

type ParkingZone struct {
	ID            uint      `gorm:"primaryKey"`
	Name          string    `gorm:"type:varchar(120);not null"`
	Type          string    `gorm:"type:varchar(30);not null"`
	TotalCapacity int       `gorm:"not null"`
	PricePerHour  float64   `gorm:"not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (ParkingZone) TableName() string {
	return "parking_zones"
}
