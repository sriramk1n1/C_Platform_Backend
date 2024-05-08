package main

import (
	"fmt"
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
	loginhandler := handlers.NewLogin(l)
	reghandler := handlers.NewRegister(l)
	seshandler := handlers.NewSession(l)
	rechandler := handlers.NewRecommend(l)
	sm := http.NewServeMux()
	sm.Handle("/", qhandler)
	sm.Handle("/question", qqhandler)
	sm.Handle("/document", dhandler)
	sm.Handle("/testcase", thandler)
	sm.Handle("/login", loginhandler)
	sm.Handle("/register", reghandler)
	sm.Handle("/session", seshandler)
	sm.Handle("/recommend", rechandler)

	godotenv.Load()
	s := &http.Server{Addr: os.Getenv("PORT"),
		Handler:      sm,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	fmt.Print("Server is running ...")
	err := s.ListenAndServe()
	fmt.Println(err, err.Error())
}
