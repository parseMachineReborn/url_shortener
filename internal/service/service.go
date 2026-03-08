package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/parseMachineReborn/url_shortener/internal/model"
)

const sliceEnd int = 7

type Repository interface {
	Add(ctx context.Context, shortenURL string, url *model.URL) error
	Get(ctx context.Context, shortURL string) (*model.URL, error)
	GetAll(ctx context.Context) (map[string]*model.URL, error)
	Delete(ctx context.Context, shortURL string) error
	IncrementRedirectCount(ctx context.Context, shortURL string) error
}

type URLService struct {
	repository Repository
}

func NewURLService(repository Repository) *URLService {
	return &URLService{
		repository: repository,
	}
}

func (s *URLService) Shorten(ctx context.Context, url string) string {
	urlBytes := []byte(url)
	hash := md5.Sum(urlBytes)
	result := fmt.Sprintf("%x", hash)

	urlModel := model.URL{
		Addr:          url,
		CreationDate:  time.Now(),
		RedirectCount: 0,
	}

	shortURL := result[:sliceEnd]
	s.repository.Add(ctx, shortURL, &urlModel)

	return shortURL
}

func (s *URLService) GetURL(ctx context.Context, shortURL string) (string, error) {
	res, err := s.repository.Get(ctx, shortURL)

	if err != nil {
		return "", err
	}

	err = s.repository.IncrementRedirectCount(ctx, shortURL)
	if err != nil {
		return "", err
	}

	return res.Addr, err
}

func (s *URLService) GetAll(ctx context.Context) (map[string]*model.URL, error) {
	return s.repository.GetAll(ctx)
}

func (s *URLService) Delete(ctx context.Context, shortURL string) error {
	return s.repository.Delete(ctx, shortURL)
}
