package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/variety-jones/polygon"
)

type mytype struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	Desc        string     `json:"desc"`
	Constraints string     `json:"constraints"`
	Input1      string     `json:"input1"`
	Output1     string     `json:"output1"`
	Input2      string     `json:"input2"`
	Output2     string     `json:"output2"`
	Input3      string     `json:"input3"`
	Output3     string     `json:"output3"`
	testcases   []testcase `json:"-"`
}
type testcase struct {
	Id     int    `json:"id"`
	Input  string `json:"input"`
	Output string `json:"output"`
}

func CreateApiObjectFromLocal() (api polygon.PolygonApi) {
	api = polygon.PolygonApi{}

	api.ApiKey = "490135045fada7ed44f119fb7ed551e8a7aeb2e7"
	api.Secret = "88dc1edf5f6a7c22e3956518627ecd7c7490545d"
	api.ProblemId = "309839"
	return api
}

func getdet(api polygon.PolygonApi, ch chan int, pno int) {
	mt := &mytype{}
	parameters := make(map[string]string)
	id, _ := strconv.Atoi(api.ProblemId)
	mt.Id = id
	statements, _ := api.ProblemStatements(parameters)
	for _, val := range statements {
		mt.Name = val.Name
		mt.Desc = val.Legend
		mt.Constraints = val.Input
		mt.Constraints += val.Output
		mt.Constraints += val.Notes
		// fmt.Println(mt)
	}

	parameters = make(map[string]string)
	parameters["testset"] = "tests"
	tests, _ := api.ProblemTests(parameters)
	// fmt.Println(tests)
	count := 0
	tsi := 0
	tso := 0
	for _, test := range tests {
		fmt.Printf("%d - %s -- %d\n", pno, mt.Name, count)
		var t testcase
		t.Id = id
		if test.UseInStatements == true && tsi <= 2 {
			tsi++
			if tsi == 1 {
				mt.Input1 += cleanstring(test.Input)
				mt.Input1 += "\n"
			} else if tsi == 2 {
				mt.Input2 += cleanstring(test.Input)
				mt.Input2 += "\n"
			} else if tsi == 3 {
				mt.Input3 += cleanstring(test.Input)
				mt.Input3 += "\n"
			}
		}
		t.Input = cleanstring(test.Input)
		parameters := make(map[string]string)
		parameters["testset"] = "tests"
		parameters["testIndex"] = strconv.Itoa(test.Index)
		answer, _ := api.ProblemTestAnswer(parameters)
		t.Output = cleanstring(answer)
		if test.UseInStatements == true && tso <= 2 {
			tso++
			if tso == 1 {
				mt.Output1 += cleanstring(answer)
				mt.Output1 += "\n"
			} else if tso == 2 {
				mt.Output2 += cleanstring(answer)
				mt.Output2 += "\n"
			} else if tso == 3 {
				mt.Output3 += cleanstring(answer)
				mt.Output3 += "\n"
			}
		}

		// fmt.Printf("%q\n", cleanstring(test.Input))
		// fmt.Printf("%q\n", cleanstring(answer))
		mt.testcases = append(mt.testcases, t)
		count++
	}
	postBody, _ := json.Marshal(mt)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:8080/question", "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	for _, testcase := range mt.testcases {
		postBody, _ := json.Marshal(testcase)
		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post("http://localhost:8080/testcase", "application/json", responseBody)
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}
		defer resp.Body.Close()
	}
	<-ch
}

func cleanstring(s string) string {
	a := strings.TrimLeft(s, " \n\r\t")
	a = strings.TrimRight(a, " \n\r\t")
	a = strings.Replace(a, "\r", "", -1)
	return a
}

func main() {
	api := CreateApiObjectFromLocal()
	parameters := make(map[string]string)
	// getdet(api)
	problems, _ := api.ProblemsList(parameters)
	// var wg sync.WaitGroup
	ch := make(chan int, 5) // Buffered channel with capacity 5
	pno := 1
	for _, problem := range problems {
		// wg.Add(1)
		api.ProblemId = strconv.Itoa(problem.Id)
		fmt.Println(problem.Name)
		ch <- 1
		go getdet(api, ch, pno)
		pno++
		// getdet(api, ch)

	}
	// wg.Wait()

}
