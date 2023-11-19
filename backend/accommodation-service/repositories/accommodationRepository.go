package repositories

import (
	"accommodation-service/models" // Import the models package
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccommodationRepository struct {
	accommodations map[primitive.ObjectID]*models.Accommodation
}

func NewAccommodationRepository() *AccommodationRepository {
	return &AccommodationRepository{
		accommodations: make(map[primitive.ObjectID]*models.Accommodation),
	}
}

func (repo *AccommodationRepository) Get(id primitive.ObjectID) (*models.Accommodation, error) {
	accommodation, found := repo.accommodations[id]
	if !found {
		return nil, errors.New("Accommodation not found")
	}
	return accommodation, nil
}

func (repo *AccommodationRepository) GetAll() []*models.Accommodation {
	var result []*models.Accommodation
	for _, acc := range repo.accommodations {
		result = append(result, acc)
	}
	return result
}

func (repo *AccommodationRepository) Create(accommodation *models.Accommodation) (*models.Accommodation, error) {
	// You may want to add database logic here in a real-world scenario
	// For now, let's assume the accommodation is directly added to the map
	repo.accommodations[accommodation.Id] = accommodation
	return accommodation, nil
}

func (repo *AccommodationRepository) Update(accommodation *models.Accommodation) error {
	// You may want to add database logic here in a real-world scenario
	// For now, let's assume the accommodation is directly updated in the map
	repo.accommodations[accommodation.Id] = accommodation
	return nil
}
