package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type TestServer struct {
}

func (t *TestServer) Meta(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"env": os.Environ(),
	})
	if err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusOK)
}

func (t *TestServer) Start() {
	http.HandleFunc("/meta", t.Meta)
	log.Println("Listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}

func main() {
	t := TestServer{}
	t.Start()
}
