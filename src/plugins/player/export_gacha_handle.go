package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/commandoperation"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/qwq233/telegram-bot-api"
	"github.com/xuri/excelize/v2"
	"io"
	"os"
	"time"
)

type PlayerOperationExport struct {
	commandoperation.OperationAbstract
}

func (_ PlayerOperationExport) hintWordForPlayerSelection() string {
	return "请选择要导出的角色"
}

func (_ PlayerOperationExport) Run(uid string, userAccount account.UserAccount, chatId int64, message *tgbotapi.Message) error {
	var userGacha []UserGacha
	res := utils.GetUserGacha(userAccount.UserNumber, uid).Scan(&userGacha)
	if res.RowsAffected == 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "不存在抽卡记录！")
		bot.Arknights.Send(sendMessage)
		return nil
	}

	sendAction := tgbotapi.NewChatAction(chatId, "upload_document")
	bot.Arknights.Send(sendAction)

	f := excelize.NewFile()
	// 设置单元格的值
	f.SetSheetRow("Sheet1", "A1", &[]interface{}{"卡池", "干员", "星级", "New", "时间"})
	f.SetColWidth("Sheet1", "A", "E", 18)
	for i, gacha := range userGacha {
		f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", i+2), &[]interface{}{gacha.PoolName, gacha.CharName, gacha.Rarity + 1, gacha.IsNew, time.Unix(gacha.Ts, 0).Format("2006-01-02 15:04:05")})
	}
	fileName := time.Now().Format("2006-01-02") + ".xlsx"
	// 根据指定路径保存文件
	if err := f.SaveAs(fileName); err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, "生成文件失败！")
		bot.Arknights.Send(sendMessage)
	}

	file, _ := os.Open(fileName)
	b, _ := io.ReadAll(file)
	file.Close()
	os.Remove(fileName)
	sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileBytes{Bytes: b, Name: fileName})
	sendDocument.Caption = "抽卡记录导出成功！"
	bot.Arknights.Send(sendDocument)
	return nil
}
