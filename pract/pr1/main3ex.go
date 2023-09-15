package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type FileType int

var numOfTypes = 3
var mu sync.Mutex

const (
	XML FileType = iota
	JSON
	XLS
)

type File struct {
	Type   FileType
	Size   int
    Id int
}

type FileEx struct {
	Type   FileType
	Size   int
    Proc   int
    Id int
}

func main() {
	queue := make(chan *File, 5)
    queueEx := make(chan *FileEx)

	go generateFiles(queue)
	numProcessors := 3
	for i := 0; i < numProcessors; i++ {
		go processFiles(queue, queueEx, i)
        fmt.Printf("Процессор %d создан!\n", i)
	}

    var fXML, fXLS, fJSON int
	for file := range queueEx {
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
        fmt.Printf("-----Proc %d-----ID: %d-----\n", file.Proc, file.Id)
		fmt.Printf("XML обработано: %d\n", fXML)
		fmt.Printf("JSON обработано: %d\n", fXLS)
		fmt.Printf("XLS обработано: %d\n", fJSON)
        fmt.Printf("Всего: %d\n", fXML + fXLS + fJSON)
	}
}

func generateFiles(queue chan *File) {
    mu.Lock()
	defer mu.Unlock()
	file := new(File)
    var id int
	for{
        id++
		delay := rand.Intn(9) + 100
		time.Sleep(time.Duration(delay) * time.Millisecond)
		file.Type = FileType(rand.Intn(numOfTypes))
        file.Size = rand.Intn(91) + 10 // Размер файла от 10 до 100
        file.Id = id
        select {
		case queue <- file:
			fmt.Printf("--->Файл добавлен в очередь (тип, размер, id): %v\n", file)
		default:
			fmt.Println("-!->Очередь полна. Ожидание...")
			queue <- file // Блокировка до того момента, пока место не освободится в очереди
            fmt.Printf("-+->Файл добавлен в очередь: %v\n", file)
		}
	}
}

func processFiles(queue chan *File, queueEx chan *FileEx, proc int) {
    fileEx := new(FileEx)
	for file := range queue {
        delay := time.Duration(file.Size*7) * time.Millisecond
		time.Sleep(delay) 
        fileEx.Type = file.Type
        fileEx.Size = file.Size
        fileEx.Proc = proc
        fileEx.Id = file.Id
        queueEx <- fileEx
	}
}