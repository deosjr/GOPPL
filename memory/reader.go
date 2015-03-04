// Start by copying from csvreader, then make adjustments
// See below for starting point / idiom reference
// http://golang.org/src/encoding/csv/reader.go

package memory

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode"
	
	"GOPPL/prolog"
)

type ParseError struct {
	Line int
	Column int
	Err error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d, column %d: %s", e.Line, e.Column, e.Err)
}

var (
	// TODO: meaningful errors
	ErrSyntaxError = errors.New("syntax error")
	ErrQueryError = errors.New("syntax error in query")
)

func (r *Reader) ThrowError(err error) error {
	return &ParseError {
		Line: r.line,
		Column: r.column,
		Err: err,
	}
}

type Reader struct {
	Comment rune
	And	rune
	Stop rune
	line int
	column int
	Last_Read rune
	r *bufio.Reader
	rulebase map[prolog.Predicate][]prolog.Rule
}

// expects UTF-8 input
func NewReader(r io.Reader) *Reader {
	return &Reader{
		Comment : '%',
		And : ',',
		Stop : '.',
		r: bufio.NewReader(r),
		rulebase : make(prolog.Data),
	}
}

func (r *Reader) ReadAll() (prolog.Data, error) {

	for {
		predicate, rule, err := r.Read()
		if err == io.EOF {
			return r.rulebase, nil
		}
		if err != nil {
			return nil, err
		}
		r.addData(predicate, rule)
	}
}

// Read returns the next full rule in a prolog file
func (r *Reader) Read() (prolog.Predicate, prolog.Rule, error) {

	// Check valid starting point
	r1, _, err := r.r.ReadRune()
	if err != nil {
		return prolog.Predicate{}, prolog.Rule{}, err
	}	
	// TODO: expand this simple check
	if r1 == '[' {
		return prolog.Predicate{}, prolog.Rule{}, r.ThrowError(ErrSyntaxError)
	}
	r.r.UnreadRune()

	term, err := r.ReadTerm()
	if err != nil {
		return prolog.Predicate{}, prolog.Rule{}, err
	}
	
	if r.Last_Read == r.Stop {
		switch term.(type){
		case prolog.Atom:
			predicate := prolog.Predicate{term.(prolog.Atom).Value, 0}
			rule := prolog.Rule{prolog.Terms{}, prolog.Terms{}}
			return predicate, rule, nil
		case prolog.Compound_Term:
			p := term.(prolog.Compound_Term)
			predicate := p.GetPredicate()
			rule := prolog.Rule{p.GetArgs(), prolog.Terms{}}
			return predicate, rule, nil
		default:
			return prolog.Predicate{}, prolog.Rule{}, err
		}
	}
	p, _ := term.(prolog.Compound_Term)
	
	err = r.readTurnstile()
	if err != nil {
		return prolog.Predicate{}, prolog.Rule{}, err
	}
	
	terms, err := r.ReadTerms()
	if err != nil {
		return prolog.Predicate{}, prolog.Rule{}, err
	}
	if r.Last_Read != r.Stop {
		ok, err := r.findNext(r.Stop, true)
		if !ok {
			return prolog.Predicate{}, prolog.Rule{}, err
		}
	}
	predicate := p.GetPredicate()
	rule := prolog.Rule{p.GetArgs(), terms}
	return predicate, rule, nil
	
}

func (r *Reader) readTurnstile() error {
	if ok, err := r.findNext(':', true); !ok || err != nil {
		return r.ThrowError(ErrSyntaxError)
	}
	if ok, err := r.findNext('-', false); !ok || err != nil {
		return r.ThrowError(ErrSyntaxError)
	}
	return nil
}

// ReadTerm returns one Term
// TODO: inline operators =, \=, is, etc
// as Term Op Term -> Op(Term, Term)
func (r *Reader) ReadTerm() (prolog.Term, error) {

	r1, err := r.skipCommentsAndSpaces()
	if err != nil {
		return nil, err
	}
	s := []rune{}
	for {
		if err != nil {
			return nil, err
		}
		if r1 == '\n' {
			r.line++
		}
		if r1 == '(' {
			if !checkValidFunctor(s) {
				return nil, r.ThrowError(ErrSyntaxError)
			}
			return r.readCompound(string(s))
		}
		if r1 == '[' {
			if len(s) > 0 {
				return nil, r.ThrowError(ErrSyntaxError)
			}	
			return r.readList()
		}
		if !checkValidAtomVar(r1) {
			if unicode.IsSpace(r1) {
				r1, err = r.skipCommentsAndSpaces()
			}
			return r.readAtomVar(s, err)
		}
		s = append(s, r1)
		r1, err = r.readRune()
	}
}

