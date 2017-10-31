cat([], L, L).
cat([H|T], L, [H|R]) :-
	cat(T, L, R).

length([], 0).
length([_|T], N) :-
	length(T, N1),
	is(N, +(N1, 1)).

sumlist([], 0).
sumlist([H|T], N) :- 
	sumlist(T, N1), is(N, +(N1, H)).

member(X, [X|_]) :- \=(X, []).
member(X, [_|T]) :- 
	member(X, T).

reverse(List, Reversed) :-
          reverse(List, [], Reversed).

reverse([], Reversed, Reversed).
reverse([Head|Tail], SoFar, Reversed) :-
          reverse(Tail, [Head|SoFar], Reversed).