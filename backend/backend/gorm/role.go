package gorm

import (
	"git.carried.ru/opcup23/backend/models"
)

func (db *DB) RoleAdd(name string, capabilities []string) error {
	caps := make([]models.Capability, len(capabilities))

	for i, _ := range caps {
		caps[i].Name = capabilities[i]
	}

	return db.Connection.Create(&models.Role{Name: name, Capabilities: caps}).Error
}

func (db *DB) RoleGet(id int) (*models.Role, error) {
	r := &models.Role{}
	err := db.Connection.Where("id = ?", id).Preload("Capabilities").Preload("Users").First(r).Error
	return r, err
}

func (db *DB) RoleGetAll() *[]*models.Role {
	roles := &[]*models.Role{}
	db.Connection.Find(roles)
	return roles
}

func (db *DB) RoleGetByName(name string) (*models.Role, error) {
	r := &models.Role{}
	err := db.Connection.Where("name = ?", name).Preload("Capabilities").Preload("Users").First(r).Error
	return r, err
}

func (db *DB) RoleUpdate(role *models.Role) error {
	db.Connection.Model(role).Association("Capabilities").Replace(role.Capabilities)
	return db.Connection.Updates(role).Error
}

func (db *DB) RoleDelete(id int) error {
	db.Connection.Where("id = ?", id).Association("Users").Clear()
	db.Connection.Where("id = ?", id).Association("Capabilities").Clear()
	return db.Connection.Delete(&models.Role{}, id).Error
}

func (db *DB) CapabilityGet(name string) (*models.Capability, error) {
	c := &models.Capability{}
	return c, db.Connection.Where("name = ?", name).First(c).Error
}

func (db *DB) CapabilityGetAll() *[]models.Capability {
	// TODO: capabilities are hard-coded, so maybe get them from the CapSet?
	caps := &[]models.Capability{}

	db.Connection.Find(caps)

	return caps
}
