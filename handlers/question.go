package handlers

import (
	"log"
	"net/http"
	"strconv"

	"first.com/data"
)

type Question struct {
	l *log.Logger
}

func NewQuestion(l *log.Logger) *Question {
	return &Question{l}
}

func (q *Question) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	userID := r.URL.Query().Get("id")
	// fmt.Printf("%#v", r.URL.Query())
	id, _ := strconv.Atoi(userID)
	question := data.GetQuestion(id)
	err := question.ToJSON(rw)
	if err != nil {
		println("ERROR: encoding json")
	}
}
