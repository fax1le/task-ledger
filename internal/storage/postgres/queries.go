package postgres

import (
	"database/sql"
	"context"
	"todo/internal/models"
	"todo/internal/utils/password"
	"todo/internal/utils/task"
)

func SelectTasks(DB *sql.DB, ctx context.Context, query_params string, args []any) ([]models.Task, error) {
	query := "SELECT title, completed, due_date, created_at, updated_at, priority, category FROM tasks" + query_params

	rows, err := DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task

		if err := rows.Scan(
			&task.Title,
			&task.Completed,
			&task.Due_date,
			&task.Created_at,
			&task.Updated_at,
			&task.Priority,
			&task.Category,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func InsertTask(DB *sql.DB, ctx context.Context, user_id int, task models.NewTask) error {
	task_utils.TrimSpace(&task)

	_, err := DB.ExecContext(ctx, "INSERT INTO tasks (user_id, title, due_date, priority, category) values ($1, $2, $3, $4, $5)",
		user_id,
		task.Title,
		task.Due_date,
		task.Priority,
		task.Category,
	)

	return err
}

func SelectTask(DB *sql.DB, ctx context.Context, user_id int, task_uuid string) (models.Task, error) {
	var task models.Task

	row := DB.QueryRowContext(ctx, "SELECT title, completed, due_date, created_at, updated_at, priority, category FROM tasks WHERE user_id = $1 AND id = $2", user_id, task_uuid)

	if err := row.Scan(
		&task.Title,
		&task.Completed,
		&task.Due_date,
		&task.Created_at,
		&task.Updated_at,
		&task.Priority,
		&task.Category,
	); err != nil {
		return task, err
	}

	return task, nil
}

func UpdateTask(DB *sql.DB, ctx context.Context, update_query string, args []any) error {
	var i int

	user_id := args[len(args)-2]
	task_uuid := args[len(args)-1]

	row := DB.QueryRowContext(ctx, "SELECT 1 FROM tasks WHERE user_id = $1 AND id = $2", user_id, task_uuid)

	if err := row.Scan(&i); err != nil {
		return err
	}

	_, err := DB.ExecContext(ctx, update_query, args...)

	return err
}

func RemoveTask(DB *sql.DB, ctx context.Context, user_id int, task_uuid string) (int64, error) {
	res, err := DB.ExecContext(ctx, "DELETE FROM tasks WHERE user_id = $1 AND id = $2", user_id, task_uuid)

	if err != nil {
		return 0, err
	}

	rows_affected, _ := res.RowsAffected()

	return rows_affected, err
}

func TaskExists(DB *sql.DB, ctx context.Context, user_id int, title string) bool {
	found := 0

	row := DB.QueryRowContext(ctx, "SELECT 1 FROM tasks WHERE user_id = $1 AND title = $2", user_id, title)

	if err := row.Scan(&found); err != nil {
		return false
	}

	return true
}

func UserExistsByEmail(DB *sql.DB, ctx context.Context, email string) bool {
	i := 0
	row := DB.QueryRowContext(ctx, "SELECT 1 FROM users WHERE email = $1", email)

	err := row.Scan(&i)

	return err == nil
}

func UserExistsByID(DB *sql.DB, ctx context.Context, id int) bool {
	i := 0
	row := DB.QueryRowContext(ctx, "SELECT 1 FROM users WHERE id = $1", id)

	err := row.Scan(&i)

	return err == nil
}

func CreateUser(DB *sql.DB, ctx context.Context, user models.User) error {
	hashed_password, err := password.Hash([]byte(user.Password))
	if err != nil {
		return err
	}

	_, err = DB.ExecContext(ctx, "INSERT INTO users (email, hashed_password) VALUES ($1, $2)", user.Email, hashed_password)
	if err != nil {
		return err
	}

	return nil
}

func GetPassword(DB *sql.DB, ctx context.Context, email string) (string, error) {
	var hashed_password string

	row := DB.QueryRowContext(ctx, "SELECT hashed_password FROM users WHERE email=$1", email)
	err := row.Scan(&hashed_password)

	return hashed_password, err
}

func GetUserID(DB *sql.DB, ctx context.Context, email string) (int, error) {
	var id int

	row := DB.QueryRowContext(ctx, "SELECT id FROM users WHERE email=$1", email)

	err := row.Scan(&id)

	return id, err
}

func SelectAllTasks(DB *sql.DB, ctx context.Context) ([]models.DBtask, error) {
	rows, err := DB.QueryContext(ctx, "SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.DBtask

	for rows.Next() {
		var task models.DBtask

		if err := rows.Scan(
			&task.ID,
			&task.User_ID,
			&task.Title,
			&task.Completed,
			&task.Due_date,
			&task.Created_at,
			&task.Updated_at,
			&task.Priority,
			&task.Category,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func SelectAllUsers(DB *sql.DB, ctx context.Context) ([]models.DBuser, error) {
	rows, err := DB.QueryContext(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.DBuser

	for rows.Next() {
		var user models.DBuser

		if err := rows.Scan(
			&user.UID,
			&user.Email,
			&user.Password,
			&user.Created_at,
			&user.Updated_at,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
