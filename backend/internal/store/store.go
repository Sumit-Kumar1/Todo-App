package store

import (
	"database/sql"
	pkgErr "errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"todoapp/internal/errors"
	"todoapp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type queries struct {
	deleteTask        string
	deleteCompleted   string
	getAllByUserID     string
	getTaskByID       string
	getChildTasks     string
	insertQuery       string
	setDone           string
	updateQuery       string
}

var sqlQueries = queries{
	deleteTask:      "DELETE FROM tasks WHERE id=$1 AND user_id=$2;",
	deleteCompleted: "DELETE FROM tasks WHERE user_id=$1 AND status='DONE';",
	getAllByUserID:   "SELECT id, parent_id, user_id, title, description, status, priority, category, due_date, created_at, updated_at FROM tasks WHERE user_id=$1;",
	getTaskByID:     "SELECT id, parent_id, user_id, title, description, status, priority, category, due_date, created_at, updated_at FROM tasks WHERE id=$1 AND user_id=$2;",
	getChildTasks:   "SELECT id, parent_id, user_id, title, description, status, priority, category, due_date, created_at, updated_at FROM tasks WHERE parent_id=$1 AND user_id=$2;",
	insertQuery:     "INSERT INTO tasks (id, parent_id, user_id, title, description, status, priority, category, due_date, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);",
	setDone:         "UPDATE tasks SET status=$1, updated_at=$2 WHERE id=$3 AND user_id=$4;",
	updateQuery:     "UPDATE tasks SET title=$1, description=$2, status=$3, priority=$4, category=$5, due_date=$6, updated_at=$7, parent_id=$10 WHERE id=$8 AND user_id=$9;",
}

type Store struct {
	DB *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{DB: db}
}

