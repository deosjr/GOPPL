
package memory

import (
	"errors"

	"GOPPL/prolog"
)

var (
	instantiationError = errors.New("arguments insufficiently instantiated")
	zeroDivisionError = errors.New("division by zero")
)

func evaluate(t prolog.Term, a prolog.Bindings) (int, error) {
	switch t.(type) {
	case prolog.Int:
		return t.(prolog.Int).IntValue(), nil
	case *prolog.Var:
		value, ok := a[t.(*prolog.Var)]
		if ok {
			return evaluate(value, a)
		} 
		return 0, instantiationError
	case prolog.Compound_Term:
		ct := t.(prolog.Compound_Term)
		if ct.Pred.Arity != 2 {
			panic("expected an arithmetic expression")
		}
		v1, err1 := evaluate(ct.Args[0], a)
		v2, err2 := evaluate(ct.Args[1], a)
		if err1 != nil {
			return 0, err1
		}
		if err2 != nil {
			return 0, err2
		}
		switch ct.Pred.Functor {
		case "+":
			return v1 + v2, nil
		case "-":
			return v1 - v2, nil
		case "*":
			return v1 * v2, nil
		case "/":
			if v2 == 0 {
				return 0, zeroDivisionError
			}
			// TODO: using ints makes this very imprecise!
			return v1 / v2, nil
		}
	}
	panic("expected an arithmetic expression")
	return 0, nil
}

func evaluateInstantiated(x prolog.Term, a prolog.Bindings) int {
	i, err := evaluate(x, a)
	if err != nil {
		// x was insufficiently instantiated
		panic(err)
	}
	return i
}

func is(terms prolog.Terms, a prolog.Bindings) prolog.Bindings {
	x, y := terms[0], terms[1]
	xassign := evaluateInstantiated(y, a)
	xvalue, err := evaluate(x, a)
	switch err {
	case instantiationError:
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

func arithmetic_equals(terms prolog.Terms, a prolog.Bindings) prolog.Bindings {
	x, y := evaluateInstantiated(terms[0], a), evaluateInstantiated(terms[1], a)
	if x == y {
		return a
	}
	return nil
}

func arithmetic_not_equals(terms prolog.Terms, a prolog.Bindings) prolog.Bindings { 
	x, y := evaluateInstantiated(terms[0], a), evaluateInstantiated(terms[1], a)
	if x != y {
		return a
	}
	return nil
}

func arithmetic_less(terms prolog.Terms, a prolog.Bindings) prolog.Bindings { 
	x, y := evaluateInstantiated(terms[0], a), evaluateInstantiated(terms[1], a)
	if x < y {
		return a
	}
	return nil
}

func arithmetic_leq(terms prolog.Terms, a prolog.Bindings) prolog.Bindings { 
	x, y := evaluateInstantiated(terms[0], a), evaluateInstantiated(terms[1], a)
	if x <= y {
		return a
	}
	return nil
}

func arithmetic_greater(terms prolog.Terms, a prolog.Bindings) prolog.Bindings {
	x, y := evaluateInstantiated(terms[0], a), evaluateInstantiated(terms[1], a)
	if x > y {
		return a
	}
	return nil
}

func arithmetic_geq(terms prolog.Terms, a prolog.Bindings) prolog.Bindings { 
	x, y := evaluateInstantiated(terms[0], a), evaluateInstantiated(terms[1], a)
	if x >= y {
		return a
	}
	return nil
}