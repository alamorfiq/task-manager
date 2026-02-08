package postgres

import (
	"context"
	"taskmgr/internal/storage"
)

func (p *Postgres) CreateTask(ctx context.Context, t storage.Task) (int64, error) {
	query := `
        INSERT INTO tasks (title, description)
        VALUES ($1, $2)
        RETURNING id;
    `
	var id int64
	err := p.pool.QueryRow(ctx, query, t.Title, t.Description).Scan(&id)
	return id, err
}

func (p *Postgres) GetTask(ctx context.Context, id int64) (storage.Task, error) {
	query := `
        SELECT id, title, description, is_done, created_at
        FROM tasks
        WHERE id = $1;
    `
	var t storage.Task
	err := p.pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.Title, &t.Description, &t.IsDone, &t.CreatedAt,
	)
	return t, err
}

func (p *Postgres) ListTasks(ctx context.Context) ([]storage.Task, error) {
	query := `
        SELECT id, title, description, is_done, created_at
        FROM tasks
        ORDER BY created_at DESC;
    `

	rows, err := p.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]storage.Task, 0)
	for rows.Next() {
		var t storage.Task
		err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.IsDone, &t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, rows.Err()
}

func (p *Postgres) MarkDone(ctx context.Context, id int64) error {
	query := `
        UPDATE tasks
        SET is_done = true
        WHERE id = $1;
    `
	_, err := p.pool.Exec(ctx, query, id)
	return err
}

func (p *Postgres) DeleteTask(ctx context.Context, id int64) error {
	query := `
        DELETE FROM tasks WHERE id = $1;
    `
	_, err := p.pool.Exec(ctx, query, id)
	return err
}
