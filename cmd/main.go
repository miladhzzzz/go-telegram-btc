package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	apiURL        = "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd"
	checkInterval = 60 * time.Second
)

type BitcoinPrice struct {
	USD float64 `json:"usd"`
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Printf("couldnt get updates : %v", err.Error())
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start":
				price, err := getBitcoinPrice()
				if err != nil {
					log.Println("Error getting Bitcoin price:", err)
					msg.Text = "Error getting Bitcoin price"
				} else {
					msg.Text = fmt.Sprintf("Welcome to BTC Price Tracker Bot!\nCurrent Bitcoin price: $%.2f\nPlease use the following commands:\n/settarget - Set your target price\n/checkprice - Check the current price", price)
				}
				bot.Send(msg)

			case "settarget":
				msg.Text = "Enter your target BTC price and specify if it's for below or above (e.g., 50000 below)"
				bot.Send(msg)
				
			case "checkprice":
				price, err := getBitcoinPrice()
				if err != nil {
					log.Println("Error getting Bitcoin price:", err)
					msg.Text = "Error getting Bitcoin price"
				} else {
					msg.Text = fmt.Sprintf("Current Bitcoin price: $%.2f", price)
				}
				bot.Send(msg)
				
			default:
				msg.Text = "I don't know that command"
				bot.Send(msg)
			}
		}

		if strings.HasPrefix(update.Message.Text, "/settarget") {
			args := strings.Split(update.Message.Text, " ")
			if len(args) != 3 {
				reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid format. Please use /settarget [price] [below/above]")
				bot.Send(reply)
				continue
			}

			targetPrice, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid price format. Please enter a valid number.")
				bot.Send(reply)
				continue
			}

			direction := args[2]
			if direction != "below" && direction != "above" {
				reply := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid direction. Please use 'below' or 'above'.")
				bot.Send(reply)
				continue
			}

			// price, err := getBitcoinPrice()
			// if err != nil {
			// 	log.Println("Error getting Bitcoin price:", err)
			// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error getting Bitcoin price")
			// 	bot.Send(msg)
			// 	continue
			// }

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Your target price is set to %.2f USD for %s the current price. You'll be notified when it's reached.", targetPrice, direction))
			bot.Send(msg)

			go trackBitcoinPrice(bot, update.Message.Chat.ID, targetPrice, direction)
		}
	}
}

func trackBitcoinPrice(bot *tgbotapi.BotAPI, chatID int64, targetPrice float64, direction string) {
	for {
		price, err := getBitcoinPrice()
		if err != nil {
			log.Println("Error getting Bitcoin price:", err)
			continue
		}

		var notify bool
		switch direction {
		case "below":
			if price <= targetPrice {
				notify = true
			}
		case "above":
			if price >= targetPrice {
				notify = true
			}
		}

		if notify {
			message := fmt.Sprintf("Bitcoin price has reached your target: $%.2f", price)
			msg := tgbotapi.NewMessage(chatID, message)
			bot.Send(msg)
			break
		}

		time.Sleep(checkInterval)
	}
}

func getBitcoinPrice() (float64, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data map[string]BitcoinPrice
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return 0, err
	}

	price := data["bitcoin"].USD
	return price, nil
}
