package local

import (
	"sync"

	"github.com/parseMachineReborn/url_shortener/internal/apperror"
	"github.com/parseMachineReborn/url_shortener/internal/model"
)

type defaultRepository struct {
	mu      sync.Mutex
	storage map[string]*model.URL
}

func NewDefaultRepository() *defaultRepository {
	return &defaultRepository{
		storage: make(map[string]*model.URL),
	}
}

func (r *defaultRepository) Add(shortenURL string, url *model.URL) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[shortenURL] = url
}

func (r *defaultRepository) Get(shortURL string) (*model.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exist := r.storage[shortURL]; !exist {
		return &model.URL{}, apperror.ErrNotFound
	}

	return r.storage[shortURL], nil
}

func (r *defaultRepository) GetAll() map[string]*model.URL {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.storage
}

func (r *defaultRepository) Delete(shortURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exist := r.storage[shortURL]; !exist {
		return apperror.ErrNotFound
	}

	delete(r.storage, shortURL)

	return nil
}
