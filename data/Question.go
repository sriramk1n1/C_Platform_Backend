package data

import (
	"context"
	"encoding/json"
	"io"

	"github.com/sriramk1n1/C_Platform_Backend/prisma"
	"github.com/sriramk1n1/C_Platform_Backend/prisma/db"
)

type Question struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Desc        string `json:"desc"`
	Constraints string `json:"constraints"`
	Input       string `json:"input"`
	Output      string `json:"output"`
}

func (q *Question) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(q)
}

func (q *Question) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(q)
}

func AddQuestion(q *Question) error {
	return insert(q)
}

func GetQuestion(id int) (*Question, error) {
	return find(id)
}

// database operation functions

func insert(q *Question) error {
	return prisma.HandleDBOperation(func(client *db.PrismaClient) error {
		ctx := context.Background()
		_, err := client.Question.CreateOne(
			db.Question.Name.Set(q.Name),
			db.Question.Desc.Set(q.Desc),
			db.Question.Input.Set(q.Input),
			db.Question.Output.Set(q.Output),
			db.Question.Constraints.Set(q.Constraints),
		).Exec(ctx)
		return err
	})
}

func find(id int) (*Question, error) {
	result := &Question{}
	return result, prisma.HandleDBOperation(func(client *db.PrismaClient) error {
		ctx := context.Background()
		res, err := client.Question.FindUnique(db.Question.ID.Equals(id)).Exec(ctx)
		constraints, _ := res.Constraints()
		result = &Question{
			Id:          res.ID,
			Name:        res.Name,
			Desc:        res.Desc,
			Constraints: constraints,
			Input:       res.Input,
			Output:      res.Output,
		}
		return err
	})
}
