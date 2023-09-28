package main

import (
	"fmt"
	"math/rand"
	"time"
)

type FileType int

var numOfTypes = 3

const (
	XML FileType = iota
	JSON
	XLS
)

type File struct {
	Type   FileType
	Size   int
}

func main() {
	queue := make(chan *File, 5)

	go generateFiles(queue)
	processFiles(queue)
}

func generateFiles(queue chan *File) {
	file := new(File)

	for {
		delay := rand.Intn(9) + 100
		time.Sleep(time.Duration(delay) * time.Millisecond)
		file.Type = FileType(rand.Intn(numOfTypes))
        file.Size = rand.Intn(91) + 10 // Размер файла от 10 до 100
        select {
		case queue <- file:
			fmt.Printf("-->Файл добавлен в очередь: %v\n", file)
		default:
			fmt.Println("-!->Очередь полна. Ожидание...")
			queue <- file // Блокировка до того момента, пока место не освободится в очереди
			fmt.Printf("-->Файл добавлен в очередь: %v\n", file)
		}
	}
}

func processFiles(queue chan *File) {
    var fXML, fXLS, fJSON int
	for file := range queue {
		switch file.Type {
		case 0:
			fXML += 1
		case 1:
			fXLS += 1
		case 2:
			fJSON += 1
		default:
			fmt.Println("Нет подобного обработчика для", file.Type)
		}
        delay := time.Duration(file.Size*7) * time.Millisecond
		time.Sleep(delay)
        fmt.Printf("------\n")
		fmt.Printf("XML обработано: %d\n", fXML)
		fmt.Printf("JSON обработано: %d\n", fXLS)
		fmt.Printf("XLS обработано: %d\n", fJSON)
        fmt.Printf("Всего: %d\n", fXML + fXLS + fJSON)
	}
}