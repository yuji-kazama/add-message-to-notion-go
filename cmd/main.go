package main

import (
	"log"
	"net/http"

	"github.com/leokite/add-message-to-notion-go"
)

func main() {
	http.HandleFunc("/", function.Function)
	log.Println("[INFO] Server listening")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("[FATAL] Failed listing server: %v", err)
	}
}