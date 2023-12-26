package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/reactivex/rxgo/v2"
)

func main() {
	temperatureObservable := rxgo.
	Interval(rxgo.WithDuration(time.Second)).
	Map(func(_ context.Context,_ interface{}) (interface{}, error) {
		return rand.Intn(16) + 15, nil
	})

	co2Observable := rxgo.
	Interval(rxgo.WithDuration(time.Second)).
	Map(func(_ context.Context,_ interface{}) (interface{}, error) {
		return rand.Intn(71) + 30, nil
	})

	temperature := temperatureObservable.Observe()
	co2 := co2Observable.Observe()

	for {
		select {
		case temp, open := <-temperature:
			if !open {
				return
			}
			fmt.Printf("temperature: %d \n", temp.V)

			co2Num, open := <-co2
			if !open {
				return
			}
			fmt.Printf("CO2: %d \n", co2Num.V)

			if temp.V.(int) > 25 && co2Num.V.(int) > 70 {
				fmt.Println("ALARM!!!")
			} else if co2Num.V.(int) > 70 {
				fmt.Println("co2 превысила норму")
			} else if temp.V.(int) > 25 {
				fmt.Println("Температура превысила норму")
			}
		}
	}
}
