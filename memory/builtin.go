
package memory

import (
	"fmt"

	"GOPPL/prolog"
)

var builtins = make( map[prolog.Predicate] prolog.Predicate)

func pred(functor string, arity int) prolog.Predicate {
	return prolog.Predicate{functor, arity}
}

func InitBuiltIns() {

	extralogical := prolog.Extralogical
	extralogical[pred("listing", 0)] = listing

	x := prolog.VarTemplate{"X"}
	y := prolog.VarTemplate{"Y"}

	//	=/2 as UNIFY(X,X)
	unify := pred("UNIFY",2)
	builtins[pred("=",2)] = unify
	addData(unify, prolog.Rule{prolog.Terms{x, x}, prolog.Terms{}})

	//	not/1, also \+ /1
	extralogical[pred("not",1)] = not
	builtins[pred("\\+",1)] = pred("not",1)

	//	\= /2 as not(UNIFY)
	notunify := pred("NOTUNIFY",2)
	builtins[pred("\\=",2)] = notunify
	addData(notunify, prolog.Rule{prolog.Terms{x, y}, prolog.Terms{prolog.Compound_Term{pred("not",1), prolog.Terms{prolog.Compound_Term{unify, prolog.Terms{x, y}}}}}})

	//	is/2 as IS
	extralogical[pred("is",2)] = is 

	// write and writeln
	extralogical[pred("write",1)] = write
	extralogical[pred("writeln",1)] = writeln

	// DCG builtin predicate (TODO: might want to just rewrite better?)
	builtins[pred("C",3)] = pred("C",3)
	addData(pred("C",3), prolog.Rule{prolog.Terms{prolog.CreateList(prolog.Terms{x},y),x,y}, prolog.Terms{}})
}

func is(terms prolog.Terms, a prolog.Bindings) prolog.Bindings {
	if len(terms) != 2 {
		return nil
	}
	x, y := terms[0], terms[1]
	xassign, err := prolog.Evaluate(y, a)
	if err != nil {
		// y was insufficiently instantiated
		panic(err)
		return nil
	}
	xvalue, err := prolog.Evaluate(x, a)
	switch err {
	case prolog.InstantiationError:
		v, ok := x.(*prolog.Var)
		if !ok {
			panic(err)
			return nil
		}
		update := make(prolog.Bindings)
		update[v] = prolog.GetInt(xassign)
		clash := prolog.UpdateAlias(a, update)
		if clash {
			return nil
		}
		return a
	case nil: // x is an expression with no vars
		if xvalue == xassign {
			return a
		}
	}
	return nil
}

// TODO: variables in terms[0] have to be bound
// deadlocks on trying to negate a true premise?
func not(terms prolog.Terms, a prolog.Bindings) prolog.Bindings {
	if len(terms) != 1 {
		return nil
	}
	node := prolog.ContinueDFS(terms, a)
	found_nothing := true
	for result := range node.Answer {
		if result.Err == prolog.Notification {
			break
		}
		found_nothing = false
	}
	if found_nothing {
		return a
	}
	return nil
}

func write(terms prolog.Terms, a prolog.Bindings) prolog.Bindings {
	fmt.Print(terms[0].SubstituteVars(a))	
	return a
}

func writeln(terms prolog.Terms, a prolog.Bindings) prolog.Bindings {
	fmt.Println(terms[0].SubstituteVars(a))	
	return a
}

// TODO: parse listing/0 (or 0 arity 'compound' terms in general)
func listing(terms prolog.Terms, a prolog.Bindings) prolog.Bindings {
	printMemory()	
	return a
}