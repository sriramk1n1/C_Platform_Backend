package handlers

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/sriramk1n1/C_Platform_Backend/data"
)

type Question struct {
	l *log.Logger
}

func NewQuestion(l *log.Logger) *Question {
	return &Question{l}
}

func (q *Question) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		rw.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS"))
		q.getQuestion(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		q.addQuestion(rw, r)
		return
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)

}

func (q *Question) getQuestion(rw http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(userID)
	question, err := data.GetQuestion(id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	err = question.ToJSON(rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func (q *Question) addQuestion(rw http.ResponseWriter, r *http.Request) {
	question := &data.Question{}
	question.FromJSON(r.Body)
	err := data.AddQuestion(question)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
}
