package service

import (
	"context"
	"errors"
	"testing"

	apperrors "github.com/shuyou-ai/shuyou-go/internal/errors"
	"github.com/shuyou-ai/shuyou-go/internal/model"
	"github.com/shuyou-ai/shuyou-go/internal/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	users map[string]*model.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{users: make(map[string]*model.User)}
}

func (m *mockUserRepository) Create(_ context.Context, user *model.User) error {
	if _, ok := m.users[user.Username]; ok {
		return repository.ErrDuplicateKey
	}
	user.PrepareCreate()
	m.users[user.Username] = user
	return nil
}

func (m *mockUserRepository) FindByID(_ context.Context, id string) (*model.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, mongo.ErrNoDocuments
}

func (m *mockUserRepository) FindByUsername(_ context.Context, username string) (*model.User, error) {
	user, ok := m.users[username]
	if !ok {
		return nil, mongo.ErrNoDocuments
	}
	return user, nil
}

func TestUserService_RegisterAndLogin(t *testing.T) {
	repo := newMockUserRepository()
	svc := NewUserService(repo)

	user, err := svc.Register(context.Background(), "alice", "alice@example.com", "secret123")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if user.Username != "alice" {
		t.Fatalf("unexpected username: %s", user.Username)
	}

	_, err = svc.Register(context.Background(), "alice", "alice2@example.com", "secret123")
	if !isAppErrorCode(err, apperrors.CodeConflict) {
		t.Fatalf("expected conflict error, got: %v", err)
	}

	loginUser, err := svc.Login(context.Background(), "alice", "secret123")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if loginUser.ID != user.ID {
		t.Fatalf("login user id mismatch")
	}

	_, err = svc.Login(context.Background(), "alice", "wrong-password")
	if !isAppErrorCode(err, apperrors.CodeUnauthorized) {
		t.Fatalf("expected unauthorized error, got: %v", err)
	}
}

func TestUserService_GetByID(t *testing.T) {
	repo := newMockUserRepository()
	svc := NewUserService(repo)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	user := &model.User{Username: "bob", Email: "bob@example.com", Password: string(hashed), Status: 1}
	_ = repo.Create(context.Background(), user)

	found, err := svc.GetByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("get by id failed: %v", err)
	}
	if found.Username != "bob" {
		t.Fatalf("unexpected username: %s", found.Username)
	}

	_, err = svc.GetByID(context.Background(), "missing-id")
	if !isAppErrorCode(err, apperrors.CodeNotFound) {
		t.Fatalf("expected not found error, got: %v", err)
	}
}

func isAppErrorCode(err error, code int) bool {
	var appErr *apperrors.Error
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}
