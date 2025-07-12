package domain

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string
	Email    string `gorm:"uniqueIndex"`
	Password string
	Role     string `gorm:"default:user"`
}

type UserRepository interface {
	FindByEmail(email string) (*User, error)
	Create(user *User) error
}

type AuthUseCase interface {
	Login(email, password string) (string, string, error)
	Register(name, email, password string) error
	RefreshToken(oldToken string) (string, error)
}
