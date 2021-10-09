package common

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/elolpuer/DiaryBot/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//isPost определяет пост ли это
//значение true означает что пост
//значение false означает что идея

//Insert добавление поста в бд
func Insert(db *sql.DB, insrt *models.PostOrIdea, isPost bool) error {
	var ID int
	var err error
	var slct []*models.PostOrIdea
	if isPost == true {
		slct, err = SelectAll(db, insrt.UserID, isPost)
	} else {
		slct, err = SelectAll(db, insrt.UserID, isPost)
	}
	if err != nil {
		return err
	}
	//если есть посты в бд то прибавляем 1 к максимальному id и получаем id для нового поста
	//если их нет то id = 1
	if len(slct) != 0 {
		ID = len(slct) + 1
	} else {
		ID = 1
	}
	var str string
	if isPost == true {
		str = "INSERT INTO posts (id, user_id, body, time_date, time_clock) VALUES ($1, $2, $3, $4, $5)"
	} else {
		str = "INSERT INTO ideas (id, user_id, body, time_date, time_clock) VALUES ($1, $2, $3, $4, $5)"
	}
	_, err = db.Exec(str, ID, insrt.UserID, insrt.Body, insrt.TimeDate, insrt.TimeClock)
	if err != nil {
		return err
	}
	return nil
}

