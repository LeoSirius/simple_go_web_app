package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// 如果浏览器输入127.0.0.1:8888/go
	// 那么r.URL的值是 “/go”
	// 这里切片去掉了最前面的"/"
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)

	// ListenAndServe返回的话，必然是error
	log.Fatal(http.ListenAndServe(":8888", nil))
}