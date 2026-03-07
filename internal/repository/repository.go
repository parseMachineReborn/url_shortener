package repository

import (
	"crypto/md5"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/parseMachineReborn/url_shortener/internal/model"
)

const sliceEnd int = 7

var ErrNotFound = errors.New("Не найден")

type Repository interface {
	Shorten(URL string) string
	Get(shortURL string) (model.URL, error)
	GetAll() map[string]model.URL
	Delete(shortURL string) error
}

type defaultRepository struct {
	mu      sync.Mutex
	storage map[string]model.URL
}

func NewDefaultRepository() *defaultRepository {
	return &defaultRepository{
		storage: make(map[string]model.URL),
	}
}

func (r *defaultRepository) Shorten(url string) string {
	r.mu.Lock()
	defer r.mu.Unlock()

	urlBytes := []byte(url)
	hash := md5.Sum(urlBytes)
	result := fmt.Sprintf("%x", hash)

	urlModel := model.URL{
		Addr:          url,
		CreationDate:  time.Now(),
		RedirectCount: 0,
	}

	shortURL := result[:sliceEnd]
	r.storage[shortURL] = urlModel

	return shortURL
}

func (r *defaultRepository) Get(shortURL string) (model.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exist := r.storage[shortURL]; !exist {
		return model.URL{}, ErrNotFound
	}

	return r.storage[shortURL], nil
}

func (r *defaultRepository) GetAll() map[string]model.URL {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.storage
}

func (r *defaultRepository) Delete(shortURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exist := r.storage[shortURL]; !exist {
		return ErrNotFound
	}

	delete(r.storage, shortURL)

	return nil
}
