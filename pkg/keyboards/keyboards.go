package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Keyboard основная клавиатура
var Keyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Add"),
		tgbotapi.NewKeyboardButton("Read"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Break"),
		tgbotapi.NewKeyboardButton("Help"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Search"),
	),
)

//KeyboardDelete кнопки для того чтобы узнать хочет пользователь удалять что-либо ил нет
var KeyboardDelete = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Yes"),
		tgbotapi.NewKeyboardButton("No"),
	),
)

//CallbackAdd меню выбора того, что добавить, пост или идею
var CallbackAdd = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Post", "Post"),
		tgbotapi.NewInlineKeyboardButtonData("Idea", "Idea"),
	),
)

//CallbackAll меню выбора того, что посмотреть, посты или идеи
var CallbackAll = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Posts", "AllPosts"),
		tgbotapi.NewInlineKeyboardButtonData("Ideas", "AllIdeas"),
	),
)

//CallbackMenuPosts меню в постах
var CallbackMenuPosts = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Start", "MenuEndPost"),
		tgbotapi.NewInlineKeyboardButtonData("<<", "MenuMinusPost"),
		tgbotapi.NewInlineKeyboardButtonData(">>", "MenuPlusPost"),
		tgbotapi.NewInlineKeyboardButtonData("End", "MenuBeginPost"),
	),
)

//CallbackMenuIdeas меню в идеях
var CallbackMenuIdeas = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Start", "MenuEndIdea"),
		tgbotapi.NewInlineKeyboardButtonData("<<", "MenuMinusIdea"),
		tgbotapi.NewInlineKeyboardButtonData(">>", "MenuPlusIdea"),
		tgbotapi.NewInlineKeyboardButtonData("End", "MenuBeginIdea"),
	),
)

//CallbackSearch меню выбора того, что удалить, посты или идеи
var CallbackSearch = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Posts", "SearchPost"),
		tgbotapi.NewInlineKeyboardButtonData("Ideas", "SearchIdea"),
	),
)

//CallbackDatePost меню выбора даты для поиска
var CallbackDatePost = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("1", "day_P1"),
		tgbotapi.NewInlineKeyboardButtonData("2", "day_P2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "day_P3"),
		tgbotapi.NewInlineKeyboardButtonData("4", "day_P4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "day_P5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "day_P6"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("7", "day_P7"),
		tgbotapi.NewInlineKeyboardButtonData("8", "day_P8"),
		tgbotapi.NewInlineKeyboardButtonData("9", "day_P9"),
		tgbotapi.NewInlineKeyboardButtonData("10", "day_P10"),
		tgbotapi.NewInlineKeyboardButtonData("11", "day_P11"),
		tgbotapi.NewInlineKeyboardButtonData("12", "day_P12"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("13", "day_P13"),
		tgbotapi.NewInlineKeyboardButtonData("14", "day_P14"),
		tgbotapi.NewInlineKeyboardButtonData("15", "day_P15"),
		tgbotapi.NewInlineKeyboardButtonData("16", "day_P16"),
		tgbotapi.NewInlineKeyboardButtonData("17", "day_P17"),
		tgbotapi.NewInlineKeyboardButtonData("18", "day_P18"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("19", "day_P19"),
		tgbotapi.NewInlineKeyboardButtonData("20", "day_P20"),
		tgbotapi.NewInlineKeyboardButtonData("21", "day_P21"),
		tgbotapi.NewInlineKeyboardButtonData("22", "day_P22"),
		tgbotapi.NewInlineKeyboardButtonData("23", "day_P23"),
		tgbotapi.NewInlineKeyboardButtonData("24", "day_P24"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("25", "day_P25"),
		tgbotapi.NewInlineKeyboardButtonData("26", "day_P26"),
		tgbotapi.NewInlineKeyboardButtonData("27", "day_P27"),
		tgbotapi.NewInlineKeyboardButtonData("28", "day_P28"),
		tgbotapi.NewInlineKeyboardButtonData("29", "day_P29"),
		tgbotapi.NewInlineKeyboardButtonData("30", "day_P30"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("31", "day_P31"),
	),
)

//CallbackDateIdea меню выбора даты для поиска
var CallbackDateIdea = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("1", "day_I1"),
		tgbotapi.NewInlineKeyboardButtonData("2", "day_I2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "day_I3"),
		tgbotapi.NewInlineKeyboardButtonData("4", "day_I4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "day_I5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "day_I6"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("7", "day_I7"),
		tgbotapi.NewInlineKeyboardButtonData("8", "day_I8"),
		tgbotapi.NewInlineKeyboardButtonData("9", "day_I9"),
		tgbotapi.NewInlineKeyboardButtonData("10", "day_I10"),
		tgbotapi.NewInlineKeyboardButtonData("11", "day_I11"),
		tgbotapi.NewInlineKeyboardButtonData("12", "day_I12"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("13", "day_I13"),
		tgbotapi.NewInlineKeyboardButtonData("14", "day_I14"),
		tgbotapi.NewInlineKeyboardButtonData("15", "day_I15"),
		tgbotapi.NewInlineKeyboardButtonData("16", "day_I16"),
		tgbotapi.NewInlineKeyboardButtonData("17", "day_I17"),
		tgbotapi.NewInlineKeyboardButtonData("18", "day_I18"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("19", "day_I19"),
		tgbotapi.NewInlineKeyboardButtonData("20", "day_I20"),
		tgbotapi.NewInlineKeyboardButtonData("21", "day_I21"),
		tgbotapi.NewInlineKeyboardButtonData("22", "day_I22"),
		tgbotapi.NewInlineKeyboardButtonData("23", "day_I23"),
		tgbotapi.NewInlineKeyboardButtonData("24", "day_I24"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("25", "day_I25"),
		tgbotapi.NewInlineKeyboardButtonData("26", "day_I26"),
		tgbotapi.NewInlineKeyboardButtonData("27", "day_I27"),
		tgbotapi.NewInlineKeyboardButtonData("28", "day_I28"),
		tgbotapi.NewInlineKeyboardButtonData("29", "day_I29"),
		tgbotapi.NewInlineKeyboardButtonData("30", "day_I30"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("31", "day_I31"),
	),
)

//CallbackMonth меню выбора даты для поиска
var CallbackMonth = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("January", "month_January"),
		tgbotapi.NewInlineKeyboardButtonData("February", "month_February"),
		tgbotapi.NewInlineKeyboardButtonData("March", "month_March"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("April", "month_April"),
		tgbotapi.NewInlineKeyboardButtonData("May", "month_May"),
		tgbotapi.NewInlineKeyboardButtonData("June", "month_June"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("July", "month_July"),
		tgbotapi.NewInlineKeyboardButtonData("August", "month_August"),
		tgbotapi.NewInlineKeyboardButtonData("September", "month_September"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("October", "month_October"),
		tgbotapi.NewInlineKeyboardButtonData("November", "month_November"),
		tgbotapi.NewInlineKeyboardButtonData("December", "month_December"),
	),
)
