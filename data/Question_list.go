package data

import (
	"context"
	"encoding/json"
	"io"

	"first.com/prisma/db"
)

type partialQuestion struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Questions []*partialQuestion

func (q *Questions) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(q)
}

func GetQuestionList() (Questions, error) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()
	ctx := context.Background()
	res, _ := client.Question.FindMany().Exec(ctx)
	qlist := Questions{}
	for _, question := range res {
		ele := &partialQuestion{
			Id:   question.ID,
			Name: question.Name,
		}
		qlist = append(qlist, ele)
	}
	return qlist, nil
}
