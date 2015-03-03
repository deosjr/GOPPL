
package memory

import (
	"strconv"

	"GOPPL/prolog"
)

var builtins = make( map[prolog.Predicate] prolog.Compound_Term)

//TODO: suppress these by default when printing memory 
func InitBuiltIns() {

	extralogical := prolog.Extralogical

	x := prolog.VarTemplate{"X"}
	anon := prolog.VarTemplate{"_"}

	//	=/2 as UNIFY(X,X)
	unify := prolog.Compound_Term{prolog.Predicate{"UNIFY",2}, prolog.Terms{x, x}}
	builtins[prolog.Predicate{"=",2}] = unify
	addData(unify.Pred, prolog.Rule{unify.Args, prolog.Terms{}})

	//TODO:
	//	not/1
	//	\=/2 as not(UNIFY)

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