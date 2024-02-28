package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"first.com/handlers"
)

func main() {

	l := log.New(os.Stdout, "--", log.LstdFlags)
	qhandler := handlers.NewQuestion_list(l)
	qqhandler := handlers.NewQuestion(l)
	dhandler := handlers.NewDocument(l)
	sm := http.NewServeMux()
	sm.Handle("/", qhandler)
	sm.Handle("/question", qqhandler)
	sm.Handle("/document", dhandler)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      sm,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	s.ListenAndServe()
}
