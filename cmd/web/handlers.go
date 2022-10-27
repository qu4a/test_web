package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

/*
Меняем сигнатуры обработчика home, чтобы он определялся как метод
структуры *application.
*/

func (app *application) index(write http.ResponseWriter, request *http.Request) {
	// Проверяется, если текущий путь URL запроса точно совпадает с шаблоном "/". Если нет, вызывается
	// функция http.NotFound() для возвращения клиенту ошибки 404.
	// Важно, чтобы мы завершили работу обработчика через return. Если мы забудем про "return", то обработчик
	// продолжит работу и выведет сообщение "Привет из SnippetBox" как ни в чем не бывало.
	if request.URL.Path != "/" {
		http.NotFound(write, request)
		return
	}

	// Инициализируем срез содержащий пути к двум файлам. Обратите внимание, что
	// файл home.page.tmpl должен быть *первым* файлом в срезе.
	files := []string{
		"ui/html/home.page.tmpl",
		"ui/html/base.layout.tmpl",
		"ui/html/footer.partial.tmpl",
	}

	//write.Write([]byte("Test"))
	// Используем функцию template.ParseFiles() для чтения файла шаблона.
	// Если возникла ошибка, мы запишем детальное сообщение ошибки и
	// используя функцию http.Error() мы отправим пользователю
	// ответ: 500 Internal Server Error (Внутренняя ошибка на сервере)
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(write, "Internal Server Error", 500)
		return
	}
	// Затем мы используем метод Execute() для записи содержимого
	// шаблона в тело HTTP ответа. Последний параметр в Execute() предоставляет
	// возможность отправки динамических данных в шаблон.
	err = ts.Execute(write, nil)
	/*
		Обновляем код для использования логгера-ошибок из структуры application.
	*/
	if err != nil {
		/*
			Поскольку обработчик home теперь является методом структуры application
			он может получить доступ к логгерам из структуры.
			Используем их вместо стандартного логгера от Go.
		*/
		app.errorLog.Println(err.Error())
		http.Error(write, "Internal Server Error", 500)
		return
	}
}

/*
Меняем сигнатуру обработчика showSnippet, чтобы он был определен как метод
структуры *application
*/
func (app *application) showSnippet(write http.ResponseWriter, request *http.Request) {
	// Извлекаем значение параметра id из URL и попытаемся
	// конвертировать строку в integer используя функцию strconv.Atoi(). Если его нельзя
	// конвертировать в integer, или значение меньше 1, возвращаем ответ
	// 404 - страница не найдена!
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(write, request)
		return
	}
	fmt.Fprintf(write, "Отображение выбранной заметки с ID %d...", id)

	//write.Write([]byte("showSnippet"))

}

// обработчик для создания заметок
func (app *application) createSnippet(write http.ResponseWriter, request *http.Request) {
	// Используем r.Method для проверки, использует ли запрос метод POST или нет. Обратите внимание,
	// что http.MethodPost является строкой и содержит текст "POST".
	if request.Method != http.MethodPost {
		write.Header().Set("Allow", http.MethodPost) //добавление заголовка к http запросу
		// Используем функцию http.Error() для отправки кода состояния 405 с соответствующим сообщением.
		http.Error(write, "Метод запрещен", 405)

		/*уже не нужены строки, т.к http.Error обрабатывает запрос и выводит тоже самое
		write.WriteHeader(405)
		write.Write([]byte("Get - метод запрещен"))
		*/
		return
	}
	/*
		Под json:
		write.Header().Set("Content-Type", "application/json")
		write.Write([]byte(`{"name":"Alex"}`))
	*/
	write.Write([]byte("Новая заметка"))

}
