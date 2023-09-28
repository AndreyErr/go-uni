package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
    // Укажите путь к файлу, который вы хотите прочитать
    filePath := "tekst.txt"

    // Чтение содержимого файла
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        fmt.Println("Ошибка чтения файла:", err)
        return
    }

    // Вывод содержимого файла в стандартный поток вывода (консоль)
    fmt.Println(string(data))
}