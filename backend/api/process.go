package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"git.carried.ru/opcup23/backend/models"
	"git.carried.ru/opcup23/backend/utils"
	"github.com/go-chi/chi/v5"
)

func (l *LoginResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (u *UserProfile) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (u *BalanceResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (u *UploadResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (p *PasswordRequest) Bind(_ *http.Request) error {
	if !utils.ValidatePassword(p.NewPassword) {
		return errors.New(ErrorBadPassword)
	}

	return nil
}

func (u *BalanceRequest) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (s *SubmissionResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (s *SubmissionRequest) Bind(_ *http.Request) error {
	return nil
}

func (u *BalanceHistory) Render(_ http.ResponseWriter, _ *http.Request) error {
	if u.History == nil {
		return errors.New("history is blank.")
	}

	return nil
}

func (u *BalanceRequest) Bind(_ *http.Request) error {
	if u.Type != models.BalOpsDeposit && u.Type != models.BalOpsWithdraw {
		return errors.New("unknown operation type.")
	}

	if u.Delta < 0 {
		u.Delta = -u.Delta
	}

	return nil
}

func (u *CredentialsRequest) Bind(_ *http.Request) error {
	return nil
}

func (r *RegistrationRequest) Bind(_ *http.Request) error {
	if !utils.ValidateEmail(r.Email) {
		return errors.New("that's not an email")
	}

	r.Username = strings.ToLower(r.Username)
	if !utils.ValidateUsername(r.Username) {
		return errors.New("username should be < 18 chars and only contain alphanumerics")
	}

	if !utils.ValidatePassword(r.Password) {
		return errors.New(ErrorBadPassword)
	}

	return nil
}

func (u *UserProfile) Bind(_ *http.Request) error {
	if !utils.ValidateName(u.FirstName) || !utils.ValidateName(u.LastName) {
		return errors.New("name is invalid.")
	}

	return nil
}

func (o *OrderRequest) Bind(_ *http.Request) error {
	if o.Price < 0 {
		return errors.New("negative order price.")
	}

	return nil
}

func (o *OrderResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (o *OrderListResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (s *Server) buildOrderListResponse(uid int, orders *[]models.Order) *OrderListResponse {
	list := make([]OrderResponse, 0)
	for _, v := range *orders {
		u, err := s.DB.UserGet(v.UserID)
		if err != nil {
			continue
		}

		sub := models.Submission{}
		if uid == v.UserID {
			sub = v.Submission
		}

		list = append(list, OrderResponse{
			ID:          v.ID,
			Alias:       u.Alias,
			BuyerRating: u.BuyerRating,
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
			Status:      v.Status,
			Due:         v.Due,
			Published:   v.Published,
			Views:       v.Views,
			OfferCount:  v.OfferCount,
			Submission:  SubmissionResponse{
				sub.Text,
				*buildAttachmentList(&sub.Attachments),
			},
		})
	}

	return &OrderListResponse{list}
}

func (s *Server) buildOfferListResponse(offers *[]models.Offer) *OfferListResponse {
	list := make([]OfferResponse, 0)
	for _, v := range *offers {
		if v.Hidden {
			continue
		}

		u, err := s.DB.UserGet(v.UserID)
		if err != nil {
			continue
		}

		o, err := s.DB.OrderGet(v.OrderID)
		if err != nil {
			continue
		}

		list = append(list, OfferResponse{
			ID:    v.ID,
			Alias: u.Alias,
			OrderAuthor: OrderAuthor{
				v.OrderID,
				o.Name,
			},
			SellerRating: u.SellerRating,
			Price:        v.Price,
			Due:          v.Due,
			Comment:      v.Comment,
			Status:       v.Status,
			Submission: SubmissionResponse{
				v.Submission.Text,
				*buildAttachmentList(&v.Submission.Attachments),
			},
		})
	}

	return &OfferListResponse{list}
}

func (o *OfferResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (resp *RedirectResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	if resp.URI != "/me" {
		resp.URI = fmt.Sprintf("/me/%s/%s", resp.URI, chi.URLParam(r, "id"))
	}

	return nil
}

func (o *OfferRequest) Bind(_ *http.Request) error {
	return nil
}

func (o *OfferListResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (r *ReviewRequest) Bind(_ *http.Request) error {
	if r.Rating < 1 || r.Rating > 5 {
		return errors.New("rating should be from 1 to 5")
	}

	return nil
}

func (n *NotificationRequest) Bind(_ *http.Request) error {
	return nil
}

func (n *NotificationResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (n *NotificationListResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func buildNotificationListResponse(notifs *[]models.Notification) *NotificationListResponse {
	resp := &NotificationListResponse{}
	list := make([]NotificationResponse, 0)
	for _, v := range *notifs {
		if v.Hidden {
			continue
		}

		if !v.Read {
			resp.UnreadCount++
		}

		list = append(list, NotificationResponse{
			ID:        v.ID,
			Message:   v.Message,
			Timestamp: v.Timestamp,
			Read:      v.Read,
		})
	}

	resp.Notifications = list

	return resp
}

func buildAttachmentList(attachments *[]models.Attachment) *[]string {
	list := make([]string, len(*attachments))
	for i, v := range *attachments {
		list[i] = v.URI
	}

	return &list
}

func (c *ChatListResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (c *ChatResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (m *MessageResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (m *MessageRequest) Bind(_ *http.Request) error {
	if len(m.Attachments) == 0 && m.Text == "" {
		return errors.New("can't send blank message")
	}

	return nil
}

func buildChatListResponse(chats *[]*models.Chat) *ChatListResponse {
	list := make([]ChatResponse, 0)
	for _, v := range *chats {
		list = append(list, ChatResponse{
			ID: v.ID,
			Name: v.Name,
		})
	}

	return &ChatListResponse{list}
}

func (s *Server) buildChatResponse(chat *models.Chat) *ChatResponse {
	resp := &ChatResponse{
		ID: chat.ID,
		Name: chat.Name,
	}

	resp.Messages = *s.buildMessageList(&chat.Messages)

	return resp
}

func (s *Server) buildMessageList(messages *[]models.Message) *[]MessageResponse {
	list := make([]MessageResponse, 0)
	for _, v := range *messages {
		u, err := s.DB.UserGet(v.UserID)
		if err != nil {
			continue
		}

		list = append(list, MessageResponse{
			Sender: u.Alias,
			Timestamp: v.Timestamp,
			Text: v.Text,
			Attachments: *buildAttachmentList(&v.Attachments),
		})
	}

	return &list
}

func (s *Server) buildReviewListResponse(reviews *[]models.Review) *ReviewListResponse {
	list := make([]ReviewResponse, 0)
	for _, v := range *reviews {
		u, err := s.DB.UserGet(v.SubmitterID)
		if err != nil {
			continue
		}

		list = append(list, ReviewResponse{
			Type:   v.Type,
			Text:   v.Text,
			Rating: v.Rating,
			Name: u.Alias,
		})
	}

	return &ReviewListResponse{list}
}

func (n *ReviewResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (n *ReviewListResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
