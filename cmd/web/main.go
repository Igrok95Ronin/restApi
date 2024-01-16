package main

import (
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Создание HTTP маршрутизатора (ServeMux)
	mux := http.NewServeMux()

	// Регистрация обработчиков для различных путей
	mux.HandleFunc("/", home)                        // Обработчик для корневого пути
	mux.HandleFunc("/snippet", showSnippet)          // Обработчик для пути /snippet
	mux.HandleFunc("/snippet/create", createSnippet) // Обработчик для пути /snippet/create

	// Создание файлового сервера с настраиваемой файловой системой
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})

	// Настройка маршрутов для статических файлов
	mux.Handle("/static", http.NotFoundHandler())                   // Возвращает ошибку 404 для пути /static
	mux.Handle("/static/", http.StripPrefix("/static", fileServer)) // Обслуживает файлы из /static/

	// Запуск HTTP-сервера
	log.Println("Запуск сервера на http://127.0.0.1:4000")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err) // Логирование критической ошибки сервера
	}
}

// neuteredFileSystem представляет файловую систему с дополнительной безопасностью
type neuteredFileSystem struct {
	fs http.FileSystem // Основная файловая система
}

// Open переопределяет стандартный метод открытия файлов
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err // Возвращает ошибку, если файл/директория не найдена
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err // Возвращает ошибку, если невозможно получить информацию о файле/директории
	}

	if s.IsDir() {
		// Проверка наличия index.html в директории
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr // Возвращает ошибку, если невозможно закрыть директорию
			}

			return nil, err // Возвращает ошибку, если index.html не найден
		}
	}

	return f, nil // Возвращает файл/директорию
}
