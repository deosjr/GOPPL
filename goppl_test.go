
package main

import (
	"testing"

	"GOPPL/memory"
	"GOPPL/prolog"
)

func evaluateQuery(t *testing.T, query prolog.Terms, testAnswers []map[string]string) {
	stack := prolog.InitStack(query)
	answer := make(chan prolog.Alias, 1)
	go prolog.DFS(stack, answer)
	for _, bindings := range testAnswers {
		alias := <- answer
		for k,v := range alias {
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

//TODO: when parsing is functioning,
//make a table driven test using filename,expected_answers
//so we have one testfunction instead of 3 identical ones

func TestPerms(t *testing.T) {
	memory.InitFromFile("tests/permutation_test.pl")
	memory.InitBuiltIns()
	x,y := &prolog.Var{"X"}, &prolog.Var{"Y"}
	h := prolog.Compound_Term{prolog.Predicate{"hardcoded2",2}, prolog.Terms{x,y}}
	query := prolog.Terms{h}
	testAnswers := []map[string]string{
		{"X":"a","Y":"a"},
		{"X":"a","Y":"b"},
		{"X":"b","Y":"a"},
		{"X":"b","Y":"b"},
	}
	evaluateQuery(t, query, testAnswers)
}

func TestExample(t *testing.T) {
	query := memory.InitExample()
	testAnswers := []map[string]string{
		{"X":"a"},
		{"X":"a"},
		{"X":"b"},
		{"X":"d"},
	}
	evaluateQuery(t, query, testAnswers)
}

func TestPeano(t *testing.T) {
	query := memory.InitPeano()
	testAnswers := []map[string]string{
		{"X":"s(s(s(s(s(0)))))"},
	}
	evaluateQuery(t, query, testAnswers)
}

func TestLists(t *testing.T) {
	query := memory.InitLists()
	testAnswers := []map[string]string{
		{"L":"[]", "X":"[1,2,3,4,5]"},
		{"L":"[1]", "X":"[2,3,4,5]"},
		{"L":"[1,2]", "X":"[3,4,5]"},
		{"L":"[1,2,3]", "X":"[4,5]"},
		{"L":"[1,2,3,4]", "X":"[5]"},
		{"L":"[1,2,3,4,5]", "X":"[]"},
	}
	evaluateQuery(t, query, testAnswers)
}	