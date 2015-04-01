/* -*-  Prolog -*- */


part([], _, _, _).

part([X|XR], P, [X|LR], H) :-
	P >= X, part(XR, P, LR, H).

part([X|XR], P, L, [X|HR]) :-
	P < X, part(XR, P, L, HR).


qsort([], _).

qsort([X|XR], S) :-
	part(XR, X, L, H),
	qsort(L, LS), qsort(H, HS),
	append(LS, [X|HS], S).

