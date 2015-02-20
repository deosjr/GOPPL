
package main

import (
	"GOPPL/prolog"
	"testing"
)

//TODO: when parsing is functioning,
//make a table driven test using filename,expected_answers
//so we have one testfunction instead of 3 identical ones

func TestPerms(t *testing.T) {
	query := prolog.InitPerms()
	stack := prolog.InitStack(query)
	answer := make(chan prolog.Alias, 1)
	go prolog.DFS(stack, answer)
	testAnswers := []map[string]string{
		{"X":"a","Y":"a"},
		{"X":"a","Y":"b"},
		{"X":"b","Y":"a"},
		{"X":"b","Y":"b"},
	}
	for _,bindings := range testAnswers {
		alias := <- answer
		for k,v := range alias {
			if v.Term_to_string() != bindings[k.String()] {
				t.Errorf(k.String(), v.Term_to_string())
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

func TestExample(t *testing.T) {
	query := prolog.InitExample()
	stack := prolog.InitStack(query)
	answer := make(chan prolog.Alias, 1)
	go prolog.DFS(stack, answer)
	testAnswers := []map[string]string{
		{"X":"a"},
		{"X":"a"},
		{"X":"b"},
		{"X":"d"},
	}
	for _,bindings := range testAnswers {
		alias := <- answer
		for k,v := range alias {
			if v.Term_to_string() != bindings[k.String()] {
				t.Errorf(k.String(), v.Term_to_string())
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

func TestPeano(t *testing.T) {
	query := prolog.InitPeano()
	stack := prolog.InitStack(query)
	answer := make(chan prolog.Alias, 1)
	go prolog.DFS(stack, answer)
	testAnswers := []map[string]string{
		{"X":"s(s(s(s(s(0)))))"},
	}
	for _,bindings := range testAnswers {
		alias := <- answer
		for k,v := range alias {
			if v.Term_to_string() != bindings[k.String()] {
				t.Errorf(k.String(), v.Term_to_string())
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