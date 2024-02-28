package data

import (
	"encoding/json"
	"io"
)

type Question struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Desc        string `json:"desc"`
	Constraints string `json:"constraints"`
	Input       string `json:"input"`
	Output      string `json:"output"`
	Done        bool   `json:"-"`
}

type Questions []*Question

var q1 = &Question{
	1,
	"Co-primes",
	`Two numbers are called co-primes or relatively primes if gcd of those two numbers is 1.
	i.e a,b are coprimes if gcd(a,b)=1.
	You are given with a prime number n, you task is to find out how many numbers in the range
	[1,n-1] are relatively prime with n.`,
	`First line of the input is an integer (1≤t≤100) is the number of test cases.
	First line of each test case consists of a prime number n (1≤n≤1000).
	Output how many relative primes are there in range [1,n-1] to number n.`,
	"3 5 7 11",
	"4 6 10",
	true,
}

var q2 = &Question{
	2,
	"Top k",
	`You are given an array of integers and an integer k.
	 You need to print top k elements of the array in sorted order.`,
	`The first line contains two integers, N (1≤N≤106) and K (1≤K≤N).
	 The second line contains N integers, separated by space, denoting the elements of the array.
	 Each element is an integer, where (−109≤ai≤109).
	 Print top k elements of the array in sorted order.`,
	`8 3
	 4 -3 5 10 2 -4 8 3`,
	"5 8 10 ",
	true,
}

var qlist = Questions{
	&Question{Id: q1.Id, Name: q1.Name},
	&Question{Id: q2.Id, Name: q2.Name},
}

var fullqlist = Questions{q1, q2}

func (q *Question) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(q)
}

func (q *Questions) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(q)
}

func (q *Question) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(q)
}

func GetQuestions() Questions {
	return qlist
}

func AddQuestions(q *Question) {
	qlist = append(qlist, q)
}

func GetQuestion(id int) *Question {
	for i := 0; i < len(fullqlist); i++ {
		if fullqlist[i].Id == id {
			return fullqlist[i]
		}
	}
	return &Question{}
}
