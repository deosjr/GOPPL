
package main

import (
	"testing"
	
	"GOPPL/prolog"
)

var unifytests = []struct {
	s string
	//TODO add alias
}{
	{"a, a."},
	{"X, 1."},
}

func TestUnification(t *testing.T) {
	for _, tt := range unifytests {
		alias := make(map[*prolog.Var]prolog.Term)
		terms := parseQuery(tt.s)
		term1, term2 := terms[0], terms[1]
		unifies, _ := term1.UnifyWith(term2, alias)
			
		if !unifies {
			t.Errorf("%s failed to unify with %s", term1.String(), term2.String())
		}

	}
}