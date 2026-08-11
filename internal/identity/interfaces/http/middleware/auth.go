package middleware

import (
	"net/http"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/infrastructure/security"
)

// Authenticate - Authorization header'dan JWT'ni oladi, tekshiradi,
// va userID/role'ni request context'iga yozadi.
func Authenticate(tokenService *security.JWTTokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: "Authorization" header'ni o'qing ("Bearer <token>" formatida)
			// TODO: header bo'sh yoki noto'g'ri formatda bo'lsa - 401 qaytaring
			// TODO: tokenService.ValidateToken(token) chaqiring
			// TODO: xato bo'lsa - 401 qaytaring
			// TODO: context.WithValue bilan userID va role'ni joylang
			// TODO: next.ServeHTTP(w, r.WithContext(...)) chaqiring
		})
	}
}
