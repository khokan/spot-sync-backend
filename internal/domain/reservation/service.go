package reservation

import (
	"errors"
	"spotsync/internal/domain/reservation/dto"
)

var (
	ErrReservationForbidden = errors.New("forbidden to modify this reservation")
	ErrReservationConflict  = errors.New("reservation cannot be cancelled")
)

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) Create(userID uint, req dto.CreateReservationRequest) (*dto.CreateReservationResponse, error) {
	reservation, err := s.repo.CreateWithCapacityLock(userID, req.ZoneID, req.LicensePlate)
	if err != nil {
		return nil, err
	}

	return &dto.CreateReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    reservation.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *service) ListMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	rows, err := s.repo.ListByUser(userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.MyReservationResponse, 0, len(rows))
	for _, row := range rows {
		result = append(result, dto.MyReservationResponse{
			ID:           row.ID,
			LicensePlate: row.LicensePlate,
			Status:       row.Status,
			Zone: dto.ZoneSummary{
				ID:   row.ZoneID,
				Name: row.ZoneName,
				Type: row.ZoneType,
			},
			CreatedAt: row.CreatedAt,
		})
	}

	return result, nil
}

func (s *service) ListAllReservations() ([]dto.AdminReservationResponse, error) {
	rows, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}

	result := make([]dto.AdminReservationResponse, 0, len(rows))
	for _, row := range rows {
		result = append(result, dto.AdminReservationResponse{
			ID:           row.ID,
			LicensePlate: row.LicensePlate,
			Status:       row.Status,
			User: dto.UserSummary{
				ID:    row.UserID,
				Name:  row.UserName,
				Email: row.UserEmail,
			},
			Zone: dto.ZoneSummary{
				ID:   row.ZoneID,
				Name: row.ZoneName,
				Type: row.ZoneType,
			},
			CreatedAt: row.CreatedAt,
		})
	}

	return result, nil
}

func (s *service) Cancel(requestingUserID uint, role string, reservationID uint) error {
	res, err := s.repo.GetByID(reservationID)
	if err != nil {
		return err
	}

	if role == "driver" && res.UserID != requestingUserID {
		return ErrReservationForbidden
	}

	if res.Status != "active" {
		return ErrReservationConflict
	}

	return s.repo.UpdateStatus(reservationID, "cancelled")
}
