\ alter the behavior of do
: do immediate
    [COMPILE] do   \ first compile do
    ' i ,          \ then print current index
    ' . ,
;

: printEachChar
0 do
    dup i +
    c@ emit cr
loop
;

s" This is a string" printEachChar