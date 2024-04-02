package db

import (
	"context"
	"database/sql"
	"errors"

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
	return nil, nil
}
func (s *PostgresTaskStore) Delete(ctx context.Context, id types.ID) error {
	return nil
}

func (s *PostgresTaskStore) Update(ctx context.Context, id types.ID, update Map) error {
	return nil
}
