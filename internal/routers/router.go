package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/handlers"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/middlewares"
)

func NewRouter(h *handlers.Handler, m middlewares.Middlewares) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(m.GzipRequest)

	r.Route("/api", func(r chi.Router) {
		r.Post("/user/register", h.SignUp)
		r.Post("/user/login", h.SignIn)

		r.Group(func(r chi.Router) {
			r.Use(m.UserAuth)

			r.Get("/user", h.GetUser)

			//r.Get("/company/", nil)
			//r.Post("/company/", nil)
			//r.Route("/company/{id}", func(r chi.Router) {
			//	r.Get("/", nil)
			//	r.Get("/devices", nil)
			//})
			//
			//r.Post("/device/", nil)
			//r.Route("/device/{id}", func(r chi.Router) {
			//	r.Get("/", nil)
			//	r.Get("/schedule", nil)
			//	r.Post("/schedule", nil)
			//	r.Delete("/schedule", nil)
			//	r.Get("/journal", nil)
			//})
		})

		//r.Group(func(r chi.Router) {
		//	r.Get("/device/info", nil)
		//	r.Get("/device/firmware", nil)
		//	r.Get("/device/schedule", nil)
		//	r.Post("/device/state", nil)
		//})
	})

	return r
}
