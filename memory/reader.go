//TODO: Start by copying from csvreader, then make adjustments
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
)

func (r *Reader) error(err error) error {
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
	last_read rune
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
		line : 1,
		rulebase : make(prolog.Data),
	}
}

func (r *Reader) ReadAll() (prolog.Data, error) {

	for {
		predicate, rule, err := r.Read()
		fmt.Println(predicate, rule)
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
		return prolog.Predicate{}, prolog.Rule{}, r.error(ErrSyntaxError)
	}
	r.r.UnreadRune()

	term, err := r.ReadTerm()
	if err != nil {
		return prolog.Predicate{}, prolog.Rule{}, err
	}
	p, _ := term.(prolog.Compound_Term)
	
	// TODO: no turnstile, just a Stop
	fact, err := r.findNext(r.Stop, true)
	if err != nil {
		return prolog.Predicate{}, prolog.Rule{}, err
	}
	if fact {
		predicate := p.GetPredicate()
		rule := prolog.Rule{p.GetArgs(), prolog.Terms{}}
		return predicate, rule, nil
	}
	
	err = r.readTurnstile()
	if err != nil {
		return prolog.Predicate{}, prolog.Rule{}, err
	}
	
	terms, err := r.ReadTerms()
	if err != nil {
		return prolog.Predicate{}, prolog.Rule{}, err
	}
	if r.last_read != r.Stop {
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
		return r.error(ErrSyntaxError)
	}
	if ok, err := r.findNext('-', false); !ok || err != nil {
		return r.error(ErrSyntaxError)
	}
	return nil
}

// ReadTerm returns one Term
func (r *Reader) ReadTerm() (prolog.Term, error) {

	r1, err := r.skipCommentsAndSpaces()
	//fmt.Printf("START READ: %v:%c\n", r1, r1)
	if err != nil {
		return nil, err
	}
	s := []rune{}
	for {
		//fmt.Printf("READ: %q + %c\n", string(s), r1)
		if err != nil {
			return nil, err
		}
		if r1 == '\n' {
			r.line++
		}
		if r1 == '(' {
			if unicode.IsUpper(s[0]) {
				return nil, r.error(ErrSyntaxError)
			}
			functor := string(s)
			args, err := r.ReadTerms()
			if err != nil {
				return nil, err
			}
			if r.last_read != ')' {
				fmt.Printf("LASTREAD %c \n", r.last_read)
				ok, err := r.findNext(')', true)
				fmt.Printf("NEXT %v %v \n", ok, err)
				if !ok {
					return nil, err
				}
			}
			predicate := prolog.Predicate{functor, len(args)}
			compound := prolog.Compound_Term{predicate, args}
			fmt.Printf("COMPOUND: %v \n", compound)
			return compound, err
		}
		if r1 == '[' {
			//TODO: Lists
		}
		//For now, only accept letters/digits as Atom/Var names
		if !unicode.IsLetter(r1) && !unicode.IsDigit(r1) {
			if unicode.IsSpace(r1) {
				r1, err = r.skipCommentsAndSpaces()
			}
			if unicode.IsUpper(s[0]) {
				fmt.Printf("VAR: %q \n", string(s))
				return &prolog.Var{string(s)}, err
			}
			fmt.Printf("ATOM: %q \n", string(s))
			return prolog.Atom{string(s)}, err
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
		ok := r.last_read == r.And
		if !ok && unicode.IsSpace(r.last_read) {
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
	r.last_read = r1
	return r1, err
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
		//fmt.Printf("COMMENT READS: %v:%c", r1, r1)
		if err != nil {
			return r1, err
		}
		if r1 == '\n' {
			r.line++
			return r1, nil
		}
	}
}

func (r *Reader) skipCommentsAndSpaces() (rune, error) {

	r1, err := r.readRune()
	//fmt.Printf("START SKIP: %v:%c\n", r1, r1)
	Skip:
	for err == nil {
		switch r1 {
		case '\n':
			r.line++
			r1, err = r.readRune()
		case r.Comment:
			//fmt.Printf("START COMMENT: %v:%c\n", r1, r1)
			r1, err = r.skipComment()
		default:
			if unicode.IsSpace(r1) {
				r1, err = r.readRune()
			} else {
				break Skip
			}
		}
		//fmt.Printf("SKIP READS: %v:%c\n", r1, r1)
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