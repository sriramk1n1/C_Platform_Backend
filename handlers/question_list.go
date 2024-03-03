package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/sriramk1n1/C_Platform_Backend/data"
)

type Question_list struct {
	l *log.Logger
}

func NewQuestionList(l *log.Logger) *Question_list {
	return &Question_list{l}
}

func (q *Question_list) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS"))
	qlist, _ := data.GetQuestionList()
	err := qlist.ToJSON(rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
