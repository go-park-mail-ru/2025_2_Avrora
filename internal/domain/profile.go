package domain

// Профиль
type ProfileInfo struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	PhotoURL  string
}

// Обновляем данные профиля
type ProfileUpdate struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Photo     []byte
}

// Смена пароля
type ProfileSecurityUpdate struct {
	OldPassword string
	NewPassword string
}
