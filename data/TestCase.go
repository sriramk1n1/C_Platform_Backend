package data

import (
	"context"
	"encoding/json"
	"io"

	"first.com/prisma"
	"first.com/prisma/db"
)

type TestCase struct {
	Id     int    `json:"id"`
	Input  string `json:"input"`
	Output string `json:"output"`
}

type TestCases []*TestCase

func (t *TestCase) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(t)
}

func AddTestCase(t *TestCase) error {
	return inserttestcase(t)
}

// database operation functions

func inserttestcase(t *TestCase) error {
	return prisma.HandleDBOperation(func(client *db.PrismaClient) error {
		ctx := context.Background()
		_, err := client.TestCase.CreateOne(
			db.TestCase.Question.Link(db.Question.ID.Equals(t.Id)),
			db.TestCase.Input.Set(t.Input),
			db.TestCase.Output.Set(t.Output),
		).Exec(ctx)
		return err
	})
}
