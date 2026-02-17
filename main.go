package main

import (
	"fmt"
	"log"
	"os"

	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var bot *tgbotapi.BotAPI

func main() {
	err1 := godotenv.Load()
	if err1 != nil {
		log.Println("Ошибка загрузки .env:", err1)
	}
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("❌ Установите переменную окружения TELEGRAM_BOT_TOKEN")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("❌ Ошибка создания бота:", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			handleCommand(update.Message)
			continue
		}
		handleMessage(update.Message)
	}

}

func handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		sendWelcome(message.Chat.ID)

	case "help":
		sendHelp(message.Chat.ID)

	case "search":
		query := message.CommandArguments()
		if query == "" {
			msg := tgbotapi.NewMessage(message.Chat.ID,
				"🔍 Использование: /search <запрос>\n"+
					"Пример: /search Queen Bohemian Rhapsody")
			bot.Send(msg)
			return
		}
		processSearch(message.Chat.ID, query)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID,
			"🤔 Неизвестная команда. Используйте /help для списка команд")
		bot.Send(msg)
	}

}
func handleMessage(message *tgbotapi.Message) {
	if strings.TrimSpace(message.Text) == "" {
		return
	}
	if len(message.Text) < 2 {
		msg := tgbotapi.NewMessage(message.Chat.ID,
			"📝 Запрос должен содержать минимум 2 символа")
		bot.Send(msg)
		return
	}
	processSearch(message.Chat.ID, message.Text)
}

func sendWelcome(chatID int64) {
	text := `🎵 *Добро пожаловать в музыкального бота!*

Я помогу найти музыку в различных сервисах.

*Как пользоваться:*
• Просто напишите название песни или исполнителя


*Доступные команды:*
/start - начать работу
/help - помощь
/search - поиск музыки`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func sendHelp(chatID int64) {
	text := `*🎵 Помощь по использованию бота*

*Основные команды:*
/start - начать работу с ботом
/help - показать эту справку
/search <запрос> - поиск музыки

*Как искать музыку:*
1. Просто напишите название песни
2. Или исполнителя и название песни
3. Бот найдет музыку в:
   • SoundCloud
   • Яндекс.Музыка
 
• Иногда поиск может занять несколько секунд`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func processSearch(chatID int64, query string) {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("🔍 Ищу: *%s*...", query))
	msg.ParseMode = "Markdown"
	bot.Send(msg)

	results, err := searchMusicSimple(query)
	if err != nil {
		log.Printf("❌ Ошибка поиска: %v", err)
		errorMsg := tgbotapi.NewMessage(chatID, "❌ Произошла ошибка при поиске. Попробуйте позже.")
		bot.Send(errorMsg)
		return
	}

	if len(results) == 0 {
		noResultsMsg := tgbotapi.NewMessage(chatID, "😔 Ничего не найдено. Попробуйте другой запрос.")
		bot.Send(noResultsMsg)
		return
	}

	response := fmt.Sprintf("🎵 *Результаты для: %s*\n\n", query)
	for _, result := range results {
		response += result + "\n"
	}
	response += "\n_Ищет музыку через открытые источники_"

	resultsMsg := tgbotapi.NewMessage(chatID, response)
	resultsMsg.ParseMode = "Markdown"
	bot.Send(resultsMsg)
}

func searchMusicSimple(query string) ([]string, error) {
	var results []string

	encodedQuery := strings.ReplaceAll(query, " ", "%20")

	results = append(results, fmt.Sprintf("  • [SoundCloud](https://soundcloud.com/search?q=%s)", encodedQuery))
	results = append(results, fmt.Sprintf("  • [Яндекс.Музыка](https://music.yandex.ru/search?text=%s)", encodedQuery))

	return results, nil
}
