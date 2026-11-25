package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

// NextDate вычисляет следующую дату задачи в формате 20060102.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата %q: %w", dstart, err)
	}

	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "d":
		return nextDateByDays(now, date, parts)
	case "y":
		return nextDateByYears(now, date, parts)
	case "w":
		return nextDateByWeekdays(now, date, parts)
	case "m":
		return nextDateByMonths(now, date, parts)
	default:
		return "", fmt.Errorf("неподдерживаемый формат правила: %q", repeat)
	}
}

// afterNow возвращает true, если date строго больше now (сравнение по дате).
func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	return time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC).
		After(time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC))
}

// nextDateByDays обрабатывает правило "d <число>".
func nextDateByDays(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("правило d должно содержать интервал в днях")
	}
	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("некорректный интервал в днях: %w", err)
	}
	if interval < 1 || interval > 400 {
		return "", fmt.Errorf("интервал в днях должен быть от 1 до 400, получено %d", interval)
	}
	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

// nextDateByYears обрабатывает правило "y".
func nextDateByYears(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 1 {
		return "", errors.New("правило y не принимает параметров")
	}
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

// nextDateByWeekdays обрабатывает правило "w <дни недели>".
func nextDateByWeekdays(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("правило w должно содержать дни недели")
	}
	var days [8]bool
	for _, s := range strings.Split(parts[1], ",") {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > 7 {
			return "", fmt.Errorf("некорректный день недели: %q", s)
		}
		days[n] = true
	}
	for {
		date = date.AddDate(0, 0, 1)
		wd := int(date.Weekday())
		if wd == 0 {
			wd = 7
		}
		if days[wd] && afterNow(date, now) {
			break
		}
	}
	return date.Format(dateFormat), nil
}

func nextDateByMonths(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", errors.New("правило m должно содержать дни месяца и опционально месяцы")
	}

	daySet := make(map[int]bool)
	for _, s := range strings.Split(parts[1], ",") {
		n, err := strconv.Atoi(s)
		if err != nil || n == 0 || n < -2 || n > 31 {
			return "", fmt.Errorf("некорректный день месяца: %q", s)
		}
		daySet[n] = true
	}

	var months [13]bool
	allMonths := true
	if len(parts) == 3 {
		allMonths = false
		for _, s := range strings.Split(parts[2], ",") {
			n, err := strconv.Atoi(s)
			if err != nil || n < 1 || n > 12 {
				return "", fmt.Errorf("некорректный месяц: %q", s)
			}
			months[n] = true
		}
	}

	const maxIter = 366 * 5
	for i := 0; i < maxIter; i++ {
		date = date.AddDate(0, 0, 1)

		if !allMonths && !months[int(date.Month())] {
			continue
		}

		day := date.Day()
		lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
		match := false
		for n := range daySet {
			switch {
			case n > 0 && day == n:
				match = true
			case n == -1 && day == lastDay:
				match = true
			case n == -2 && day == lastDay-1:
				match = true
			}
			if match {
				break
			}
		}

		if match && afterNow(date, now) {
			return date.Format(dateFormat), nil
		}
	}
	return "", errors.New("не найдена подходящая дата в указанных месяцах")
}

// nextDateHandler обрабатывает GET /api/nextdate?now=&date=&repeat=
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "некорректный параметр now: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, next)
}
