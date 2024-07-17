package actions

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Enotisi/go_final_project/internal/models"
)

var db *sql.DB

const (
	DateTemplate    = "20060102"
	dataBaseFields  = "id, date, title, comment, repeat"
	baseSearchlimit = 10
)

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
		err = repeatByYear(now, &startDate)
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
		query = fmt.Sprintf("SELECT %s FROM scheduler ORDER BY date LIMIT :limit", dataBaseFields)
	} else if date, err := time.Parse("02.01.2006", search); err == nil {
		dateSearch = date.Format(DateTemplate)
		query = fmt.Sprintf("SELECT %s FROM scheduler WHERE date = :date ORDER BY date LIMIT :limit", dataBaseFields)
	} else {
		textSearch = "%" + search + "%"
		query = fmt.Sprintf("SELECT %s FROM scheduler WHERE title LIKE :text OR comment LIKE :text ORDER BY date LIMIT :limit", dataBaseFields)
	}

	rows, err := db.Query(query,
		sql.Named("date", dateSearch),
		sql.Named("text", textSearch),
		sql.Named("limit", baseSearchlimit),
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

	sqlSearch := fmt.Sprintf("SELECT %s FROM scheduler WHERE id = :id", dataBaseFields)

	row := db.QueryRow(sqlSearch,
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

	if task.Id == "" {
		return errors.New("не указан идентификатор")
	}
	_, err := GetTaskById(task.Id)
	if err != nil {
		return err
	}
	if check {
		err = checkTaskData(&task)
		if err != nil {
			return err
		}
	}

	_, err = db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
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
		err := DeleteTaskById(task.Id)
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
	repeatData := nowData

	if task.Repeat != "" {
		newDate, err := NextDate(repeatData, repeatData.Format(DateTemplate), task.Repeat)
		if err != nil {
			return err
		}

		repeatData = newDate
	}

	if task.Date == "" {
		taskDate = nowData
	} else {
		dateParse, err := time.Parse(DateTemplate, task.Date)
		if err != nil {
			return errors.New("некорректная дата")
		}

		if dateParse.Format(DateTemplate) < nowData.Format(DateTemplate) {
			taskDate = repeatData
		} else {
			taskDate = dateParse
		}
	}

	task.Date = taskDate.Format(DateTemplate)

	return nil
}
