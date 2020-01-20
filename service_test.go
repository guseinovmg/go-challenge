package main

import (
	"context"
	"golang.org/x/text/language"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var translator = newRandomTranslator(
	100*time.Millisecond,
	150*time.Millisecond,
	0.0,
)

var failedTranslator = newRandomTranslator(
	100*time.Millisecond,
	150*time.Millisecond,
	1.0,
)

func TestSingle(t *testing.T) {
	ctx := context.Background()
	rand.Seed(time.Now().UTC().UnixNano())
	s := NewService(translator)
	_, err := s.translator.Translate(ctx, language.English, language.Japanese, "test")

	if err != nil {
		t.Errorf("Got error: \"%v\"", err)
	}
}

func TestCanceledContext(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	s := NewService(translator)
	canceledContext, cancelFunc := context.WithCancel(context.Background())
	cancelFunc()
	_, err := s.Translate(canceledContext, language.English, language.Japanese, "test")
	if err == nil {
		t.Error("Must return error")
	}
}

func TestFailedTranslator(t *testing.T) {
	ctx := context.Background()
	rand.Seed(time.Now().UTC().UnixNano())
	s := NewService(failedTranslator)
	_, err := s.Translate(ctx, language.English, language.Japanese, "test")
	if err == nil {
		t.Error("Must return error")
	}
}

func TestSequence(t *testing.T) {
	ctx := context.Background()
	rand.Seed(time.Now().UTC().UnixNano())
	s := NewService(translator)
	res, _ := s.Translate(ctx, language.English, language.Japanese, "test")
	res2, _ := s.Translate(ctx, language.English, language.Japanese, "test")

	if res != res2 {
		t.Errorf("Results does not match: got \"%s\" and \"%s\"", res, res2)
	}
}

func TestConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	var res [5]string

	ctx := context.Background()
	rand.Seed(time.Now().UTC().UnixNano())
	s := NewService(translator)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(s *Service, ctx context.Context, str *string) {
			defer wg.Done()

			res, err := s.Translate(ctx, language.English, language.Japanese, "test")

			if err == nil {
				*str = res
			}
		}(s, ctx, &res[i])
	}

	wg.Wait()

	for i := 1; i < 5; i++ {
		if res[0] != res[i] {
			t.Errorf("Results does not match: %v", res)
			break
		}
	}
}
