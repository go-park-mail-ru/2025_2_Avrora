package domain

import "time"

// Профиль
type ProfileInfo struct {
	ID        int
	UserID    int // связываем пользователя с профилем
	FirstName string
	LastName  string
	Email     string
	Phone     string
	PhotoURL  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Обновляем данные профиля
type ProfileUpdate struct {
	FirstName string
	LastName  string
	Phone     string
	Photo     []byte
}

// Смена пароля
type ProfileSecurityUpdate struct {
	OldPassword string
	NewPassword string
}
