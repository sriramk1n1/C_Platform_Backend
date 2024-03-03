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

	if r.Method == http.MethodGet {
		q.getQuestion(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		q.addQuestion(rw, r)
		return
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)

}

func (q *Question) getQuestion(rw http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(userID)
	question := data.GetQuestion(id)
	err := question.ToJSON(rw)
	if err != nil {
		println("ERROR: encoding json")
	}
}

func (q *Question) addQuestion(rw http.ResponseWriter, r *http.Request) {
	question := &data.Question{}
	question.FromJSON(r.Body)
	data.AddQuestion(question)
}
