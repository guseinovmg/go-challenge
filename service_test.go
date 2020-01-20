package main

import (
	"context"
	"golang.org/x/text/language"
	"math/rand"
	"testing"
	"time"
)

func TestInc(t *testing.T) {
	ctx := context.Background()
	rand.Seed(time.Now().UTC().UnixNano())
	translator := newRandomTranslator(
		100*time.Millisecond,
		150*time.Millisecond,
		0.0,
	)
	s := NewService(translator)
	_, err := s.Translate(ctx, language.English, language.Japanese, "test")
	if err != nil {
		t.Error(err)
	}
	canceledContext, cancelFunc := context.WithCancel(ctx)
	cancelFunc()
	_, err = s.Translate(canceledContext, language.English, language.Japanese, "test2")
	if err == nil {
		t.Error("Must return error")
	}
	failedTranslator := newRandomTranslator(
		100*time.Millisecond,
		150*time.Millisecond,
		1.0,
	)
	s = NewService(failedTranslator)
	_, err = s.Translate(ctx, language.English, language.Japanese, "test3")
	if err == nil {
		t.Error("Must return error")
	}
}
