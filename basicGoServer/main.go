package main

import (
	"fmt"
	"log"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/hello" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not supported..", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello")
}

func FormHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err %v", err)
		return
	}
	fmt.Fprintf(w, "POST successful..")

	name := r.FormValue("name")
	address := r.FormValue("address")

	fmt.Fprintf(w, "Name = %v\n", name)
	fmt.Fprintf(w, "address = %v\n", address)

}
func main() {
	fileServer := http.FileServer((http.Dir("./static")))
	http.Handle("/", fileServer)
	http.HandleFunc("/hello", HelloHandler)
	http.HandleFunc("/form", FormHandler)

	fmt.Printf(("Starting Servrt on port 8080\n"))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
