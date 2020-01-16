package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/text/language"
)

func main() {
	ctx := context.Background()
	rand.Seed(time.Now().UTC().UnixNano())
	s := NewService()
	fmt.Println(s.Translate(ctx, language.English, language.Japanese, "test"))
	fmt.Println(s.Translate(ctx, language.English, language.Japanese, "test"))
	time.Sleep(time.Second * 40)
	ctx, cancelFunc := context.WithCancel(ctx)
	cancelFunc()
	fmt.Println(s.Translate(ctx, language.English, language.Japanese, "test"))
	ctx = context.Background()
	fmt.Println(s.Translate(ctx, language.English, language.Japanese, "test"))
}
