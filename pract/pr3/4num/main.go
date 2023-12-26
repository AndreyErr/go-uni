package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/reactivex/rxgo/v2"
)

type File struct {
	Type string
	Size int
}

func workerXML(wg *sync.WaitGroup, file File) {
	defer wg.Done()
	fmt.Println("Received XML file:", file)
	time.Sleep(time.Duration(file.Size*7) * time.Millisecond)
	fmt.Println("Processed XML file:", file)
}

func workerJSON(wg *sync.WaitGroup, file File) {
	defer wg.Done()
	fmt.Println("Received JSON file:", file)
	time.Sleep(time.Duration(file.Size*7) * time.Millisecond)
	fmt.Println("Processed JSON file:", file)
}

func workerXLS(wg *sync.WaitGroup, file File) {
	defer wg.Done()
	fmt.Println("Received XLS file:", file)
	time.Sleep(time.Duration(file.Size*7) * time.Millisecond)
	fmt.Println("Processed XLS file:", file)
}

func main() {
	ch := make(chan rxgo.Item)
	fileTypes := []string{"XML", "JSON", "XLS"}
	var wg sync.WaitGroup

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Duration(rand.Intn(900)+100) * time.Millisecond)
			ch <- rxgo.Of(File{
				Type: fileTypes[rand.Intn(len(fileTypes))],
				Size: rand.Intn(91) + 10,
			})
		}
		close(ch)
	}()

	observable := rxgo.FromChannel(ch, rxgo.WithBufferedChannel(5))

	for item := range observable.Observe() {
		if item.Error() {
			fmt.Println("Received an error:", item.E.Error())
		} else {
			file := item.V.(File)
			wg.Add(1)
			switch file.Type {
			case "XML":
				go workerXML(&wg, file)
			case "JSON":
				go workerJSON(&wg, file)
			case "XLS":
				go workerXLS(&wg, file)
			default:
				fmt.Println("Unknown file type:", file.Type)
				wg.Done()
			}
		}
	}

	wg.Wait()
}
