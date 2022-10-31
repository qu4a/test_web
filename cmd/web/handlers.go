package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"test_web/pkg/models"
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
		app.notFound(write) //Использование помощника notFound()
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
		//app.errorLog.Println(err.Error())
		app.serverError(write, err) //Использование помощника serverError
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
		//app.errorLog.Println(err.Error())
		app.serverError(write, err)
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
		app.notFound(write) //Использование помощника notFound()
		return
	}
	// Вызываем метода Get из модели Snipping для извлечения данных для
	// конкретной записи на основе её ID. Если подходящей записи не найдено,
	// то возвращается ответ 404 Not Found (Страница не найдена).
	s, err := app.snippet.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNotRecord) {
			app.notFound(write)
		} else {
			app.serverError(write, err)
		}
		return
	}
	//отображаем весь вывод на странице
	fmt.Fprintf(write, "%v", s)

	//write.Write([]byte("showSnippet"))

}

// обработчик для создания заметок
func (app *application) createSnippet(write http.ResponseWriter, request *http.Request) {
	// Используем r.Method для проверки, использует ли запрос метод POST или нет. Обратите внимание,
	// что http.MethodPost является строкой и содержит текст "POST".
	if request.Method != http.MethodPost {
		write.Header().Set("Allow", http.MethodPost) //добавление заголовка к http запросу
		// Используем функцию http.Error() для отправки кода состояния 405 с соответствующим сообщением.
		app.clientError(write, http.StatusMethodNotAllowed)

		/*уже не нужены строки, т.к http.Error обрабатывает запрос и выводит тоже самое
		write.WriteHeader(405)
		write.Write([]byte("Get - метод запрещен"))
		*/
		return
	}
	title := "История про улитку"
	content := "Улитка выползла из раковины,\nвытянула рожки,\nи опять подобрала их."
	expires := "7"
	id, err := app.snippet.Insert(title, content, expires)
	if err != nil {
		app.serverError(write, err)
		return
	}
	/*
		Под json:
		write.Header().Set("Content-Type", "application/json")
		write.Write([]byte(`{"name":"Alex"}`))
	*/
	//write.Write([]byte("Новая заметка"))
	// Перенаправляем пользователя на соответствующую страницу заметки.
	http.Redirect(write, request, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

}
