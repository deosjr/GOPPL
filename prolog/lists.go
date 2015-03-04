
package prolog

// List Predicate = {"LIST", 2}
type List interface {
	Compound
	head() Term
	tail() Term
	isEmpty() bool
}

type Cons struct {
	Compound_Term
	h Term
	t Term 		// can be variable, not always a List
}

type Nil struct {
	Compound_Term
}

var Empty_List List = Nil{}

func (c Cons) head() Term {
	return c.h
}

func (n Nil) head() Term {
	panic("Attempted to get head of []")
}

func (c Cons) tail() Term {
	return c.t
}

func (n Nil) tail() Term {
	panic("Attempted to get tail of []")
}

func (c Cons) isEmpty() bool {
	return false
}

func (n Nil) isEmpty() bool {
	return true
}

func CreateList(heads Terms, tail Term) Term {
	list := tail
	for i := len(heads)-1; i >= 0; i-- {
		list = Cons{Compound_Term{Predicate{"LIST",2}, Terms{heads[i], list}}, heads[i], list}
	}
	return list
}
