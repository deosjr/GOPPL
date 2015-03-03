
package memory

import (
	"strconv"

	"GOPPL/prolog"
)

var builtins = make( map[prolog.Predicate] prolog.Predicate)

//TODO: suppress these by default when printing memory 
func InitBuiltIns() {

	extralogical := prolog.Extralogical

	x := prolog.VarTemplate{"X"}
	y := prolog.VarTemplate{"Y"}
	anon := prolog.VarTemplate{"_"}

	//	=/2 as UNIFY(X,X)
	unify := prolog.Predicate{"UNIFY",2}
	builtins[prolog.Predicate{"=",2}] = unify
	addData(unify, prolog.Rule{prolog.Terms{x, x}, prolog.Terms{}})

	//	not/1, also \+ /1
	extralogical[prolog.Predicate{"not",1}] = not
	builtins[prolog.Predicate{"\\+",1}] = prolog.Predicate{"not",1}

	//	\= /2 as not(UNIFY)
	notunify := prolog.Predicate{"NOTUNIFY",2}
	builtins[prolog.Predicate{"\\=",2}] = notunify
	addData(notunify, prolog.Rule{prolog.Terms{x, y}, prolog.Terms{prolog.Compound_Term{prolog.Predicate{"not",1}, prolog.Terms{prolog.Compound_Term{unify, prolog.Terms{x, y}}}}}})

	//	is/2 as IS
	extralogical[prolog.Predicate{"is",2}] = is 

	// TODO: is this definition necessary?
	// Lists as LIST/2 using prolog.Atom EMPTYLIST as [] and RESERVED as end of list
	list := prolog.Predicate{"LIST",2}
	
	// LIST([], RESERVED)
	addData(list, prolog.Rule{prolog.Terms{prolog.Atom{"EMPTYLIST"}, prolog.Atom{"RESERVED"}}, prolog.Terms{}})
	
	// LIST(_, LIST(_,_))
	tlist := prolog.CreateList(prolog.Terms{anon, anon}, prolog.Empty_List)
	addData(list, prolog.Rule{prolog.Terms{anon, tlist}, prolog.Terms{}})

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
		update[v] = prolog.Atom{strconv.FormatInt(xassign, 10)}
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