package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func (s *Server) RegisterAPI() {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(cors.AllowAll().Handler)
	mux.Use(render.SetContentType(render.ContentTypeJSON))

	mux.Group(func(r chi.Router) {
		r.Use(s.Authenticator)

		r.Route("/uploads", func(r chi.Router) {
			r.Post("/", s.UploadFile)
			r.Get("/{uploadName}", s.GetUpload)
		})

		r.Route("/me", func(r chi.Router) {
			r.Get("/", s.GetUserProfile)
			r.Get("/offers", s.GetUserOffers)
			r.Patch("/", s.UpdateUserProfile)
			r.Patch("/password", s.UpdatePassword)
			r.Get("/reviews", s.GetUserReviews)

			r.Route("/orders", func(r chi.Router) {
				r.Get("/", s.GetUserOrders)
				r.Route("/{id:[0-9]+}", func(r chi.Router) {
					r.Use(s.OrderCtx)
					r.Get("/", s.GetOrder)
					r.Group(func(r chi.Router) {
						r.Use(s.OrderAccessController)
						r.Patch("/", s.EditOrder)
						r.Delete("/", s.CancelOrder)
						r.Get("/offers", s.GetOrderOffers)
					})
				})
			})

			r.Route("/notifications", func(r chi.Router) {
				r.Get("/", s.FetchNotifications)
				r.With(s.NotificationCtx).Patch("/{id:[0-9]+}", s.UpdateNotification)
			})
		})

		r.Route("/balance", func(r chi.Router) {
			r.Get("/", s.GetBalance)
			r.Get("/history", s.GetBalanceHistory)
			r.Post("/", s.UpdateBalance)
		})

		r.Route("/orders", func(r chi.Router) {
			r.With(s.RoleController(CapCreateBuyOrders)).Post("/", s.CreateOrder)
			r.Route("/{id:[0-9]+}", func(r chi.Router) {
				r.Use(s.OrderGeneralCtx)
				r.Get("/", s.GetOrder)
				r.With(s.RoleController(CapCreateResponses)).Post("/offers", s.CreateOffer)
			})
		})

		r.Route("/reviews", func(r chi.Router) {
			r.Post("/", s.SubmitReview)
		})

		r.Route("/offers/{id:[0-9]+}", func(r chi.Router) {
			r.Use(s.OfferCtx)
			r.Use(s.OfferAccessController)
			r.Get("/", s.GetOffer)
			r.Patch("/", s.UpdateOffer)
			r.Delete("/", s.CancelOffer)
		})

		r.With(s.UserCtx).Post("/u/{alias:[a-zA-Z0-9]+}/message", s.SendMessage)
		r.Route("/chats", func(r chi.Router) {
			r.Get("/", s.GetChats)
			r.Route("/{id:[0-9]+}", func(r chi.Router) {
				r.Use(s.ChatCtx)
				r.Get("/", s.GetChatHistory)
				r.Post("/", s.SendMessage)
			})
		})
	})

	mux.With(s.Registrar).Post("/register", s.LoginUser)
	mux.Post("/login", s.LoginUser)

	mux.Route("/u/{alias:[a-zA-Z0-9]+}", func(r chi.Router) {
		r.Use(s.UserCtx)
		r.Get("/", s.GetUserProfile)
		r.Get("/reviews", s.GetUserReviews)
	})

	mux.Get("/orders", s.GetOrderList)

	s.Mux = mux
}
