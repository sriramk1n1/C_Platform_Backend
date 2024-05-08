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
	Input1      string `json:"input1"`
	Output1     string `json:"output1"`
	Input2      string `json:"input2"`
	Output2     string `json:"output2"`
	Input3      string `json:"input3"`
	Output3     string `json:"output3"`
	Category    string `json:"category"`
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
			db.Question.ID.Set(q.Id),
			db.Question.Name.Set(q.Name),
			db.Question.Desc.Set(q.Desc),
			db.Question.Input1.Set(q.Input1),
			db.Question.Output1.Set(q.Output1),
			db.Question.Input2.Set(q.Input2),
			db.Question.Output2.Set(q.Output2),
			db.Question.Input3.Set(q.Input3),
			db.Question.Output3.Set(q.Output3),
			db.Question.Constraints.Set(q.Constraints),
			db.Question.Category.Set(q.Category),
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
		Input1, _ := res.Input1()
		Output1, _ := res.Output1()
		Input2, _ := res.Input2()
		Output2, _ := res.Output2()
		Input3, _ := res.Input3()
		Output3, _ := res.Output3()
		Category, _ := res.Category()

		result = &Question{

			Id:          res.ID,
			Name:        res.Name,
			Desc:        res.Desc,
			Constraints: constraints,
			Input1:      Input1,
			Output1:     Output1,
			Input2:      Input2,
			Output2:     Output2,
			Input3:      Input3,
			Output3:     Output3,
			Category:    Category,
		}
		return err
	})
}
