package data

import (
	"encoding/json"
	"io"
)

type Question struct {
	Id   int    `json:"id"`
	Desc string `json:"desc"`
	Done bool   `json:"-"`
}

type Questions []*Question

var qlist = Questions{
	{1, "Two Sum", true},
	{2, "Three Sum", false},
}

func (q *Question) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(q)
}

func (q *Questions) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(q)
}

func (q *Question) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(q)
}

func GetQuestions() Questions {
	return qlist
}

func AddQuestions(q *Question) {
	qlist = append(qlist, q)
}

func GetQuestion(id int) *Question {
	for i := 0; i < len(qlist); i++ {
		if qlist[i].Id == id {
			return qlist[i]
		}
	}
	return &Question{}
}
