package services

import (
	"accommodation-service/models"
	"accommodation-service/repositories"
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

func (s *AccommodationService) GetAccommodation(id primitive.ObjectID) (*models.Accommodation, error) {
	return s.accommodationRepo.Get(id)
}

func (s *AccommodationService) GetAllAccommodations() []*models.Accommodation {
	return s.accommodationRepo.GetAll()
}

func (s *AccommodationService) CreateAccommodation(accommodation *models.Accommodation) (*models.Accommodation, error) {
	return s.accommodationRepo.Create(accommodation)
}

func (s *AccommodationService) UpdateAccommodation(accommodation *models.Accommodation) error {
	return s.accommodationRepo.Update(accommodation)
}
