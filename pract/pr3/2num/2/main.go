package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/reactivex/rxgo/v2"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	digits := []rune("0123456789")

	letterObservable := rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		for i := 0; i < 1000; i++ {
			next <- rxgo.Of(string(letters[rand.Intn(len(letters))]))
		}
	}})

	digitObservable := rxgo.Create([]rxgo.Producer{func(ctx context.Context, next chan<- rxgo.Item) {
		for i := 0; i < 1000; i++ {
			next <- rxgo.Of(string(digits[rand.Intn(len(digits))]))
		}
	}})

	letterChan := letterObservable.Observe()
	digitChan := digitObservable.Observe()

	for i := 0; i < 1000; i++ {
		select {
		case letterItem := <-letterChan:
			digitItem := <-digitChan
			fmt.Println(letterItem.V.(string) + digitItem.V.(string))
		}
	}
}
