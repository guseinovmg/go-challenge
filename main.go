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
	fmt.Println(s.translator.Translate(ctx, language.English, language.Japanese, "test"))
}
