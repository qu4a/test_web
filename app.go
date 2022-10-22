package main

import (
	"log"
	"net/http"
)

func index(write http.ResponseWriter, request *http.Request) {
	write.Write([]byte("Test"))

}
func showSnippet(write http.ResponseWriter, request *http.Request) {
	write.Write([]byte("showSnippet"))

}
func createSnippet(write http.ResponseWriter, request *http.Request) {
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
