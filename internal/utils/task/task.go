package task_utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"todo/internal/models"
)

func GetTaskUUID(path string) string {
	extracted, _ := strings.CutPrefix(path, "/tasks/")

	return extracted
}

func ValidString(s string) bool {
	for _, c := range s {
		if !(('!' <= c && c <= '~') || (c == ' ')) {
			return false
		}

	}

	return true
}

var allowedOrderBy = map[string]string{
	"title":           "title",
	"title desc":      "title desc",
	"priority":        "priority",
	"priority desc":   "priority desc",
	"completed":       "completed",
	"completed desc":  "completed desc",
	"due_date":        "due_date",
	"due_date desc":  "due_date desc",
	"category":        "category",
	"category desc":   "category desc",
	"created_at":      "created_at",
	"created_at desc": "created_at desc",
	"updated_at":      "updated_at",
	"updated_at desc": "updated_at desc",
}

func GetDynamicQuery(user_id int, r *http.Request) (string, []any, error) {
	condition_query := " WHERE user_id = $1 AND"
	operation_query := ""
	query := ""
	args := []interface{}{}
	arg_ind := 2

	args = append(args, user_id)

	query_params := GetQueryParams(r)

	if query_params["completed"] != "" {
		param := query_params["completed"]
		param = strings.ToLower(param)
		param = strings.TrimSpace(param)

		if param != "false" && param != "true" {
			return "", args, errors.New("completed param not a bool value")
		}

		completed_str := fmt.Sprintf(" completed = $%d", arg_ind)
		condition_query += completed_str + " AND"

		completed_bool, err := strconv.ParseBool(param)

		if err != nil {
			return "", args, err
		}
		args = append(args, completed_bool)
		arg_ind++
	}

	if query_params["category"] != "" {
		param := query_params["category"]
		param = strings.TrimSpace(param)

		category_str := fmt.Sprintf(" category = $%d", arg_ind)
		condition_query += category_str + " AND"
		args = append(args, param)
		arg_ind++
	}

	if query_params["due"] != "" {
		param := query_params["due"]
		param = strings.TrimSpace(param)

		layout := "2006-01-02 15:04:05"

		date, err := time.Parse(layout, param)

		if err != nil {
			return "", args, err
		}

		due_str := fmt.Sprintf(" due_date <= $%d", arg_ind)
		condition_query += due_str + " AND"
		args = append(args, date)
		arg_ind++
	}

	if query_params["search"] != "" {
		param := query_params["search"]
		param = strings.TrimSpace(param)

		search_str := fmt.Sprintf(" title = $%d", arg_ind)
		condition_query += search_str + " AND"
		args = append(args, param)
		arg_ind++
	}

	if query_params["priority"] != "" {
		param := query_params["priority"]
		param = strings.ToLower(param)
		param = strings.TrimSpace(param)

		if param != "low" && param != "medium" && param != "high" {
			return "", args, errors.New("priority param not in ('low', 'medium', 'high')")
		}

		priority_str := fmt.Sprintf(" priority = $%d", arg_ind)
		condition_query += priority_str + " AND"
		args = append(args, param)
		arg_ind++
	}

	condition_query, _ = strings.CutSuffix(condition_query, " AND")

	if query_params["sort"] != "" {
		param := query_params["sort"]
		param = strings.ToLower(param)
		param = strings.TrimSpace(param)

		if _, ok := allowedOrderBy[param]; !ok {
			return "", args, errors.New("sort param not allowed")
		}

		sort_str := fmt.Sprintf(" ORDER BY %s", param)
		operation_query += sort_str
	}

	if query_params["limit"] != "" {
		param := query_params["limit"]
		param = strings.TrimSpace(param)

		limit, err := strconv.Atoi(param)

		if err != nil {
			return "", args, errors.New("limit param not a number")
		}

		if limit < 0 {
			return "", args, errors.New("limit must be positive")
		}

		limit_str := fmt.Sprintf(" LIMIT $%d", arg_ind)
		operation_query += limit_str
		args = append(args, limit)
		arg_ind++
	}

	if condition_query == " WHERE" {
		condition_query = ""
	}

	query += condition_query + operation_query

	return query, args, nil
}

func GetQueryParams(r *http.Request) map[string]string {
	return map[string]string{
		"completed": r.URL.Query().Get("completed"),
		"category":  r.URL.Query().Get("category"),
		"due":       r.URL.Query().Get("due"),
		"search":    r.URL.Query().Get("search"),
		"sort":      r.URL.Query().Get("sort"),
		"limit":     r.URL.Query().Get("limit"),
		"priority":  r.URL.Query().Get("priority"),
	}
}

func TrimSpace(task *models.NewTask) {
	task.Title = strings.TrimSpace(task.Title)
	task.Due_date = strings.TrimSpace(task.Due_date)
	task.Priority = strings.TrimSpace(task.Priority)
	task.Category = strings.TrimSpace(task.Category)
}

func GetUpdateQuery(user_id int, task_uuid string, update_task models.UpdateTask) (string, []any) {
	update_query := "UPDATE tasks SET "
	args := []any{}
	arg_ind := 1

	if update_task.Title != nil {
		update_query += fmt.Sprintf("title = $%d, ", arg_ind)

		args = append(args, *update_task.Title)
		arg_ind++
	}

	if update_task.Due_date != nil {
		update_query += fmt.Sprintf("due_date = $%d, ", arg_ind)

		args = append(args, *update_task.Due_date)
		arg_ind++
	}

	if update_task.Priority != nil {
		update_query += fmt.Sprintf("priority = $%d, ", arg_ind)

		args = append(args, *update_task.Priority)
		arg_ind++
	}

	if update_task.Category != nil {
		update_query += fmt.Sprintf("category = $%d, ", arg_ind)

		args = append(args, *update_task.Category)
		arg_ind++
	}

	if update_task.Completed != nil {
		update_query += fmt.Sprintf("completed = $%d, ", arg_ind)

		args = append(args, *update_task.Completed)
		arg_ind++
	}

	if update_query == "UPDATE tasks SET " {
		update_query = ""
	} else {
		update_query += fmt.Sprintf("updated_at = $%d WHERE user_id = $%d AND id = $%d", arg_ind, arg_ind+1, arg_ind+2)
		args = append(args, time.Now(), user_id, task_uuid)
		arg_ind += 2
	}

	return update_query, args
}
