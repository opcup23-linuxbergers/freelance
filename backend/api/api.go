package api

import (
	"io"
	"net/http"
	"os"
	"fmt"

	"git.carried.ru/opcup23/backend/models"
	"git.carried.ru/opcup23/backend/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Server struct {
	models.DB
	*chi.Mux

	DevMode     bool
	TokenSecret []byte
}

func (s *Server) LoginUser(w http.ResponseWriter, r *http.Request) {
	data := &CredentialsRequest{}

	if v := r.Context().Value("login_info"); v != nil {
		data = v.(*CredentialsRequest)
	} else {
		err := render.Bind(r, data)
		if err != nil {
			s.SendError(&w, r, err, http.StatusBadRequest, ErrorBadJSON)
			return
		}
	}

	u, err := s.DB.UserGetByEmail(data.Email)
	if err != nil {
		s.SendError(&w, r, err, http.StatusBadRequest, ErrorBadCredentials)
		return
	}

	if !utils.VerifyCredentials(u.Hash, data.Password) {
		s.SendError(&w, r, nil, http.StatusBadRequest, ErrorBadCredentials)
		return
	}

	lr := &LoginResponse{}

	lr.Token, err = utils.GenerateNewToken(u.ID, s.TokenSecret)
	if err != nil {
		s.SendError(&w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}

	render.Render(w, r, lr)
}

func (s *Server) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		u   = &models.User{}
	)

	if userCtx := r.Context().Value("auth_data"); userCtx != nil {
		u, err = s.DB.UserGet(userCtx.(*models.User).ID)
	} else {
		u, err = s.DB.UserGetByAlias(chi.URLParam(r, "alias"))
	}

	if err != nil {
		s.SendError(&w, r, err, http.StatusNotFound, ErrorNotFound)
		return
	}

	render.Render(w, r, &UserProfile{
		Alias:        u.Alias,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		About:        u.About,
		SellerRating: u.SellerRating,
		BuyerRating:  u.BuyerRating,
	})
}

