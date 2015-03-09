split(L,0,[],L).
split([X|Xs],N,[X|Ys],Zs) :- 	
	>(N,0), 
	is(N1, -(N,1)), 
	split(Xs,N1,Ys,Zs).

% from http://www.ic.unicamp.br/~meidanis/courses/mc336/2009s2/prolog/problemas/index.html
% This is very very very slow compared to SWI-Prolog!
% Try the following goal:   ?- do([2,3,5,7,11]).

append([],L,L). 
append([H|T],L2,[H|L3]) :- append(T,L2,L3).

length([], 0).
length([_|T], N) :-
	length(T, N1),
	is(N, +(N1, 1)).

equation(L,LT,RT) :-
   split(L,LL,RL),              % decompose the list L
   term(LL,LT),                 % construct the left term
   term(RL,RT),                 % construct the right term
   =:=(LT, RT).                 % evaluate and compare the terms

term([X],X).                    % a number is a term in itself
term(L,T) :-                    % general case: binary term
   split(L,LL,RL),              % decompose the list L
   term(LL,LT),                 % construct the left term
   term(RL,RT),                 % construct the right term
   binterm(LT,RT,T).            % construct combined binary term

binterm(LT,RT,+(LT,RT)).
binterm(LT,RT,-(LT,RT)).
binterm(LT,RT,*(LT,RT)).
binterm(LT,RT,/(LT,RT)) :- =\=(RT, 0).   % avoid division by zero

split(L,L1,L2) :- append(L1,L2,L), \=(L1,[]), \=(L2,[]).

do(L) :- 
   equation(L,LT,RT),
   writeln([LT, =, RT]),
   false.