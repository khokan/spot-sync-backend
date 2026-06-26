package zone

import (
	"spotsync/internal/domain/zone/dto"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) Create(req dto.CreateZoneRequest) (*dto.CreateZoneResponse, error) {
	zone := ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.repo.Create(&zone); err != nil {
		return nil, err
	}

	return &dto.CreateZoneResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     zone.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *service) Update(id uint, req dto.UpdateZoneRequest) (*dto.ZoneResponse, error) {
	zone, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		zone.Name = *req.Name
	}
	if req.Type != nil {
		zone.Type = *req.Type
	}
	if req.TotalCapacity != nil {
		zone.TotalCapacity = *req.TotalCapacity
	}
	if req.PricePerHour != nil {
		zone.PricePerHour = *req.PricePerHour
	}

	if err := s.repo.Update(zone); err != nil {
		return nil, err
	}

	zoneWithAvailability, err := s.repo.GetWithAvailabilityByID(zone.ID)
	if err != nil {
		return nil, err
	}

	return mapZone(zoneWithAvailability), nil
}

func (s *service) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *service) List() ([]dto.ZoneResponse, error) {
	zones, err := s.repo.ListWithAvailability()
	if err != nil {
		return nil, err
	}

	result := make([]dto.ZoneResponse, 0, len(zones))
	for _, z := range zones {
		mapped := mapZone(&z)
		result = append(result, *mapped)
	}

	return result, nil
}

func (s *service) GetByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.repo.GetWithAvailabilityByID(id)
	if err != nil {
		return nil, err
	}
	return mapZone(zone), nil
}

func mapZone(zone *ZoneWithAvailability) *dto.ZoneResponse {
	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		PricePerHour:   zone.PricePerHour,
		AvailableSpots: zone.AvailableSpots,
		CreatedAt:      zone.CreatedAt,
		UpdatedAt:      zone.UpdatedAt,
	}
}
