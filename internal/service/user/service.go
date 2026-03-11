package user

import (
	"context"

	"github.com/parseMachineReborn/url_shortener/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Register(ctx context.Context, email string, passHash string) error
	GetUser(ctx context.Context, email string) (model.User, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Register(ctx context.Context, email string, pass string) error {
	hashedPassByte, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	hashedPass := string(hashedPassByte)

	return s.repository.Register(ctx, email, hashedPass)
}

func (s *Service) LogIn(ctx context.Context, email string, pass string) (int, error) {
	user, err := s.repository.GetUser(ctx, email)
	if err != nil {
		return -1, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(pass)); err != nil {
		return -1, err
	}

	return int(user.ID), nil
}
