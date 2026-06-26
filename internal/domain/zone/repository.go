package zone

import (
	"errors"

	"gorm.io/gorm"
)

var ErrZoneNotFound = errors.New("zone not found")

type ZoneWithAvailability struct {
	ID             uint
	Name           string
	Type           string
	TotalCapacity  int
	PricePerHour   float64
	AvailableSpots int
	CreatedAt      string
	UpdatedAt      string
}

type Repository interface {
	Create(zone *ParkingZone) error
	GetByID(id uint) (*ParkingZone, error)
	Update(zone *ParkingZone) error
	Delete(id uint) error
	ListWithAvailability() ([]ZoneWithAvailability, error)
	GetWithAvailabilityByID(id uint) (*ZoneWithAvailability, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(zone *ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *repository) GetByID(id uint) (*ParkingZone, error) {
	var zone ParkingZone
	err := r.db.First(&zone, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}
	return &zone, nil
}

func (r *repository) Update(zone *ParkingZone) error {
	return r.db.Save(zone).Error
}

func (r *repository) Delete(id uint) error {
	result := r.db.Delete(&ParkingZone{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrZoneNotFound
	}
	return nil
}

func (r *repository) ListWithAvailability() ([]ZoneWithAvailability, error) {
	var rows []ZoneWithAvailability
	err := r.db.Table("parking_zones z").
		Select(`
			z.id,
			z.name,
			z.type,
			z.total_capacity,
			z.price_per_hour,
			GREATEST(z.total_capacity - COUNT(r.id), 0) AS available_spots,
			TO_CHAR(z.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS created_at,
			TO_CHAR(z.updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS updated_at
		`).
		Joins("LEFT JOIN reservations r ON r.zone_id = z.id AND r.status = ?", "active").
		Group("z.id, z.name, z.type, z.total_capacity, z.price_per_hour, z.created_at, z.updated_at").
		Order("z.id ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *repository) GetWithAvailabilityByID(id uint) (*ZoneWithAvailability, error) {
	var row ZoneWithAvailability
	err := r.db.Table("parking_zones z").
		Select(`
			z.id,
			z.name,
			z.type,
			z.total_capacity,
			z.price_per_hour,
			GREATEST(z.total_capacity - COUNT(r.id), 0) AS available_spots,
			TO_CHAR(z.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS created_at,
			TO_CHAR(z.updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') AS updated_at
		`).
		Joins("LEFT JOIN reservations r ON r.zone_id = z.id AND r.status = ?", "active").
		Where("z.id = ?", id).
		Group("z.id, z.name, z.type, z.total_capacity, z.price_per_hour, z.created_at, z.updated_at").
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, ErrZoneNotFound
	}
	return &row, nil
}
