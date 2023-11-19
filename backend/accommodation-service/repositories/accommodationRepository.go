package repositories

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccommodationRepository struct {
	// You can add MongoDB client or any other necessary dependencies here
	// ...

	accommodations map[primitive.ObjectID]*Accommodation
}

func NewAccommodationRepository() *AccommodationRepository {
	return &AccommodationRepository{
		accommodations: make(map[primitive.ObjectID]*Accommodation),
	}
}

func (repo *AccommodationRepository) Get(id primitive.ObjectID) (*Accommodation, error) {
	accommodation, found := repo.accommodations[id]
	if !found {
		return nil, errors.New("Accommodation not found")
	}
	return accommodation, nil
}

func (repo *AccommodationRepository) GetAll() []*Accommodation {
	var result []*Accommodation
	for _, acc := range repo.accommodations {
		result = append(result, acc)
	}
	return result
}

func (repo *AccommodationRepository) Create(accommodation *Accommodation) (*Accommodation, error) {
	// You may want to add database logic here in a real-world scenario
	// For now, let's assume the accommodation is directly added to the map
	repo.accommodations[accommodation.Id] = accommodation
	return accommodation, nil
}

func (repo *AccommodationRepository) Update(accommodation *Accommodation) error {
	// You may want to add database logic here in a real-world scenario
	// For now, let's assume the accommodation is directly updated in the map
	repo.accommodations[accommodation.Id] = accommodation
	return nil
}
