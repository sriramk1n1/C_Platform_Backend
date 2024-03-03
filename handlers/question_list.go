package handlers

import (
	"log"
	"net/http"

	"first.com/data"
)

type Question_list struct {
	l *log.Logger
}

func NewQuestionList(l *log.Logger) *Question_list {
	return &Question_list{l}
}

func (q *Question_list) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	qlist, _ := data.GetQuestionList()
	err := qlist.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Failed to encode json", http.StatusInternalServerError)
	}
}
