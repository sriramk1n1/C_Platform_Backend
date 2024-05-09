package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/sriramk1n1/C_Platform_Backend/data"
	"github.com/sriramk1n1/C_Platform_Backend/prisma"
	"github.com/sriramk1n1/C_Platform_Backend/prisma/db"
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
	Points   int    `json:"points"`
}

func (r *Result) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(r)
}

func executeWithTimeout(input string) (string, error, string) {
	cmd := exec.Command("./a.out")
	cmd.Stdin = strings.NewReader(input)
	var output, sterr bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &sterr

	if err := cmd.Start(); err != nil {
		return "", err, sterr.String()
	}

	timer := time.NewTimer(time.Second)

	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	select {
	case <-timer.C:
		if err := cmd.Process.Kill(); err != nil {
			return "", err, ""
		}
		return "", fmt.Errorf("Execution timed out"), "Time Limit Exceeded"
	case err := <-done:
		if err != nil {
			return "", err, sterr.String()
		}
	}

	return output.String(), nil, ""
}

func (d *Document) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS"))
	skip := r.FormValue("skip")

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	if err := db.Ping(); err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	defer db.Close()
	id, _ := strconv.Atoi(r.FormValue("id"))
	if skip == "true" {
		var p int
		db.QueryRow("SELECT CAST(level * 0.3 AS INT) FROM question where id = $1;", id).Scan(&p)
		w.Write([]byte(strconv.Itoa(p)))
		fmt.Println("skip", p)
		return
	}
	err = r.ParseMultipartForm(10 << 20)
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
	cmd := exec.Command("g++", filename)
	var compileerror bytes.Buffer
	cmd.Stderr = &compileerror
	err = cmd.Run()
	result := &Result{}
	if err != nil {
		result.Warnings = compileerror.String()
		result.Accepted = false
		result.Points = 0
		os.Remove("a.out")
		os.Remove(filename)
		result.ToJSON(w)
		return
	}
	t, _ := getTestCases(id)
	result.Total = len(t)
	pass := 0
	for _, i := range t {

		output, err, sterr := executeWithTimeout(i.Input)
		if err != nil {
			result.Warnings = sterr
			result.Accepted = false
			result.Passed = pass
			result.Points = 0
			result.ToJSON(w)
			os.Remove("a.out")
			os.Remove(filename)
			return
		}
		if cleanstring(output) == i.Output {
			pass++
		} else {
			result.Accepted = false
			result.Points = 0
			result.Passed = pass
			result.ToJSON(w)
			os.Remove("a.out")
			os.Remove(filename)
			return
		}
	}
	os.Remove("a.out")
	os.Remove(filename)
	db.QueryRow("SELECT CAST(level * 0.5 AS INT) FROM question where id = $1;", id).Scan(&result.Points)
	result.Accepted = true
	result.Passed = pass

	result.ToJSON(w)
}
func cleanstring(s string) string {
	a := strings.TrimLeft(s, " \n")
	a = strings.TrimRight(a, " \n")
	return a
}
func getTestCases(id int) (data.TestCases, error) {
	t := data.TestCases{}
	return t, prisma.HandleDBOperation(func(client *db.PrismaClient) error {
		ctx := context.Background()
		res, err := client.Testcase.FindMany(db.Testcase.Qid.Equals(id)).Exec(ctx)
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
