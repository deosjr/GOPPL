
package main

import (
	"testing"
	
	"GOPPL/prolog"
)

var unifytests = []struct {
	s string
	alias map[string]string
	result map[string]string
}{
	{"a, a", nil, nil},
	{"X, 1", nil, map[string]string{"X":"1"}},
	{"X, Y", map[string]string{"X":"a", "Y":"a"}, nil},
	{"[X], [1]", nil, map[string]string{"X":"1"}},
	{"[X,Y], [1,2]", nil, map[string]string{"X":"1", "Y":"2"}},
	{"p([0,1], [1]), p([0|A], A)", nil, map[string]string{"A":"[1]"}},
	// TODO: change UpdateAlias:
	// unifies but clashes: {B:[1], C:[]} and {B:[1|C]}
	{"p(A, [1|C]), p([0|B], B)", map[string]string{"A":"[0,1]", "C":"[]"}, map[string]string{"A":"[0,1]", "C":"[]", "B":"[1]"}},
}

func TestDoesUnify(t *testing.T) {
	for _, tt := range unifytests {
		terms := parseQuery(tt.s + ".")
		term1, term2 := terms[0], terms[1]
		alias := createTestAlias(terms, tt.alias)
		unifies, resulting_bindings := term1.UnifyWith(term2, alias)
			
		if !unifies {
			t.Errorf("%s failed to unify with %s", term1.String(), term2.String())
		}

		for k,v := range resulting_bindings {
			if _, contains := tt.result[k.String()]; !contains {
				t.Errorf("Out of scope variable %s in alias", k.String())
			} else if v.String() != tt.result[k.String()] {
				t.Errorf("%s bound to %s, not %s", k.String(), v.String(), tt.result[k.String()])
			}
			delete(tt.result, k.String())
		}
		if len(tt.result) > 0 {
			t.Errorf("Unbound input variables: %v", tt.result)
		}

	}
}

var notunifytests = []struct {
	s string
	alias map[string]string
}{
	{"a, b", nil},
	{"[X], [1,2]", nil},
	{"X, Y", map[string]string{"X":"a", "Y":"b"}},
}

func TestDoesNotUnify(t *testing.T) {
	for _, tt := range notunifytests {
		terms := parseQuery(tt.s + ".")
		term1, term2 := terms[0], terms[1]
		alias := createTestAlias(terms, tt.alias)
		unifies, _ := term1.UnifyWith(term2, alias)
			
		if unifies {
			t.Errorf("%s unified with %s", term1.String(), term2.String())
		}
	}
}

func createTestAlias(terms prolog.Terms, alias map[string]string) prolog.Bindings {
	test_alias := prolog.Bindings{}
	if alias == nil {
		return test_alias
	}
	vars := prolog.VarsInTermArgs(terms)
	for v, b := range alias {
		variable := findVarInTerms(v, vars)
		bind := parseQuery(b + ".")[0]
		test_alias[variable] = bind
	}
	return test_alias
}

func findVarInTerms(vs string, vars []*prolog.Var) *prolog.Var {
	for _, v := range vars {
		if v.Name == vs {
			return v
		}
	}
	panic("Error in unit test: variable " + vs + " not in terms!")
}