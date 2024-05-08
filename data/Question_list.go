package data

import (
	"context"
	"encoding/json"
	"io"

	"github.com/sriramk1n1/C_Platform_Backend/prisma"
	"github.com/sriramk1n1/C_Platform_Backend/prisma/db"
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
		res, _ := client.Question.FindMany().OrderBy(
			db.Question.Level.Order(db.SortOrderAsc),
		).Exec(ctx)
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
