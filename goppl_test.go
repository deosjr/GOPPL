
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
//TODO: parse queries instead of constructing them

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
	memory.InitFromFile("tests/example_test.pl")
	memory.InitBuiltIns()
	qx := &prolog.Var{"X"}
	px := prolog.Compound_Term{prolog.Predicate{"p",1}, prolog.Terms{qx}}
	query := prolog.Terms{px}
	testAnswers := []map[string]string{
		{"X":"a"},
		{"X":"a"},
		{"X":"b"},
		{"X":"d"},
	}
	evaluateQuery(t, query, testAnswers)
}

func TestPeano(t *testing.T) {
	memory.InitFromFile("tests/peano_test.pl")
	memory.InitBuiltIns()
	s := prolog.Predicate{"s",1}
	x := &prolog.Var{"X"}
	//query := prolog.Terms{prolog.Compound_Term{prolog.Predicate{"int",1}, prolog.Terms{x}}}
	s2 := prolog.Compound_Term{s, prolog.Terms{prolog.Compound_Term{s, prolog.Terms{prolog.Atom{"0"}}}}}
	s3 := prolog.Compound_Term{s, prolog.Terms{prolog.Compound_Term{s, prolog.Terms{prolog.Compound_Term{s, prolog.Terms{prolog.Atom{"0"}}}}}}}
	sum := prolog.Compound_Term{prolog.Predicate{"sum",3}, prolog.Terms{s2,s3,x}}
	query := prolog.Terms{sum}
	testAnswers := []map[string]string{
		{"X":"s(s(s(s(s(0)))))"},
	}
	evaluateQuery(t, query, testAnswers)
}

func TestLists(t *testing.T) {
	memory.InitFromFile("tests/lists_test.pl")
	memory.InitBuiltIns()
	list := prolog.Predicate{"LIST",2}
	l12345 := prolog.List{prolog.Compound_Term{list, prolog.Terms{prolog.Atom{"1"}, prolog.List{prolog.Compound_Term{list, prolog.Terms{prolog.Atom{"2"}, prolog.List{prolog.Compound_Term{list, prolog.Terms{prolog.Atom{"3"}, prolog.List{prolog.Compound_Term{list, prolog.Terms{prolog.Atom{"4"}, prolog.List{prolog.Compound_Term{list, prolog.Terms{prolog.Atom{"5"}, prolog.Empty_List}}}}}}}}}}}}}}}
	lx := &prolog.Var{"L"}
	x := &prolog.Var{"X"}
	cat := prolog.Compound_Term{prolog.Predicate{"cat",3}, prolog.Terms{lx,x,l12345}}
	query := prolog.Terms{cat}
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