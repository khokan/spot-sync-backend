package reservation

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrReservationNotFound = errors.New("reservation not found")
	ErrZoneFull            = errors.New("zone is full")
	ErrDuplicateLicensePlate = errors.New("duplicate active reservation for license plate in this zone")
)

type ReservationProjection struct {
	ID           uint
	UserID       uint
	UserName     string
	UserEmail    string
	ZoneID       uint
	ZoneName     string
	ZoneType     string
	PricePerHour float64
	LicensePlate string
	Status       string
	CreatedAt    string
}

type Repository interface {
	CreateWithCapacityLock(userID, zoneID uint, licensePlate string) (*Reservation, error)
	GetByID(id uint) (*Reservation, error)
	UpdateStatus(id uint, status string) error
	ListByUser(userID uint) ([]ReservationProjection, error)
	ListAll() ([]ReservationProjection, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

type zoneLockView struct {
	ID            uint
	TotalCapacity int
}

func (r *repository) CreateWithCapacityLock(userID, zoneID uint, licensePlate string) (*Reservation, error) {
	var created Reservation
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var zone zoneLockView
		if err := tx.Table("parking_zones").Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrReservationNotFound
			}
			return err
		}

		var activeCount int64
		if err := tx.Table("reservations").Where("zone_id = ? AND status = ?", zoneID, "active").Count(&activeCount).Error; err != nil {
			return err
		}

		if int(activeCount) >= zone.TotalCapacity {
			return ErrZoneFull
		}

		var duplicateCount int64
		if err := tx.Table("reservations").
			Where("zone_id = ? AND status = ? AND license_plate = ?", zoneID, "active", licensePlate).
			Count(&duplicateCount).Error; err != nil {
			return err
		}

		if duplicateCount > 0 {
			return ErrDuplicateLicensePlate
		}

		created = Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       "active",
		}

		if err := tx.Create(&created).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *repository) GetByID(id uint) (*Reservation, error) {
	var res Reservation
	err := r.db.First(&res, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}
	return &res, nil
}

func (r *repository) UpdateStatus(id uint, status string) error {
	result := r.db.Model(&Reservation{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrReservationNotFound
	}
	return nil
}

func (r *repository) ListByUser(userID uint) ([]ReservationProjection, error) {
	var rows []ReservationProjection
	err := r.db.Table("reservations r").
		Select(`
			r.id,
			r.user_id,
			r.zone_id,
			z.name as zone_name,
			z.type as zone_type,
			z.price_per_hour,
			r.license_plate,
			r.status,
			TO_CHAR(r.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at
		`).
		Joins("JOIN parking_zones z ON z.id = r.zone_id").
		Where("r.user_id = ?", userID).
		Order("r.created_at DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *repository) ListAll() ([]ReservationProjection, error) {
	var rows []ReservationProjection
	err := r.db.Table("reservations r").
		Select(`
			r.id,
			r.user_id,
			u.name as user_name,
			u.email as user_email,
			r.zone_id,
			z.name as zone_name,
			z.type as zone_type,
			z.price_per_hour,
			r.license_plate,
			r.status,
			TO_CHAR(r.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at
		`).
		Joins("JOIN users u ON u.id = r.user_id").
		Joins("JOIN parking_zones z ON z.id = r.zone_id").
		Order("r.created_at DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}
