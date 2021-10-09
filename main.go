package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	// "log"
	// "math/rand"

	// "os"

	"github.com/elolpuer/DiaryBot/pkg/common"
	"github.com/elolpuer/DiaryBot/pkg/ideas"
	"github.com/elolpuer/DiaryBot/pkg/keyboards"
	"github.com/elolpuer/DiaryBot/pkg/posts"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var db *sql.DB
var startText = "If you do not display the buttons, click /start\nAdd - to add a post/idea\nRead - to read all posts/ideas\nBreak - exit the function\nSearch - search by date\nHelp - help\nTo update (delete) a post (idea), you need to click on its ID and select what you need. Note that when deleting a post (idea), the ID of the subsequent ones is reduced by 1.\n/deleteAllPosts deletes all user posts from the diary.\n/deleteAllIdeas deletes all the user's ideas from the diary."

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env", err.Error())
	}
	Token := os.Getenv("Token")
	PgHost := os.Getenv("PgHost")
	PgPort := os.Getenv("PgPort")
	PgUser := os.Getenv("PgUser")
	PgPass := os.Getenv("PgPass")
	PgDB := os.Getenv("PgDB")
	RedisHost := os.Getenv("RedisHost")
	RedisPort := os.Getenv("RedisPort")
	RedisPassword := os.Getenv("RedisPassword")
	SSLmode := os.Getenv("SSLmode")
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s ", PgHost, PgPort, PgUser, PgPass, PgDB, SSLmode))
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Postgres has been connected.")
	rdb := redis.NewClient(&redis.Options{
		Addr:     RedisHost + ":" + RedisPort,
		Password: RedisPassword, // no password set
		DB:       0,             // use default DB
	})
	fmt.Println("Redis has been connected.")

	bot, err := tgbotapi.NewBotAPI(Token)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.CallbackQuery != nil {
			if update.CallbackQuery.Data == "Post" {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Write your text or click Break"))
				_, err := bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
				for update := range updates {
					if update.Message == nil {
						continue
					}
					if update.Message.Text == "Break" {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("OK, don't add anything"))
						_, err := bot.Send(msg)
						if err != nil {
							log.Fatal(err)
						}
						break
					}
					err := posts.CreatePost(db, update.Message.From.ID, update.Message.Chat.ID, update.Message.Text, update.Message.Time(), bot)
					if err != nil {
						log.Fatal(err)
					}
					break
				}
			}
			if update.CallbackQuery.Data == "Idea" {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Write your text or click Break"))
				_, err := bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
				for update := range updates {
					if update.Message == nil {
						continue
					}
					if update.Message.Text == "Break" {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("OK, don't add anything"))
						_, err := bot.Send(msg)
						if err != nil {
							log.Fatal(err)
						}
						break
					}
					err := ideas.CreateIdea(db, update.Message.From.ID, update.Message.Chat.ID, update.Message.Text, update.Message.Time(), bot)
					if err != nil {
						log.Fatal(err)
					}
					break
				}
			}
			if update.CallbackQuery.Data == "MenuPlusPost" {
				err := posts.ChangePostsStatePlusOne(rdb, update.CallbackQuery.From.ID)
				if err != nil {
					log.Fatal(err)
				}
				delMsg := tgbotapi.DeleteMessageConfig{
					ChatID:    update.CallbackQuery.Message.Chat.ID,
					MessageID: update.CallbackQuery.Message.MessageID,
				}
				_, err = bot.DeleteMessage(delMsg)
				if err != nil {
					log.Fatal(err)
				}
				err = posts.AllPosts(db, rdb, update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID, bot)
				if err != nil {
					log.Fatal(err)
				}
			}
			if update.CallbackQuery.Data == "MenuMinusPost" {
				err := posts.ChangePostsStateMinusOne(rdb, update.CallbackQuery.From.ID)
				if err != nil {
					log.Fatal(err)
				}
				delMsg := tgbotapi.DeleteMessageConfig{
					ChatID:    update.CallbackQuery.Message.Chat.ID,
					MessageID: update.CallbackQuery.Message.MessageID,
				}
				_, err = bot.DeleteMessage(delMsg)
				if err != nil {
					log.Fatal(err)
				}
				err = posts.AllPosts(db, rdb, update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID, bot)
				if err != nil {
					log.Fatal(err)
				}
			}
			if update.CallbackQuery.Data == "MenuEndPost" {
				err := posts.ChangePostsStateEnd(rdb, update.CallbackQuery.From.ID)
				if err != nil {
					log.Fatal(err)
				}
				delMsg := tgbotapi.DeleteMessageConfig{
					ChatID:    update.CallbackQuery.Message.Chat.ID,
					MessageID: update.CallbackQuery.Message.MessageID,
				}
				_, err = bot.DeleteMessage(delMsg)
				if err != nil {
					log.Fatal(err)
				}
				err = posts.AllPosts(db, rdb, update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID, bot)
				if err != nil {
					log.Fatal(err)
				}
			}
			if update.CallbackQuery.Data == "MenuBeginPost" {
				err := posts.ChangePostsStateBegin(db, rdb, update.CallbackQuery.From.ID)
				if err != nil {
					log.Fatal(err)
				}
				delMsg := tgbotapi.DeleteMessageConfig{
					ChatID:    update.CallbackQuery.Message.Chat.ID,
					MessageID: update.CallbackQuery.Message.MessageID,
				}
				_, err = bot.DeleteMessage(delMsg)
				if err != nil {
					log.Fatal(err)
				}
				err = posts.AllPosts(db, rdb, update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID, bot)
				if err != nil {
					log.Fatal(err)
				}
			}
			if update.CallbackQuery.Data == "AllPosts" {
				err := posts.ChangePostsStateEnd(rdb, update.CallbackQuery.From.ID)
				if err != nil {
					log.Fatal(err)
				}
				err = posts.AllPosts(db, rdb, update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID, bot)
				if err != nil {
					log.Fatal(err)
				}
			}
			if update.CallbackQuery.Data == "MenuPlusIdea" {
				err := ideas.ChangeIdeasStatePlusOne(rdb, update.CallbackQuery.From.ID)
				if err != nil {
					log.Fatal(err)
				}
				delMsg := tgbotapi.DeleteMessageConfig{
					ChatID:    update.CallbackQuery.Message.Chat.ID,
					MessageID: update.CallbackQuery.Message.MessageID,
				}
				_, err = bot.DeleteMessage(delMsg)
				if err != nil {
					log.Fatal(err)
				}
				err = ideas.AllIdeas(db, rdb, update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID, bot)
				if err != nil {
					log.Fatal(err)
				}
			}
			if update.CallbackQuery.Data == "MenuMinusIdea" {
				err := ideas.ChangeIdeasStateMinusOne(rdb, update.CallbackQuery.From.ID)
				if err != nil {
					log.Fatal(err)
				}
				delMsg := tgbotapi.DeleteMessageConfig{
					ChatID:    update.CallbackQuery.Message.Chat.ID,
					MessageID: update.CallbackQuery.Message.MessageID,
				}
				_, err = bot.DeleteMessage(delMsg)
				if err != nil {
					log.Fatal(err)
				}
				err = ideas.AllIdeas(db, rdb, update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID, bot)
				if err != nil {
					log.Fatal(err)
				}
			}

			if update.CallbackQuery.Data == "MenuEndIdea" {
				err := ideas.ChangeIdeasStateEnd(rdb, update.CallbackQuery.From.ID)
				if err != nil {
					log.Fatal(err)
				}
				delMsg := tgbotapi.DeleteMessageConfig{
					ChatID:    update.CallbackQuery.Message.Chat.ID,
					MessageID: update.CallbackQuery.Message.MessageID,
				}
				_, err = bot.DeleteMessage(delMsg)
				if err != nil {
					log.Fatal(err)
				}
				err = ideas.AllIdeas(db, rdb, update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID, bot)
				if err != nil {
					log.Fatal(err)
				}
			}
			if update.CallbackQuery.Data == "MenuBeginIdea" {
				err := ideas.ChangeIdeasStateBegin(db, rdb, update.CallbackQuery.From.ID)
				if err != nil {
					log.Fatal(err)
				}
				delMsg := tgbotapi.DeleteMessageConfig{
					ChatID:    update.CallbackQuery.Message.Chat.ID,
					MessageID: update.CallbackQuery.Message.MessageID,
				}
				_, err = bot.DeleteMessage(delMsg)
				if err != nil {
					log.Fatal(err)
				}
				err = ideas.AllIdeas(db, rdb, update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID, bot)
				if err != nil {
					log.Fatal(err)
				}
			}
			if update.CallbackQuery.Data == "AllIdeas" {
				err := ideas.ChangeIdeasStateEnd(rdb, update.CallbackQuery.From.ID)
				if err != nil {
					log.Fatal(err)
				}
				err = ideas.AllIdeas(db, rdb, update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID, bot)
				if err != nil {
					log.Fatal(err)
				}

			}
			if update.CallbackQuery.Data == "SearchPost" {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Select the date"))
				msg.ReplyMarkup = keyboards.CallbackDatePost
				_, err := bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
			}
			if update.CallbackQuery.Data == "SearchIdea" {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Date"))
				msg.ReplyMarkup = keyboards.CallbackDateIdea
				_, err := bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
			}
			var dateOutside string
			if strings.Contains(update.CallbackQuery.Data, "day_P") {
				indxP := strings.Index(update.CallbackQuery.Data, "P")
				dateOutside = update.CallbackQuery.Data[indxP+1:]
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Select the month"))
				msg.ReplyMarkup = keyboards.CallbackMonth
				_, err := bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
				for update := range updates {
					if update.CallbackQuery != nil {
						if strings.Contains(update.CallbackQuery.Data, "month_") {
							indx := strings.Index(update.CallbackQuery.Data, "_")
							m := update.CallbackQuery.Data[indx+1:]
							err := posts.SearchPost(db, update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID, dateOutside, m, bot)
							if err != nil {
								log.Fatal(err)
							}
							break
						}

					}
					if update.Message != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Press again")
						_, err := bot.Send(msg)
						if err != nil {
							log.Fatal(err)
						}
						break
					}
				}
			}
			if strings.Contains(update.CallbackQuery.Data, "day_I") {
				indxI := strings.Index(update.CallbackQuery.Data, "I")
				dateOutside = update.CallbackQuery.Data[indxI+1:]
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Select the month"))
				msg.ReplyMarkup = keyboards.CallbackMonth
				_, err := bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
				for update := range updates {
					if update.CallbackQuery != nil {
						if strings.Contains(update.CallbackQuery.Data, "month_") {
							indx := strings.Index(update.CallbackQuery.Data, "_")
							m := update.CallbackQuery.Data[indx+1:]
							err := ideas.SearchIdea(db, update.CallbackQuery.From.ID, update.CallbackQuery.Message.Chat.ID, dateOutside, m, bot)
							if err != nil {
								log.Fatal(err)
							}
							break
						}

					}
					if update.Message != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Press again")
						_, err := bot.Send(msg)
						if err != nil {
							log.Fatal(err)
						}
						break
					}
				}
			}
			if strings.Contains(update.CallbackQuery.Data, "DeletePost_") {
				indx := strings.Index(update.CallbackQuery.Data, "_")
				ID := update.CallbackQuery.Data[indx+1:]
				intID, err := strconv.Atoi(ID)
				if err != nil {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Something went wrong"))
					_, err := bot.Send(msg)
					if err != nil {
						log.Fatal(err)
					}
					log.Fatal(err)
				}
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Are we definitely deleting it?"))
				msg.ReplyMarkup = keyboards.KeyboardDelete
				_, err = bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
				for update := range updates {
					if update.Message == nil {
						continue
					}
					if update.Message.Text == "Yes" {
						num, err := common.Delete(db, intID, update.Message.From.ID, true)
						if num == 1 || err != nil {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Something went wrong"))
							msg.ReplyMarkup = keyboards.Keyboard
							_, err := bot.Send(msg)
							if err != nil {
								log.Fatal(err)
							}
						} else {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Deleted"))
							msg.ReplyMarkup = keyboards.Keyboard
							_, err := bot.Send(msg)
							if err != nil {
								log.Fatal(err)
							}
						}
						break
					}
					if update.Message.Text == "No" {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("OK, we don't delete anything"))
						msg.ReplyMarkup = keyboards.Keyboard
						_, err := bot.Send(msg)
						if err != nil {
							log.Fatal(err)
						}
						break
					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Send Yes or No"))
						msg.ReplyMarkup = keyboards.KeyboardDelete
						_, err := bot.Send(msg)
						if err != nil {
							log.Fatal(err)
						}
					}
				}
			}
			if strings.Contains(update.CallbackQuery.Data, "DeleteIdea_") {
				indx := strings.Index(update.CallbackQuery.Data, "_")
				ID := update.CallbackQuery.Data[indx+1:]
				intID, err := strconv.Atoi(ID)
				if err != nil {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Something went wrong"))
					_, err := bot.Send(msg)
					if err != nil {
						log.Fatal(err)
					}
					log.Fatal(err)
				}
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Are we definitely deleting it?"))
				msg.ReplyMarkup = keyboards.KeyboardDelete
				_, err = bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
				for update := range updates {
					if update.Message == nil {
						continue
					}
					if update.Message.Text == "Yes" {
						num, err := common.Delete(db, intID, update.Message.From.ID, false)
						if num == 1 || err != nil {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Something went wrong"))
							msg.ReplyMarkup = keyboards.Keyboard
							_, err := bot.Send(msg)
							if err != nil {
								log.Fatal(err)
							}
						} else {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Deleted"))
							msg.ReplyMarkup = keyboards.Keyboard
							_, err := bot.Send(msg)
							if err != nil {
								log.Fatal(err)
							}
						}
						break
					}
					if update.Message.Text == "No" {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("OK, we don't delete anything"))
						msg.ReplyMarkup = keyboards.Keyboard
						_, err := bot.Send(msg)
						if err != nil {
							log.Fatal(err)
						}
						break
					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Send Yes or No"))
						msg.ReplyMarkup = keyboards.KeyboardDelete
						_, err := bot.Send(msg)
						if err != nil {
							log.Fatal(err)
						}
					}
				}
			}
			callBack := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			callBack.ShowAlert = false
			_, err := bot.AnswerCallbackQuery(callBack)
			if err != nil {
				log.Fatal(err)
			}
		}
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, startText)
			msg.ReplyMarkup = keyboards.Keyboard
			_, err := bot.Send(msg)
			if err != nil {
				log.Fatal(err)
			}
			err = common.AddToStart(db, update.Message.From.ID)
			if err != nil {
				log.Fatal(err)
			}

		case "/deleteAllPosts":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Are you sure you want to permanently delete all posts?")
			msg.ReplyMarkup = keyboards.KeyboardDelete
			_, err = bot.Send(msg)
			if err != nil {
				log.Fatal(err)
			}
			for update := range updates {
				if update.Message == nil {
					continue
				}
				if update.Message.Text == "Yes" {
					err := posts.DeleteAllPosts(db, update.Message.From.ID)
					if err != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Something went wrong"))
						msg.ReplyMarkup = keyboards.Keyboard
						_, err := bot.Send(msg)
						if err != nil {
							log.Fatal(err)
						}
					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Deleted"))
						msg.ReplyMarkup = keyboards.Keyboard
						_, err := bot.Send(msg)
						if err != nil {
							log.Fatal(err)
						}
					}
					break
				}
				if update.Message.Text == "No" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("OK, we don't delete anything"))
					msg.ReplyMarkup = keyboards.Keyboard
					_, err := bot.Send(msg)
					if err != nil {
						log.Fatal(err)
					}
					break
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Send Yes or No"))
					msg.ReplyMarkup = keyboards.KeyboardDelete
					_, err := bot.Send(msg)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		case "/deleteAllIdeas":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Are you sure you want to permanently delete all ideas?")
			msg.ReplyMarkup = keyboards.KeyboardDelete
			_, err = bot.Send(msg)
			if err != nil {
				log.Fatal(err)
			}
			for update := range updates {
				if update.Message == nil {
					continue
				}
				if update.Message.Text == "Yes" {
					err := ideas.DeleteAllIdeas(db, update.Message.From.ID)
					if err != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Something went wrong"))
						msg.ReplyMarkup = keyboards.Keyboard
						_, err := bot.Send(msg)
						if err != nil {
							log.Fatal(err)
						}
					} else {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Deleted"))
						msg.ReplyMarkup = keyboards.Keyboard
						_, err := bot.Send(msg)
						if err != nil {
							log.Fatal(err)
						}
					}
					break
				}
				if update.Message.Text == "No" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("OK, we don't delete anything"))
					msg.ReplyMarkup = keyboards.Keyboard
					_, err := bot.Send(msg)
					if err != nil {
						log.Fatal(err)
					}
					break
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Send Yes or No"))
					msg.ReplyMarkup = keyboards.KeyboardDelete
					_, err := bot.Send(msg)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		case "Help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, startText)
			_, err := bot.Send(msg)
			if err != nil {
				log.Fatal(err)
			}
		case "Add":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Choose what to add"))
			msg.ReplyMarkup = keyboards.CallbackAdd
			_, err := bot.Send(msg)
			if err != nil {
				log.Fatal(err)
			}
		case "Read":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Choose what to watch"))
			msg.ReplyMarkup = keyboards.CallbackAll
			_, err := bot.Send(msg)
			if err != nil {
				log.Fatal(err)
			}
		case "Search":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Search by date\nChoose what to look for"))
			msg.ReplyMarkup = keyboards.CallbackSearch
			_, err := bot.Send(msg)
			if err != nil {
				log.Fatal(err)
			}
		}
		//смотрим выключает ли в себя запрос /p_ или /i_, если включает, то отдаем пост/идею и клавиатуру.
		if strings.Contains(update.Message.Text, "/p_") && len(update.Message.Text) <= 7 {
			indx := strings.Index(update.Message.Text, "_")
			ID := update.Message.Text[indx+1:]
			err := posts.CallbackOnePost(db, update.Message.From.ID, update.Message.Chat.ID, ID, bot)
			if err != nil {
				log.Fatal(err)
			}
		}
		if strings.Contains(update.Message.Text, "/i_") && len(update.Message.Text) <= 7 {
			indx := strings.Index(update.Message.Text, "_")
			ID := update.Message.Text[indx+1:]
			err := ideas.CallbackOneIdea(db, update.Message.From.ID, update.Message.Chat.ID, ID, bot)
			if err != nil {
				log.Fatal(err)
			}
		}
		if (len(update.Message.Text) > 20) && (strings.Contains(update.Message.Text, "@godofdiarybot update")) {
			err := common.UpdatePostOrIdea(db, update.Message.Text, update.Message.From.ID, update.Message.Chat.ID, bot)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
