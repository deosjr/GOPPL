
package prolog

//import "fmt"

// List Predicate = {"LIST", 2}
type List struct {
	Compound_Term
}

func (l List) String() string{ 
	if l.isEmpty() {
		return "[]"
	}
	tail := l.tail()
	switch tail.(type){
	case List:
		ltail := tail.(List)
		if ltail.isEmpty() {
			return "[" + l.head().String() + "]"
		}
		rec := ltail.String()
		return  "[" + l.head().String() + "," + rec[1:len(rec)-1] + "]"
	}
	// TODO: doesnt take into account var X halfway string thats otherwise grounded
	return  "[" + l.head().String() + "|" + tail.String() + "]"
}

func (l List) head() Term {
	t := l.args[0]
	if l.isEmpty() {
		panic("Attempted to get head of []")	//TODO: better solution
	}
	return t
}

func (l List) tail() Term {
	if l.isEmpty() {
		panic("Attempted to get tail of []")	//TODO: better solution
	}
	return l.args[1]
}

func (l List) isEmpty() bool {
	t := l.args[0]
	switch t.(type) {
	case Atom:
		if t.(Atom).value == "EMPTYLIST" {
			return true
		}
	}
	return false
}

