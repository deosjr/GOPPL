package main

import (
	"testing"

	"GOPPL/memory"
	"GOPPL/prolog"
)

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
		t.Errorf("Too many answers for query %s", query)
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
		t.Errorf("Too many answers for query %s", query)
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
	query = parseQuery("updateAliasProblem(X).")
	testAnswers = []map[string]string{
		{"X": "b"},
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
	query = parseQuery("p(X,Y).")
	testAnswers = []map[string]string{
		{"X": "1", "Y":"2"},
	}
	evaluateQuery(t, query, testAnswers)
	query = parseQuery("test(EEN).")
	testAnswers = []map[string]string{
		{"EEN": "1"},
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
	testAnswers = []map[string]string{
		{"X": "0"},
		{"X": "s(0)"},
		{"X": "s(s(0))"},
	}
	evaluateNonterminatingQuery(t, parseQuery("int(X)."), testAnswers)
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

func TestDifferenceLists(t *testing.T) {
	memory.InitFromFile("tests/difference_lists_test.pl")
	memory.InitBuiltIns()
	evaluateQueryTrue(t, parseQuery("pal([0],[])."))
	evaluateQueryTrue(t, parseQuery("pal([1,0,1],[])."))
	evaluateQueryTrue(t, parseQuery("pal([1,1,1,1,1],[])."))
	query := parseQuery("pal(S, [1,0,1,0,1], []).")
	testAnswers :=  []map[string]string{
		{"S":"s(s(0))"},
	}
	evaluateQuery(t, query, testAnswers)
}

func TestArithmetic(t *testing.T) {
	memory.InitFromFile("tests/arithmetic_test.pl")
	memory.InitBuiltIns()
	query := parseQuery("split([a,b,c,d,e,f,g,h,i,j],3,L1,L2).")
	testAnswers :=  []map[string]string{
		{"L1":"[a,b,c]", "L2":"[d,e,f,g,h,i,j]"},
	}
	evaluateQuery(t, query, testAnswers)
}