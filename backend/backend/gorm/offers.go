package gorm

import (
	"errors"
	"time"

	"git.carried.ru/opcup23/backend/models"
)

func (db *DB) OfferGet(offerid int) (*models.Offer, error) {
	offer := &models.Offer{}
	return offer, db.Connection.Where("id = ?", offerid).
		Preload("Submission.Attachments").
		First(offer).Error
}

func (db *DB) OfferAdd(userid int, orderid int, comment string, due int64, price float64) error {
	count := db.Connection.Where("order_id = ? AND user_id = ?", orderid, userid).First(&models.Offer{}).RowsAffected
	if count != 0 {
		return errors.New("you already submitted your offer to this order.")
	}

	if price < 0 || due < time.Now().Unix() {
		return errors.New("negative price value or due date is in past")
	}

	err := db.Connection.Create(&models.Offer{
		UserID:  userid,
		OrderID: orderid,
		Status:  models.OffSubmitted,
		Comment: comment,
		Due:     due,
		Price:   price,
	}).Error
	if err != nil {
		return err
	}

	count = db.Connection.Where("order_id = ?", orderid).Find(&[]models.Offer{}).RowsAffected
	db.Connection.Model(&models.Order{}).Where("id = ?", orderid).Update("offer_count", count)

	return nil
}

func (db *DB) OfferChangeStatus(offer *models.Offer, status string) error {
	return db.Connection.Model(offer).Update("status", status).Error
}

func (db *DB) OfferGetAll(orderid int) *[]models.Offer {
	list := make([]models.Offer, 0)
	db.Connection.Where("order_id = ?", orderid).Preload("Submission.Attachments").Find(&list)

	return &list
}

func (db *DB) OfferHide(offer *models.Offer, hidden bool) {
	db.Connection.Model(offer).Update("hidden", hidden)
}

func (db *DB) OfferUpdateSubmission(offer *models.Offer, text string, uris []string) error {
	if text == "" {
		return errors.New("please comment your submission")
	}

	sub := &models.Submission{
		OrderID: offer.OrderID,
		OfferID: offer.ID,
		Text:    text,
	}

	if offer.Submission.ID != 0 {
		sub.ID = offer.Submission.ID
	}

	err := db.Connection.Save(sub).Error
	if err != nil {
		return errors.New("failed saving submission")
	}

	db.Connection.Model(offer).Association("Submission").Replace(sub)
	if err != nil {
		return errors.New("failed updating order")
	}

	attachments := make([]models.Attachment, len(uris))
	for i := range uris {
		attachments[i].URI = uris[i]
	}

	db.Connection.Model(sub).Association("Attachments").Replace(attachments)

	db.Connection.Model(&models.Order{ID: offer.OrderID}).Association("Submission").Replace(sub)

	return nil
}

func (db *DB) OfferUpdate(offer *models.Offer, comment string, price float64, due int64) error {
	if price < 0 || due < time.Now().Unix() {
		return errors.New("negative price value or due date is in past")
	}

	return db.Connection.Model(offer).Updates(models.Offer{
		Comment: comment,
		Price:   price,
		Due:     due,
	}).Error
}

func (db *DB) OfferDelete(offerid int) error {
	offer, err := db.OfferGet(offerid)
	if err != nil {
		return err
	}

	if offer.Status != models.OffSubmitted {
		return errors.New("you can't cancel this offer")
	}

	orderid := offer.OrderID

	err = db.Connection.Delete(&offer).Error
	if err != nil {
		return err
	}

	count := db.Connection.Where("order_id = ?", orderid).Find(&[]models.Offer{}).RowsAffected
	db.Connection.Model(&models.Order{}).Where("id = ?", orderid).Update("offer_count", count)

	return nil
}
