int(0).
int(s(M)) :- int(M).

sum(0, M, M).
sum(s(N), M, s(K)) :-
	sum(N, M, K).