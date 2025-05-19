package main

import (
	"log"
	"net/http"
)

const PORT = ":1450"

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/run", runHandler)
	log.Printf("Server on %s", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
