package repository

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Repository struct {
	Conn *pgxpool.Pool
}

func (r *Repository) Create(ctx context.Context, t Order) error {
	query := `INSERT INTO orders (id, json_data) SELECT CAST($1 AS VARCHAR), $2 WHERE NOT EXISTS(SELECT id FROM orders WHERE id = $1)`
	_, err := r.Conn.Exec(ctx, query, t.Id, t)
	return err
}

func (r *Repository) FindAll(ctx context.Context) (t []Order, err error) {
	query := `SELECT json_data FROM orders`
	rows, err := r.Conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	buff := make([]Order, 0)

	for rows.Next() {
		var s string
		var d Order
		err := rows.Scan(&s)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(s), &d)
		if err != nil {
			log.Fatalf("%v", err)
		}
		buff = append(buff, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return buff, nil
}
