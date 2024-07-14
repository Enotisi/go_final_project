package actions

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/Enotisi/go_final_project/internal/models"
)

var db *sql.DB

const DateTemplate = "20060102"

func InitAction(dataBase *sql.DB) {
	db = dataBase
}

func NextDate(now time.Time, date string, repeat string) (time.Time, error) {

	startDate, err := time.Parse(DateTemplate, date)

	if err != nil {
		return now, err
	}

	if repeat == "" {
		return now, errors.New("правило не может быть пустым")
	}

	switch []rune(repeat)[0] {
	case 'd':
		err = repeatByDays(now, &startDate, repeat)
	case 'y':
		startDate = startDate.AddDate(1, 0, 0)
	case 'w':
		err = repeatByWeek(now, &startDate, repeat)
	case 'm':
		err = repeatByMonthDay(now, &startDate, repeat)
	default:
		return now, errors.New("недопустимый символ")
	}

	if err != nil {
		return now, err
	}

	return startDate, nil
}

func CreateTask(taskData models.Task) (int, error) {

	err := checkTaskData(&taskData)

	if err != nil {
		return 0, err
	}

	res, err := db.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`,
		sql.Named("date", taskData.Date),
		sql.Named("title", taskData.Title),
		sql.Named("comment", taskData.Comment),
		sql.Named("repeat", taskData.Repeat),
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func TasksList(search string) ([]models.Task, error) {

	var query string
	var dateSearch string
	var textSearch string

	if search == "" {
		query = "SELECT * FROM scheduler ORDER BY date"
	} else if date, err := time.Parse("02.01.2006", search); err == nil {
		dateSearch = date.Format(DateTemplate)
		query = "SELECT * FROM scheduler WHERE date = :date ORDER BY date"
	} else {
		textSearch = "%" + search + "%"
		query = "SELECT * FROM scheduler WHERE title LIKE :text OR comment LIKE :text ORDER BY date"
	}

	rows, err := db.Query(query,
		sql.Named("date", dateSearch),
		sql.Named("text", textSearch),
	)

	if err != nil {
		return []models.Task{}, err
	}
	defer rows.Close()

	tasks := make([]models.Task, 0, 10)

	for rows.Next() {
		task := models.Task{}
		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)

		if err != nil {
			return []models.Task{}, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func GetTaskById(id string) (models.Task, error) {

	task := models.Task{}

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id",
		sql.Named("id", id),
	)

	if err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {

		if err == sql.ErrNoRows {
			return task, errors.New("задача не найдена")
		} else {
			return task, err
		}
	}

	return task, nil
}

func UpdateTask(task models.Task, check bool) error {

	if check {
		err := checkTaskData(&task)
		if err != nil {
			return err
		}
	}

	_, err := db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.Id),
	)

	return err
}

func DoneTask(id string) error {

	task, err := GetTaskById(id)

	if err != nil {
		return err
	}

	if task.Repeat != "" {
		newDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return err
		}

		task.Date = newDate.Format(DateTemplate)
		err = UpdateTask(task, false)
		return err
	} else {
		err := DeleteTaskById(strconv.Itoa(task.Id))
		return err
	}
}

func DeleteTaskById(id string) error {

	task, err := GetTaskById(id)

	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", task.Id),
	)

	return err
}

func checkTaskData(task *models.Task) error {

	if task.Title == "" {
		return errors.New("не указан заголовок")
	}

	var taskDate time.Time
	nowData := time.Now()

	if task.Repeat != "" {
		newDate, err := NextDate(nowData, nowData.Format(DateTemplate), task.Repeat)
		if err != nil {
			return err
		}

		nowData = newDate
	}

	if task.Date == "" {
		taskDate = nowData
	} else {
		dateParse, err := time.Parse(DateTemplate, task.Date)
		if err != nil {
			return err
		}

		if dateParse.Before(time.Now()) {
			taskDate = nowData
		} else {
			taskDate = dateParse
		}
	}

	task.Date = taskDate.Format(DateTemplate)

	return nil
}
