package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var uploadPath = "/app/data"

func initUploadPath() error {
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		err := os.MkdirAll(uploadPath, 0755)
		if err != nil {
			return fmt.Errorf("не удалось создать папку для загрузки: %v", err)
		}
	}
	return nil
}

func saveFile(file io.Reader, fileName string) (string, error) {
	fileName = uuid.New().String() + "_" + fileName
	filePath := filepath.Join(uploadPath, fileName)
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("не удалось сохранить файл: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", fmt.Errorf("не удалось сохранить файл: %v", err)
	}

	return fileName, nil
}

func header() string {
	header := `<!DOCTYPE html>
	<html lang='en'>
	<head>
	  <meta charset='UTF-8'>
	  <meta name='viewport' content='width=device-width, initial-scale=1.0'>
	  <title>Архив гос документов</title>
	</head>
	<body>
	  <h1>Архив гос документов</h1>`
	return string(header)
}

func showFiles(w http.ResponseWriter) {
	defer fmt.Fprintf(w, `
		<br><a href="/">на главную</a>
		</body>
		</html>`)

	files, err := os.ReadDir(uploadPath)
	if err != nil {
		http.Error(w, "Не удалось получить файлы", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, `
		<form enctype="multipart/form-data" action="/upload" method="post">
			<input type="file" name="file">
			<br>
			<input type="submit" value="Загрузить">
		</form>
		<h3>Файлы</h3>`)

	if len(files) != 0 {
		for _, file := range files {
			fileInfo, err := file.Info()
			if err != nil {
				http.Error(w, "Не удалось получить информацию о файле", http.StatusInternalServerError)
				return
			}
			n := strings.SplitN(fileInfo.Name(), "_", 2)[1]
			fmt.Fprintf(w, "%s | <a href=\"/show/%s\" target=\"_blank\">Показать</a> | <a href=\"/download/%s\" target=\"_blank\">Скачать</a><br>", n, fileInfo.Name(), fileInfo.Name())
		}
	} else {
		fmt.Fprintf(w, "Пусто")
	}
}

func main() {
	initUploadPath()

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, header())
		fmt.Fprintf(w, `Сервер: %s`, os.Getenv("SERVER_NAME"))
		if r.Method == http.MethodPost {
			file, header, err := r.FormFile("file")
			if err != nil {
				http.Error(w, "Не удалось загрузить файл", http.StatusBadRequest)
				showFiles(w)
				return
			}
			defer file.Close()

			_, err = file.Seek(0, io.SeekStart)
			if err != nil {
				http.Error(w, "Не удалось прочитать файл", http.StatusBadRequest)
				showFiles(w)
				return
			}

			_, err = saveFile(file, header.Filename)
			if err != nil {
				http.Error(w, "Не удалось сохранить файл", http.StatusBadRequest)
				showFiles(w)
				return
			}

			fmt.Fprintf(w, `<p>Файл успешно загружен</p>`)
			showFiles(w)
		} else {
			http.Error(w, "Недопустимый метод HTTP", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, header())
		fmt.Fprintf(w, `Сервер: %s`, os.Getenv("SERVER_NAME"))
		if r.Method == http.MethodGet {
			showFiles(w)
		} else {
			http.Error(w, "Недопустимый метод HTTP", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/show/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fileName := r.URL.Path[len("/show/"):]
			filePath := filepath.Join(uploadPath, fileName)
			http.ServeFile(w, r, filePath)
		} else {
			http.Error(w, "Недопустимый метод HTTP", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fileName := r.URL.Path[len("/download/"):]
			filePath := filepath.Join(uploadPath, fileName)
	
			file, err := os.Open(filePath)
			if err != nil {
				http.Error(w, "Файл не найден", http.StatusNotFound)
				return
			}
			defer file.Close()
	
			fileInfo, err := file.Stat()
			if err != nil {
				http.Error(w, "Ошибка чтения информации о файле", http.StatusInternalServerError)
				return
			}
	
			w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Content-Length", fmt.Sprint(fileInfo.Size()))
	
			http.ServeContent(w, r, fileName, fileInfo.ModTime(), file)
		} else {
			http.Error(w, "Недопустимый метод HTTP", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}
