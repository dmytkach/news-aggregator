package storage

import (
	"errors"
	"sync"
)

type Storage interface {
	Set(key, value string)
	Get(key string) (string, error)
	GetAll() map[string]string
	Delete(key string) error
	SetCachedNews(key string, news interface{}) error
	GetCachedNews(key string) (interface{}, error)
}

type memoryStorage struct {
	data       map[string]string
	newsCache  map[string]interface{}
	cacheMutex sync.RWMutex
}

func NewMemoryStorage() *memoryStorage {
	return &memoryStorage{
		data:      make(map[string]string),
		newsCache: make(map[string]interface{}),
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

func (s *memoryStorage) GetAll() map[string]string {
	return s.data
}

func (s *memoryStorage) Delete(key string) error {
	delete(s.data, key)
	return nil
}

func (s *memoryStorage) SetCachedNews(key string, news interface{}) error {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()
	s.newsCache[key] = news
	return nil
}

func (s *memoryStorage) GetCachedNews(key string) (interface{}, error) {
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
	for source := range s.data {
		sources = append(sources, source)
	}
	return sources
}

//
//type fileStorage struct {
//	data       map[string]string
//	filePath   string
//	newsCache  map[string]interface{}
//	cacheMutex sync.RWMutex
//}
//
//func NewFileStorage(filePath string) (*fileStorage, error) {
//	s := &fileStorage{
//		data:      make(map[string]string),
//		filePath:  filePath,
//		newsCache: make(map[string]interface{}),
//	}
//	if err := s.load(); err != nil {
//		return nil, err
//	}
//	return s, nil
//}
//
//func (s *fileStorage) load() error {
//	file, err := os.Open(s.filePath)
//	if err != nil {
//		if os.IsNotExist(err) {
//			return nil // File does not exist, initialize with empty data
//		}
//		return err
//	}
//	defer file.Close()
//	return json.NewDecoder(file).Decode(&s.data)
//}
//
//func (s *fileStorage) save() error {
//	file, err := os.Create(s.filePath)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//	return json.NewEncoder(file).Encode(s.data)
//}
//
//func (s *fileStorage) Set(key, value string) error {
//	s.data[key] = value
//	return s.save()
//}
//
//func (s *fileStorage) Get(key string) (string, error) {
//	value, ok := s.data[key]
//	if !ok {
//		return "", errors.New("key not found")
//	}
//	return value, nil
//}
//

//
//func (s *fileStorage) Delete(key string) error {
//	delete(s.data, key)
//	return s.save()
//}
//
//func (s *fileStorage) SetCachedNews(key string, news interface{}) error {
//	s.cacheMutex.Lock()
//	defer s.cacheMutex.Unlock()
//	s.newsCache[key] = news
//	return s.saveCache()
//}
//
//func (s *fileStorage) GetCachedNews(key string) (interface{}, error) {
//	s.cacheMutex.RLock()
//	defer s.cacheMutex.RUnlock()
//	news, ok := s.newsCache[key]
//	if !ok {
//		return nil, errors.New("cached news not found")
//	}
//	return news, nil
//}
//
//func (s *fileStorage) saveCache() error {
//	cacheFile, err := os.Create(s.filePath + ".cache")
//	if err != nil {
//		return err
//	}
//	defer cacheFile.Close()
//	return json.NewEncoder(cacheFile).Encode(s.newsCache)
//}
//
//func (s *fileStorage) loadCache() error {
//	cacheFile, err := os.Open(s.filePath + ".cache")
//	if err != nil {
//		if os.IsNotExist(err) {
//			return nil // Cache file does not exist, initialize with empty data
//		}
//		return err
//	}
//	defer cacheFile.Close()
//	return json.NewDecoder(cacheFile).Decode(&s.newsCache)
//}
