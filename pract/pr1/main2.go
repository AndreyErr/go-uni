package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// Вычисление квадрата
func sq(num int, resultChan chan int) {
	delay := time.Duration(rand.Intn(5)+1) * time.Second
	time.Sleep(delay)

	result := int(math.Pow(float64(num), 2))

	resultChan <- result
}

func main() {
	var wg sync.WaitGroup

	for {
		var input int
		fmt.Print("Введите число: ")
		_, err := fmt.Scanf("%d\n", &input)
		if err != nil {
			fmt.Println("Ошибка ввода.")
			continue
		}
		resultChan := make(chan int)
		wg.Add(1)
		go func(input int){
			defer wg.Done()
			sq(input, resultChan)
		}(input)
		go func() {
			result := <-resultChan
			fmt.Printf("Результат: %d\n", result)
		}()
	}
	// Ожидание завершения всех горутин
	wg.Wait()
}