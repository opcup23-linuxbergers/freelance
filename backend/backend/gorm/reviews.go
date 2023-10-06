package gorm

import (
	"errors"

	"git.carried.ru/opcup23/backend/models"
)

func (db *DB) ReviewAdd(orderid int, submitter int, userid int, rtype string, text string, rating int) error {
	u, err := db.UserGet(userid)
	if err != nil {
		return errors.New("user not found?")
	}

	db.UserRecalculateRating(u, rtype, rating)
	return db.Connection.Create(&models.Review{
		SubmitterID: submitter,
		Type:    rtype,
		Text:    text,
		Rating:  rating,
		UserID:  userid,
		OrderID: orderid,
	}).Error
}

func (db *DB) ReviewExist(orderid int, userid int) bool {
	review := &models.Review{}

	db.Connection.Where("order_id = ? AND submitter_id = ?", orderid, userid).First(review)
	if review.ID != 0 {
		return true
	}

	return false
}
