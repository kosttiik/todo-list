package tasks

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kosttiik/todo-list/ds"
)

func NormalizeToDate(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func IsLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("Не указан repeat")
	}

	nextDate, err := time.ParseInLocation(ds.DateFormat, date, time.UTC)
	if err != nil {
		return "", fmt.Errorf("Неверный формат даты: %v", err)
	}

	now = NormalizeToDate(now.UTC())
	nextDate = NormalizeToDate(nextDate.UTC())

	repeatRule := strings.Fields(repeat)

	switch repeatRule[0] {
	case "d":
		if len(repeatRule) < 2 {
			return "", errors.New("Не указано количество дней")
		}

		days, err := strconv.Atoi(repeatRule[1])
		if err != nil {
			return "", fmt.Errorf("Неверное значение для дней: %w", err)
		}

		if days > 400 {
			return "", errors.New("Количество дней не должно превышать 400")
		}

		nextDate = nextDate.AddDate(0, 0, days)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		return nextDate.Format(ds.DateFormat), nil

	case "y":
		nextDate = nextDate.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}

		if nextDate.Month() == time.February && nextDate.Day() == 29 && !IsLeapYear(nextDate.Year()) {
			nextDate = time.Date(nextDate.Year(), time.March, 1, 0, 0, 0, 0, time.UTC)
		}
		return nextDate.Format(ds.DateFormat), nil

	default:
		return "", fmt.Errorf("Неподдерживаемый формат правила: %s", repeatRule[0])
	}
}

func HandleNextDate(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	if nowStr == "" || dateStr == "" || repeatStr == "" {
		http.Error(w, "Отсутствуют обязательные параметры", http.StatusBadRequest)
		return
	}

	now, err := time.ParseInLocation(ds.DateFormat, nowStr, time.UTC)
	if err != nil {
		http.Error(w, "Неверный формат параметра now", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}
