package mysql

import (
	"solid/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepo() *UserRepository {
	return &UserRepository{
		db: nil, // ここではDB接続の詳細は省略しています。
	}
}

func (u *UserRepository) GetByID(id string) (*domain.User, error) {
	return MockUserHelper(), nil
}

func (u *UserRepository) CreateUser(user *domain.User) error {
	return nil
}

// もともとは分離したテスト用ファイルで、便宜上ここに定義する
func MockUserHelper() *domain.User {
	return &domain.User{
		ID:          "1",
		Name:        "Test User",
		Email:       "test@example.com",
		Password:    "123",
		IsAnonymous: false,
	}
}
