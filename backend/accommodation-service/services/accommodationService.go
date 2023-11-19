// services/accommodation_service.go
package services

import (
	"github.com/yourusername/accommodation-service/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccommodationService struct {
	accommodationRepo *repositories.AccommodationRepository
}

func NewAccommodationService(repo *repositories.AccommodationRepository) *AccommodationService {
	return &AccommodationService{
		accommodationRepo: repo,
	}
}

func (s *AccommodationService) GetAccommodation(id primitive.ObjectID) (*repositories.Accommodation, error) {
	return s.accommodationRepo.Get(id)
}

func (s *AccommodationService) GetAllAccommodations() []*repositories.Accommodation {
	return s.accommodationRepo.GetAll()
}

func (s *AccommodationService) CreateAccommodation(accommodation *repositories.Accommodation) (*repositories.Accommodation, error) {
	return s.accommodationRepo.Create(accommodation)
}

func (s *AccommodationService) UpdateAccommodation(accommodation *repositories.Accommodation) error {
	return s.accommodationRepo.Update(accommodation)
}
