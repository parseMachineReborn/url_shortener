package url

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/parseMachineReborn/url_shortener/internal/model"
)

const sliceEnd int = 7

type Repository interface {
	Add(ctx context.Context, shortenURL string, url *model.URL, userId int) error
	Get(ctx context.Context, shortURL string, userId int) (*model.URL, error)
	GetAll(ctx context.Context, userId int) (map[string]*model.URL, error)
	Delete(ctx context.Context, shortURL string, userId int) error
	IncrementRedirectCount(ctx context.Context, shortURL string) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Shorten(ctx context.Context, url string) (string, error) {
	urlBytes := []byte(url)
	hash := md5.Sum(urlBytes)
	result := fmt.Sprintf("%x", hash)

	urlModel := model.URL{
		Addr:          url,
		CreationDate:  time.Now(),
		RedirectCount: 0,
	}

	shortURL := result[:sliceEnd]
	userId := int(ctx.Value("user_id").(float64))
	err := s.repository.Add(ctx, shortURL, &urlModel, userId)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (s *Service) GetURL(ctx context.Context, shortURL string) (string, error) {
	userId := int(ctx.Value("user_id").(float64))
	res, err := s.repository.Get(ctx, shortURL, userId)

	if err != nil {
		return "", err
	}

	err = s.repository.IncrementRedirectCount(ctx, shortURL)
	if err != nil {
		return "", err
	}

	return res.Addr, err
}

func (s *Service) GetAll(ctx context.Context) (map[string]*model.URL, error) {
	userId := int(ctx.Value("user_id").(float64))
	return s.repository.GetAll(ctx, userId)
}

func (s *Service) Delete(ctx context.Context, shortURL string) error {
	userId := int(ctx.Value("user_id").(float64))
	return s.repository.Delete(ctx, shortURL, userId)
}
