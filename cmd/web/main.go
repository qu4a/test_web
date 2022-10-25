package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	/*
		Создаем новый флаг командной строки, значение по умолчанию: ":4000".
		Добавляем небольшую справку, объясняющая, что содержит данный флаг.
		Значение флага будет сохранено в переменной addr.
	*/
	addr := flag.String("addr", ":8080", "Сетевой адрес HTTP")
	/*
		Мы вызываем функцию flag.Parse() для извлечения флага из командной строки.
		Она считывает значение флага из командной строки и присваивает его содержимое
		переменной. Вам нужно вызвать ее *до* использования переменной addr
		иначе она всегда будет содержать значение по умолчанию ":4000".
		Если есть ошибки во время извлечения данных - приложение будет остановлено.
	*/
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	/*Инициализируем FileServer, он будет обрабатывать
	  HTTP-запросы к статическим файлам из папки "./ui/static".
	  Обратите внимание, что переданный в функцию http.Dir путь
	  является относительным корневой папке проекта*/
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})

	// Используем функцию mux.Handle() для регистрации обработчика для
	// всех запросов, которые начинаются с "/static/". Мы убираем
	// префикс "/static" перед тем как запрос достигнет http.FileServer
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	/*
		Значение, возвращаемое функцией flag.String(), является указателем на значение
		из флага, а не самим значением. Нам нужно убрать ссылку на указатель
		то есть перед использованием добавьте к нему префикс *. Обратите внимание, что мы используем
		функцию log.Printf() для записи логов в журнал работы нашего приложения.
	*/
	log.Printf("Запуск сервера на %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}

/*
В данном коде мы создаем настраиваемый тип neuteredFileSystem, который включает в себя http.FileSystem.
Затем мы создаем метод Open(), который вызывается каждый раз, когда http.FileServer получает запрос.

В методе Open() мы открываем вызываемый путь. Используя метод IsDir() мы проверим если вызываемый путь является папкой
или нет. Если это папка, то с помощью метода Stat("index.html") мы проверим если файл index.html существует внутри данной папки.

Если файл index.html не существует, то метод вернет ошибку os.ErrNotExist (которая, в свою очередь, будет преобразована через
http.FileServer в ответ 404 страница не найдена). Мы также вызываем метод Close() для закрытия только, что открытого index.html файла, чтобы избежать утечки файлового дескриптора.
*/
type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
