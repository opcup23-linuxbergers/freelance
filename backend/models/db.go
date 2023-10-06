package models

type DB interface {
	Setup() error

	UserAdd(string, []byte, string) error
	UserGet(int) (*User, error)
	UserGetByAlias(string) (*User, error)
	UserGetByEmail(string) (*User, error)
	UserGetAll() *[]User
	UserUpdate(*User, string, string, string, string) error
	UserUpdatePassword(*User, []byte) error
	UserUpdateBalance(*User, string, float64) error
	UserRecalculateRating(*User, string, int)
	UserDelete(int) error

	RoleAdd(string, []string) error
	RoleGet(int) (*Role, error)
	RoleGetByName(string) (*Role, error)
	RoleGetAll() *[]*Role
	RoleUpdate(*Role) error
	RoleDelete(int) error
	CapabilityGet(string) (*Capability, error)
	CapabilityGetAll() *[]Capability

	OrderAdd(int, string, string, float64, int64, []string) error
	OrderGet(int) (*Order, error)
	OrderProcessView(int, int)
	OrderUpdate(*Order, string, string, int, string, []string) error
	OrderGetAll() *[]Order
	OrderSearch(string, string, string) *[]Order
	OrderDelete(int) error

	OfferAdd(int, int, string, int64, float64) error
	OfferGet(int) (*Offer, error)
	OfferChangeStatus(*Offer, string) error
	OfferHide(*Offer, bool)
	OfferUpdate(*Offer, string, float64, int64) error
	OfferUpdateSubmission(*Offer, string, []string) error
	OfferGetAll(int) *[]Offer
	OfferDelete(int) error

	ReviewAdd(int, int, int, string, string, int) error
	ReviewExist(int, int) bool

	ChatAdd(string, []int) error
	ChatGet(int) (*Chat, error)
	ChatPostMessage(*Chat, int, string, []string) error

	NotificationGet(int) (*Notification, error)
	NotificationAdd(userid int, message string) error
	NotificationUpdate(*Notification, bool, bool) error

	UploadCreate()
	UploadGetTotal() int
}
