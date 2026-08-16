package middleware

import (
	"context"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/domain"
)

type ContextKey string

const (
	UserIDKey ContextKey = "user_id"
	RoleKey   ContextKey = "role"
)

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

func RoleFromContext(ctx context.Context) (domain.Role, bool) {
	role, ok := ctx.Value(RoleKey).(domain.Role)
	return role, ok
}
