package middleware

type ContextKey string

const (
	UserIDKey ContextKey = "user_id"
	RoleKey   ContextKey = "role"
)