func scanTask(scanner interface{ Scan(...any) error }, task *models.Task) (*sql.NullString, error) {
	var parentID sql.NullString

	err := scanner.Scan(&task.ID, &parentID, &task.UserID, &task.Title, &task.Description, &task.Status,
		&task.Priority, &task.Category, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if parentID.Valid {
		task.ParentID = &parentID.String
	}

	return &parentID, nil
}

func (s *Store) GetAll(ctx *gin.Context, userID *uuid.UUID) ([]models.Task, error) {
	var res = make([]models.Task, 0)

	rows, err := s.DB.QueryContext(ctx, sqlQueries.getAllByUserID, *userID)
	if err != nil {
		if pkgErr.Is(err, sql.ErrNoRows) {
			return res, nil
		}

		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			slog.ErrorContext(ctx, "error while closing sql rows", slog.String("error", err.Error()))
		}
	}(rows)

	for rows.Next() {
		var task models.Task
		if _, err := scanTask(rows, &task); err != nil {
			return nil, err
		}

		res = append(res, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	slog.DebugContext(ctx, "get all tasks", slog.String("user", userID.String()))

	return res, nil
}

func (s *Store) Create(ctx *gin.Context, task *models.Task) error {
	logger := models.GetLoggerFromCtx(ctx)

	var parentID sql.NullString
	if task.ParentID != nil {
		parentID = sql.NullString{String: *task.ParentID, Valid: true}
	}

	res, err := s.DB.ExecContext(ctx, sqlQueries.insertQuery, task.ID, parentID, task.UserID,
		task.Title, task.Description, task.Status, task.Priority, task.Category, task.DueDate, task.CreatedAt)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.ErrInsertFailed
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task added successfully",
		slog.String("task", task.ID),
	)

	return nil
}

func (s *Store) Update(ctx *gin.Context, task *models.Task) error {
	logger := models.GetLoggerFromCtx(ctx)

	var parentID sql.NullString
	if task.ParentID != nil {
		parentID = sql.NullString{String: *task.ParentID, Valid: true}
	}

	res, err := s.DB.ExecContext(ctx, sqlQueries.updateQuery, task.Title, task.Description,
		task.Status, task.Priority, task.Category, task.DueDate, task.UpdatedAt, task.ID, task.UserID, parentID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.ErrTaskNotFound
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task updated successfully",
		slog.String("task", task.ID))

	return nil
}

func (s *Store) Delete(ctx *gin.Context, id string, userID *uuid.UUID) error {
	logger := models.GetLoggerFromCtx(ctx)

	res, err := s.DB.ExecContext(ctx, sqlQueries.deleteTask, id, *userID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.ErrTaskNotFound
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task deleted successfully", slog.String("task", id))

	return nil
}

func (s *Store) DeleteCompleted(ctx *gin.Context, userID *uuid.UUID) error {
	logger := models.GetLoggerFromCtx(ctx)

	_, err := s.DB.ExecContext(ctx, sqlQueries.deleteCompleted, *userID)
	if err != nil {
		return err
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "completed tasks deleted", slog.String("user", userID.String()))

	return nil
}

func (s *Store) MarkDone(ctx *gin.Context, id string, userID *uuid.UUID) error {
	logger := models.GetLoggerFromCtx(ctx)

	res, err := s.DB.ExecContext(ctx, sqlQueries.setDone, "DONE", time.Now(), id, *userID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.ErrTaskNotFound
	}

	logger.LogAttrs(ctx, slog.LevelDebug, "task marked as done", slog.String("task", id))
	return nil
}

func (s *Store) GetTaskByID(ctx *gin.Context, taskID string, userID *uuid.UUID) (*models.Task, error) {
	var task models.Task

	row := s.DB.QueryRowContext(ctx, sqlQueries.getTaskByID, taskID, *userID)
	if _, err := scanTask(row, &task); err != nil {
		if pkgErr.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrTaskNotFound
		}

		return nil, err
	}

	// Fetch child tasks
	childTasks, err := s.GetChildTasks(ctx, taskID, userID)
	if err != nil {
		return nil, err
	}
	task.ChildTasks = childTasks

	return &task, nil
}

func (s *Store) GetChildTasks(ctx *gin.Context, taskID string, userID *uuid.UUID) ([]models.Task, error) {
	var res = make([]models.Task, 0)

	rows, err := s.DB.QueryContext(ctx, sqlQueries.getChildTasks, taskID, *userID)
	if err != nil {
		if pkgErr.Is(err, sql.ErrNoRows) {
			return res, nil
		}

		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			slog.ErrorContext(ctx, "error while closing sql rows", slog.String("error", err.Error()))
		}
	}(rows)

	for rows.Next() {
		var task models.Task
		if _, err := scanTask(rows, &task); err != nil {
			return nil, err
		}

		res = append(res, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Store) GetFilteredTasks(ctx *gin.Context, userID *uuid.UUID, status, priority, category, search string) ([]models.Task, error) {
	var (
		res    = make([]models.Task, 0)
		args   = []any{*userID}
		wheres = []string{"user_id=$1"}
		argIdx = 2
	)

	if status != "" {
		wheres = append(wheres, fmt.Sprintf("status=$%d", argIdx))
		args = append(args, status)
		argIdx++
	}

	if priority != "" {
		wheres = append(wheres, fmt.Sprintf("priority=$%d", argIdx))
		args = append(args, priority)
		argIdx++
	}

	if category != "" {
		wheres = append(wheres, fmt.Sprintf("category=$%d", argIdx))
		args = append(args, category)
		argIdx++
	}

	if search != "" {
		wheres = append(wheres, fmt.Sprintf("title ILIKE $%d", argIdx))
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query := "SELECT id, parent_id, user_id, title, description, status, priority, category, due_date, created_at, updated_at FROM tasks WHERE " +
		strings.Join(wheres, " AND ") + " ORDER BY created_at DESC;"

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		if pkgErr.Is(err, sql.ErrNoRows) {
			return res, nil
		}

		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			slog.ErrorContext(ctx, "error while closing sql rows", slog.String("error", err.Error()))
		}
	}(rows)

	for rows.Next() {
		var task models.Task
		if _, err := scanTask(rows, &task); err != nil {
			return nil, err
		}

		res = append(res, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
