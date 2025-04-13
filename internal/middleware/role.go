package middleware

import (
	"net/http"

	"pvz-service/internal/handler/pkg/response"
)

const ErrForbidden = "forbidden"

func RequireRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctxRole, ok := r.Context().Value("role").(string)
			if !ok {
				response.WriteError(w, ErrForbidden, http.StatusForbidden)
				return
			}

			for _, role := range allowedRoles {
				if ctxRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			response.WriteError(w, ErrForbidden, http.StatusForbidden)
		})
	}
}
