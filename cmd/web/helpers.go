package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

/*
Помощник serverError записывает сообщение об ошибке в errorLog и
затем отправляет пользователю ответ 500 "Внутренняя ошибка сервера".
*/
func (app *application) serverError(write http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace) //использовать логгер Output() и установив глубину вызова на 2

	http.Error(write, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

/*
Помощник clientError отправляет определенный код состояния и соответствующее описание
пользователю. Мы будем использовать это в следующий уроках, чтобы отправлять ответы вроде 400 "Bad
Request", когда есть проблема с пользовательским запросом.
*/
func (app *application) clientError(write http.ResponseWriter, status int) {
	http.Error(write, http.StatusText(status), status)
}

/*
Мы также реализуем помощник notFound. Это простоудобная оболочка вокруг clientError,
которая отправляет пользователю ответ "404 Страница не найдена".
*/
func (app *application) notFound(write http.ResponseWriter) {
	app.clientError(write, http.StatusNotFound)
}

// слздаем метод render, чтобы было можно легко отображать шаблоны из кэша
func (app *application) render(write http.ResponseWriter, request *http.Request, name string, td *templateData) {
	// Извлекаем соответствующий набор шаблонов из кэша в зависимости от названия страницы
	// (например, 'home.page.tmpl'). Если в кэше нет записи запрашиваемого шаблона, то
	// вызывается вспомогательный метод serverError(), который мы создали ранее.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(write, fmt.Errorf("Шаблон %s не существует!", name))
		return
	}
	// Рендерим файлы шаблона, передавая динамические данные из переменной `td`.
	err := ts.Execute(write, td)
	if err != nil {
		app.serverError(write, err)
	}
}
