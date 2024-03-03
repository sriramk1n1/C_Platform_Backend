package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sriramk1n1/C_Platform_Backend/handlers"
)

func main() {
	l := log.New(os.Stdout, "--", log.LstdFlags)
	qhandler := handlers.NewQuestionList(l)
	qqhandler := handlers.NewQuestion(l)
	dhandler := handlers.NewDocument(l)
	thandler := handlers.NewTestCase(l)
	sm := http.NewServeMux()
	sm.Handle("/", qhandler)
	sm.Handle("/question", qqhandler)
	sm.Handle("/document", dhandler)
	sm.Handle("/testcase", thandler)

	s := &http.Server{Addr: ":8080",
		Handler:      sm,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	godotenv.Load()
	s.ListenAndServe()
}
