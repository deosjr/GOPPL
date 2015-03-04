package main

import (
	"testing"

	"GOPPL/memory"
	"GOPPL/prolog"
)

// TODO: evaluate nonterminating queries, by comparing the first X results from answer

func evaluateQuery(t *testing.T, query prolog.Terms, testAnswers []map[string]string) {
	node :=  prolog.StartDFS(query)
	for _, bindings := range testAnswers {
		result, open := <- node.Answer
		if !open {
			t.Errorf("Not enough answers")
		}
		for k, v := range result.Alias {
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
		node.Notify()
	}
	if result, open := <- node.Answer; open && result.Err != prolog.Notification {
		t.Errorf("Too many answers")
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
	query = parseQuery("length([a, [b, c, d], e, [f | g], h], N).")
	testAnswers =  []map[string]string{
		{"N":"5"},
	}
	evaluateQuery(t, query, testAnswers)
}

/** TODO: evaluateQueryTrue
func TestDifferenceLists(t *testing.T) {
	memory.InitFromFile("tests/difference_lists_test.pl")
	memory.InitBuiltIns()
	evaluateQueryTrue(t, parseQuery("pal([0],[])."))
	evaluateQueryTrue(t, parseQuery("pal([1,0,1],[])."))
	evaluateQueryTrue(t, parseQuery("pal([1,1,1,1,1],[])."))
}
*/