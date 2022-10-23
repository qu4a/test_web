package main

import (
	"log"
	"net/http"
)

func main() {
	//http.HandleFunc()
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	//http.HandleFunc("/hello", viewHandler)
	log.Println("run server")
	err := http.ListenAndServe("localhost:8080", mux)
	log.Fatal(err)
}
