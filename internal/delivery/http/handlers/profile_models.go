package handlers

type Profile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
}

type ProfileUpdate struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
}

type UpdateEmail struct {
	Email string `json:"email"`
}

type SecurityUpdate struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}
