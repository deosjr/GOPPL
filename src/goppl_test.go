package main

import (
	"fmt"
	"testing"

	"memory"
	"prolog"
)

func TestQueries(t *testing.T) {
	for _, tt := range []struct {
		file    string
		query   string
		answers []map[string]string
	}{
		{
			file:  "permutation_test.pl",
			query: "hardcoded2(X,Y).",
			answers: []map[string]string{
				{"X": "a", "Y": "a"},
				{"X": "a", "Y": "b"},
				{"X": "b", "Y": "a"},
				{"X": "b", "Y": "b"},
			},
		},
		{
			file:  "permutation_test.pl",
			query: "updateAliasProblem(X).",
			answers: []map[string]string{
				{"X": "b"},
			},
		},
		{
			file:  "example_test.pl",
			query: "p(X).",
			answers: []map[string]string{
				{"X": "a"},
				{"X": "a"},
				{"X": "b"},
				{"X": "d"},
			},
		},
		{
			file:  "example_test.pl",
			query: "p(X,Y).",
			answers: []map[string]string{
				{"X": "1", "Y": "2"},
			},
		},
		{
			file:  "example_test.pl",
			query: "test(EEN).",
			answers: []map[string]string{
				{"EEN": "1"},
			},
		},
		{
			file:  "peano_test.pl",
			query: "sum(s(s(0)), s(s(s(0))), X).",
			answers: []map[string]string{
				{"X": "s(s(s(s(s(0)))))"},
			},
		},
		{
			file:  "lists_test.pl",
			query: "cat(L, X, [1,2,3,4,5]).",
			answers: []map[string]string{
				{"L": "[]", "X": "[1,2,3,4,5]"},
				{"L": "[1]", "X": "[2,3,4,5]"},
				{"L": "[1,2]", "X": "[3,4,5]"},
				{"L": "[1,2,3]", "X": "[4,5]"},
				{"L": "[1,2,3,4]", "X": "[5]"},
				{"L": "[1,2,3,4,5]", "X": "[]"},
			},
		},
		{
			file:  "lists_test.pl",
			query: "length([a, [b, c, d], e, [f | g], h], N).",
			answers: []map[string]string{
				{"N": "5"},
			},
		},
		{
			file:  "difference_lists_test.pl",
			query: "pal(S, [1,0,1,0,1], []).",
			answers: []map[string]string{
				{"S": "s(s(0))"},
			},
		},
		{
			file:  "arithmetic_test.pl",
			query: "split([a,b,c,d,e,f,g,h,i,j],3,L1,L2).",
			answers: []map[string]string{
				{"L1": "[a,b,c]", "L2": "[d,e,f,g,h,i,j]"},
			},
		},
	} {
		memory.Reset()
		memory.InitFromFile(fmt.Sprintf("tests/%s", tt.file))
		memory.InitBuiltIns()
		q := parseQuery(tt.query)
		evaluateQuery(t, q, tt.answers)
	}
}

func TestNonterminatingQueries(t *testing.T) {
	for _, tt := range []struct {
		file    string
		query   string
		answers []map[string]string
	}{
		{
			file:  "peano_test.pl",
			query: "int(X).",
			answers: []map[string]string{
				{"X": "0"},
				{"X": "s(0)"},
				{"X": "s(s(0))"},
			},
		},
	} {
		memory.Reset()
		memory.InitFromFile(fmt.Sprintf("tests/%s", tt.file))
		memory.InitBuiltIns()
		q := parseQuery(tt.query)
		evaluateNonterminatingQuery(t, q, tt.answers)
	}
}

func TestTrueQueries(t *testing.T) {
	for _, tt := range []struct {
		file  string
		query string
	}{
		{
			file:  "difference_lists_test.pl",
			query: "pal([0],[]).",
		},
		{
			file:  "difference_lists_test.pl",
			query: "pal([1,0,1],[]).",
		},
		{
			file:  "difference_lists_test.pl",
			query: "pal([1,1,1,1,1],[]).",
		},
	} {
		memory.Reset()
		memory.InitFromFile(fmt.Sprintf("tests/%s", tt.file))
		memory.InitBuiltIns()
		q := parseQuery(tt.query)
		evaluateQueryTrue(t, q)
	}
}

func evaluateQuery(t *testing.T, query prolog.Terms, testAnswers []map[string]string) {
	node := prolog.EmptyDFS(query)
	for _, bindings := range testAnswers {
		alias := node.GetAnswer()
		if alias == nil {
			t.Errorf("Not enough answers")
		}
		for k, v := range alias {
			if _, contains := bindings[k.String()]; !contains {
				t.Errorf("Out of scope variable %s in alias", k.String())
			} else if v.String() != bindings[k.String()] {
				t.Errorf("%s bound to %s, not %s", k.String(), v.String(), bindings[k.String()])
			}
			delete(bindings, k.String())
		}
		if len(bindings) > 0 {
			t.Errorf("Unbound input variables: %v", bindings)
		}
	}
	alias := node.GetAnswer()
	if alias != nil {
		t.Errorf("Too many answers for query %s: %v", query, alias)
	}
}

func evaluateQueryTrue(t *testing.T, query prolog.Terms) {
	node := prolog.EmptyDFS(query)
	alias := node.GetAnswer()
	if alias == nil {
		t.Errorf("Query evaluated to false: %s", query)
	}
	if len(alias) > 0 {
		t.Errorf("Bindings found: %v", alias)
	}
	alias = node.GetAnswer()
	if alias != nil {
		t.Errorf("Too many answers for query %s: %v", query, alias)
	}
}

func evaluateNonterminatingQuery(t *testing.T, query prolog.Terms, testAnswers []map[string]string) {
	node := prolog.EmptyDFS(query)
	for _, bindings := range testAnswers {
		alias := node.GetAnswer()
		if alias == nil {
			t.Errorf("Not enough answers")
		}
		for k, v := range alias {
			if _, contains := bindings[k.String()]; !contains {
				t.Errorf("Out of scope variable %s in alias", k.String())
			} else if v.String() != bindings[k.String()] {
				t.Errorf("%s bound to %s, not %s", k.String(), v.String(), bindings[k.String()])
			}
			delete(bindings, k.String())
		}
		if len(bindings) > 0 {
			t.Errorf("Unbound input variables: %v", bindings)
		}
	}
}
