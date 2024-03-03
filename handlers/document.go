package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"first.com/data"
	"first.com/prisma"
	"first.com/prisma/db"
)

type Document struct {
	l *log.Logger
}

func NewDocument(l *log.Logger) *Document {
	return &Document{l}
}

type Result struct {
	Accepted bool   `json:"accepted"`
	Warnings string `json:"warnings"`
	Time     string `json:"Time"`
	Total    int    `json:"total"`
	Passed   int    `json:"passed"`
}

func (r *Result) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(r)
}

func (d *Document) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File not found in form data", http.StatusBadRequest)
		return
	}
	defer file.Close()
	filename := r.FormValue("fileName")
	outFile, err := os.Create(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Print(filename)
	cmd := exec.Command("g++", filename)
	var compileerror bytes.Buffer
	cmd.Stderr = &compileerror
	err = cmd.Run()
	result := &Result{}
	if err != nil {
		result.Warnings = compileerror.String()
		result.Accepted = false
		result.ToJSON(w)
		return
	}
	t, _ := getTestCases()
	result.Total = len(t)
	pass := 0
	for _, i := range t {
		cmd := exec.Command("./a.out")
		cmd.Stdin = strings.NewReader(i.Input)
		var rterror, output bytes.Buffer
		cmd.Stdout = &output
		cmd.Stderr = &rterror
		err := cmd.Run()
		if err != nil {
			result.Warnings = rterror.String()
			result.Accepted = false
			result.Passed = pass
			result.ToJSON(w)
			return
		}
		fmt.Println(output.String(), "|", i.Output)
		if strings.TrimSpace(output.String()) == strings.TrimSpace(i.Output) {
			pass++
		} else {
			fmt.Printf("failed")
			result.Accepted = false
			result.Passed = pass
			result.ToJSON(w)
			return
		}
	}
	os.Remove("a.out")
	os.Remove(filename)
	os.Remove(filename)
	result.Accepted = true
	result.Passed = pass
	result.ToJSON(w)
}

func getTestCases() (data.TestCases, error) {
	t := data.TestCases{}
	return t, prisma.HandleDBOperation(func(client *db.PrismaClient) error {
		ctx := context.Background()
		res, err := client.TestCase.FindMany(db.TestCase.Qid.Equals(1)).Exec(ctx)
		for _, i := range res {
			cur := &data.TestCase{
				Input:  i.Input,
				Output: i.Output,
			}
			t = append(t, cur)
		}
		return err
	})
}
