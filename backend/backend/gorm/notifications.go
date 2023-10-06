package gorm

import (
	"time"

	"git.carried.ru/opcup23/backend/models"
)

func (db *DB) NotificationAdd(userid int, message string) error {
	return db.Connection.Create(&models.Notification{
		UserID:    userid,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}).Error
}

func (db *DB) NotificationGet(id int) (*models.Notification, error) {
	notif := &models.Notification{}
	return notif, db.Connection.Where("id = ?", id).First(notif).Error
}

func (db *DB) NotificationUpdate(notif *models.Notification, hidden bool, read bool) error {
	return db.Connection.Model(notif).Updates(models.Notification{
		Hidden: hidden,
		Read:   read,
	}).Error
}
