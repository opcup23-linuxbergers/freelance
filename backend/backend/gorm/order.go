package gorm

import (
	"errors"
	"time"
	"fmt"

	"git.carried.ru/opcup23/backend/models"
	"git.carried.ru/opcup23/backend/utils"
)

func (db *DB) OrderAdd(userid int, name string, desc string, price float64, due int64, uris []string) error {
	if due < time.Now().Unix() {
		return errors.New("due date is off")
	}

	attachments := make([]models.Attachment, len(uris))
	for i := range uris {
		attachments[i].URI = uris[i]
	}

	return db.Connection.Create(&models.Order{
		UserID:      userid,
		Name:        name,
		Description: desc,
		Price:       price,
		Due:         due,
		Published:   time.Now().Unix(),
		Status:      models.OrdAvailable,
		Attachments: attachments,
	}).Error
}

func (db *DB) OrderGet(orderid int) (*models.Order, error) {
	order := &models.Order{}
	err := db.Connection.Where("id = ?", orderid).
		Preload("Submission.Attachments").
		Preload("Attachments").
		First(order).Error
	if err != nil {
		return nil, err
	}

	db.Connection.Where("order_id = ?", order.ID).Order("id desc").Find(&order.Offers)

	return order, nil
}

func (db *DB) OrderProcessView(orderid int, userid int) {
	db.Connection.Create(&models.View{UserID: userid, OrderID: orderid})
	count := db.Connection.Where("order_id = ?", orderid).Find(&[]models.View{}).RowsAffected
	db.Connection.Model(&models.Order{}).Where("id = ?", orderid).Update("views", count)
}

func (db *DB) OrderUpdate(order *models.Order, name string, desc string, offerid int, status string, uris []string) error {
	updated := models.Order{
		Name:        name,
		Description: desc,
	}

	attachments := make([]models.Attachment, len(uris))
	for i := range uris {
		attachments[i].URI = uris[i]
	}

	if offerid != 0 {
		if order.Status != models.OrdAvailable {
			return errors.New("can't change seller: order is in progress.")
		}

		offer, err := db.OfferGet(offerid)
		if err != nil || offer.ID == 0 || offer.OrderID != order.ID {
			return errors.New("offer doesn't exist")
		}

		updated.Status = models.OrdInProgress
		updated.ContractorID = offer.UserID
		updated.Price = offer.Price

		u, err := db.UserGet(order.UserID)
		if err != nil {
			return errors.New("buyer doesn't exist anymore?")
		}

		switch {
		case order.Price > updated.Price:
			db.UserUpdateBalance(u, models.BalOpsRefund, order.Price-updated.Price)
		case order.Price < updated.Price:
			err := db.UserUpdateBalance(u, models.BalOpsExpense, updated.Price-order.Price)
			if err != nil {
				return errors.New("this offer is too expensive for you.")
			}
		}

		db.OfferChangeStatus(offer, models.OffConfirmed)
		db.NotificationAdd(offer.UserID, fmt.Sprintf("Order #%d. Buyer accepted your offer.", order.ID))
		db.ChatAdd(order.Name, []int{offer.UserID, order.UserID})
	}

	if status == models.OrdDone {
		if order.Status != models.OrdInProgress {
			return errors.New("don't lie to yourself.. it's not done.")
		}

		offer := &models.Offer{}
		db.Connection.Where("user_id = ?", order.ContractorID).First(offer)

		err := db.OfferChangeStatus(offer, models.OffDone)
		if err != nil {
			return err
		}

		err = db.Connection.Model(order).Update("status", status).Error
		if err != nil {
			return err
		}


		err = db.NotificationAdd(offer.UserID, fmt.Sprintf("Order #%d. Your submission has been accepted.", order.ID))
		if err != nil {
			return errors.New("failed creating notification")
		}

		return nil
	}

	err := db.Connection.Model(order).Updates(updated).Error
	if err != nil {
		return errors.New("failed updating order")
	}

	db.Connection.Model(order).Association("Attachments").Replace(attachments)

	return nil
}

func (db *DB) OrderGetAll() *[]models.Order {
	list := make([]models.Order, 0)
	db.Connection.Order("id desc").Preload("Submission.Attachments").Preload("Attachments").Find(&list)

	return &list
}

func (db *DB) OrderSearch(query string, filter string, order string) *[]models.Order {
	list := make([]models.Order, 0)
	re := utils.FormatQuery(query)
	if filter != "" && order != "" {
		db.Connection.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", re, re).
			Order(fmt.Sprintf("%s %s", filter, order)).Find(&list)
	} else {
		db.Connection.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", re, re).Order("id desc").Find(&list)
	}

	return &list
}

func (db *DB) OrderDelete(orderid int) error {
	order, err := db.OrderGet(orderid)
	if err != nil {
		return err
	}

	if order.Status != models.OrdAvailable {
		return errors.New("you can't cancel this order")
	}

	db.Connection.Unscoped().Model(&order).Association("Offers").Unscoped().Clear()
	db.Connection.Model(&order).Update("status", models.OrdCancelled)

	return nil
}
