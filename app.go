package main

import (
	"log"
	"net/http"
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
	write.Write([]byte("showSnippet"))

}

// обработчик для создания заметок
func createSnippet(write http.ResponseWriter, request *http.Request) {
	// Используем r.Method для проверки, использует ли запрос метод POST или нет. Обратите внимание,
	// что http.MethodPost является строкой и содержит текст "POST".
	if request.Method != http.MethodPost {
		write.WriteHeader(405)
		write.Write([]byte("Get - метод запрещен"))
		return
	}
	write.Write([]byte("createSnippet"))

}
func main() {
	//http.HandleFunc()
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	//http.HandleFunc("/hello", viewHandler)
	err := http.ListenAndServe("localhost:8080", mux)
	log.Fatal(err)
}
