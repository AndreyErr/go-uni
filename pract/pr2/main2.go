package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/otiai10/copy"
)


func main() {
    copyFile := "tekst.txt"
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	fmt.Printf("Затраты памяти начальные: %d bytes\n", mem.Alloc)

    // 1 СПОСОБ FileInput
	start := time.Now()
	copyFileDest := "copyDir/tekst1.txt"
	readText, err := ioutil.ReadFile(copyFile)
	if err != nil {
        fmt.Println("1 Ошибка чтения файла:", err)
        return
    }
	err = ioutil.WriteFile(copyFileDest, readText, 0644)
    if err != nil {
        fmt.Println("2 Ошибка записи файла:", err)
        return
    }
	duration := time.Since(start).Seconds()
	fmt.Printf("Время выполнения FileInput: %f секунд\n", duration)

	runtime.ReadMemStats(&mem)
	fmt.Printf("Затраты памяти для FileInput: %d bytes\n", mem.Alloc)

    // 2 СПОСОБ FileChannal
	start = time.Now()
	copyFileDest = "copyDir/tekst2.txt"
	middle := make([]byte, 1024)

	textRead2, err := os.Open(copyFile)
	if err != nil {
        fmt.Println("3 Ошибка чтения файла:", err)
        return
    }
	defer textRead2.Close()
	newFile2, err := os.Create(copyFileDest)
	if err != nil {
        fmt.Println("4 Ошибка создания нового файла:", err)
        return
    }
	defer newFile2.Close()
	for {
		middleText, err := textRead2.Read(middle)
		if err != nil && err != io.EOF {
            fmt.Println("5 Ошибка чтения файла:", err)
            return
        }
		if middleText == 0 {
			break
		}
        _, err = newFile2.Write(middle[:middleText])
        if err != nil {
            fmt.Println("6 Ошибка записи файла:", err)
            return
        }
	}
	duration = time.Since(start).Seconds()
	fmt.Printf("Время выполнения FileChannal: %f секунд\n", duration)
	runtime.ReadMemStats(&mem)
	fmt.Printf("Затраты памяти для FileChannal: %d bytes\n", mem.Alloc)

    // 3 СПОСОБ Apache
	start = time.Now()
    duplicate3 := "copyDir/tekst3.txt"

    _ = copy.Copy(copyFile, duplicate3)
	duration = time.Since(start).Seconds()
	fmt.Printf("Время выполнения Apache: %f секунд\n", duration)
	runtime.ReadMemStats(&mem)
	fmt.Printf("Затраты памяти для Apache: %d bytes\n", mem.Alloc)
}