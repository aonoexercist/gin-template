package user

import (
	"context"
	"errors"

	"github.com/yourname/gin-template/internal/platform/apperror"
)

type Service struct{ repo Repository }

func NewService(repo Repository) *Service { return &Service{repo: repo} }

// Service methods return *apperror.Error so handlers stay thin.

func (s *Service) Create(ctx context.Context, req CreateRequest) (*User, error) {
	u := &User{Name: req.Name, Email: req.Email}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, mapErr(err)
	}
	return u, nil
}

func (s *Service) Get(ctx context.Context, id uint) (*User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, mapErr(err)
	}
	return u, nil
}

func (s *Service) List(ctx context.Context, q ListQuery) ([]User, int64, error) {
	users, total, err := s.repo.List(ctx, q.PageSize, (q.Page-1)*q.PageSize)
	if err != nil {
		return nil, 0, mapErr(err)
	}
	return users, total, nil
}

func (s *Service) Update(ctx context.Context, id uint, req UpdateRequest) (*User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, mapErr(err)
	}
	u.Name, u.Email = req.Name, req.Email
	if err := s.repo.Update(ctx, u); err != nil {
		return nil, mapErr(err)
	}
	return u, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return mapErr(s.repo.Delete(ctx, id))
}

func mapErr(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, ErrNotFound):
		return apperror.NotFound("user not found")
	case errors.Is(err, ErrDuplicate):
		return apperror.Conflict("email already in use")
	default:
		return apperror.Internal(err)
	}
}
