
package memory

import (
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
	extralogical[prolog.Predicate{"IS",2}] = is 

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
/*	if len(terms) != 2 {
		return nil
	}
	x, y := terms[0], terms[1]
	xassign, err := evaluate(y, a)
	if err != nil {
		// y was insufficiently instantiated
		return nil
	}
	xvalue, err := evaluate(x, a)
	switch err {
	case xisavar:
		update := make(alias)
		update[x] = xassign
		clash := updateAlias(a, update)
		if clash {
			return nil
		}
	case nil: // x is an expression with no vars

	}
	
*/	return nil
}

func evaluate(t prolog.Term, a prolog.Bindings) (int64, error) {
	return 0, nil
}