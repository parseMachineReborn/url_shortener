package service

import (
	"github.com/parseMachineReborn/url_shortener/internal/model"
	"github.com/parseMachineReborn/url_shortener/internal/repository"
)

type URLService struct {
	repository repository.Repository
}

func NewURLService(repository repository.Repository) *URLService {
	return &URLService{
		repository: repository,
	}
}

func (s *URLService) Shorten(url string) string {
	return s.repository.Shorten(url)
}

func (s *URLService) GetURL(shortURL string) (string, error) {
	res, err := s.repository.Get(shortURL)

	if err != nil {
		return "", err
	}

	res.RedirectCount++

	return res.Addr, err
}

func (s *URLService) GetAll() map[string]model.URL {
	return s.repository.GetAll()
}

func (s *URLService) Delete(shortURL string) error {
	return s.repository.Delete(shortURL)
}
