package main

import (
	"fmt"
	"io/ioutil"
)

func main(){
    filePath := "tekst.txt"

    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        fmt.Println("Ошибка чтения файла:", err)
        return
    }



    var checksum uint16 = 0xFFFF // Инициализируем контрольную сумму
	for _, b := range data { // Для каждого байта
		checksum ^= uint16(b) // XOR байта с суммой

        // Для каждого бита 
        // Если младший бит равен 1, то контрольная сумма сдвигается на 1 бит и выполняется XOR с полиномом 0xA001
		for i := 0; i < 8; i++ {
			if (checksum & 0x0001) == 1 {
				checksum = (checksum >> 1) ^ 0xA001
			} else {
				checksum = checksum >> 1
			}
		}
	}

    fmt.Printf("16-битная контрольная сумма: %d \n", checksum)

}