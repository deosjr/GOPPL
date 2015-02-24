
%% PEANO

int(0).
int(s(M)) :- sum(N,M,K).
	
sum(0, M, M).
sum(s(N),M,s(K)) :-
	sum(N,M,K).
	
%int(0) .								% Stop not immediately following is a syntax error atm	
	
%whyWouldYou(X) :- int(X). DoThis(0).	% Newlines are ignored! Parser cant handle this right now
%zeroArguments.							% For now, all rules start with a compound term

cat([], L, L).
cat([H|T], L, [H|R]) :-
	cat(T, L, R).