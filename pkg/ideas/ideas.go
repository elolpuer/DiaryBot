package ideas

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/elolpuer/DiaryBot/pkg/common"
	"github.com/elolpuer/DiaryBot/pkg/keyboards"
	"github.com/elolpuer/DiaryBot/pkg/models"
	"github.com/go-redis/redis/v8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var ctx = context.Background()

//CreateIdea создает пост
func CreateIdea(db *sql.DB, userID int, chatID int64, text string, time time.Time, bot *tgbotapi.BotAPI) error {
	var idea = new(models.PostOrIdea)
	idea.UserID = userID
	idea.Body = text
	idea.TimeDate = common.MakeStringDate(time)
	idea.TimeClock = common.MakeStringClock(time)
	err := common.Insert(db, idea, false)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Idea added"))
	_, err = bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

//AllIdeas отдает 10 идей через сообщение в зависимости от state
func AllIdeas(db *sql.DB, rdb *redis.Client, userID int, chatID int64, bot *tgbotapi.BotAPI) error {
	state, err := GetIdeasState(rdb, userID)
	if err != nil {
		return err
	}
	ideas, err := common.SelectAll(db, userID, false)
	if err != nil {
		return err
	}
	var ideasMore bool
	if len(ideas) > 10 {
		ideasMore = true
		maxID := len(ideas)
		maxState := maxID / 10
		stateInt, err := strconv.ParseInt(state, 10, 64)
		if err != nil {
			return err
		}
		if int(stateInt) < 0 {
			err := ChangeIdeasStateEnd(rdb, userID)
			if err != nil {
				return err
			}
			stateInt = 0
		}
		if int(stateInt) > maxState {
			err := ChangeIdeasStateBegin(db, rdb, userID)
			if err != nil {
				return err
			}
			stateInt -= 1
		}
		if len(ideas)-10*(int(stateInt)+1) < 0 {
			ideas = ideas[:len(ideas)-10*int(stateInt)]
		} else {
			ideas = ideas[len(ideas)-10*(int(stateInt)+1) : len(ideas)-10*int(stateInt)]
		}
	}
	if ideas == nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("You don't have any ideas in the diary yet to add, click Add"))
		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
	} else {
		var str string
		for _, i := range ideas {
			tm := i.TimeDate + " " + i.TimeClock
			IDstr := "/i_" + fmt.Sprint(i.ID)
			if (len(i.Body) > 100) && (cap([]rune(i.Body)) > 100) {
				i.Body = string([]rune(i.Body)[:100]) + "...[Click on the ID]"
			}
			arr := []string{IDstr, i.Body, tm}
			str = str + strings.Join(arr, "\n") + "\n\n"
		}
		str = str + "To read/update/delete a idea, click on its /i_id"
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s", str))
		if ideasMore == true {
			msg.ReplyMarkup = keyboards.CallbackMenuIdeas
		}
		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

//DeleteAllIdeas удаляет все идеи этого пользователя
func DeleteAllIdeas(db *sql.DB, userID int) error {
	_, err := db.Exec("DELETE FROM ideas WHERE user_id=$1", userID)
	if err != nil {
		return err
	}
	return nil
}

//SearchIdea поиск идеи по дате
func SearchIdea(db *sql.DB, userID int, chatID int64, dateOutside string, month string, bot *tgbotapi.BotAPI) error {
	search, err := common.Search(db, userID, dateOutside, month, false)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Something went wrong"))
		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
		return err
	}
	var count int
	if len(search) > 10 {
		count = len(search) / 10
	}
	if search == nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("No ideas on this day"))
		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
	} else {
		//собираем строку из идей по 10 или меньше. И как отдельные сообщения отправляем пользователю.
		//count нужен чтобы отделять в массиве по 10шт
		for i := count; i >= 0; i-- {
			var str string
			var arr []*models.PostOrIdea
			if len(search)-10*(i+1) < 0 {
				arr = search[:len(search)-10*i]
			} else {
				arr = search[len(search)-10*(i+1) : len(search)-10*i]
			}
			for _, s := range arr {
				tm := s.TimeDate + " " + s.TimeClock
				IDstr := "/i_" + fmt.Sprint(s.ID)
				if len(s.Body) > 100 && (cap([]rune(s.Body)) > 100) {
					s.Body = string([]rune(s.Body)[:100]) + "...[Click on the]"
				}
				arr := []string{IDstr, s.Body, tm}
				str = str + strings.Join(arr, "\n") + "\n\n"
			}
			str = str + "To read/update/delete a idea, click on its /i_id"
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s", str))
			_, err = bot.Send(msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//CallbackOneIdea возвращает одну идею
func CallbackOneIdea(db *sql.DB, userID int, chatID int64, ID string, bot *tgbotapi.BotAPI) error {
	content, err := common.ReturnOne(db, ID, userID, false)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Something went wrong"))
		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
		return err
	}
	txt := fmt.Sprintf("updateI(%s)\n%s", ID, content.Body)
	var callbackDelUpIdea = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Delete", fmt.Sprintf("DeleteIdea_%s", ID)),
			tgbotapi.InlineKeyboardButton{
				Text:                         "Update",
				SwitchInlineQueryCurrentChat: &txt,
			},
		),
	)
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s\n%s %s", content.Body, content.TimeDate, content.TimeClock))
	msg.ReplyMarkup = callbackDelUpIdea
	_, err = bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

//ChangeIdeasStatePlusOne меняет(+1) state в дб для того чтобы сделать меню вывода постов, по 10 последних шт
func ChangeIdeasStatePlusOne(rdb *redis.Client, userID int) error {
	state, err := rdb.Get(ctx, fmt.Sprintf("%d_ideas_state", userID)).Result()
	if err != nil {
		return err
	}
	stateInt, err := strconv.ParseInt(state, 10, 64)
	if err != nil {
		return err
	}
	err = rdb.Set(ctx, fmt.Sprintf("%d_ideas_state", userID), fmt.Sprintf("%d", stateInt+1), 0).Err()
	if err != nil {
		return err
	}
	return nil
}

//ChangeIdeasStateMinusOne меняет(-1) state в дб для того чтобы сделать меню вывода постов, по 10 последних шт
func ChangeIdeasStateMinusOne(rdb *redis.Client, userID int) error {
	state, err := rdb.Get(ctx, fmt.Sprintf("%d_ideas_state", userID)).Result()
	if err != nil {
		return err
	}
	stateInt, err := strconv.ParseInt(state, 10, 64)
	if err != nil {
		return err
	}
	err = rdb.Set(ctx, fmt.Sprintf("%d_ideas_state", userID), fmt.Sprintf("%d", stateInt-1), 0).Err()
	if err != nil {
		return err
	}
	return nil
}

//ChangeIdeasStateEnd меняет state в начало
func ChangeIdeasStateEnd(rdb *redis.Client, userID int) error {
	err := rdb.Set(ctx, fmt.Sprintf("%d_ideas_state", userID), fmt.Sprintf("%d", 0), 0).Err()
	if err != nil {
		return err
	}
	return nil
}

//ChangeIdeasStateBegin меняет state на последний
func ChangeIdeasStateBegin(db *sql.DB, rdb *redis.Client, userID int) error {
	maxID, err := common.ReturnOneMax(db, userID, false)
	if err != nil {
		return err
	}
	maxState := maxID / 10
	err = rdb.Set(ctx, fmt.Sprintf("%d_ideas_state", userID), fmt.Sprintf("%d", maxState), 0).Err()
	if err != nil {
		return err
	}
	return nil
}

//GetIdeasState получаем state из дб
func GetIdeasState(rdb *redis.Client, userID int) (string, error) {
	state, err := rdb.Get(ctx, fmt.Sprintf("%d_ideas_state", userID)).Result()
	if err == redis.Nil {
		err := createIdeaState(rdb, userID)
		if err != nil {
			return "", err
		}
		return "0", nil
	} else if err != nil {
		return "", err
	}
	return state, nil
}

//создаем state в redis
func createIdeaState(rdb *redis.Client, userID int) error {
	err := rdb.Set(ctx, fmt.Sprintf("%d_ideas_state", userID), "0", 0).Err()
	if err != nil {
		return err
	}
	return nil
}
