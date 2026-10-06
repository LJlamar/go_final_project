package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		now = time.Now()
	}

	if dateStr == "" {
		http.Error(w, "date parameter is crucial, no argument is given", http.StatusBadRequest)
		return
	}
	response, err := NextDate(now, dateStr, repeatStr)

	if err != nil {
		log.Printf("Error with calculating next date: %v", err)
		http.Error(w, fmt.Sprintf("next date was not calculated: %v", err), http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}
}

// Getting valid "not past" time: true if 'date" is after "today"
func afterNow(givenTime time.Time, currentTime time.Time) bool {
	return givenTime.After(currentTime)
}

// Calculating year and month's day values for upcoming date
func calculatedYearAndDay(timeValue time.Time, month int, value string) ([2]int, error) {
	//Setting array with year and date
	arr := [2]int{0, 0}

	//Current month's(m) index correction, based on arriving of new year
	if month > 12 {
		month = month - 12
	}
	//Getting year, varies with current month's(m) calculated index
	y := timeValue.Year()
	if month < int(timeValue.Month()) {
		y += 1
	}

	//Getting amount of days in current calculated month(m)
	//Reminder - for instance, func time.Date(2027, 1, 0, ...)
	//would return in result 0.01.2027 -> 31.12.2026,
	//so trying to get t later, without month's correction would return
	//amount of previous month's days
	t := time.Date(y, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC)
	daysAmountInMonth := t.Day()

	//Checking for type of given task for adding dates:
	//number of the day taken by substracting value in task from current month's quantity or
	if strings.HasPrefix(value, "-") {
		repeatDaysSubstracted, err := strconv.Atoi(strings.TrimLeft(value, "-"))
		if err != nil {
			return arr, err
		}
		//Returning error if "repeat" value is out of task's bounds:
		if repeatDaysSubstracted < 1 || repeatDaysSubstracted > 2 {
			return arr, errors.New("given date value is out of month's bounds")
		}
		//Checking if given day is suitable for current month
		//then adding it to resulting array:

		//If given month is equal to current month and
		if month == int(timeValue.Month()) {
			//if current day is less than resulted repeat date then
			//add resulted repeat date to array
			if timeValue.Day() < daysAmountInMonth-(repeatDaysSubstracted-1) {
				arr[0] = y
				arr[1] = daysAmountInMonth - (repeatDaysSubstracted - 1)
			}
			//Adding repeated day in default(not current month) cases:
		} else {
			arr[0] = y
			arr[1] = daysAmountInMonth - (repeatDaysSubstracted - 1)
		}

		//number of the day in month by it's order
	} else {
		repeatedDayInMonth, err := strconv.Atoi(value)
		if err != nil {
			//return arr, err
		}
		//Returning error if "repeat" value is out of task's bounds:
		if repeatedDayInMonth < 1 || repeatedDayInMonth > 31 {
			return arr, errors.New("given date value is out of month's bounds")
		}
		//Checking if given day is suitable for current month
		//then adding it to resulting array:
		if repeatedDayInMonth <= daysAmountInMonth {
			//If given month is equal to current month and
			if month == int(timeValue.Month()) {
				//if repeated day in this month is more than current day
				//we add it to array
				if repeatedDayInMonth > timeValue.Day() {
					arr[0] = y
					arr[1] = repeatedDayInMonth
				}
				//Adding repeated day in default(not current month) cases:
			} else {
				arr[0] = y
				arr[1] = repeatedDayInMonth
			}
		}
	}
	return arr, nil
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", errors.New("the repeat parameter's argument is empty string")
	}
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", errors.New("given date argument is not valid, reqiured format: YYYYMMDD")
	}

	switch strings.Count(repeat, " ") {
	case 0:
		if repeat != "y" {
			return "", errors.New("given repeat argument is not valid")
		} else {
			for {
				date = date.AddDate(1, 0, 0)
				if afterNow(date, now) {
					break
				}
			}
		}
	case 1:
		//Splitting repeat string into array
		repeatValues := strings.Split(repeat, " ")

		switch repeatValues[0] {
		default:
			return "", errors.New("given timetype argument is not valid, required: day/week/month")

		case "d":
			//Getting repeat value
			repeatDayIndex, err := strconv.Atoi(repeatValues[1])
			if err != nil {
				return "", err
			}

			//Array for multiple days:
			//var dateList []string

			if repeatDayIndex < 1 || repeatDayIndex > 400 {
				return "", errors.New("given day value is out of required bounds")
			}
			for {
				date = date.AddDate(0, 0, repeatDayIndex)
				if afterNow(date, now) {
					break
				}
			}
			//}
			//sort.Strings(dateList)
			//return dateList[0], nil //this string is for ONE VALUE

			//dateString := strings.Join(dateList, ";")            //this string is for ALL VALUES
			//return dateString, nil                               //this string is for ALL VALUES

		case "w":
			//Defining latest date among given values
			var latestDate time.Time

			if afterNow(date, now) {
				latestDate = date
			} else {
				latestDate = now
			}
			//Getting day of the week from current date and
			//setting it from 0-6 to 1-7, due to task's demands

			//Default case - new Sunday (7 instead of 0)
			currentWeekdayInNumber := 7
			//If day of the week isn't Sunday:
			if int(latestDate.Weekday()) != 0 {
				currentWeekdayInNumber = int(latestDate.Weekday())
			}
			//New variable for finding closest suitable day
			gap := 0

			repeatDaysIndexes := strings.Split(repeatValues[1], ",")

			for _, value := range repeatDaysIndexes {
				v, err := strconv.Atoi(value)
				if err != nil {
					return "", err
				}
				//Cheking for wrong values:
				if v < 1 || v > 7 {
					return "", errors.New("given week's day value is out of bounds")
				}

				g := 0
				switch {
				case currentWeekdayInNumber > v:
					g = (7 - currentWeekdayInNumber) + v
				case currentWeekdayInNumber < v:
					g = v - currentWeekdayInNumber
				case currentWeekdayInNumber == v:
					g = 7
				}
				if gap == 0 || g < gap {
					gap = g
				}
			}

			date = latestDate.AddDate(0, 0, gap)

		case "m":
			//Defining latest date among given values
			var latestDate time.Time

			if afterNow(date, now) {
				latestDate = date
			} else {
				latestDate = now
			}
			//Case of multiple days, to pick closest one in current month, and all dates in next one:
			var dateList []string
			repeatDaysIndexes := strings.Split(repeatValues[1], ",")

			//Adding dates in every month, for next one year
			for i := 0; i < 12; i++ {
				//Adding dates for each given month
				for _, value := range repeatDaysIndexes {

					//Getting month of new date, varies with indexed month(i)
					m := int(latestDate.Month()) + i

					arr, err := calculatedYearAndDay(latestDate, m, value)
					if err != nil {
						return "", err
					}

					//Current month's(m) index correction, based on arriving of new year
					month := m
					if month > 12 {
						month = month - 12
					}

					if arr[1] != 0 {
						d := time.Date(arr[0], time.Month(month), arr[1], 0, 0, 0, 0, time.UTC)

						dateList = append(dateList, d.Format("20060102"))
					}
				}
			}
			sort.Strings(dateList)
			//fmt.Printf("Список дат помесячно: \n%v\n", dateList)
			//fmt.Printf("Первая дата: %v\n", dateList[0])
			return dateList[0], nil //this string is for ONE VALUE
			//dateString := strings.Join(dateList, ";")            //this string is for ALL VALUES
			//return dateString, nil                               //this string is for ALL VALUES
		}

	case 2:
		//Defining latest date among given values
		var latestDate time.Time

		if afterNow(date, now) {
			latestDate = date
		} else {
			latestDate = now
		}

		repeatValues := strings.Split(repeat, " ")

		repeatDaysIndexes := strings.Split(repeatValues[1], ",")

		repeatMonthsIndexes := strings.Split(repeatValues[2], ",")

		var dateList []string

		//Adding dates in every month, for next one year
		for _, month := range repeatMonthsIndexes {
			//Adding dates for each given month
			for _, value := range repeatDaysIndexes {

				//Getting month of new date, varies with indexed month(i)
				m, err := strconv.Atoi(month)
				if err != nil {
					return "", err
				}

				if m < 1 || m > 12 {
					return "", errors.New("given month's value is out of reasonable bounds")
				}
				arr, err := calculatedYearAndDay(latestDate, m, value)
				if err != nil {
					return "", err
				}
				if arr[1] != 0 {
					d := time.Date(arr[0], time.Month(m), arr[1], 0, 0, 0, 0, time.UTC)
					dateList = append(dateList, d.Format("20060102")) //this string is for ALL VALUES
				}
			}
		}
		sort.Strings(dateList)
		return dateList[0], nil //this string is for ONE VALUE
		//dateString := strings.Join(dateList, ";") //this string is for ALL VALUES
		//return dateString, nil //this string is for ALL VALUES
	}

	return date.Format("20060102"), nil
}
