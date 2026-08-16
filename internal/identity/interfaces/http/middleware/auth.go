package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/security"
	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/pkg/response"
)

// Authenticate - Authorization header'dan JWT'ni oladi, tekshiradi,
// va userID/role'ni request context'iga yozadi.
func Authenticate(tokenService *security.JWTTokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: "Authorization" header'ni o'qing ("Bearer <token>" formatida)
			authHeader := r.Header.Get("Authorization")

			// TODO: header bo'sh yoki noto'g'ri formatda bo'lsa - 401 qaytaring
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				response.Error(w, http.StatusUnauthorized, "Authorization header is required")
				return
			}

			// TODO: tokenService.ValidateToken(token) chaqiring
			userID, role, err := tokenService.ValidateToken(authHeader[7:])
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			// TODO: context.WithValue bilan userID va role'ni joylang
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDKey, userID)
			ctx = context.WithValue(ctx, RoleKey, role)
			r = r.WithContext(ctx)

			// TODO: next.ServeHTTP(w, r.WithContext(...)) chaqiring
			next.ServeHTTP(w, r)
		})
	}
}
