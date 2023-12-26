package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Последовательная сумма
func sequentialSum(arr []int) int {
	sum := 0
	for _, val := range arr {
		sum += val
		time.Sleep(1 * time.Millisecond)
	}
	return sum
}

// Сумма через горутины с разделением массивов
func concurrentSumGur(arr []int) int {
	numberThreads := runtime.NumCPU()
	fmt.Printf("numCPU:  %d\n", numberThreads)
	numberThreads += 44
	// numberThreads := 70
	chunkSize := len(arr) / numberThreads
	var wg sync.WaitGroup
	// Канал для сбора результатов из горутин
	resultCh := make(chan int, numberThreads)
	// Запуск горутин для суммирования частей массива
	for i := 0; i < numberThreads; i++ {
		wg.Add(1)
		go func(chunk []int) {
			defer wg.Done()
			resultCh <- sequentialSum(chunk)
		}(arr[i*chunkSize : (i+1)*chunkSize])
	}

	// Ожидание завершения всех горутин
	wg.Wait()

	close(resultCh)

	// Сбор результатов из канала
	totalSum := 0
	for sum := range resultCh {
		totalSum += sum
	}
	return totalSum
}

func main() {
	arrs := 100
	arr := make([]int, arrs)
	for i := 0; i < arrs; i++ {
		arr[i] = i
	}

	startTime := time.Now()
	result := sequentialSum(arr)
	endTime := time.Since(startTime)
	fmt.Printf("Последовательно: Сумма = %d, Время выполнения = %v\n", result, endTime)

	startTime = time.Now()
	result = concurrentSumGur(arr)
	endTime = time.Since(startTime)
	fmt.Printf("С использованием горутин 1 попытка: Сумма = %d, Время выполнения = %v\n", result, endTime)

	startTime = time.Now()
	result = concurrentSumGur(arr)
	endTime = time.Since(startTime)
	fmt.Printf("С использованием горутин 2 попытка: Сумма = %d, Время выполнения = %v\n", result, endTime)
}