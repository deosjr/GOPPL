%:- consult(tests/example_test.pl)		% Include other .pl files, not parsed yet

%% PARSING TODO's

int(0).
int(s(M)) :- int(M).
	
sum(0, M, M).
sum(s(N),M,s(K)) :-
	sum(N,M,K).
	
%int(0) .								% Stop not immediately following is a syntax error
	
whyWouldYou(X) :- int(X). doThis(0).	% Newlines are ignored, no problem.
%zeroArguments.							% For now, all rules start with a compound term

cat([], L, L).
cat([H|T], L, [H|R]) :-
	cat(T, L, R).
	
singleton(A) :- int(A), int(B).			% Singleton variables should be a syntax error!
	
isIgnored(0).
this(X) :- isIgnored(X)					% No Stop at EOF means this rule is simply
										% thrown away. Very annoying!