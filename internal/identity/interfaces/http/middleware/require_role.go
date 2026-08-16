package middleware

import (
	"net/http"

	"github.com/JavascriptDev347/uzum-clone-with-ddd.git/internal/identity/domain"
)

func RequireRole(allowed ...domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := RoleFromContext(r.Context())
			if !ok {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			for _, allowedRole := range allowed {
				if role == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
