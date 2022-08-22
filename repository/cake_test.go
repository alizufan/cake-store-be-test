package repository_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zufzuf/cake-store/repository"
	"github.com/zufzuf/cake-store/schema"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

var cake = schema.Cake{
	ID:          1,
	Title:       "Test Title",
	Description: "Test Description",
	Rating:      7,
	Image:       "https://cdn.lorem.space/images/movie/.cache/150x220/totoro-1988.jpeg",
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
}

var cakes = []schema.Cake{
	cake,
	{
		ID:          2,
		Title:       "Test Title 2",
		Description: "Test Description 2",
		Rating:      8,
		Image:       "https://cdn.lorem.space/images/movie/.cache/150x220/totoro-1988.jpeg",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
}

func Test_Cake_Repository_Find(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id",
		"title",
		"description",
		"rating",
		"image",
		"created_at",
		"updated_at",
	}).AddRow(
		cake.ID,
		cake.Title,
		cake.Description,
		cake.Rating,
		cake.Image,
		cake.CreatedAt,
		cake.UpdatedAt,
	)
	mock.ExpectQuery("SELECT * FROM cakes WHERE id = ? LIMIT 1").WithArgs(cake.ID).WillReturnRows(rows)

	repo := &repository.Cake{DB: db}
	res, err := repo.Find(context.Background(), cake.ID)
	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func Test_Cake_Repository_Find_All(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	tests := []struct {
		Name   string
		Query  string
		Filter repository.FindAllFilter
		Result []schema.Cake
	}{
		{
			Name:   "No_Filter",
			Query:  "SELECT * FROM cakes ORDER BY title ASC, rating ASC",
			Result: cakes,
		},
		{
			Name:  "Title_Filter",
			Query: "SELECT * FROM cakes WHERE title LIKE ? ORDER BY title ASC, rating ASC",
			Filter: repository.FindAllFilter{
				Title: "Test Title",
			},
			Result: cakes,
		},
		{
			Name:  "Description_Filter",
			Query: "SELECT * FROM cakes WHERE description LIKE ? ORDER BY title ASC, rating ASC",
			Filter: repository.FindAllFilter{
				Description: "Test Description",
			},
			Result: cakes,
		},
		{
			Name:  "All_Filter",
			Query: "SELECT * FROM cakes WHERE title LIKE ? AND description LIKE ? ORDER BY title ASC, rating ASC",
			Filter: repository.FindAllFilter{
				Title:       "Test Title",
				Description: "Test Description",
			},
			Result: cakes,
		},
		{
			Name:  "Record_Not_Found",
			Query: "SELECT * FROM cakes ORDER BY title ASC, rating ASC",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			rows := sqlmock.NewRows([]string{
				"id",
				"title",
				"description",
				"rating",
				"image",
				"created_at",
				"updated_at",
			})

			for _, c := range test.Result {
				rows = rows.AddRow(
					c.ID,
					c.Title,
					c.Description,
					c.Rating,
					c.Image,
					c.CreatedAt,
					c.UpdatedAt,
				)
			}

			args := []driver.Value{}
			if test.Filter.IsValidTitle() {
				args = append(args, "%"+test.Filter.Title+"%")
			}
			if test.Filter.IsValidDescription() {
				args = append(args, "%"+test.Filter.Description+"%")
			}

			mock.ExpectQuery(test.Query).WithArgs(args...).WillReturnRows(rows)

			repo := &repository.Cake{DB: db}
			res, err := repo.FindAll(context.Background(), &test.Filter)

			if test.Name == "Record_Not_Found" {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, repository.ErrRecordNotFound)
				return
			}

			assert.NotNil(t, res)
			assert.NoError(t, err)
		})
	}
}

func Test_Cake_Repository_Insert(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	expec := 1
	result := sqlmock.NewResult(int64(expec), 1)

	mock.ExpectExec("INSERT INTO cakes (title, description, rating, image, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)").
		WithArgs(
			cake.Title,
			cake.Description,
			cake.Rating,
			cake.Image,
			cake.CreatedAt,
			cake.UpdatedAt,
		).WillReturnResult(result)
	repo := &repository.Cake{DB: db}
	err := repo.Insert(context.Background(), &cake)

	assert.Equal(t, expec, cake.ID)
	assert.NoError(t, err)
}

func Test_Cake_Repository_Update(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	result := sqlmock.NewResult(0, 1)
	mock.ExpectExec("UPDATE cakes SET title=?, description=?, rating=?, image=?, updated_at=? WHERE id = ?").
		WithArgs(
			cake.Title,
			cake.Description,
			cake.Rating,
			cake.Image,
			cake.UpdatedAt,
			cake.ID,
		).WillReturnResult(result)

	repo := &repository.Cake{DB: db}
	err := repo.Update(context.Background(), &cake)
	assert.NoError(t, err)
}

func Test_Cake_Repository_Delete(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	result := sqlmock.NewResult(0, 1)
	mock.ExpectExec("DELETE FROM cakes WHERE id = ?").WithArgs(cake.ID).WillReturnResult(result)

	repo := &repository.Cake{DB: db}
	err := repo.Delete(context.Background(), cake.ID)
	assert.NoError(t, err)
}
