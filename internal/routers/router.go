package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/handlers"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/middlewares"
)

func NewRouter(handler *handlers.Handler, m middlewares.Middlewares) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(m.GzipRequest)

	r.Route("/", func(r chi.Router) {
	})

	return r
}
