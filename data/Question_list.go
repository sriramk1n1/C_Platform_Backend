package data

import (
	"context"
	"encoding/json"
	"io"

	"first.com/prisma"
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
	questions := Questions{}
	return questions, prisma.HandleDBOperation(func(client *db.PrismaClient) error {
		ctx := context.Background()
		res, _ := client.Question.FindMany().Exec(ctx)
		for _, question := range res {
			ele := &partialQuestion{
				Id:   question.ID,
				Name: question.Name,
			}
			questions = append(questions, ele)
		}
		return nil
	})
}
