package main

import (
	"fmt"
	"net/http"
	"os"
)

type MacVMProvider struct{}

func NewMacVMProvider() *MacVMProvider {
	return &MacVMProvider{}
}

func (p *MacVMProvider) Serve() error {
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Creating VM...")
		// Example call to macvmagt service
		// macvmagtURL := os.Getenv("MACVMAGT_URL")
		// http.Post(macvmagtURL+"/create", "application/json", r.Body)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Deleting VM...")
		// Example: DELETE via macvmagt API
		w.WriteHeader(http.StatusOK)
	})

	port := os.Getenv("PLUGIN_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("MacVM plugin server listening on :%s\n", port)
	return http.ListenAndServe(":"+port, nil)
}
