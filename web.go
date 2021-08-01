package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed web/*
var embedFs embed.FS

func startWebserver(coordinator Coordinator) {

	webRoot, err := fs.Sub(embedFs, "web")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(coordinator, w, r)
	})
	http.Handle("/", http.FileServer(http.FS(webRoot)))
	log.Fatal(http.ListenAndServe(":9855", nil))
}
