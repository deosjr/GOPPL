%:- consult(tests/example_test.pl)		% Include other .pl files, not parsed yet

% PARSING TODO's
%% PARSING NOTES

int(0).
int(s(M)) :- int(M).

dontCare(_).
	
difflist([],X,X).						% Only works one way, append([],[],X).
difflist([H|T],[H|HDiff],TDiff)			% not append([], X, []).
	:- difflist(T,HDiff,TDiff).
	
append(A,B,C) :-
		difflist(A,A1,A2),
		difflist(B,B1,B2),
		append(A1,A2,B1,B2,C1,[]),
		difflist(C,C1,[]).
		
append(A,B,B,C,A,C).
	
%int(0) .								%% Stop not immediately following is a syntax error
	
whyWouldYou(X) :- int(X). doThis(0).	%% Newlines are ignored, no problem.
%zeroArguments.							% For now, all rules start with a compound term
	
canWeDoThis(X, Y) :-					%% Lists and predicates over multiple lines
	areInteger([						%% parse just fine.
		X,
		Y
	]).
	
areInteger([]).							
areInteger([H|T]) :-					
	int(H),
	areInteger(T).
	
singleton(A) :- int(A), int(B).			% Singleton variables should be a syntax error!
	
isIgnored(0).
this(X) :- isIgnored(X)					% No Stop at EOF means this rule is simply
										% thrown away. Very annoying!