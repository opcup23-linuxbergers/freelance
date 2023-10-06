package models

type User struct {
	ID                int
	Alias             string `gorm:"unique"`
	Email             string `gorm:"unique"`
	Hash              []byte
	FirstName         string
	LastName          string
	About             string
	BuyerRating       float64
	SellerRating      float64
	Balance           float64
	BalanceOperations []BalanceOperation
	Roles             []*Role `gorm:"many2many:user_roles"`
	Chats             []*Chat  `gorm:"many2many:user_chats"`
	Notifications     []Notification
	Orders            []Order
	Offers            []Offer
	Reviews           []Review
}

type Chat struct {
	ID       int
	Name     string
	Users    []*User `gorm:"many2many:user_chats"`
	Messages []Message
}

type Message struct {
	ID         int
	UserID     int
	ChatID     int
	Text       string
	Timestamp  int64
	Attachments []Attachment `gorm:"polymorphic:Source;"`
}

type Notification struct {
	ID        int
	UserID    int
	Timestamp int64
	Message   string
	Hidden    bool
	Read      bool
}

const (
	OrdDone       = "done"
	OrdAvailable  = "available"
	OrdInProgress = "in_progress"
	OrdCancelled  = "cancelled"
)

type Order struct {
	ID           int
	Name         string
	Description  string
	Price        float64
	Views        int64
	OfferCount   int64
	Status       string
	UserID       int
	ContractorID int
	Submission   Submission
	Offers       []Offer
	Due          int64
	Published    int64
	Attachments  []Attachment `gorm:"polymorphic:Source;"`
}

type Submission struct {
	ID          int
	OfferID     int
	OrderID     int
	Text        string
	Attachments []Attachment `gorm:"polymorphic:Source;"`
}

type View struct {
	UserID  int `gorm:"primaryKey;autoIncrement:false"`
	OrderID int `gorm:"primaryKey;autoIncrement:false"`
}

const (
	OffSubmitted = "submitted"
	OffConfirmed = "confirmed"
	OffDone      = "done"
)

type Offer struct {
	ID         int
	OrderID    int
	BuyerID    int
	UserID     int
	Due        int64
	Price      float64
	Comment    string
	Status     string
	Hidden     bool
	Submission Submission
}

const (
	RevBuyer  = "buyer"
	RevSeller = "seller"
)

type Review struct {
	ID          int
	Type        string
	UserID      int
	SubmitterID int
	Rating      int
	Text        string
	OrderID     int
}

const (
	BalOpsDeposit  = "deposit"
	BalOpsWithdraw = "withdraw"
	BalOpsIncome   = "income"
	BalOpsExpense  = "expense"
	BalOpsRefund   = "refund"
)

type BalanceOperation struct {
	ID        int
	UserID    int
	Type      string
	Delta     float64
	Timestamp int64
}

type Upload struct {
	ID int
}

type Attachment struct {
	ID         int
	URI        string
	SourceType string
	SourceID   string
}

const (
	RoleBuyerStr  = "Buyer"
	RoleSellerStr = "Seller"
)

type Role struct {
	ID           int
	Name         string       `gorm:"unique"`
	Users        []*User      `gorm:"many2many:user_roles"`
	Capabilities []Capability `gorm:"many2many:role_capabilities"`
}

type Capability struct {
	ID   int
	Name string `gorm:"unique"`
}
