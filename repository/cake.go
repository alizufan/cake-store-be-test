package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/rotisserie/eris"
	"github.com/zufzuf/cake-store/libs/util"
	"github.com/zufzuf/cake-store/schema"
)

var (
	ErrSQLRowsNill = eris.New("sql rows is nil")
	ErrResultNill  = eris.New("result is nil")

	ErrFilterNill     = eris.New("filter is nil")
	ErrRecordNill     = eris.New("record is nil")
	ErrRecordNotFound = eris.New("record not found")
)

type Cake struct {
	DB *sql.DB
}

func (s *Cake) Find(ctx context.Context, id int) (*schema.Cake, error) {
	if id <= 0 {
		return nil, ErrRecordNotFound
	}

	query := "SELECT * FROM cakes WHERE id = ? LIMIT 1"
	log.Print(query)

	rows, err := s.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, eris.Wrap(err, "find cake by id, an error occurred")
	}

	res := schema.Cake{}
	if err := s.retrieveRow(rows, &res); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, eris.Wrap(err, "find cake by id, an error occurred")
	}

	if res.ID <= 0 {
		return nil, ErrRecordNotFound
	}

	return &res, nil
}

type FindAllFilter struct {
	Title       string
	Description string
}

func (f *FindAllFilter) IsValidTitle() bool {
	return len(f.Title) > 0
}

func (f *FindAllFilter) IsValidDescription() bool {
	return len(f.Description) > 0
}

func (s *Cake) FindAll(ctx context.Context, fil *FindAllFilter) ([]schema.Cake, error) {
	if fil == nil {
		return nil, ErrFilterNill
	}

	q := util.NewQuery()
	if fil.IsValidTitle() {
		q.Where("title LIKE ?", "%"+fil.Title+"%")
	}

	if fil.IsValidDescription() {
		q.Where("description LIKE ?", "%"+fil.Description+"%")
	}

	where, args := q.Build()
	query := "SELECT * FROM cakes " + where + " ORDER BY title ASC, rating ASC"
	log.Print(query)

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, eris.Wrap(err, "find cakes, an error occurred")
	}

	res := []schema.Cake{}
	if err := s.retrieveRows(rows, &res); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, eris.Wrap(err, "find cakes, an error occurred")
	}

	if len(res) == 0 {
		return nil, ErrRecordNotFound
	}

	return res, nil
}

func (s *Cake) Insert(ctx context.Context, rec *schema.Cake) error {
	if rec == nil {
		return ErrRecordNill
	}

	query := "INSERT INTO cakes (title, description, rating, image, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	log.Print(query)

	res, err := s.DB.ExecContext(ctx, query, rec.Title, rec.Description, rec.Rating, rec.Image, rec.CreatedAt, rec.UpdatedAt)
	if err != nil {
		return eris.Wrap(err, "insert cake, an error occurred")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return eris.Wrap(err, "insert cake, an error occurred")
	}
	rec.ID = int(id)

	return nil
}

func (s *Cake) Update(ctx context.Context, rec *schema.Cake) error {
	if rec == nil {
		return ErrRecordNill
	}
	if rec.ID <= 0 {
		return ErrRecordNotFound
	}

	query := "UPDATE cakes SET title=?, description=?, rating=?, image=?, updated_at=? WHERE id = ?"
	log.Print(query)

	if _, err := s.DB.ExecContext(ctx, query,
		rec.Title,
		rec.Description,
		rec.Rating,
		rec.Image,
		rec.UpdatedAt,
		rec.ID,
	); err != nil {
		return err
	}

	return nil
}

func (s *Cake) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrRecordNotFound
	}

	query := "DELETE FROM cakes WHERE id = ?"
	log.Print(query)

	if _, err := s.DB.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}

func (s *Cake) retrieveRow(rows *sql.Rows, res *schema.Cake) error {
	if rows == nil {
		return ErrSQLRowsNill
	}
	if res == nil {
		return ErrResultNill
	}

	for rows.Next() {
		if err := rows.Scan(
			&res.ID,
			&res.Title,
			&res.Description,
			&res.Rating,
			&res.Image,
			&res.CreatedAt,
			&res.UpdatedAt,
		); err != nil {
			return err
		}
	}

	return rows.Err()
}

func (s *Cake) retrieveRows(rows *sql.Rows, res *[]schema.Cake) error {
	if rows == nil {
		return ErrSQLRowsNill
	}
	if res == nil {
		return ErrResultNill
	}

	for rows.Next() {
		var o schema.Cake
		if err := rows.Scan(
			&o.ID,
			&o.Title,
			&o.Description,
			&o.Rating,
			&o.Image,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			util.ResetSlice(res)
			return err
		}

		*res = append(*res, o)
	}

	if err := rows.Err(); err != nil {
		util.ResetSlice(res)
		return err
	}

	return nil
}
