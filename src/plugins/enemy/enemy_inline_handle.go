package enemy

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	gonanoid "github.com/matoous/go-nanoid/v2"
	tgbotapi "github.com/qwq233/telegram-bot-api"
	"strings"
)

func InlineEnemy(update tgbotapi.Update) error {
	_, name, _ := strings.Cut(update.InlineQuery.Query, "敌人-")
	enemyList := utils.GetEnemiesByName(name)
	var inlineQueryResults []interface{}
	for k, v := range enemyList {
		id, _ := gonanoid.New(32)
		queryResult := tgbotapi.InlineQueryResultArticle{
			ID:          id,
			Type:        "article",
			Title:       k,
			Description: "查询" + k,
			ThumbURL:    v,
			InputMessageContent: tgbotapi.InputTextMessageContent{
				Text: "/enemy " + k,
			},
		}
		inlineQueryResults = append(inlineQueryResults, queryResult)
	}
	answerInlineQuery := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		Results:       inlineQueryResults,
		CacheTime:     0,
	}
	bot.Arknights.Send(answerInlineQuery)
	return nil
}
