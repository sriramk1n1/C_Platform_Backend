package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sriramk1n1/C_Platform_Backend/data"
)

type TestCase struct {
	l *log.Logger
}

func NewTestCase(l *log.Logger) *TestCase {
	return &TestCase{l}
}

func (t *TestCase) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS"))
	testcase := &data.TestCase{}
	testcase.FromJSON(r.Body)
	fmt.Printf("%#v", testcase)
	err := data.AddTestCase(testcase)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
}
