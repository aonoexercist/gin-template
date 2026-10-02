package user

import (
	"context"
	"errors"
	"testing"

	"github.com/yourname/gin-template/internal/platform/apperror"
)

// fakeRepo is an in-memory Repository for unit tests.
type fakeRepo struct {
	users  map[uint]*User
	nextID uint
}

func newFakeRepo() *fakeRepo { return &fakeRepo{users: map[uint]*User{}, nextID: 1} }

func (f *fakeRepo) Create(_ context.Context, u *User) error {
	for _, e := range f.users {
		if e.Email == u.Email {
			return ErrDuplicate
		}
	}
	u.ID = f.nextID
	f.nextID++
	f.users[u.ID] = u
	return nil
}
func (f *fakeRepo) GetByID(_ context.Context, id uint) (*User, error) {
	if u, ok := f.users[id]; ok {
		return u, nil
	}
	return nil, ErrNotFound
}
func (f *fakeRepo) List(context.Context, int, int) ([]User, int64, error) { return nil, 0, nil }
func (f *fakeRepo) Update(context.Context, *User) error                    { return nil }
func (f *fakeRepo) Delete(_ context.Context, id uint) error {
	if _, ok := f.users[id]; !ok {
		return ErrNotFound
	}
	delete(f.users, id)
	return nil
}

func TestCreate_DuplicateEmail(t *testing.T) {
	svc := NewService(newFakeRepo())
	req := CreateRequest{Name: "Ana", Email: "ana@example.com"}

	if _, err := svc.Create(context.Background(), req); err != nil {
		t.Fatalf("first create: %v", err)
	}
	_, err := svc.Create(context.Background(), req)

	var appErr *apperror.Error
	if !errors.As(err, &appErr) || appErr.Code != "conflict" {
		t.Fatalf("expected conflict error, got %v", err)
	}
}

func TestGet_NotFound(t *testing.T) {
	svc := NewService(newFakeRepo())
	_, err := svc.Get(context.Background(), 99)

	var appErr *apperror.Error
	if !errors.As(err, &appErr) || appErr.Code != "not_found" {
		t.Fatalf("expected not_found error, got %v", err)
	}
}
