package middleware

import (
	"net/http"

	"pvz-service/internal/handler/pkg"
)

func RequireRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctxRole, ok := r.Context().Value("role").(string)
			if !ok {
				pkg.WriteError(w, "forbidden", http.StatusForbidden)
				return
			}

			for _, role := range allowedRoles {
				if ctxRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			pkg.WriteError(w, "forbidden", http.StatusForbidden)
		})
	}
}
