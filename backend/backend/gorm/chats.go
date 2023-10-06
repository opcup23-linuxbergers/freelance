package gorm

import (
	"errors"

	"git.carried.ru/opcup23/backend/models"
)

func (db *DB) ChatAdd(topic string, users []int) error {
	chat := &models.Chat{Name: topic}
	err := db.Connection.Create(chat).Error
	if err != nil {
		return errors.New("failed to create a chat")
	}

	db.Connection.Model(chat).Association("Users").
		Append([]models.User{models.User{ID: users[0]}, models.User{ID: users[1]}})

	return nil
}

func (db *DB) ChatGet(chatid int) (*models.Chat, error) {
	chat := &models.Chat{}
	err := db.Connection.Where("id = ?", chatid).Preload("Messages.Attachments").First(chat).Error
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (db *DB) ChatPostMessage(chat *models.Chat, senderid int, message string, uris []string) error {
	msg := &models.Message{
		ChatID: chat.ID,
		UserID: senderid,
		Text: message,
	}

	attachments := make([]models.Attachment, len(uris))
	for i := range uris {
		attachments[i].URI = uris[i]
	}

	err := db.Connection.Create(msg).Error
	if err != nil {
		return err
	}

	db.Connection.Model(msg).Association("Attachments").
		Append(attachments)

	db.Connection.Model(chat).Association("Messages").
		Append(msg)

	return nil
}
