package model

type LoginResponse struct {
	Jwt      string `json:"jwt"`
	ExpireAt int64  `json:"expire_at"`
}

type UserResponse struct {
	ID             uint64  `json:"id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	Email          string  `json:"email"`
	HashedPassword []byte  `json:"-"`
	RoleID         uint64  `json:"role_id"`
	CreatedUser    uint64  `json:"created_user"`
	CreatedAt      string  `json:"created_at"`
	UpdatedUser    *uint64 `json:"updated_user"`
	UpdatedAt      *string `json:"updated_at"`
	DeletedUser    *uint64 `json:"deleted_user"`
	DeletedAt      *string `json:"-"`
}
