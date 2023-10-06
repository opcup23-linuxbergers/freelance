package api

type CredentialsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type BalanceRequest struct {
	Type      string  `json:"operation_type"`
	Delta     float64 `json:"delta"`
	Timestamp int64   `json:"timestamp"`
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}

type BalanceHistory struct {
	History []BalanceRequest `json:"history"`
}

type PasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type RedirectResponse struct {
	URI string `json:"uri"`
}

type UserProfile struct {
	Alias        string  `json:"alias"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	About        string  `json:"about"`
	SellerRating float64 `json:"seller_rating"`
	BuyerRating  float64 `json:"buyer_rating"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type OrderRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	OfferID     int      `json:"offer_id"`
	Status      string   `json:"status"`
	Due         int64    `json:"due"`
	Attachments []string `json:"attachments"`
}

type OrderResponse struct {
	ID          int                `json:"order_id"`
	Alias       string             `json:"user"`
	BuyerRating float64            `json:"rating"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Price       float64            `json:"price"`
	Status      string             `json:"status"`
	Due         int64              `json:"due"`
	Published   int64              `json:"published"`
	Views       int64              `json:"views"`
	OfferCount  int64              `json:"offer_count"`
	Submission  SubmissionResponse `json:"final_comment"`
	Attachments []string           `json:"attachments"`
}

type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
}

type OfferRequest struct {
	Comment    string            `json:"comment"`
	Price      float64           `json:"price"`
	Due        int64             `json:"due"`
	Hidden     bool              `json:"hidden"`
	Submission SubmissionRequest `json:"final_comment"`
}

type OrderAuthor struct {
	OrderID   int    `json:"id,omitempty"`
	OrderName string `json:"name,omitempty"`
}

type OfferResponse struct {
	ID           int                `json:"id"`
	OrderAuthor  OrderAuthor        `json:"order"`
	Alias        string             `json:"user"`
	SellerRating float64            `json:"rating"`
	Comment      string             `json:"comment"`
	Status       string             `json:"status,omitempty"`
	Due          int64              `json:"due"`
	Price        float64            `json:"price"`
	Submission   SubmissionResponse `json:"final_comment,omitempty"`
}

type UploadResponse struct {
	URI string `json:"uri"`
}

type SubmissionRequest struct {
	Text        string   `json:"text"`
	Attachments []string `json:"attachments"`
}

type SubmissionResponse struct {
	Text        string   `json:"text"`
	Attachments []string `json:"attachments"`
}

type OfferListResponse struct {
	Offers []OfferResponse `json:"offers"`
}

type ReviewRequest struct {
	OrderID int    `json:"order_id"`
	Rating  int    `json:"rating"`
	Text    string `json:"text"`
}

type ReviewResponse struct {
	Type   string `json:"review_type"`
	Rating int    `json:"rating"`
	Text   string `json:"text"`
	Name string `json:"name"`
}

type ReviewListResponse struct {
	Reviews []ReviewResponse `json:"reviews"`
}

type NotificationRequest struct {
	Hidden bool `json:"hidden"`
	Read   bool `json:"read"`
}

type NotificationResponse struct {
	ID        int    `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Message   string `json:"message"`
	Read      bool   `json:"read"`
}

type NotificationListResponse struct {
	UnreadCount int `json:"unread_count"`
	Notifications []NotificationResponse `json:"notifications"`
}

type MessageResponse struct {
	Sender string `json:"sender"`
	Timestamp int64 `json:"timestamp"`
	Text string `json:"text"`
	Attachments []string `json:"attachments"`
}

type MessageRequest struct {
	Text string `json:"text"`
	Attachments []string `json:"attachments"`
}

type ChatResponse struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Messages []MessageResponse `json:"messages"`
}

type ChatListResponse struct {
	Chats []ChatResponse `json:"chats"`
}
