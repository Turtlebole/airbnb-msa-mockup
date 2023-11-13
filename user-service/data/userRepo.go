package data

// Transformed the only implementation we had into an interface so we can standardize the way our code works
// and allow for easy swapping between multiple implementations.
type UsersRepo interface {
	GetAll() Users
	GetUsers() Users
	AddUser(u *User)
	PutUser(u *User, id int) error
	DeleteUser(id int) error
}
