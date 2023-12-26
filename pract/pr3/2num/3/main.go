package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/reactivex/rxgo/v2"
)

func main() {
	numbers := rxgo.Range(0, 4).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return rand.Int(), nil
	}).Skip(3)

	for item := range numbers.Observe() {
		fmt.Println(item.V)
	}
}