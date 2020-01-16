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
const cashDurationSeconds = 30

func init() {
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
	default:
	}
	key := fmt.Sprintf("from %s, to %s  data %s", from, to, data)
	fmt.Println(key)
	var err error = nil
	for i := 1; i <= retries; i++ {
		item, ok := cash[key]
		if ok {
			return item.data, nil
		}
		fmt.Println(cash, ok, i)
		var translation string
		translation, err = s.translator.Translate(ctx, from, to, data)
		fmt.Println(err, i)
		if err == nil {
			cash[key] = cashItem{data: translation, timestamp: time.Now()}
			return translation, err
		}
		for j := 0; j < i*i; j++ {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(time.Second):
				item, ok := cash[key]
				if ok {
					return item.data, nil
				}
			}
		}
	}
	return "", err
}

func NewService() *Service {
	t := newRandomTranslator(
		100*time.Millisecond,
		500*time.Millisecond,
		0.4,
	)

	return &Service{
		translator: t,
	}
}
