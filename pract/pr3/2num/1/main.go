package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/reactivex/rxgo/v2"
)

func main() {
	numbers := rxgo.Range(0, 1000).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		num := rand.Intn(1001)
		return num * num, nil
	})

	for item := range numbers.Observe() {
		fmt.Println(item.V)
	}
}