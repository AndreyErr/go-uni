package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type fileState struct {
	Size    int64
	Content string
	IsDir   bool
}

func main() {
	// Задаем начальный каталог, который будем мониторить
	rootPath := "123"
	lastFileStates := make(map[string]fileState)

	for {
		// Просматриваем все файлы и папки внутри каталога
		err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("Ошибка при поиске файлов и папок:", err)
				return err
			}

			// Игнорируем текущий каталог
			if path == rootPath {
				return nil
			}

			state, exists := lastFileStates[path]

			if !exists {
				if info.IsDir() {
					// Если папка ранее не существовала, выводим ее как созданную
					fmt.Printf("Создана новая папка: %s\n", path)
					lastFileStates[path] = fileState{IsDir: true}
				} else {
					// Если файл ранее не существовал, выводим его как новый
					fmt.Printf("Создан новый файл: %s\n", path)
					content := getFileContent(path)
					size := getFileSize(path)
					lastFileStates[path] = fileState{Size: size, Content: content}
				}
			} else if info.IsDir() {
				// Проверяем, если это папка, то обновляем состояние
				lastFileStates[path] = fileState{IsDir: true}
			} else {
				// Сравниваем текущее содержимое файла с предыдущим
				content := getFileContent(path)
				size := getFileSize(path)

				if state.Content != content {
					fmt.Printf("Файл изменен: %s\n", path)
					printFileChanges(path, state.Content, content)
				}

				if state.Size != size {
					fmt.Printf("Размер файла изменен: %s\n", path)
					fmt.Printf("Размер до: %d байт, Размер после: %d байт\n", state.Size, size)
				}

				lastFileStates[path] = fileState{Size: size, Content: content}
			}

			return nil
		})

		if err != nil {
			fmt.Print(err)
		}

		// Проверяем, были ли удалены файлы или папки
		for path, state := range lastFileStates {
			_, exists := os.Stat(path)
			if os.IsNotExist(exists) {
				if state.IsDir {
					// Если папка была удалена, выводим сообщение об удалении папки
					fmt.Printf("Папка удалена: %s\n", path)
				} else {
					// Если файл был удален, выводим его контрольную сумму и размер
					fmt.Printf("Файл удален: %s\n", path)
					fmt.Printf("16-битная контрольная сумма: %s\n", calculateChecksum([]byte(state.Content)))
					fmt.Printf("Размер файла: %d байт\n", state.Size)
				}
				delete(lastFileStates, path)
			}
		}

		// Засыпаем на некоторое время перед следующей проверкой
		time.Sleep(5 * time.Second)
	}
}

func getFileContent(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Ошибка чтения файла %s: %v\n", filePath, err)
		return ""
	}
	return string(data)
}

func getFileSize(filePath string) int64 {
	info, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("Ошибка при получении информации о файле %s: %v\n", filePath, err)
		return 0
	}
	return info.Size()
}

func calculateChecksum(data []byte) string {
	var checksum uint16 = 0xFFFF // Инициализируем контрольную сумму
	for _, b := range data {      // Для каждого байта
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
	return fmt.Sprintf("%d", checksum)
}

func printFileChanges(path string, oldContent, newContent string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(oldContent, newContent, false)

	fmt.Printf("Изменения в файле %s:\n", path)
	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			fmt.Printf("- %s\n", diff.Text)
		case diffmatchpatch.DiffInsert:
			fmt.Printf("+ %s\n", diff.Text)
		}
	}
}