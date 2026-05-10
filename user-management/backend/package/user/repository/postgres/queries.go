package postgres

const (
	GetByEmail    = "SELECT * FROM users WHERE email = $1"
	GetByID       = "SELECT * FROM users WHERE id = $1"
	RegisterUser  = "INSERT INTO users (first_name, last_name, email, hashed_password, role_id, created_user, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING  id"
	UpdateUser    = "UPDATE users SET first_name = $1, last_name = $2, updated_user = $3, updated_at = now() WHERE id = $4"
	ResetPassword = "UPDATE users SET hashed_password = $1, updated_user = $2, updated_at = now() WHERE id = $3"
	ResetEmail    = "UPDATE users SET email = $1, updated_user = $2, updated_at = now() WHERE id = $3"
)