// ReadTerms returns a list of Terms, which where And-separated
func (r *Reader) ReadTerms() (prolog.Terms, error) {
	
	terms := prolog.Terms{}
	t, err := r.ReadTerm()
	if err != nil {
		return nil, err
	}
	terms = append(terms, t)
	for {
		ok := r.Last_Read == r.And
		if !ok && unicode.IsSpace(r.Last_Read) {
			ok, err = r.findNext(r.And, true)
			if err != nil {
				return nil, err
			}
		}
		if !ok {
			break
		}
		
		t, err := r.ReadTerm()
		if err != nil {
			return nil, err
		}
		terms = append(terms, t)		
	}
	return terms, err
	
}

func (r *Reader) checkBuiltin(ct prolog.Compound_Term, err error) (prolog.Term, error) {
	if builtin, ok := builtins[ct.Pred]; ok {
		ct.Pred = builtin
	}
	return ct, err
}

func (r *Reader) readCompound(functor string) (prolog.Term, error) {
	args, err := r.ReadTerms()
	if err != nil {
		return nil, err
	}
	if r.Last_Read != ')' {
		ok, err := r.findNext(')', true)
		if !ok {
			return nil, err
		}
	}
	_, err = r.readRune()
	predicate := prolog.Predicate{functor, len(args)}
	compound := prolog.Compound_Term{predicate, args}
	return r.checkBuiltin(compound, err)
}

func (r *Reader) readList() (prolog.Term, error) {
	args, err := r.ReadTerms()
	if len(args) == 0 {		
		if r.Last_Read != ']' {
			return nil, r.ThrowError(ErrSyntaxError)
		}
		_, err = r.readRune()
		return prolog.Empty_List, nil
	}
	switch r.Last_Read {
	case ']':
		_, err = r.readRune()
		return prolog.CreateList(args, prolog.Empty_List), err
	case '|':
		tail, err := r.ReadTerm()
		if err != nil {
			return nil, err
		}
		if r.Last_Read != ']' {
			return nil, r.ThrowError(ErrSyntaxError)
		}
		_, err = r.readRune()
		return prolog.CreateList(args, tail), err
	default:
		return nil, r.ThrowError(ErrSyntaxError)
	}
	return nil, err
}

func (r *Reader) readAtomVar(s []rune, err error) (prolog.Term, error) {
	if len(s) == 0 {
		return nil, r.ThrowError(ErrSyntaxError)
	}
	if unicode.IsUpper(s[0]) || s[0] == '_' {
		return prolog.VarTemplate{string(s)}, err
	}
	return prolog.Atom{string(s)}, err
}

func (r *Reader) readRune() (rune, error) {
	r1, _, err := r.r.ReadRune()
	
	if r1 == '\r' {
		r1, _, err = r.r.ReadRune()
		if err == nil && r1 != '\n' {
			r.r.UnreadRune()
			r1 = '\r'	// Should never happen, right?
		}
	}
	r.column++
	r.Last_Read = r1
	return r1, err
}

// TODO: '=+-*/\\' added for simplicity for now, need extended check
func checkValidAtomVar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '=' || r == '+' || r == '-' || r == '*' || r == '/' || r == '\\'
}

func checkValidFunctor(s []rune) bool {
	return len(s) != 0 && !unicode.IsUpper(s[0])
}

// findNext throws an error unless the next rune
// is equal to c
func (r *Reader) findNext(c rune, skip bool) (bool, error) {

	var r1 rune
	var err error
	if skip {
		r1, err = r.skipCommentsAndSpaces()
	} else {
		r1, err = r.readRune()
	}
	
	if err != nil {
		return false, err
	}
	if r1 == c {
		return true, nil
	} else {
		return false, nil
	}
}

func (r *Reader) skipComment() (rune, error) {
	for {
		r1, err := r.readRune()
		if err != nil {
			return r1, err
		}
		if r1 == '\n' {
			r.line++
			r.column = -1
			return r1, nil
		}
	}
}

func (r *Reader) skipCommentsAndSpaces() (rune, error) {

	r1, err := r.readRune()
	Skip:
	for err == nil {
		switch r1 {
		case '\n':
			r.line++
			r.column = -1
			r1, err = r.readRune()
		case r.Comment:
			r1, err = r.skipComment()
		default:
			if unicode.IsSpace(r1) {
				r1, err = r.readRune()
			} else {
				break Skip
			}
		}
	}
	return r1, err
}

func (r *Reader) addData(pred prolog.Predicate, rule prolog.Rule) {
	if value, ok := r.rulebase[pred]; ok {
		r.rulebase[pred] = append(value, rule)
	} else {
		r.rulebase[pred] = []prolog.Rule{rule}
	}
}