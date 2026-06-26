package reservation

import "time"

type Reservation struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"index;not null"`
	ZoneID       uint      `gorm:"index;not null"`
	LicensePlate string    `gorm:"type:varchar(15);not null"`
	Status       string    `gorm:"type:varchar(20);not null;default:active;index"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (Reservation) TableName() string {
	return "reservations"
}
