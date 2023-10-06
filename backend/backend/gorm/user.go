package gorm

import (
	"errors"
	"time"

	"git.carried.ru/opcup23/backend/models"
)

func (db *DB) UserAdd(alias string, hash []byte, email string) error {
	rb, err := db.RoleGetByName(models.RoleBuyerStr)
	if err != nil {
		return err
	}

	rs, err := db.RoleGetByName(models.RoleSellerStr)
	if err != nil {
		return err
	}

	return db.Connection.Create(&models.User{
		Alias: alias,
		Email: email,
		Hash:  hash,
		Roles: []*models.Role{rb, rs},
	}).Error
}

func (db *DB) UserGet(id int) (*models.User, error) {
	u := &models.User{}
	err := db.Connection.Where("id = ?", id).
		Preload("Roles.Capabilities").
		Preload("BalanceOperations").
		Preload("Chats.Messages.Attachments").
		First(u).Error
	if err != nil {
		return nil, err
	}

	db.Connection.Where("user_id = ?", u.ID).Order("id desc").Find(&u.Offers)
	db.Connection.Where("user_id = ?", u.ID).Order("id desc").Find(&u.Orders)
	db.Connection.Where("user_id = ?", u.ID).Order("id desc").Find(&u.Reviews)
	db.Connection.Where("user_id = ?", u.ID).Order("id desc").Find(&u.BalanceOperations)
	db.Connection.Where("user_id = ?", u.ID).Order("id desc").Find(&u.Notifications)

	return u, err
}

func (db *DB) UserGetByAlias(alias string) (*models.User, error) {
	u := &models.User{}
	err := db.Connection.Where("alias = ?", alias).First(u).Error

	return u, err
}

func (db *DB) UserGetByEmail(email string) (*models.User, error) {
	u := &models.User{}
	err := db.Connection.Where("email = ?", email).First(u).Error

	return u, err
}

func (db *DB) UserGetAll() *[]models.User {
	u := &[]models.User{}
	db.Connection.Find(u)
	return u
}

func (db *DB) UserUpdate(u *models.User, alias string, firstname string, lastname string, about string) error {
	return db.Connection.Model(u).Updates(models.User{
		Alias:     alias,
		FirstName: firstname,
		LastName:  lastname,
		About:     about,
	}).Error
}

func (db *DB) UserUpdatePassword(u *models.User, hash []byte) error {
	return db.Connection.Model(u).Update("hash", hash).Error
}

func (db *DB) UserUpdateBalance(u *models.User, op string, d float64) error {
	if d == 0 {
		return nil
	}

	switch op {
	case models.BalOpsDeposit:
		fallthrough
	case models.BalOpsRefund:
		fallthrough
	case models.BalOpsIncome:
		u.Balance += d
	case models.BalOpsWithdraw:
		fallthrough
	case models.BalOpsExpense:
		if u.Balance < d {
			return errors.New("you don't have enough money.")
		}
		u.Balance -= d
	}

	if err := db.Connection.Model(u).Update("balance", u.Balance).Error; err != nil {
		return err
	}

	return db.Connection.Create(&models.BalanceOperation{
		UserID:    u.ID,
		Type:      op,
		Delta:     d,
		Timestamp: time.Now().Unix(),
	}).Error
}

func (db *DB) UserDelete(id int) error {
	db.Connection.Model(&models.User{ID: id}).Association("Roles").Clear()
	return db.Connection.Delete(&models.User{}, id).Error
}

func (db *DB) UserRecalculateRating(u *models.User, rtype string, r int) {
	count := float64(db.Connection.Where("user_id = ? AND type = ?", u.ID, rtype).Find(&[]models.Review{}).RowsAffected)
	switch rtype {
	case models.RevSeller:
		u.SellerRating = (u.SellerRating*count + float64(r)) / (count + 1)
	case models.RevBuyer:
		u.BuyerRating = (u.BuyerRating*count + float64(r)) / (count + 1)
	}
	db.Connection.Updates(u)
}
