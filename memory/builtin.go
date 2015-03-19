
package memory

import (
	"fmt"
	"strconv"
	"strings"

	"GOPPL/prolog"
)

var builtins = make( map[prolog.Predicate] prolog.Predicate)

func pred(functor string, arity int) prolog.Predicate {
	return prolog.Predicate{functor, arity}
}

func InitBuiltIns() {

	extralogical := prolog.Extralogical
	extralogical[pred("listing", 0)] = listingAll
	extralogical[pred("listing", 1)] = listing

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

	// write and writeln
	extralogical[pred("write",1)] = write
	extralogical[pred("writeln",1)] = writeln

	// DCG builtin predicate (TODO: might want to just rewrite better?)
	builtins[pred("C",3)] = pred("C",3)
	addData(pred("C",3), prolog.Rule{prolog.Terms{prolog.CreateList(prolog.Terms{x},y),x,y}, prolog.Terms{}})

	// true and false
	extralogical[pred("true",0)] = trueFunc
	extralogical[pred("false",0)] = falseFunc
	extralogical[pred("fail",0)] = falseFunc

	// arithmetics
	extralogical[pred("is",2)] = is 
	extralogical[pred("=:=",2)] = arithmetic_equals 
	extralogical[pred("=\\=",2)] = arithmetic_not_equals 
	extralogical[pred("<",2)] = arithmetic_less 
	extralogical[pred("=<",2)] = arithmetic_leq 
	extralogical[pred(">",2)] = arithmetic_greater 
	extralogical[pred(">=",2)] = arithmetic_geq 
}

// TODO: variables in terms[0] have to be bound
// deadlocks on trying to negate a true premise?
func not(terms prolog.Terms, a prolog.Bindings) prolog.Bindings {
	node := prolog.BoundDFS(terms, a)
	if node.GetAnswer() == nil {
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

func trueFunc(terms prolog.Terms, a prolog.Bindings) prolog.Bindings {
	return a
}

func falseFunc(terms prolog.Terms, a prolog.Bindings) prolog.Bindings {
	return nil
}

func listing(terms prolog.Terms, a prolog.Bindings) prolog.Bindings {
	switch terms[0].(type) {
	case prolog.Atom:
		input := terms[0].(prolog.Atom).Value()
		if split := strings.Split(input, "/"); len(split) == 2 {
			i, err := strconv.Atoi(split[1])
			if err != nil {
				panic("Wrong argument for listing/1")
			}
			printTermInMemory(prolog.Predicate{split[0], i})
			return a
		}
	}
	panic("Wrong argument for listing/1")
	return nil	
}

func listingAll(terms prolog.Terms, a prolog.Bindings) prolog.Bindings {
	printMemory()	
	return a
}