
package prolog

// List Predicate = {"LIST", 2}
type List struct {
	Compound_Term
}

var Empty_List List = List{Compound_Term{Predicate{"LIST",2}, Terms{Atom{"EMPTYLIST"}, Atom{"RESERVED"}}}}

func (l List) head() Term {
	t := l.Args[0]
	if l.isEmpty() {
		panic("Attempted to get head of []")	//TODO: better solution
	}
	return t
}

func (l List) tail() Term {
	if l.isEmpty() {
		panic("Attempted to get tail of []")	//TODO: better solution
	}
	return l.Args[1]
}

func (l List) isEmpty() bool {
	t := l.Args[0]
	switch t.(type) {
	case Atom:
		if t.(Atom).Value == "EMPTYLIST" {
			return true
		}
	}
	return false
}

