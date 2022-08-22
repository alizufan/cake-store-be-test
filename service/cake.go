package service

import (
	"context"
	"time"

	"github.com/rotisserie/eris"
	"github.com/zufzuf/cake-store/repository"
	"github.com/zufzuf/cake-store/schema"
)

var (
	ErrRequestNil = eris.New("request is nil")
)

type CakeRepository interface {
	Find(ctx context.Context, id int) (*schema.Cake, error)
	FindAll(ctx context.Context, fil *repository.FindAllFilter) ([]schema.Cake, error)
	Insert(ctx context.Context, rec *schema.Cake) error
	Update(ctx context.Context, rec *schema.Cake) error
	Delete(ctx context.Context, id int) error
}

type Cake struct {
	Repo CakeRepository
}

func (s *Cake) Find(ctx context.Context, id int) (*schema.Cake, error) {
	return s.Repo.Find(ctx, id)
}

type FindAllRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (s *Cake) FindAll(ctx context.Context, req *FindAllRequest) ([]schema.Cake, error) {
	if req == nil {
		return nil, ErrRequestNil
	}

	fil := repository.FindAllFilter{
		Title:       req.Title,
		Description: req.Description,
	}

	return s.Repo.FindAll(ctx, &fil)
}

type CakeRequest struct {
	ID          int     `json:"-"`
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Rating      float64 `json:"rating" validate:"required"`
	Image       string  `json:"image" validate:"required"`
}

func (s *Cake) Insert(ctx context.Context, req *CakeRequest) error {
	if req == nil {
		return ErrRequestNil
	}

	timeNow := time.Now()
	rec := schema.Cake{
		Title:       req.Title,
		Description: req.Description,
		Rating:      req.Rating,
		Image:       req.Image,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}
	if err := s.Repo.Insert(ctx, &rec); err != nil {
		return err
	}
	req.ID = rec.ID

	return nil
}

func (s *Cake) Update(ctx context.Context, req *CakeRequest) error {
	if req == nil {
		return ErrRequestNil
	}

	if _, err := s.Repo.Find(ctx, req.ID); err != nil {
		return err
	}

	rec := schema.Cake{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		Rating:      req.Rating,
		Image:       req.Image,
		UpdatedAt:   time.Now(),
	}
	return s.Repo.Update(ctx, &rec)
}

func (s *Cake) Delete(ctx context.Context, id int) error {
	if _, err := s.Repo.Find(ctx, id); err != nil {
		return err
	}
	return s.Repo.Delete(ctx, id)
}
