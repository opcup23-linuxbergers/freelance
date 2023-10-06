package api

import (
	"context"
	"net/http"
	"strconv"

	"git.carried.ru/opcup23/backend/models"
	"git.carried.ru/opcup23/backend/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// could be easily extended to add more (paid) features
const (
	CapCreateBuyOrders = "create_buy_orders"
	CapCreateResponses = "create_responses"
)

func (s *Server) RoleController(cap string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value("auth_data").(*models.User)

			if can(user, cap) {
				h.ServeHTTP(w, r)
			} else {
				http.NotFound(w, r)
			}
		})
	}
}

func can(u *models.User, c string) bool {
	for _, v := range u.Roles {
		for _, w := range v.Capabilities {
			if w.Name == c {
				return true
			}
		}
	}
	return false
}

func (s *Server) Authenticator(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := utils.GetClaimsFromToken(r.Header.Get("Authorization"), s.TokenSecret)
		if err != nil {
			s.SendError(&w, r, err, http.StatusForbidden, ErrorInvalidToken)
			return
		}

		u, err := s.DB.UserGet(int(claims["iss"].(float64)))
		if err != nil {
			s.SendError(&w, r, err, http.StatusForbidden, ErrorInvalidToken)
			return
		}

		ctx := context.WithValue(r.Context(), "auth_data", u)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) UserCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := s.DB.UserGetByAlias(chi.URLParam(r, "alias"))
		if err != nil {
			s.SendError(&w, r, err, http.StatusForbidden, ErrorInvalidToken)
			return
		}

		if u := r.Context().Value("auth_data"); u != nil {
			if u.(*models.User).ID == user.ID {
				render.Status(r, http.StatusPermanentRedirect)
				render.Render(w, r, &RedirectResponse{"/me"})
				return
			}
		}

		ctx := context.WithValue(r.Context(), "user_data", user)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) Registrar(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := &RegistrationRequest{}

		err := render.Bind(r, data)
		if err != nil {
			s.SendError(&w, r, err, http.StatusBadRequest, err.Error())
			return
		}

		passwordHash, err := utils.SaltPassword(data.Password)
		if err != nil {
			s.SendError(&w, r, err, http.StatusInternalServerError, ErrorInternal)
			return
		}

		err = s.DB.UserAdd(data.Username, passwordHash, data.Email)
		if err != nil {
			s.SendError(&w, r, err, http.StatusConflict, ErrorUsernameTaken)
			return
		}

		credentials := &CredentialsRequest{data.Email, data.Password}

		ctx := context.WithValue(r.Context(), "login_info", credentials)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) OrderCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.SendError(&w, r, err, http.StatusInternalServerError, ErrorInternal)
			return
		}

		ord, err := s.DB.OrderGet(id)
		if err != nil {
			s.SendError(&w, r, err, http.StatusNotFound, ErrorNotFound)
			return
		}

		uid := r.Context().Value("auth_data").(*models.User).ID
		if uid != ord.UserID {
			s.SendError(&w, r, nil, http.StatusNotFound, ErrorNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "order_data", ord)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) OrderGeneralCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.SendError(&w, r, err, http.StatusInternalServerError, ErrorInternal)
			return
		}

		ord, err := s.DB.OrderGet(id)
		if err != nil {
			s.SendError(&w, r, err, http.StatusNotFound, ErrorNotFound)
			return
		}

		if u := r.Context().Value("auth_data"); u != nil {
			if u.(*models.User).ID == ord.UserID {
				render.Status(r, http.StatusPermanentRedirect)
				render.Render(w, r, &RedirectResponse{"orders"})
				return
			}

			if u.(*models.User).ID != ord.ContractorID {
				ord.Attachments = nil
			}
		}

		ord.Submission = models.Submission{}	


		ctx := context.WithValue(r.Context(), "order_data", ord)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// TODO: Maybe there is a generic approach for resource ownership verification?
func (s *Server) OrderAccessController(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order := r.Context().Value("order_data").(*models.Order)
		user := r.Context().Value("auth_data").(*models.User)

		if order.Status == models.OrdDone || order.Status == models.OrdCancelled || order.UserID != user.ID {
			s.SendError(&w, r, nil, http.StatusForbidden, "can't edit this order.")
			return
		}

		h.ServeHTTP(w, r)
	})
}

func (s *Server) OfferCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offerid, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.SendError(&w, r, err, http.StatusInternalServerError, ErrorInternal)
			return
		}

		offer, err := s.DB.OfferGet(offerid)
		if err != nil {
			s.SendError(&w, r, err, http.StatusNotFound, ErrorNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "offer_data", offer)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) OfferAccessController(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offer := r.Context().Value("offer_data").(*models.Offer)
		user := r.Context().Value("auth_data").(*models.User)

		order, err := s.DB.OrderGet(offer.OrderID)
		if err != nil {
			s.SendError(&w, r, err, http.StatusConflict, err.Error())
			s.DB.OfferDelete(offer.ID)
			return
		}

		if offer.UserID != user.ID && order.UserID != user.ID {
			s.SendError(&w, r, nil, http.StatusForbidden, "can't edit this offer.")
			return
		}

		h.ServeHTTP(w, r)
	})
}

func (s *Server) ChatCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chatid, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.SendError(&w, r, err, http.StatusInternalServerError, ErrorInternal)
			return
		}

		chat, err := s.DB.ChatGet(chatid)
		if err != nil {
			s.SendError(&w, r, err, http.StatusNotFound, ErrorNotFound)
			return
		}

		uid := r.Context().Value("auth_data").(*models.User).ID
		for _, v := range chat.Users {
			if v.ID == uid {
				ctx := context.WithValue(r.Context(), "chat_data", chat)
				h.ServeHTTP(w, r.WithContext(ctx))
			}
		}

		return
	})
}

func (s *Server) NotificationCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notifid, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			s.SendError(&w, r, err, http.StatusInternalServerError, ErrorInternal)
			return
		}

		notification, err := s.DB.NotificationGet(notifid)
		if err != nil {
			s.SendError(&w, r, err, http.StatusNotFound, ErrorNotFound)
			return
		}

		uid := r.Context().Value("auth_data").(*models.User).ID
		if uid != notification.UserID {
			s.SendError(&w, r, nil, http.StatusForbidden, "not your notif")
			return
		}

		ctx := context.WithValue(r.Context(), "notification_data", notification)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
