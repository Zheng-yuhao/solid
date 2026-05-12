package domain

type User struct {
	ID          string
	Name        string
	Email       string
	Password    string
	IsAnonymous bool
}

func NewUser(id string, name string, email string, password string, isAnonymous bool) *User {
	return &User{
		ID:          id,
		Name:        name,
		Email:       email,
		Password:    password,
		IsAnonymous: isAnonymous,
	}
}
