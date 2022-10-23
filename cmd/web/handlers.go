package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// с обработчиком главноей страницы
func index(write http.ResponseWriter, request *http.Request) {
	// Проверяется, если текущий путь URL запроса точно совпадает с шаблоном "/". Если нет, вызывается
	// функция http.NotFound() для возвращения клиенту ошибки 404.
	// Важно, чтобы мы завершили работу обработчика через return. Если мы забудем про "return", то обработчик
	// продолжит работу и выведет сообщение "Привет из SnippetBox" как ни в чем не бывало.
	if request.URL.Path != "/" {
		http.NotFound(write, request)
		return
	}
	write.Write([]byte("Test"))
}
func showSnippet(write http.ResponseWriter, request *http.Request) {
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
func createSnippet(write http.ResponseWriter, request *http.Request) {
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
	write.Header().Set("Content-Type", "application/json")
	write.Write([]byte(`{"name":"Alex"}`))

}
