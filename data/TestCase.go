package data

import (
	"context"
	"encoding/json"
	"io"

	"github.com/sriramk1n1/C_Platform_Backend/prisma"
	"github.com/sriramk1n1/C_Platform_Backend/prisma/db"
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
		_, err := client.Testcase.CreateOne(
			db.Testcase.Input.Set(t.Input),
			db.Testcase.Output.Set(t.Output),
			db.Testcase.Question.Link(db.Question.ID.Equals(t.Id)),
		).Exec(ctx)
		return err
	})
}
