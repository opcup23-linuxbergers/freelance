package gorm

import (
	"git.carried.ru/opcup23/backend/models"
)

func (db *DB) UploadCreate() {
	db.Connection.Create(&models.Upload{})
}

func (db *DB) UploadGetTotal() int {
	return int(db.Connection.Find(&[]models.Upload{}).RowsAffected)
}
