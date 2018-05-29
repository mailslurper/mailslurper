package cache

import (
	"time"
)

type ICacheService interface {
	Delete(key string)
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, timeout time.Duration)
}
