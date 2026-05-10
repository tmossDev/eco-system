package mysql

const (
	GetByEmail    = "SELECT * FROM users WHERE email = ?"
	GetByID       = "SELECT * FROM users WHERE id = ?"
	RegisterUser  = "INSERT INTO users (first_name, last_name, email, hashed_password, role_id, created_user, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	UpdateUser    = "UPDATE users SET first_name = ?, last_name = ?, updated_user = ?, updated_at = now() WHERE id = ?"
	ResetPassword = "UPDATE users SET hashed_password = ?, updated_user = ? WHERE id = ?"
	ResetEmail    = "UPDATE users SET email = ?, updated_user = ? WHERE id = ?"
)
