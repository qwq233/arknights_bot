package system

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/qwq233/telegram-bot-api"
)

func RegulationHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)

	if bot.Arknights.IsAdmin(chatId, userId) {
		replyToMessage := update.Message.ReplyToMessage
		if replyToMessage != nil {
			replyMessageId := replyToMessage.MessageID
			var joined utils.GroupJoined
			utils.GetJoinedByChatId(chatId).Scan(&joined)
			joined.Reg = replyMessageId
			bot.DBEngine.Table("group_joined").Save(&joined)
			sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("消息[%d](https://t.me/%s/%d)已设置为群规！", replyMessageId, replyToMessage.Chat.UserName, replyMessageId))
			sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
			msg, err := bot.Arknights.Send(sendMessage)
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		}
		return nil
	}

	sendMessage := tgbotapi.NewMessage(chatId, "无使用权限！")
	sendMessage.ReplyToMessageID = messageId
	msg, err := bot.Arknights.Send(sendMessage)
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
	return nil
}
