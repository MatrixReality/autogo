package main

import (
	"net/http"
)

func main() {
	static := http.FileServer(http.Dir("./static"))
	http.Handle("/", static)

	http.ListenAndServe(":8081", nil)

	//http.HandleFunc("/", handler)
	//http.ListenAndServe(":8080", nil)
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "./index.html")
// }
