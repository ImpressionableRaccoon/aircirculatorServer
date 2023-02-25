package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/service"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/utils"
)

func (m *Middlewares) DeviceAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, token, ok := r.BasicAuth()
		if !ok {
			m.h.HTTPJSONError(w, service.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		deviceID, err := uuid.Parse(id)
		if err != nil || deviceID == uuid.Nil {
			m.h.HTTPJSONError(w, service.ErrWrongDeviceID.Error(), http.StatusUnauthorized)
			return
		}

		device, err := m.s.AuthDevice(r.Context(), deviceID, token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), utils.ContextKey("device"), device)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
