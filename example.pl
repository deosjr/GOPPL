
%% PEANO

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