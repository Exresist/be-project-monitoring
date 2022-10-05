package api

import (
	"net/http"

	"be-project-monitoring/internal/domain/model"
)

func (s *server) authMiddleware(toAllow ...model.UserRole) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token := r.Header.Get("Authorization")

			err := s.svc.VerifyToken(ctx, token, toAllow...)
			if err != nil {
				s.response.ErrorWithCode(w, err, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
