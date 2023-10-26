package main

import (
	"encryption_bot/internal/config"
	"encryption_bot/internal/encryption"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	config, err := config.LoadConfiguration("data.json")
	if err != nil {
		log.Panicln("Не удалось получить данные json")
	}

	// используя токен создаем новый инстанс бота
	bot, err := tgbotapi.NewBotAPI(config.TgBotToken)
	if err != nil {
		log.Panic("не удалось запустить бота", err)
	}

	log.Printf("Авторизация аккаунта: %s", bot.Self.UserName)

	// updater - структура с конфигом для получения апдейтов
	updater := tgbotapi.NewUpdate(0)
	updater.Timeout = 60

	// создаем канал в который будут прилетать новые сообщения
	updates, err := bot.GetUpdatesChan(updater)
	if err != nil {
		log.Panicln("Не приходят апгрейды")
	}

	// в канал updates прилетают структуры типа Update
	// вычитываем их и обрабатываем
	for update := range updates {
		reply := ""
		// логируем от кого какое сообщение пришло
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// преобразуем комманду пользователя (Находим ключ и текст для шифрации)
		switch update.Message.Command() {
		case "help":
			_, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, config.Messagess.Help))
			if err != nil {
				log.Println("Не удалось отправить сообщение с помощью")
				continue
			}
		}

		text := update.Message.Text
		task := strings.Split(text, " ")
		if strings.Count(text, " ") < 2 {
			continue
		}

		resultTask := task[2]
		// Преобразуем подстроку ключа в байтовый формат
		rot, err := strconv.Atoi(task[1])
		if err != nil {
			_, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, config.Messagess.InvalidDataFormat))
			if err != nil {
				log.Println("Не удалось отправить сообщение о неверном формате пользователю")
				continue
			}

			_, err = bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID))
			if err != nil {
				log.Println("Не удалось удалить сообщение пользователя")
				continue
			}
			continue
		}
		rotByte := byte(rot)

		// Свитч на обработку комманд.
		// Комманда - сообщение, начинающееся с "/".
		switch update.Message.Command() {
		case "decrypt":
			if strings.Count(text, " ") < 2 {
				_, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, config.Messagess.InvalidFormat))
				if err != nil {
					log.Println("Не удалось отправить сообщение")
					continue
				}
				continue
			}

			reply = encryption.UnCesar(resultTask, rotByte)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

			// Отправляем расшифрованный текст пользователю.
			_, err := bot.Send(msg)
			if err != nil {
				log.Println("Не удалось отправить расшифрованный текст пользователю")
				continue
			}
		case "encrypt":
			if strings.Count(text, " ") < 2 {
				_, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, config.Messagess.InvalidDataFormat))
				if err != nil {
					log.Println("Не удалось отправить зашифрованный текст пользователю")
					continue
				}
			}

			reply = encryption.Cesar(resultTask, rotByte)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

			// Отправляем зашифрованный текст пользователю.
			_, err = bot.Send(msg)
			if err != nil {
				log.Println("Не удалось отправить зашифрованный текст пользователю")
			}

			_, err = bot.DeleteMessage(tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID))
			if err != nil {
				log.Println("Не удалось удалить сообщение пользователя")
				continue
			}
		}
	}
}
