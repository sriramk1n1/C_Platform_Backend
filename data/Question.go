package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"first.com/prisma/db"
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
	err := insert(q)
	if err != nil {
		return err
	}
	return nil
}

func GetQuestion(id int) *Question {
	res, _ := find(id)
	return res
}

// database operations
func insert(q *Question) error {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return err
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	client.Question.CreateOne(
		db.Question.Name.Set(q.Name),
		db.Question.Desc.Set(q.Desc),
		db.Question.Input.Set(q.Input),
		db.Question.Output.Set(q.Output),
		db.Question.Constraints.Set(q.Constraints),
	).Exec(ctx)
	fmt.Print("Inserted in to database")
	return nil
}

func find(id int) (*Question, error) {
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

	res, _ := client.Question.FindUnique(db.Question.ID.Equals(id)).Exec(ctx)
	constraints, _ := res.Constraints()
	result := &Question{
		Id:          res.ID,
		Name:        res.Name,
		Desc:        res.Desc,
		Constraints: constraints,
		Input:       res.Input,
		Output:      res.Output,
	}
	return result, nil
}
