% Family tree program. Facts hold an invented family of 4 generations, and rules derive the other relationships.
% Run the demo with `swipl -g demo -t halt family.pl`, or load the file with `swipl family.pl` and type queries.
%
% Assumptions:
% - The family has no remarriage, so siblings always share both parents.
% - A person who marries into the family has no parents in the facts.
% - A query helper that finds no answer returns [], not false.

% Facts

parent(george, david).
parent(george, susan).
parent(george, mark).
parent(helen, david).
parent(helen, susan).
parent(helen, mark).
parent(david, emma).
parent(david, jack).
parent(laura, emma).
parent(laura, jack).
parent(susan, olivia).
parent(paul, olivia).
parent(emma, lily).
parent(emma, noah).
parent(ryan, lily).
parent(ryan, noah).

male(george).
male(david).
male(mark).
male(paul).
male(jack).
male(ryan).
male(noah).

female(helen).
female(susan).
female(laura).
female(emma).
female(olivia).
female(lily).

% Required rules

grandparent(X, Y) :- parent(X, Z), parent(Z, Y).

% X \= Y comes last, so both people are bound when the test runs.
sibling(X, Y) :- parent(P, X), parent(P, Y), X \= Y.

cousin(X, Y) :- parent(PX, X), parent(PY, Y), sibling(PX, PY).

% Recursive rules

% Base case: a child is a descendant. Recursive step: a descendant of a child is a descendant.
% The recursion stops because the family tree has no cycle.
descendant(D, A) :- parent(A, D).
descendant(D, A) :- parent(A, C), descendant(D, C).

ancestor(A, D) :- descendant(D, A).

% Query helpers

child(C, P) :- parent(P, C).

% setof/3 sorts the answers and removes the duplicates that come from two shared parents.
% setof/3 fails when it finds no answer, so each helper returns [] in that case.
children_of(Person, Children) :-
    (   setof(C, child(C, Person), Children0) -> Children = Children0 ; Children = [] ).

siblings_of(Person, Siblings) :-
    (   setof(S, sibling(Person, S), Siblings0) -> Siblings = Siblings0 ; Siblings = [] ).

descendants_of(Person, Descendants) :-
    (   setof(D, descendant(D, Person), Descendants0) -> Descendants = Descendants0 ; Descendants = [] ).

% Demo

sample_query("parent(george, david)").
sample_query("female(olivia)").
sample_query("children_of(david, Children)").
sample_query("children_of(mark, Children)").
sample_query("grandparent(george, emma)").
sample_query("setof(G, grandparent(G, lily), Grandparents)").
sample_query("siblings_of(emma, Siblings)").
sample_query("siblings_of(david, Siblings)").
sample_query("siblings_of(olivia, Siblings)").
sample_query("cousin(emma, olivia)").
sample_query("cousin(emma, jack)").
sample_query("setof(C, cousin(olivia, C), Cousins)").
sample_query("ancestor(george, noah)").
sample_query("descendants_of(emma, Descendants)").
sample_query("descendants_of(george, Descendants)").

demo :-
    format("=== Family Tree in Prolog ===~n"),
    forall(sample_query(Text), run_query(Text)).

% The query text is parsed at run time, so the printed query and the goal that runs are the same text.
run_query(Text) :-
    term_string(Goal, Text, [variable_names(Bindings)]),
    format("~n?- ~s.~n", [Text]),
    (   call(Goal)
    ->  print_answer(Bindings)
    ;   format("false.~n")
    ).

% As at the SWI-Prolog prompt, a variable that stays unbound, such as G in setof/3, is not printed.
print_answer(Bindings) :-
    include([_ = Value]>>nonvar(Value), Bindings, Bound),
    (   Bound == []
    ->  format("true.~n")
    ;   maplist(binding_text, Bound, Texts),
        atomic_list_concat(Texts, ',\n', Answer),
        format("~w.~n", [Answer])
    ).

binding_text(Name = Value, Text) :-
    format(string(Text), "~w = ~W", [Name, Value, [quoted(true), spacing(next_argument)]]).
