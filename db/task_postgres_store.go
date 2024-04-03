package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ficontini/gotasks/types"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostgresTaskStore struct {
	db *sql.DB
}

func NewPostgresTaskStore() (*PostgresTaskStore, error) {
	db, err := sql.Open("postgres", DBURI)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresTaskStore{
		db: db,
	}, nil
}

func (s *PostgresTaskStore) InsertTask(ctx context.Context, task *types.Task) (*types.Task, error) {
	query := `INSERT INTO tasks 
	(title, description, due_date, completed, projects)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id`
	err := s.db.QueryRowContext(ctx, query,
		task.Title,
		task.Description,
		task.DueDate,
		task.Completed,
		pq.Array(task.Projects)).Scan(&task.ID)

	if err != nil {
		return nil, err
	}
	return task, nil
}
func (s *PostgresTaskStore) GetTaskByID(ctx context.Context, id types.ID) (*types.Task, error) {
	intId, err := id.Int()
	if err != nil {
		return nil, err
	}
	query := `
	SELECT * 
	FROM tasks
	WHERE id = $1`
	var task types.Task
	err = s.db.QueryRowContext(ctx, query, intId).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.Completed,
		pq.Array(&task.Projects),
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorNotFound
		}
		return nil, err
	}
	return &task, nil
}
func (s *PostgresTaskStore) GetTasks(ctx context.Context, filter Map, pagination *Pagination) ([]*types.Task, error) {
	var tasks []*types.Task
	query := fmt.Sprintf("SELECT * FROM tasks ORDER BY id %s", pagination.getQuery())
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		task, err := scanIntoTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
func (s *PostgresTaskStore) Delete(ctx context.Context, id types.ID) error {
	_, err := s.db.Exec("delete from tasks where id=$1", id)
	return err
}

// task/:id/complete
func (s *PostgresTaskStore) Update(ctx context.Context, id types.ID, update Map) error {

	return nil
}
func scanIntoTask(rows *sql.Rows) (*types.Task, error) {
	var task types.Task
	err := rows.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.Completed,
		pq.Array(&task.Projects))
	return &task, err

}
