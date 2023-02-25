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

			r.Get("/company", h.GetUserCompanies)
			r.Post("/company", h.AddCompany)
			r.Route("/company/{id}", func(r chi.Router) {
				r.Get("/", h.GetCompany)
				r.Get("/devices", h.GetCompanyDevices)
			})

			r.Post("/device/", h.AddDevice)
			r.Route("/device/{id}", func(r chi.Router) {
				r.Get("/", h.GetDevice)
				r.Get("/schedule", h.GetSchedules)
				r.Post("/schedule", h.AddSchedules)
				r.Delete("/schedule/{schedule_id}", h.DeleteSchedule)
				r.Get("/journal", h.GetJournal)
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(m.DeviceAuth)

			r.Get("/device/info", h.GetDeviceInfo)
			r.Get("/device/schedule", h.GetDeviceSchedule)
			//r.Post("/device/state", nil)
		})

		r.Get("/device/firmware", h.GetFirmware)
	})

	return r
}
