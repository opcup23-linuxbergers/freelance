package gorm

import (
	"fmt"

	"git.carried.ru/opcup23/backend/api"
	"git.carried.ru/opcup23/backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	Connection *gorm.DB
}

func Open(dsn string, dbtype string) (*DB, error) {
	var (
		err error
		db  = &DB{}
	)

	switch dbtype {
	case "postgres":
		db.Connection, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})
	case "sqlite":
		db.Connection, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	case "mysql":
		db.Connection, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		return nil, err
	}

	db.Connection.AutoMigrate(models.Role{}, models.Capability{}, models.User{}, models.BalanceOperation{},
		models.Attachment{}, models.Offer{}, models.Order{}, models.Review{}, models.View{}, models.Submission{},
		models.Chat{}, models.Message{}, models.Notification{}, models.Upload{})

	return db, nil
}

func (db *DB) Setup() error {
	var err error

	defaultBuyerCapabilities := []string{
		api.CapCreateBuyOrders,
	}

	defaultSellerCapabilities := []string{
		api.CapCreateResponses,
	}

	err = db.RoleAdd(models.RoleBuyerStr, defaultBuyerCapabilities)
	if err != nil {
		return fmt.Errorf("role 'Buyer' already exist")
	}

	err = db.RoleAdd(models.RoleSellerStr, defaultSellerCapabilities)
	if err != nil {
		return fmt.Errorf("role 'Seller' already exist")
	}

	return nil
}
