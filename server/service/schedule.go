package service

import (
	"sync"
	"time"
)

var (
	fetchInterval = time.Hour
	mu            sync.Mutex
)

func GetFetchInterval() time.Duration {
	mu.Lock()
	defer mu.Unlock()
	return fetchInterval
}

func SetFetchInterval(interval time.Duration) {
	mu.Lock()
	defer mu.Unlock()
	fetchInterval = interval
}
