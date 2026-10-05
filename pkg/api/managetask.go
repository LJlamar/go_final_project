package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/LJlamar/go_final_project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	data, err := db.GetTask(id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			writeJson(w, map[string]string{"error": "задача не найдена"})
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			writeJson(w, map[string]string{"error": "ошибка базы данных"})
			return
		}
	}
	writeJson(w, data)

}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task *db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "Некорректный формат JSON"})
		return
	}

	//Checking for non-empty title
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	now := time.Now()

	//Checking if given date is empty. If so - use current date.
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	//Checking for right value of given date
	_, err = time.Parse("20060102", task.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "Дата указана в отличном от требуемого (YYYYMMDD) формате"})
		return
	}

	if task.Repeat == "" {
		now := now.Format("20060102")
		if task.Date < now {
			task.Date = now
		}
	} else if task.Date < now.Format("20060102") {

		newDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Printf("Error with adding/editing event: %v", err)

			w.WriteHeader(http.StatusBadRequest)
			writeJson(w, map[string]string{"error": "Поле repeat заполнено неверно"})
			return
		}
		task.Date = newDate
	}

	//Executing database command
	result, err := db.AddTask(task)

	//Returning id or error
	if err != nil {
		//Creating map to serialise into JSON
		responseData := map[string]string{
			//Structure prior to JSON:
			"error": err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, responseData)
	} else {
		idStr := strconv.Itoa(int(result))
		responseData := map[string]string{
			//Structure prior to JSON:
			"id": idStr,
		}
		w.WriteHeader(http.StatusOK)
		writeJson(w, responseData)
	}
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task *db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "Некорректный формат JSON"})
		return
	}

	//Checking for non-empty title
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	//Checking date
	now := time.Now()
	//Checking if given date is empty. If so - use current date's value.
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	//Checking for right formatting and value of given date:
	_, err = time.Parse("20060102", task.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "Дата указана в отличном от требуемого (YYYYMMDD) формате"})
		return
	}
	//Check if new given date value differs from previous one in database,
	//with supplement function
	dateChanged, err := db.ChangedDate(task)

	if dateChanged == false && err != nil {
		//Creating map to serialise into JSON
		responseData := map[string]string{
			//Structure prior to JSON:
			"error": err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, responseData)
		return

		//Calculating new valid date based on given new value:
	} else if dateChanged == true && err == nil {
		if task.Repeat == "" {
			now := now.Format("20060102")
			if task.Date < now {
				task.Date = now
			}
		} else {
			newDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				log.Printf("Error with adding/editing event: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				writeJson(w, map[string]string{"error": "Поле repeat заполнено неверно"})
				return
			}
			task.Date = newDate
		}
	}

	err = db.UpdateTask(task)

	//Returning empty value or error
	if err != nil {
		//Creating map to serialise into JSON
		responseData := map[string]string{
			//Structure prior to JSON:
			"error": err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, responseData)
	} else {
		writeJson(w, map[string]interface{}{})
	}
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "Внутренняя ошибка сервера"})
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
		log.Printf("Error with adding/editing event: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, map[string]interface{}{})
}

func clearTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "Внутренняя ошибка сервера"})
		return
	}

	//Getting task by id
	task, err := db.GetTask(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": "Внутренняя ошибка сервера"})
		return
	}
	// Returning empty value or error
	//Updating task's date if it has repeats
	if task.Repeat != "" {
		newDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			log.Printf("Error with adding/editing event: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			writeJson(w, map[string]string{"error": "Поле repeat заполнено неверно"})
			return
		}
		task.Date = newDate

		err = db.UpdateDate(task)
		if err != nil {
			//Creating map to serialise into JSON
			responseData := map[string]string{
				//Structure prior to JSON:
				"error": err.Error(),
			}
			w.WriteHeader(http.StatusInternalServerError)
			writeJson(w, responseData)
		} else {
			writeJson(w, map[string]interface{}{})
		}
		//Deleting completed task from database
	} else {
		err := db.DeleteTask(task.ID)
		if err != nil {
			log.Printf("Error with adding/editing event: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		writeJson(w, map[string]interface{}{})
	}
}
