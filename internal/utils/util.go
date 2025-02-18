package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate  вычисляет следующую дату для задачи в соответствии с указанным  правилом
func NextDate(now time.Time, date string, repeat string) (string, error) {

	//Проверка на пустую строку в колонке repeat
	if repeat == "" {
		return "", errors.New("не определено правило повторения")
	}

	//Парсинг исходного времени, от которго начинается отсчет повторений

	parseDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("некорректный формат времени: %s", err.Error())
	}
	switch {
	case repeat == "y":
		nextDate := parseDate.AddDate(1, 0, 0)
		for nextDate.Before(now) {

			nextDate = nextDate.AddDate(1, 0, 0)
		}
		return nextDate.Format("20060102"), nil
	case len(repeat) > 2 && string(repeat[0]) == "d":
		strdays, _ := strings.CutPrefix(repeat, "d ")
		days, err := strconv.Atoi(strdays)
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("неправильное правило повторения")
		}
		nextDate := parseDate.AddDate(0, 0, days)
		for nextDate.Before(now) {

			nextDate = nextDate.AddDate(0, 0, days)
		}
		return nextDate.Format("20060102"), nil
	default:
		return "", errors.New("неподдерживаемое правило повторений")
	}

}
