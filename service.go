package main

import (
	"context"
	"fmt"
	"golang.org/x/text/language"
	"time"
)

type cashItem struct {
	timestamp time.Time
	data      string
}

// Service is a Translator user.
type Service struct {
	translator Translator
}

var cash = make(map[string]cashItem)

const retries = 3
const cashDurationSeconds = 15

func init() {
	//remove outdated data in cash every 5 sec
	var cashChecker = time.Tick(time.Second * 5)
	go func() {
		for timer := range cashChecker {
			fmt.Println("clear cash ", timer)
			for k, v := range cash {
				if time.Since(v.timestamp) < time.Second*cashDurationSeconds {
					delete(cash, k)
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
	fmt.Println(key)
	var err error = nil
	for i := 1; i <= retries; i++ {
		if cashedTranslation, ok := cash[key]; ok {
			return cashedTranslation.data, nil
		}
		translation, err := s.translator.Translate(ctx, from, to, data)
		fmt.Println(err, i)
		if err == nil {
			cash[key] = cashItem{data: translation, timestamp: time.Now()}
			return translation, err
		}
		//retry after delay if request is unsuccessful
		for j := 0; j < i*i; j++ {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
				// check every second if another same request put result in cash
			case <-time.After(time.Second):
				if cashedTranslation, ok := cash[key]; ok {
					return cashedTranslation.data, nil
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
