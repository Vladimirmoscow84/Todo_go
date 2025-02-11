package repeatrule

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// NextDate  вычисляет следующую дату для задачи в соответствии с указанным  правилом
func NextDate(now time.Time, date string, repeat string) (string, error) {

	//Проверка на пустую строку в колонке repeat
	if repeat == "" {
		return "", errors.New("Не определено правило повторения")
	}

	//Парсинг исходного времени, от которго начинается отсчет повторений
	parseDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}
	switch {
	case repeat == "y":
		nextDate := parseDate.Add(1, 0, 0)
		if nextDate.Before(now) {
			return "", errors.New("Несоотвтествие следующей даты относительо текущей")
		}
		return nextDate.Format("20060102"), nil
	case len(repeat) > 2 && string(repeat[0]) == "d":
		strdays, _ := strings.CutPrefix(repeat, "d ")
		days, err := strconv.Atoi(strdays)
		if err != nil && days > 400 {
			return "", errors.New("Неправильное правило  повторения")
		}
		nextDate := parseDate.Add(0, 0, days)
		if nextDate.Before(now) {
			return "", errors.New("Несоотвтествие следующей даты относительо текущей")
		}
		return nextDate.Format("20060102"), nil
	default:
		return "", errors.New("Неподдерживаемое правило повторений")
	}

}
