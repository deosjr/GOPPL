
%% Peano

int(0).
int(s(M)) :- sum(N,M,K).
	
sum(0, M, M).
sum(s(N),M,s(K)) :-
	sum(N,M,K).