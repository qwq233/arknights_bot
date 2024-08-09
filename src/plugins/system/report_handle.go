package system

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"bytes"
	"fmt"
	tgbotapi "github.com/qwq233/telegram-bot-api"
)

// ReportHandle 举报
func ReportHandle(update tgbotapi.Update) error {
	message := update.Message
	chatId := message.Chat.ID

	message.Delete()

	if message.ReplyToMessage != nil {
		replyToMessage := message.ReplyToMessage
		replyMessageId := replyToMessage.MessageID
		target := replyToMessage.From.ID
		name := replyToMessage.From.FullName()

		if bot.Arknights.IsAdmin(chatId, target) {
			sendMessage := tgbotapi.NewMessage(chatId, "无法举报管理员！")
			msg, err := bot.Arknights.Send(sendMessage)
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			return nil
		}

		// 获取全部管理员
		getAdmins := tgbotapi.ChatAdministratorsConfig{
			ChatConfig: tgbotapi.ChatConfig{
				ChatID: chatId,
			},
		}

		var buttons [][]tgbotapi.InlineKeyboardButton

		var text bytes.Buffer
		text.WriteString(fmt.Sprintf("被举报人：[%s](tg://user?id=%d)\n", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, name), target))
		text.WriteString(fmt.Sprintf("消息存放：[%d](https://t.me/%s/%d)", replyMessageId, replyToMessage.Chat.UserName, replyMessageId))
		charAdmins, _ := bot.Arknights.GetChatAdministrators(getAdmins)
		var admins []int64
		for _, admin := range charAdmins {
			if !admin.User.IsBot {
				admins = append(admins, admin.User.ID)
			}
		}

		for _, admin := range admins {
			text.WriteString(fmt.Sprintf("[\u200b](tg://user?id=%d) ", admin))
		}

		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🚫封禁", fmt.Sprintf("report,%s,%d,%d", "BAN", target, replyMessageId)),
			tgbotapi.NewInlineKeyboardButtonData("❌关闭", fmt.Sprintf("report,%s,%d,%d", "CLOSE", target, replyMessageId)),
		))

		inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
			buttons...,
		)

		sendMessage := tgbotapi.NewMessage(chatId, text.String())
		sendMessage.ReplyMarkup = inlineKeyboardMarkup
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		bot.Arknights.Send(sendMessage)
	}

	return nil
}
