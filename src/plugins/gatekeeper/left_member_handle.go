package gatekeeper

import (
	tgbotapi "github.com/qwq233/telegram-bot-api"
)

func LeftMemberHandle(update tgbotapi.Update) error {
	update.Message.Delete()
	return nil
}
