package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func Tasks(limit int, searchStr string) ([]*Task, error) {

	if db == nil {
		return nil, fmt.Errorf("`error: value of db variable equals nil! Database is not connected`")
	}
	//Tasks
	tasksList := make([]*Task, 0, limit)

	//Base query without filter applied
	query := "SELECT * FROM scheduler ORDER BY date LIMIT ?"

	//Arguments
	var args []interface{}

	//Checking if search filter not empty
	if searchStr != "" {
		//Defining type of search - by context or by date:
		if strings.Count(searchStr, ".") == 2 {
			//Checking for valid date input and converting date to universal type
			filter, err := time.Parse("02.01.2006", searchStr)
			if err != nil {
				return tasksList, errors.New("ошибка поиска по дате: некорректные данные")
			}
			//Converting date into type suitable for database
			searchStr = filter.Format("20060102")
			fmt.Println(searchStr)

			query = "SELECT * FROM scheduler WHERE date = ? LIMIT ?"
			args = append(args, searchStr, limit)
			/*
				rows, err := db.Query(query, searchStr, limit)
				if err != nil {
					return tasksList, err
				}
				defer rows.Close()
			*/
			/*for rows.Next() {
				//Пересмотри потом для чего это:
				task := &Task{}
				err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
				if err != nil {
					return tasksList, err
				}
				//Adding current entry to our slice:
				tasksList = append(tasksList, task)
				if err := rows.Err(); err != nil {
					return tasksList, err
				}
			}*/
		} else {
			query = "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER by date LIMIT ?"
			searchStr = "%" + searchStr + "%"
			args = append(args, searchStr, searchStr, limit)
			/*
				rows, err := db.Query(query, searchStr, searchStr, limit)
				if err != nil {
					return tasksList, err
				}
				defer rows.Close()
			*/
			/*for rows.Next() {
				//Пересмотри потом для чего это:
				task := &Task{}
				err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
				if err != nil {
					return tasksList, err
				}
				//Adding current entry to our slice:
				tasksList = append(tasksList, task)
				if err := rows.Err(); err != nil {
					return tasksList, err
				}
			}*/
		}
	} else {
		args = append(args, limit)
		/*
			//Executing base query:
			rows, err := db.Query(query, limit)
			if err != nil {
				return tasksList, err
			}
			defer rows.Close()
		*/
	}

	//fmt.Printf("Выполняем запрос: %s\n", query)
	//fmt.Printf("Аргументы: %v\n", args)

	rows, err := db.Query(query, args...)
	if err != nil {
		return tasksList, err
	}
	defer rows.Close()

	for rows.Next() {
		//Пересмотри потом для чего это:
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasksList, err
		}
		//Adding current entry to our slice:
		tasksList = append(tasksList, task)
		if err := rows.Err(); err != nil {
			return tasksList, err
		}
	}
	//fmt.Printf("Итоговый список в слайсе: %v\n", tasksList)
	return tasksList, nil
}

// Supportive function for operation of task's not crucial postitons updating:
func ChangedDate(task *Task) (bool, error) {
	if db == nil {
		return false, fmt.Errorf("`error: value of db variable equals nil! Database is not connected`")
	}

	query := "SELECT date FROM scheduler WHERE ID = ?"

	var t Task

	idInt, err := strconv.Atoi(task.ID)
	if err != nil {
		return false, errors.New("ошибка выборки по идентификатору: некорректные исходные данные")
	}

	err = db.QueryRow(query, idInt).Scan(&t.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, sql.ErrNoRows
		} else {
			return false, err
		}
	}
	if t.Date != task.Date {
		return true, nil
	}
	return false, nil
}

func GetTask(id string) (*Task, error) {
	if db == nil {
		return nil, fmt.Errorf("`error: value of db variable equals nil! Database is not connected`")
	}

	query := "SELECT * FROM scheduler WHERE ID = ?"

	var t Task

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("ошибка выборки по идентификатору: некорректные исходные данные")
	}

	err = db.QueryRow(query, idInt).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		} else {
			return nil, err
		}
	}
	return &t, nil

}

func AddTask(task *Task) (int64, error) {
	if db == nil {
		return 0, fmt.Errorf(`error: value of db variable equals nil! Database is not connected`)
	}
	var id int64

	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func UpdateTask(task *Task) error {
	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func UpdateDate(task *Task) error {
	// параметры пропущены, не забудьте указать WHERE
	query := "UPDATE scheduler SET date = ? WHERE id = ?"
	res, err := db.Exec(query, task.Date, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`error: incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {
	if db == nil {
		return fmt.Errorf(`error: value of db variable equals nil! Database is not connected`)
	}

	query := "DELETE FROM scheduler WHERE ID = ?"

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return errors.New(`error of converting id: incorrect data`)
	}

	_, err = db.Exec(query, idInt)
	if err != nil {
		return err

	}
	return nil

}
