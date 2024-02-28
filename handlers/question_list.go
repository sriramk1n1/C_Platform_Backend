package handlers

import (
	"log"
	"net/http"

	"first.com/data"
)

type Question_list struct {
	l *log.Logger
}

func NewQuestion_list(l *log.Logger) *Question_list {
	return &Question_list{l}
}

func (q *Question_list) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodGet {
		q.getQuestions(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		q.addQuestions(rw, r)
		return
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (q *Question_list) getQuestions(rw http.ResponseWriter, r *http.Request) {
	qlist := data.GetQuestions()
	// bqlist, err := json.Marshal(qlist)
	err := qlist.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Failed to encode json", http.StatusInternalServerError)
	}
}

func (q *Question_list) addQuestions(rw http.ResponseWriter, r *http.Request) {
	qo := &data.Question{}
	err := qo.FromJSON(r.Body)
	if err != nil {
	}
	// q.l.Println(err)
	// q.l.Printf("Received %#v", qo)
	data.AddQuestions(qo)
}
