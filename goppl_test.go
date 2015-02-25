package main

import (
	"testing"

	"GOPPL/memory"
	"GOPPL/prolog"
)

// TODO: evaluate nonterminating queries, by comparing the first X results from answer

func evaluateQuery(t *testing.T, query prolog.Terms, testAnswers []map[string]string) {
	empty, answer := prolog.GetInit()
	go prolog.DFS(query, empty, answer)
	for _, bindings := range testAnswers {
		alias := <-answer
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
	_, open := <-answer
	if open {
		t.Errorf("Channel still open!")
	}
}

func TestPerms(t *testing.T) {
	memory.InitFromFile("tests/permutation_test.pl")
	memory.InitBuiltIns()
	query := parseQuery("hardcoded2(X,Y).")
	testAnswers := []map[string]string{
		{"X": "a", "Y": "a"},
		{"X": "a", "Y": "b"},
		{"X": "b", "Y": "a"},
		{"X": "b", "Y": "b"},
	}
	evaluateQuery(t, query, testAnswers)
}

func TestExample(t *testing.T) {
	memory.InitFromFile("tests/example_test.pl")
	memory.InitBuiltIns()
	query := parseQuery("p(X).")
	testAnswers := []map[string]string{
		{"X": "a"},
		{"X": "a"},
		{"X": "b"},
		{"X": "d"},
	}
	evaluateQuery(t, query, testAnswers)
}

func TestPeano(t *testing.T) {
	memory.InitFromFile("tests/peano_test.pl")
	memory.InitBuiltIns()
	query := parseQuery("sum(s(s(0)), s(s(s(0))), X).")
	testAnswers := []map[string]string{
		{"X": "s(s(s(s(s(0)))))"},
	}
	evaluateQuery(t, query, testAnswers)
}

func TestLists(t *testing.T) {
	memory.InitFromFile("tests/lists_test.pl")
	memory.InitBuiltIns()
	query := parseQuery("cat(L, X, [1,2,3,4,5]).")
	testAnswers := []map[string]string{
		{"L": "[]", "X": "[1,2,3,4,5]"},
		{"L": "[1]", "X": "[2,3,4,5]"},
		{"L": "[1,2]", "X": "[3,4,5]"},
		{"L": "[1,2,3]", "X": "[4,5]"},
		{"L": "[1,2,3,4]", "X": "[5]"},
		{"L": "[1,2,3,4,5]", "X": "[]"},
	}
	evaluateQuery(t, query, testAnswers)
}