func (s *Server) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	data := &UserProfile{}

	err := render.Bind(r, data)
	if err != nil {
		s.SendError(&w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	u := r.Context().Value("auth_data").(*models.User)

	err = s.DB.UserUpdate(u, data.Alias, data.FirstName, data.LastName, data.About)
	if err != nil {
		s.SendError(&w, r, nil, http.StatusConflict, ErrorNotFound)
		return
	}

	render.Render(w, r, &UserProfile{
		u.Alias,
		u.FirstName,
		u.LastName,
		u.About,
		u.SellerRating,
		u.BuyerRating,
	})
}

func (s *Server) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	data := &PasswordRequest{}

	err := render.Bind(r, data)
	if err != nil {
		s.SendError(&w, r, err, http.StatusBadRequest, ErrorBadJSON)
		return
	}

	u := r.Context().Value("auth_data").(*models.User)
	if !utils.VerifyCredentials(u.Hash, data.OldPassword) {
		s.SendError(&w, r, nil, http.StatusForbidden, ErrorBadCredentials)
		return
	}

	hash, err := utils.SaltPassword(data.NewPassword)
	if err != nil {
		s.SendError(&w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}

	err = s.DB.UserUpdatePassword(u, hash)
	if err != nil {
		s.SendError(&w, r, nil, http.StatusConflict, ErrorNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetBalance(w http.ResponseWriter, r *http.Request) {
	u, err := s.DB.UserGet(r.Context().Value("auth_data").(*models.User).ID)
	if err != nil {
		s.SendError(&w, r, err, http.StatusNotFound, ErrorNotFound)
		return
	}

	render.Render(w, r, &BalanceResponse{u.Balance})
}

func (s *Server) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	data := &BalanceRequest{}

	err := render.Bind(r, data)
	if err != nil {
		s.SendError(&w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	u := r.Context().Value("auth_data").(*models.User)

	err = s.DB.UserUpdateBalance(u, data.Type, data.Delta)
	if err != nil {
		s.SendError(&w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetBalanceHistory(w http.ResponseWriter, r *http.Request) {
	ops := r.Context().Value("auth_data").(*models.User).BalanceOperations

	hist := make([]BalanceRequest, len(ops))
	for i, v := range ops {
		hist[i] = BalanceRequest{v.Type, v.Delta, v.Timestamp}
	}

	render.Render(w, r, &BalanceHistory{hist})
}

func (s *Server) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("auth_data").(*models.User)
	render.Render(w, r, s.buildOrderListResponse(u.ID, &u.Orders))
}

func (s *Server) GetOrderList(w http.ResponseWriter, r *http.Request) {
	orders := &[]models.Order{}
	raw := r.URL.Query()
	q := raw.Get("q")
	switch {
	case raw.Get("sort") == "offer_count_asc":
		orders = s.DB.OrderSearch(q, "offer_count", "asc")
	case raw.Get("sort") == "budget_asc":
		orders = s.DB.OrderSearch(q, "price", "asc")
	case raw.Get("sort") == "offer_count_desc":
		orders = s.DB.OrderSearch(q, "offer_count", "desc")
	case raw.Get("sort") == "budget_desc":
		orders = s.DB.OrderSearch(q, "price", "desc")
	default:
		orders = s.DB.OrderSearch(q, "", "")
	}

	render.Render(w, r, s.buildOrderListResponse(0, orders))
}

func (s *Server) CreateOrder(w http.ResponseWriter, r *http.Request) {
	data := &OrderRequest{}

	err := render.Bind(r, data)
	if err != nil {
		s.SendError(&w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	if data.Name == "" {
		s.SendError(&w, r, err, http.StatusBadRequest, "blank order name")
		return
	}

	u := r.Context().Value("auth_data").(*models.User)
	err = s.DB.UserUpdateBalance(u, models.BalOpsExpense, data.Price)
	if err != nil {
		s.SendError(&w, r, err, http.StatusBadRequest, ErrorLogic)
		return
	}

	err = s.DB.OrderAdd(u.ID, data.Name, data.Description, data.Price, data.Due, data.Attachments)
	if err != nil {
		s.SendError(&w, r, err, http.StatusBadRequest, err.Error())
		s.DB.UserUpdateBalance(u, models.BalOpsRefund, data.Price)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) GetOrder(w http.ResponseWriter, r *http.Request) {
	ord := r.Context().Value("order_data").(*models.Order)

	if ctx := r.Context().Value("auth_data"); ctx != nil {
		uid := ctx.(*models.User).ID
		s.DB.OrderProcessView(ord.ID, uid)
	}

	u, err := s.DB.UserGet(ord.UserID)
	if err != nil {
		s.SendError(&w, r, nil, http.StatusConflict, ErrorNotFound)
		return
	}

	render.Render(w, r, &OrderResponse{
		ord.ID,
		u.Alias,
		u.BuyerRating,
		ord.Name,
		ord.Description,
		ord.Price,
		ord.Status,
		ord.Due,
		ord.Published,
		ord.Views,
		ord.OfferCount,
		SubmissionResponse{
			ord.Submission.Text,
			*buildAttachmentList(&ord.Submission.Attachments),
		},
		*buildAttachmentList(&ord.Attachments),
	})
}

func (s *Server) EditOrder(w http.ResponseWriter, r *http.Request) {
	ord := r.Context().Value("order_data").(*models.Order)

	data := &OrderRequest{}

	err := render.Bind(r, data)
	if err != nil {
		s.SendError(&w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	err = s.DB.OrderUpdate(ord, data.Name, data.Description, data.OfferID, data.Status, data.Attachments)
	if err != nil {
		s.SendError(&w, r, err, http.StatusConflict, ErrorNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) CancelOrder(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("order_data").(*models.Order).ID

	err := s.DB.OrderDelete(id)
	if err != nil {
		s.SendError(&w, r, err, http.StatusConflict, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetOrderOffers(w http.ResponseWriter, r *http.Request) {
	offers := r.Context().Value("order_data").(*models.Order).Offers
	render.Render(w, r, s.buildOfferListResponse(&offers))
}

func (s *Server) GetUserOffers(w http.ResponseWriter, r *http.Request) {
	offers := r.Context().Value("auth_data").(*models.User).Offers
	render.Render(w, r, s.buildOfferListResponse(&offers))
}

func (s *Server) GetOffer(w http.ResponseWriter, r *http.Request) {
	offer := r.Context().Value("offer_data").(*models.Offer)

	u, err := s.DB.UserGet(offer.UserID)
	if err != nil {
		s.SendError(&w, r, err, http.StatusConflict, ErrorUserNotFound)
		return
	}

	order, err := s.DB.OrderGet(offer.OrderID)
	if err != nil {
		s.SendError(&w, r, err, http.StatusConflict, "order not found")
		return
	}

	render.Render(w, r, &OfferResponse{
		offer.ID,
		OrderAuthor{
			offer.OrderID,
			order.Name,
		},
		u.Alias,
		u.SellerRating,
		offer.Comment,
		offer.Status,
		offer.Due,
		offer.Price,
		SubmissionResponse{
			offer.Submission.Text,
			*buildAttachmentList(&offer.Submission.Attachments),
		},
	})
}

func (s *Server) CreateOffer(w http.ResponseWriter, r *http.Request) {
	order := r.Context().Value("order_data").(*models.Order)
	if order.Status != models.OrdAvailable {
		s.SendError(&w, r, nil, http.StatusForbidden, "too late bruv, you can't send your offers anymore.")
		return
	}

	user := r.Context().Value("auth_data").(*models.User)

	if order.UserID == user.ID {
		s.SendError(&w, r, nil, http.StatusBadRequest, "can't create offers for your own orders.")
		return
	}

	data := &OfferRequest{}
	err := render.Bind(r, data)
	if err != nil {
		s.SendError(&w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	err = s.DB.OfferAdd(user.ID, order.ID, data.Comment, data.Due, data.Price)
	if err != nil {
		s.SendError(&w, r, err, http.StatusConflict, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) UpdateOffer(w http.ResponseWriter, r *http.Request) {
	offer := r.Context().Value("offer_data").(*models.Offer)
	uid := r.Context().Value("auth_data").(*models.User).ID

	data := &OfferRequest{}
	err := render.Bind(r, data)
	if err != nil {
		s.SendError(&w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	if offer.UserID != uid {
		s.DB.OfferHide(offer, data.Hidden)
		w.WriteHeader(http.StatusOK)
		return
	}

	if offer.Status == models.OffConfirmed {
		err := s.DB.OfferUpdateSubmission(offer, data.Submission.Text, data.Submission.Attachments)
		if err != nil {
			s.SendError(&w, r, err, http.StatusBadRequest, err.Error())
			return
		}

		order, err := s.DB.OrderGet(offer.OrderID)
		if err != nil {
			s.SendError(&w, r, err, http.StatusConflict, "order not founnd")
			return
		}

		err = s.DB.NotificationAdd(order.UserID, fmt.Sprintf("Order #%d. Seller updated his final comment.", order.ID))
		if err != nil {
			s.SendError(&w, r, err, http.StatusInternalServerError, "failed creating notification")
			return
		}

		return
	}

	err = s.DB.OfferUpdate(offer, data.Comment, data.Price, data.Due)
	if err != nil {
		s.SendError(&w, r, err, http.StatusConflict, ErrorNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) CancelOffer(w http.ResponseWriter, r *http.Request) {
	offer := r.Context().Value("offer_data").(*models.Offer)
	uid := r.Context().Value("auth_data").(*models.User).ID

	if offer.UserID != uid {
		s.SendError(&w, r, nil, http.StatusForbidden, "can't touch this")
		return
	}

	err := s.DB.OfferDelete(offer.ID)
	if err != nil {
		s.SendError(&w, r, err, http.StatusConflict, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) SubmitReview(w http.ResponseWriter, r *http.Request) {
	data := &ReviewRequest{}

	err := render.Bind(r, data)
	if err != nil {
		s.SendError(&w, r, err, http.StatusBadRequest, ErrorBadJSON)
		return
	}

	user := r.Context().Value("auth_data").(*models.User)

	order, err := s.DB.OrderGet(data.OrderID)
	if err != nil {
		s.SendError(&w, r, err, http.StatusNotFound, "order not found")
		return
	}

	if order.Status != models.OrdDone || s.DB.ReviewExist(order.ID, user.ID) {
		s.SendError(&w, r, nil, http.StatusForbidden, "can't post review to the order.")
		return
	}

	switch user.ID {
	case order.ContractorID:
		err := s.DB.ReviewAdd(order.ID, user.ID, order.UserID, models.RevBuyer, data.Text, data.Rating)
		if err != nil {
			s.SendError(&w, r, err, http.StatusInternalServerError, ErrorInternal)
			return
		}
		s.DB.UserUpdateBalance(user, models.BalOpsIncome, order.Price)
		err = s.DB.NotificationAdd(order.UserID, fmt.Sprintf("Order #%d. Contractor reviewed your order.", order.ID))
		if err != nil {
			s.SendError(&w, r, err, http.StatusInternalServerError, "failed creating notification")
			return
		}
	case order.UserID:
		err := s.DB.ReviewAdd(order.ID, user.ID, order.ContractorID, models.RevSeller, data.Text, data.Rating)
		if err != nil {
			s.SendError(&w, r, err, http.StatusInternalServerError, ErrorInternal)
			return
		}
		err = s.DB.NotificationAdd(order.ContractorID, fmt.Sprintf("Order #%d. Seller reviewed your work.", order.ID))
		if err != nil {
			s.SendError(&w, r, err, http.StatusInternalServerError, "failed creating notification")
			return
		}
	default:
		s.SendError(&w, r, nil, http.StatusForbidden, "can't post review to the order.")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	if ctx := r.Context().Value("user_data"); ctx != nil {
		render.Render(w, r, s.buildReviewListResponse(&ctx.(*models.User).Reviews))
	} else {
		reviews := r.Context().Value("auth_data").(*models.User).Reviews
		render.Render(w, r, s.buildReviewListResponse(&reviews))
	}
}

func (s *Server) GetChats(w http.ResponseWriter, r *http.Request) {
	chats := r.Context().Value("auth_data").(*models.User).Chats
	render.Render(w, r, buildChatListResponse(&chats))
}

func (s *Server) GetChatHistory(w http.ResponseWriter, r *http.Request) {
	chat := r.Context().Value("chat_data").(*models.Chat)
	render.Render(w, r, &ChatResponse{
		chat.ID,
		chat.Name,
		*s.buildMessageList(&chat.Messages),
	})
}

func (s *Server) SendMessage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("auth_data").(*models.User)
	chat := r.Context().Value("chat_data").(*models.Chat)

	data := &MessageRequest{}
	err := render.Bind(r, data)
	if err != nil {
		s.SendError(&w, r, nil, http.StatusBadRequest, err.Error())
		return
	}

	err = s.DB.ChatPostMessage(chat, user.ID, data.Text, data.Attachments)
	if err != nil {
		s.SendError(&w, r, err, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) FetchNotifications(w http.ResponseWriter, r *http.Request) {
	notifs := r.Context().Value("auth_data").(*models.User).Notifications

	render.Render(w, r, buildNotificationListResponse(&notifs))
}

func (s *Server) UpdateNotification(w http.ResponseWriter, r *http.Request) {
	notif := r.Context().Value("notification_data").(*models.Notification)

	data := &NotificationRequest{}
	err := render.Bind(r, data)
	if err != nil {
		s.SendError(&w, r, err, http.StatusBadRequest, ErrorBadJSON)
		return
	}

	err = s.DB.NotificationUpdate(notif, data.Hidden, data.Read)
	if err != nil {
		s.SendError(&w, r, err, http.StatusConflict, "failed to update notification.")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) UploadFile(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, utils.MaxFileSize)
	if err := r.ParseMultipartForm(utils.MaxFileSize); err != nil {
		s.SendError(&w, r, err, http.StatusBadRequest, "file is too large")
		return
	}

	ifile, head, err := r.FormFile("uploadedfile")
	if err != nil {
		s.SendError(&w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}
	defer ifile.Close()

	filename := utils.GenerateUploadName(head.Filename, s.DB.UploadGetTotal())

	ofile, err := os.Create("./data/" + filename)
	if err != nil {
		s.SendError(&w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}
	defer ofile.Close()

	if _, err := io.Copy(ofile, ifile); err != nil {
		s.SendError(&w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}

	s.DB.UploadCreate()

	uri := "/" + filename
	render.Render(w, r, &UploadResponse{URI: uri})
}

func (s *Server) GetUpload(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("./data/" + chi.URLParam(r, "uploadName"))
	if err != nil {
		s.SendError(&w, r, err, http.StatusNotFound, ErrorNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(file)
}
