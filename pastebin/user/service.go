package user

import (
	"solid/domain"
)

type UserService struct {
	repo UserRepository // replace with interface if unit test is needed
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

type UserRepository interface {
	GetByID(id string) (*domain.User, error)
	CreateUser(user *domain.User) error
}

// ここのレシーバーはController層から呼び出されることを想定しています→Controller層にinterfaceを定義→レシーバー先はたとえばUserControllerにすることを想定
func (u *UserService) GetByID(id string) (*domain.User, error) {
	return u.repo.GetByID(id)
}

func (u *UserService) CreateUser(user *domain.User) error {
	return u.repo.CreateUser(user)
}
