
%% PEANO

int(0).
int(s(M)) :- sum(N,M,K).
	
sum(0, M, M).
sum(s(N),M,s(K)) :-
	sum(N,M,K).
	
whyWouldYou(X) :- int(X). DoThis(0).	% Newlines are ignored!
%zeroArguments.							% For now, all rules start with a compound term

%cat([], L, L).							% Lists are hard, let's start easy
%cat([H|T], L, [H|R]) :-
%	cat(T, L, R).