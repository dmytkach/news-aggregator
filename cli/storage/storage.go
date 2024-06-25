package storage

import (
	"errors"
	"news-aggregator/internal/entity"
	"sync"
)

type Storage interface {
	Set(key, value string)
	Get(key string) (string, error)
	GetAll() map[string][]entity.News
	Delete(key string) error
	AddNewsToCache(key string, news []entity.News)
	GetCachedNews(key string) (news []entity.News, err error)
	SaveMapToCache(newsMap map[string][]entity.News)
	AvailableSources() []string
}

type memoryStorage struct {
	data       map[string]string
	newsCache  map[string][]entity.News
	cacheMutex sync.RWMutex
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		data:      make(map[string]string),
		newsCache: make(map[string][]entity.News),
	}
}

func (s *memoryStorage) Set(key, value string) {
	s.data[key] = value
}

func (s *memoryStorage) Get(key string) (string, error) {
	value, ok := s.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return value, nil
}

func (s *memoryStorage) GetAll() map[string][]entity.News {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	c := make(map[string][]entity.News)
	for key, news := range s.newsCache {
		c[key] = append([]entity.News{}, news...)
	}
	return c
}

func (s *memoryStorage) Delete(key string) error {
	delete(s.data, key)
	return nil
}

func (s *memoryStorage) AddNewsToCache(key string, news []entity.News) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	for _, newArticle := range news {
		if !s.newsExistsInCache(key, newArticle) {
			s.newsCache[key] = append(s.newsCache[key], newArticle)
		}
	}
}
func (s *memoryStorage) SaveMapToCache(newsMap map[string][]entity.News) {
	for source, news := range newsMap {
		s.AddNewsToCache(source, news)
	}
}

func (s *memoryStorage) newsExistsInCache(key string, newsItem entity.News) bool {
	for _, item := range s.newsCache[key] {
		if item.Link == newsItem.Link {
			return true
		}
	}
	return false
}

func (s *memoryStorage) GetCachedNews(key string) ([]entity.News, error) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()
	news, ok := s.newsCache[key]
	if !ok {
		return nil, errors.New("cached news not found")
	}
	return news, nil
}

// AvailableSources returns all the available registered is storage sources.
func (s *memoryStorage) AvailableSources() []string {
	sources := make([]string, 0, len(s.data))
	for source := range s.newsCache {
		sources = append(sources, source)
	}
	return sources
}
