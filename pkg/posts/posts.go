package posts

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

//CreatePost создает пост
func CreatePost(db *sql.DB, userID int, chatID int64, text string, time time.Time, bot *tgbotapi.BotAPI) error {
	var post = new(models.PostOrIdea)
	post.UserID = userID
	post.Body = text
	post.TimeDate = common.MakeStringDate(time)
	post.TimeClock = common.MakeStringClock(time)
	err := common.Insert(db, post, true)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Something went wrong"))
		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
		return err
	}
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Post added"))
	_, err = bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

//AllPosts отдает 10 постов через сообщение в зависимости от state
func AllPosts(db *sql.DB, rdb *redis.Client, userID int, chatID int64, bot *tgbotapi.BotAPI) error {
	//берем state из redis
	//state отвечает за меню, он начинается с 0, что означает последнюю страницу с последними постами
	//также если state перешел за какую то черту, будь то 0 или больше максимального, то мы его соответственно
	//возвращаем в 0 и максимальный-1
	state, err := GetPostsState(rdb, userID)
	if err != nil {
		return err
	}
	posts, err := common.SelectAll(db, userID, true)
	if err != nil {
		return err
	}
	var postsMore bool
	if len(posts) > 10 {
		postsMore = true
		maxID := len(posts)
		maxState := maxID / 10
		//перевод state в int64 из string
		stateInt, err := strconv.ParseInt(state, 10, 64)
		if err != nil {
			return err
		}
		if int(stateInt) < 0 {
			err := ChangePostsStateEnd(rdb, userID)
			if err != nil {
				return err
			}
			stateInt = 0
		}
		if int(stateInt) > maxState {
			err := ChangePostsStateBegin(db, rdb, userID)
			if err != nil {
				return err
			}
			stateInt -= 1
		}
		//это сделано чтобы не было проблемы с отрицательными значениями индекса массива
		if len(posts)-10*(int(stateInt)+1) < 0 {
			posts = posts[:len(posts)-10*int(stateInt)]
		} else {
			posts = posts[len(posts)-10*(int(stateInt)+1) : len(posts)-10*int(stateInt)]
		}
	}
	if posts == nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("You don't have any posts in your diary yet, to add, click Add"))
		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
	} else {
		var str string
		for _, p := range posts {
			tm := p.TimeDate + " " + p.TimeClock
			IDstr := "/p_" + fmt.Sprint(p.ID)
			//тут вылетала ошибка с тем что cap слишком маленький для массива
			//решил ее таким образом
			if (len(p.Body) > 100) && (cap([]rune(p.Body)) > 100) {
				p.Body = string([]rune(p.Body)[:100]) + "...[Click on the ID]"
			}
			arr := []string{IDstr, p.Body, tm}
			str = str + strings.Join(arr, "\n") + "\n\n"
		}
		str = str + "To read/update/delete a post, click on its /p_id"
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s", str))
		if postsMore == true {
			msg.ReplyMarkup = keyboards.CallbackMenuPosts
		}
		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

//DeleteAllPosts удаляет все посты этого пользователя
func DeleteAllPosts(db *sql.DB, userID int) error {
	_, err := db.Exec("DELETE FROM posts WHERE user_id=$1", userID)
	if err != nil {
		return err
	}
	return nil
}

//SearchPost ищет пост по дате
func SearchPost(db *sql.DB, userID int, chatID int64, dateOutside string, month string, bot *tgbotapi.BotAPI) error {
	search, err := common.Search(db, userID, dateOutside, month, true)
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
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("There are no posts on this day"))
		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
	} else {
		//собираем строку из постов по 10 или меньше. И как отдельные сообщения отправляем пользователю.
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
				IDstr := "/p_" + fmt.Sprint(s.ID)
				if len(s.Body) > 100 && (cap([]rune(s.Body)) > 100) {
					s.Body = string([]rune(s.Body)[:100]) + "...[Click on the ID]"
				}
				arr := []string{IDstr, s.Body, tm}
				str = str + strings.Join(arr, "\n") + "\n\n"
			}
			str = str + "To read/update/delete a post, click on its /p_id"
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s", str))
			_, err = bot.Send(msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//CallbackOnePost возвращает клавиатуру с постом
func CallbackOnePost(db *sql.DB, userID int, chatID int64, ID string, bot *tgbotapi.BotAPI) error {
	content, err := common.ReturnOne(db, ID, userID, true)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Something went wrong"))
		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
		return err
	}
	txt := fmt.Sprintf("updateP(%s)\n%s", ID, content.Body)
	var callbackDelUpPost = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Delete", fmt.Sprintf("DeletePost_%s", ID)),
			tgbotapi.InlineKeyboardButton{
				Text:                         "Update",
				SwitchInlineQueryCurrentChat: &txt,
			},
		),
	)
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s\n%s %s", content.Body, content.TimeDate, content.TimeClock))
	msg.ReplyMarkup = callbackDelUpPost
	_, err = bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

//ChangePostsStatePlusOne меняет(+1) state в дб для того чтобы сделать меню вывода постов, по 10 последних шт
func ChangePostsStatePlusOne(rdb *redis.Client, userID int) error {
	state, err := rdb.Get(ctx, fmt.Sprintf("%d_posts_state", userID)).Result()
	if err != nil {
		return err
	}
	stateInt, err := strconv.ParseInt(state, 10, 64)
	if err != nil {
		return err
	}
	err = rdb.Set(ctx, fmt.Sprintf("%d_posts_state", userID), fmt.Sprintf("%d", stateInt+1), 0).Err()
	if err != nil {
		return err
	}
	return nil
}

//ChangePostsStateMinusOne меняет(-1) state в дб для того чтобы сделать меню вывода постов, по 10 последних шт
func ChangePostsStateMinusOne(rdb *redis.Client, userID int) error {
	state, err := rdb.Get(ctx, fmt.Sprintf("%d_posts_state", userID)).Result()
	if err != nil {
		return err
	}
	stateInt, err := strconv.ParseInt(state, 10, 64)
	if err != nil {
		return err
	}
	err = rdb.Set(ctx, fmt.Sprintf("%d_posts_state", userID), fmt.Sprintf("%d", stateInt-1), 0).Err()
	if err != nil {
		return err
	}
	return nil
}

//ChangePostsStateEnd меняет state в начало
func ChangePostsStateEnd(rdb *redis.Client, userID int) error {
	err := rdb.Set(ctx, fmt.Sprintf("%d_posts_state", userID), fmt.Sprintf("%d", 0), 0).Err()
	if err != nil {
		return err
	}
	return nil
}

//ChangePostsStateBegin меняет state на последний
func ChangePostsStateBegin(db *sql.DB, rdb *redis.Client, userID int) error {
	maxID, err := common.ReturnOneMax(db, userID, true)
	if err != nil {
		return err
	}
	maxState := maxID / 10
	err = rdb.Set(ctx, fmt.Sprintf("%d_posts_state", userID), fmt.Sprintf("%d", maxState), 0).Err()
	if err != nil {
		return err
	}
	return nil
}

//GetPostsState получаем state из дб
func GetPostsState(rdb *redis.Client, userID int) (string, error) {
	state, err := rdb.Get(ctx, fmt.Sprintf("%d_posts_state", userID)).Result()
	if err == redis.Nil {
		err := createPostState(rdb, userID)
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
func createPostState(rdb *redis.Client, userID int) error {
	err := rdb.Set(ctx, fmt.Sprintf("%d_posts_state", userID), "0", 0).Err()
	if err != nil {
		return err
	}
	return nil
}
