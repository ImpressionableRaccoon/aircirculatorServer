package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/utils"
)

func (m *Middlewares) UserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			m.h.HttpJSONError(w, "empty authorization header", http.StatusUnauthorized)
			return
		}

		splitted := strings.Split(authorization, " ")
		if len(splitted) != 2 {
			m.h.HttpJSONError(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := splitted[1]

		user, err := m.s.ParseToken(r.Context(), token)
		if err != nil {
			m.h.HttpJSONError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), utils.ContextKey("user"), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
