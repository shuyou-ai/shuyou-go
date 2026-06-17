package service

import (
	"context"
	"errors"

	apperrors "github.com/shuyou-ai/shuyou-go/internal/errors"
	"github.com/shuyou-ai/shuyou-go/internal/model"
	"github.com/shuyou-ai/shuyou-go/internal/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, username, email, password string) (*model.User, error)
	Login(ctx context.Context, username, password string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, username, email, password string) (*model.User, error) {
	if _, err := s.repo.FindByUsername(ctx, username); err == nil {
		return nil, apperrors.New(apperrors.CodeConflict, "username already exists")
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, apperrors.Wrap(err, apperrors.CodeInternalError, "query user failed")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeInternalError, "hash password failed")
	}

	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hashed),
		Status:   1,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return nil, apperrors.New(apperrors.CodeConflict, "user already exists")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeInternalError, "create user failed")
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, username, password string) (*model.User, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperrors.New(apperrors.CodeUnauthorized, "invalid username or password")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeInternalError, "query user failed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, apperrors.New(apperrors.CodeUnauthorized, "invalid username or password")
	}

	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperrors.New(apperrors.CodeNotFound, "user not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeInternalError, "query user failed")
	}
	return user, nil
}
