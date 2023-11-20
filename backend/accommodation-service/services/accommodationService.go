package services

import (
	"accommodation-service/repositories"
)

type AccommodationService struct {
	accommodations repositories.AccommodationRepo
}

func NewAccommodationService(accommodations repositories.AccommodationRepo) (AccommodationService, error) {
	return AccommodationService{
		accommodations: accommodations,
	}, nil
}

//func (s AccommodationService) Create(ctx context.Context, ownerId, content string) (repositories.Post, error) {
//	authAny := ctx.Value("auth")
//	if authAny == nil {
//		return repositories.Post{}, repositories.ErrUnauthorized()
//	}
//	authenticated := authAny.(*repositories.User)
//	if authenticated == nil {
//		return repositories.Post{}, repositories.ErrUnauthorized()
//	}
//	ownerUuid, err := uuid.Parse(ownerId)
//	if err != nil {
//		return repositories.Post{}, repositories.ErrUnauthorized()
//	}
//	owner := repositories.User{Id: ownerUuid}
//	if !owner.Equals(*authenticated) {
//		return repositories.Post{}, repositories.ErrUnauthorized()
//	}
//	post := repositories.Post{
//		Owner:   owner,
//		Content: content,
//		Likes:   make([]domain.User, 0),
//	}
//	return s.posts.Create(post)
//}