//SelectAll выдает все посты в бд по этому пользователю
func SelectAll(db *sql.DB, ID int, isPost bool) ([]*models.PostOrIdea, error) {
	var str string
	if isPost == true {
		str = "SELECT * FROM posts WHERE user_id=$1"
	} else {
		str = "SELECT * FROM ideas WHERE user_id=$1"
	}
	rows, err := db.Query(str, ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	slct := make([]*models.PostOrIdea, 0)
	for rows.Next() {
		s := new(models.PostOrIdea)
		err := rows.Scan(&s.ID, &s.UserID, &s.Body, &s.TimeDate, &s.TimeClock)
		if err != nil {
			return nil, err
		}
		slct = append(slct, s)
	}
	if len(slct) == 0 {
		return nil, nil
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return sort(slct), nil
}

//Delete удаляет пост из бд
func Delete(db *sql.DB, ID int, userID int, isPost bool) (int64, error) {
	fmt.Println(ID, userID)
	var strDelete string
	var strUpdate string
	//удаляем пост и обновляем id у постов, у которых id больше чем у удаленного
	//id больших равно id - 1
	if isPost == true {
		strDelete = "DELETE FROM posts WHERE id=$1 AND user_id=$2"
		strUpdate = "UPDATE posts SET id=id-1 WHERE id>$1 AND user_id=$2"
	} else {
		strDelete = "DELETE FROM ideas WHERE id=$1 AND user_id=$2"
		strUpdate = "UPDATE ideas SET id=id-1 WHERE id>$1 AND user_id=$2"
	}
	var err error
	res, err := db.Exec(strDelete, ID, userID)
	if err != nil {
		return 0, err
	}
	num, err := res.RowsAffected()
	fmt.Println(num)
	if num == 0 {
		return 1, err
	}
	if err != nil {
		return 0, err
	}
	_, err = db.Exec(strUpdate, ID, userID)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

//UpdatePostOrIdea ...
//Берем между скобками ID поста или идеи чтобы его обновить
//txtWhat показывает что это пост или идея
//txtId - ID
//txtT - текст раньше, до обновления.
func UpdatePostOrIdea(db *sql.DB, text string, userID int, chatID int64, bot *tgbotapi.BotAPI) error {
	s1 := strings.Index(text, "(")
	s2 := strings.Index(text, ")")
	txtWhat := text[s1-1]
	txtID := text[s1+1 : s2]
	txtT := text[s2+2:]
	var err error
	if string(txtWhat) == "P" {
		err = update(db, txtID, txtT, userID, true)
	}
	if string(txtWhat) == "I" {
		err = update(db, txtID, txtT, userID, false)
	}
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Something went wrong"))
		_, err = bot.Send(msg)
		if err != nil {
			return err
		}
		return err
	}
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s", "Updated"))
	_, err = bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

//update изменяет пост в бд
func update(db *sql.DB, ID string, text string, userID int, isPost bool) error {
	var str string
	if isPost == true {
		str = "UPDATE posts SET body=$1 WHERE id=$2 AND user_id=$3"
	} else {
		str = "UPDATE ideas SET body=$1 WHERE id=$2 AND user_id=$3"
	}
	_, err := db.Exec(str, text, ID, userID)
	if err != nil {
		return err
	}
	return nil
}

//ReturnOne для update
func ReturnOne(db *sql.DB, ID string, userID int, isPost bool) (*models.PostOrIdea, error) {
	content := new(models.PostOrIdea)
	var str string
	if isPost == true {
		str = "SELECT body,time_date,time_clock FROM posts WHERE id=$1 AND user_id=$2"
	} else {
		str = "SELECT body,time_date,time_clock FROM ideas WHERE id=$1 AND user_id=$2"
	}
	row := db.QueryRow(str, ID, userID)

	if err := row.Scan(&content.Body, &content.TimeDate, &content.TimeClock); err == sql.ErrNoRows {
		return nil, err
	}

	return content, nil

}

//ReturnOneMax для state
func ReturnOneMax(db *sql.DB, userID int, isPost bool) (int, error) {
	var maxID int
	var str string
	if isPost == true {
		str = "SELECT MAX(id) FROM posts WHERE user_id=$1"
	} else {
		str = "SELECT MAX(id) FROM ideas WHERE user_id=$1"
	}
	row := db.QueryRow(str, userID)

	if err := row.Scan(&maxID); err == sql.ErrNoRows {
		return 0, err
	}

	return maxID, nil

}

//Search поиск поста по дате
func Search(db *sql.DB, userID int, date, month string, isPost bool) ([]*models.PostOrIdea, error) {
	str := date + " " + month + " " + "2021"
	var query string
	if isPost == true {
		query = "SELECT * FROM posts WHERE user_id=$1 AND time_date=$2"
	} else {
		query = "SELECT * FROM ideas WHERE user_id=$1 AND time_date=$2"
	}
	rows, err := db.Query(query, userID, str)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	search := make([]*models.PostOrIdea, 0)
	for rows.Next() {
		s := new(models.PostOrIdea)
		err := rows.Scan(&s.ID, &s.UserID, &s.Body, &s.TimeDate, &s.TimeClock)
		if err != nil {
			return nil, err
		}
		search = append(search, s)
	}
	if len(search) == 0 {
		return nil, nil
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return sort(search), nil
}

//MakeStringDate делает из времение(даты) строку
func MakeStringDate(time time.Time) string {
	y, m, d := time.Date()
	var month string
	switch fmt.Sprint(m) {
	case "January":
		month = "January"
	case "February":
		month = "February"
	case "March":
		month = "March"
	case "April":
		month = "April"
	case "May":
		month = "May"
	case "June":
		month = "June"
	case "July":
		month = "July"
	case "August":
		month = "August"
	case "September":
		month = "September"
	case "October":
		month = "October"
	case "November":
		month = "November"
	case "December":
		month = "December"
	}
	s := []string{fmt.Sprint(d), month, fmt.Sprint(y)}
	return strings.Join(s, " ")
}

//MakeStringClock делает из времение(времени) строку
func MakeStringClock(time time.Time) string {
	var sm string
	nums := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	h, m, _ := time.Clock()
	for _, num := range nums {
		if fmt.Sprint(m) == num {
			sm = "0" + fmt.Sprint(m)
			s := []string{fmt.Sprint(h), sm}
			return strings.Join(s, ":")
		}
	}
	s := []string{fmt.Sprint(h), fmt.Sprint(m)}
	return strings.Join(s, ":")
}

//сортирует по id от меньшего к большему
func sort(arr []*models.PostOrIdea) []*models.PostOrIdea {
	for i := 0; i < len(arr); i++ {
		for j := len(arr) - 1; j > i; j-- {
			if arr[j-1].ID > arr[j].ID {
				x := arr[j-1]
				arr[j-1] = arr[j]
				arr[j] = x
			}
		}
	}
	return arr
}

func AddToStart(db *sql.DB, userID int) error {
	_, err := db.Exec("INSERT INTO start (user_id) VALUES ($1)", userID)
	if err != nil {
		return err
	}
	return nil
}
