package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
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
	outFile, err := os.Create("./files/" + r.FormValue("fileName"))
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
	res := &Result{Accepted: true, Warnings: "none", Time: "null"}

	res.ToJSON(w)
}
