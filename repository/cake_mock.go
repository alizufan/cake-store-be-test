package repository

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/zufzuf/cake-store/schema"
)

type CakeMock struct {
	mock.Mock
}

func (m *CakeMock) Find(ctx context.Context, id int) (*schema.Cake, error) {
	args := m.Called(ctx, id)
	res, _ := args.Get(0).(*schema.Cake)
	return res, args.Error(1)
}

func (m *CakeMock) FindAll(ctx context.Context, fil *FindAllFilter) ([]schema.Cake, error) {
	args := m.Called(ctx, fil)
	return args.Get(0).([]schema.Cake), args.Error(1)
}

func (m *CakeMock) Insert(ctx context.Context, rec *schema.Cake) error {
	rec.CreatedAt = time.Time{}
	rec.UpdatedAt = time.Time{}
	args := m.Called(ctx, rec)
	return args.Error(0)
}

func (m *CakeMock) Update(ctx context.Context, rec *schema.Cake) error {
	rec.CreatedAt = time.Time{}
	rec.UpdatedAt = time.Time{}
	args := m.Called(ctx, rec)
	return args.Error(0)
}

func (m *CakeMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
