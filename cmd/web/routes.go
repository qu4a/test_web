package main

import "net/http"

/*
Рефактор кода. Переносим объявление маршрутов
*/
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.index)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	/*Инициализируем FileServer, он будет обрабатывать
	  HTTP-запросы к статическим файлам из папки "./ui/static".
	  Обратите внимание, что переданный в функцию http.Dir путь
	  является относительным корневой папке проекта*/
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Используем функцию mux.Handle() для регистрации обработчика для
	// всех запросов, которые начинаются с "/static/". Мы убираем
	// префикс "/static" перед тем как запрос достигнет http.FileServer
	//mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
