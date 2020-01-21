package main

import (
	"context"
	"fmt"
	"golang.org/x/text/language"
	"time"
)

type cacheItem struct {
	timestamp time.Time
	data      string
}

// Service is a Translator user.
type Service struct {
	translator Translator
}

var cache = make(map[string]cacheItem)

const retries = 3
const cacheDurationSeconds = 15

func init() {
	//remove outdated data in cache every 5 sec
	var cacheChecker = time.Tick(time.Second * 5)
	go func() {
		for timer := range cacheChecker {
			fmt.Println("clear cache ", timer)
			for k, v := range cache {
				if time.Since(v.timestamp) < time.Second*cacheDurationSeconds {
					delete(cache, k)
				}
			}
		}
	}()
}

func (s *Service) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default: // default case prevents blocking
	}
	key := fmt.Sprintf("from %s, to %s  data %s", from, to, data)
	var err error = nil
	for i := 1; i <= retries; i++ {
		if cachedTranslation, ok := cache[key]; ok {
			return cachedTranslation.data, nil
		}
		var translation string
		translation, err = s.translator.Translate(ctx, from, to, data)
		if err == nil {
			cache[key] = cacheItem{data: translation, timestamp: time.Now()}
			return translation, err
		}
		if i == retries {
			break
		}
		//retry after delay if request is unsuccessful
		for j := 0; j < i*i; j++ {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
				// check every second if another same request put result in cache
			case <-time.After(time.Second):
				if cachedTranslation, ok := cache[key]; ok {
					return cachedTranslation.data, nil
				}
			}
		}
	}
	return "", err
}

func NewService(t Translator) *Service {
	return &Service{
		translator: t,
	}
}
