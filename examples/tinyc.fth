\ From https://groups.google.com/g/comp.lang.forth/c/lBYFfVJ1qhc/m/BxvaTHCX_JgJ

(*

A translation of Marc Feeley's Tiny C compiler.
It is somewhat interesting because of the embedded virtual machine.

I have added two debugging helpers ( TRUE TO DEBUG? ):

1. a utility to show the parsed text
2. a decompiling tracer

Feel free to alter this code. It is really a stretch
to call this a "C compiler" :-)

It is probably possible to rewrite this code without using the
virtual machine part, as Forth can do this interpretively.

Run with MAIN <type some text>^D

Stop input by typing ^D (the EOI character).

It is assumed that your Forth has ENUM (Wil Baden's toolchest)
and STRUCT (gForth distribution). Words like DOC ENDDOC $"
?ALLOCATE =CELL CLEAR are hopefully obvious.
*)

\ Copyright (C) 2001 by Marc Feeley, All Rights Reserved.

NEEDS -miscutil
NEEDS -structs

ANEW -tiny-C

DOC
(*
This is a compiler for the Tiny-C language. Tiny-C is a
considerably stripped down version of C and it is meant as a
pedagogical tool for learning about compilers. The integer global
variables "a" to "z" are predefined and initialized to zero, and it
is not possible to declare new variables. The compiler reads the
program from standard input and prints out the value of the
variables that are not zero. The grammar of Tiny-C in EBNF is:

<program> ::= <statement>
<statement> ::= "if" <paren_expr> <statement> |
"if" <paren_expr> <statement> "else" <statement> |
"while" <paren_expr> <statement> |
"do" <statement> "while" <paren_expr> ";" |
"{" { <statement> } "}" |
<expr> ";" |
";"
<paren_expr> ::= "(" <expr> ")"
<expr> ::= <test> | <id> "=" <expr>
<test> ::= <sum> | <sum> "<" <sum>
<sum> ::= <term> | <sum> "+" <term> | <sum> "-" <term>
<term> ::= <id> | <int> | <paren_expr>
<id> ::= "a" | "b" | "c" | "d" | ... | "z"
<int> ::= <an_unsigned_decimal_integer>

Here are a few invocations of the compiler:

% echo "a=b=c=2<3;" | ./a.out
a = 1
b = 1
c = 1
% echo "{ i=1; while (i<100) i=i+i; }" | ./a.out
i = 128
% echo "{ i=125; j=100; while (i-j) if (i<j) j=j-i; else i=i-j; }" | ./a.out
i = 25
j = 25
% echo "{ i=1; do i=i+10; while (i<50); }" | ./a.out
i = 51
% echo "{ i=1; while ((i=i+10)<50) ; }" | ./a.out
i = 51
% echo "{ i=7; if (i<5) x=1; if (i<10) y=2; }" | ./a.out
i = 7
y = 2

The compiler does a minimal amount of error checking to help
highlight the structure of the compiler.
*)
ENDDOC

\ ---------------------------------------------------------------------------

\ Lexer.

0 ENUM lex:
lex: DO_SYM lex: ELSE_SYM lex: IF_SYM lex: WHILE_SYM
lex: LBRA lex: RBRA lex: LPAR lex: RPAR
lex: PLUS lex: MINUS lex: LESS
lex: SEMI lex: EQUAL lex: INT
lex: ID lex: EOI

$" while" $" if" $" else" $" do" CREATE cwords , , , , 0 ,

BL VALUE ch
0 VALUE sym
0 VALUE int_val

CREATE id_name #100 CHARS ALLOT

: ?syntax-error ( flag -- ) ABORT" syntax error" ;

: next_ch ( -- ) EKEY TO ch ;

: numbers? ( -- )
0 TO int_val \ missing overflow check
BEGIN ch '0' '9' 1+ WITHIN
WHILE int_val #10 * ch '0' - + TO int_val
next_ch INT TO sym
REPEAT ;

: var? ( -- )
id_name C0! \ missing overflow check
BEGIN ch 'a' 'z' 1+ WITHIN ch '_' = OR
WHILE 'OF ch 1 id_name PLACE+ next_ch
REPEAT
0 TO sym
BEGIN cwords sym CELL[] @
WHILE cwords sym CELL[] @ COUNT id_name COUNT COMPARE
WHILE 1 +TO sym
cwords sym CELL[] @
0= IF id_name C@ 1 <> ?syntax-error ID TO sym EXIT
ENDIF
REPEAT THEN ;

TRUE VALUE DEBUG?

: (next_sym) ( -- )
BEGIN
CASE ch
BL OF next_ch ENDOF
^M OF next_ch ENDOF
^D OF EOI TO sym EXIT ENDOF
'{' OF next_ch LBRA TO sym EXIT ENDOF
'}' OF next_ch RBRA TO sym EXIT ENDOF
'(' OF next_ch LPAR TO sym EXIT ENDOF
')' OF next_ch RPAR TO sym EXIT ENDOF
'+' OF next_ch PLUS TO sym EXIT ENDOF
'-' OF next_ch MINUS TO sym EXIT ENDOF
'<' OF next_ch LESS TO sym EXIT ENDOF
';' OF next_ch SEMI TO sym EXIT ENDOF
'=' OF next_ch EQUAL TO sym EXIT ENDOF
'0' '9' 1+ WITHIN
IF numbers? EXIT
ELSE ch 'a' 'z' 1+ WITHIN 0= ?syntax-error
var? EXIT
ENDIF
ENDCASE
AGAIN ;

: next_sym ( -- )
(next_sym) DEBUG? 0= ?EXIT
CR ." NEXT_SYM = "
CASE sym
DO_SYM OF ." do " ENDOF
ELSE_SYM OF ." else " ENDOF
IF_SYM OF ." if " ENDOF
WHILE_SYM OF ." while " ENDOF
EOI OF ." EOI " ENDOF
LBRA OF ." { " ENDOF
RBRA OF ." } " ENDOF
LPAR OF ." ( " ENDOF
RPAR OF ." )" ENDOF
PLUS OF ." + " ENDOF
MINUS OF ." - " ENDOF
LESS OF ." < " ENDOF
SEMI OF ." ; " ENDOF
EQUAL OF ." = " ENDOF
INT OF int_val DEC. ENDOF
ID OF id_name .$ ENDOF
ENDCASE ;

\ ---------------------------------------------------------------------------

\ Parser.

0 ENUM parse:
parse: VAR parse: CST
parse: ADD parse: SUB parse: LT parse: SET
parse: IF1 parse: IF2 parse: WHILE1 parse: DO1
parse: EMPTY
parse: SEQ
parse: EXPR1
parse: PROG

STRUCT
CELL% FIELD kind
CELL% FIELD o1
CELL% FIELD o2
CELL% FIELD o3
CELL% FIELD val
END-STRUCT node%

: new_node ( k -- node )
node% %ALLOCATE? LOCAL x
( k ) x kind !
x ;

DEFER paren_expr ( -- node ) \ forward declaration

\ <term> ::= <id> | <int> | <paren_expr>
: term ( -- node )
0 LOCAL x
sym ID = IF VAR new_node TO x
id_name CHAR+ C@ 'a' - x val !
next_sym x EXIT
ENDIF
sym INT = IF CST new_node TO x
int_val x val !
next_sym x EXIT
ENDIF
paren_expr ;

\ <sum> ::= <term> | <sum> "+" <term> | <sum> "-" <term>
: sum ( -- node )
0 LOCAL t
term LOCAL x
BEGIN sym PLUS = sym MINUS = OR
WHILE x TO t
sym PLUS = IF ADD ELSE SUB ENDIF new_node TO x
next_sym
t x o1 !
term x o2 !
REPEAT x ;

\ <test> ::= <sum> | <sum> "<" <sum>
: test ( -- node )
0 LOCAL t
sum LOCAL x
sym LESS = IF x TO t
LT new_node TO x
next_sym
t x o1 !
sum x o2 !
ENDIF x ;

\ <expr> ::= <test> | <id> "=" <expr>
: expr ( -- node ) RECURSIVE
sym ID <> IF test EXIT ENDIF
0 LOCAL t
test LOCAL x
x kind @ VAR = sym EQUAL =
AND IF x TO t
SET new_node TO x
next_sym
t x o1 !
expr x o2 !
ENDIF x ;

\ <paren_expr> ::= "(" <expr> ")"
:NONAME ( -- node )
0 LOCAL x
sym LPAR <> ?syntax-error
next_sym
expr TO x
sym RPAR <> ?syntax-error
next_sym
x ; IS paren_expr

: statement RECURSIVE ( -- node )
0 LOCAL x
0 LOCAL t
\ "if" <paren_expr> <statement>
sym IF_SYM = IF IF1 new_node TO x next_sym
paren_expr x o1 !
statement x o2 !
\ ... "else" <statement>
sym ELSE_SYM = IF IF2 x kind ! next_sym statement x o3 ! ENDIF
x EXIT
ENDIF

\ "while" <paren_expr> <statement>
sym WHILE_SYM = IF WHILE1 new_node TO x
next_sym
paren_expr x o1 !
statement x o2 !
x EXIT
ENDIF

\ "do" <statement> "while" <paren_expr> ";"
sym DO_SYM = IF DO1 new_node TO x
next_sym statement x o1 !
sym WHILE_SYM <> ?syntax-error
next_sym paren_expr x o2 !
sym SEMI <> ?syntax-error
next_sym
x EXIT
ENDIF

\ ";"
sym SEMI = IF EMPTY new_node TO x next_sym x EXIT ENDIF

\ "{" { <statement> } "}"
sym LBRA = IF EMPTY new_node TO x next_sym
BEGIN sym RBRA <>
WHILE x TO t SEQ new_node TO x
t x o1 !
statement x o2 !
REPEAT
next_sym x EXIT
ENDIF

\ <expr> ";"
EXPR1 new_node TO x
expr x o1 !
sym SEMI <> ?syntax-error
next_sym x ;

0 VALUE root

\ <program> ::= <statement>
: PROGRAM ( -- node )
PROG new_node TO root
next_sym
statement root o1 !
sym EOI <> ?syntax-error
root ;

\ ---------------------------------------------------------------------------

\ Code generator.

0 ENUM code:
code: IFETCH
code: ISTORE
code: IPUSH
code: IPOP
code: IADD
code: ISUB
code: ILT
code: JZ
code: JNZ
code: JMP
code: HALT

CREATE object #1000 CELLS ALLOT

object VALUE =here

: g ( code -- ) =here ! =CELL +TO =here ; \ missing overflow check
: hole ( -- here ) =here =CELL +TO =here ;
: fix ( 'src 'dst -- ) OVER - SWAP ! ; \ missing overflow check

: C ( node -- ) RECURSIVE
0 0 LOCALS| p1 p2 x |
CASE x kind @
VAR OF IFETCH g x val @ g ENDOF
CST OF IPUSH g x val @ g ENDOF
ADD OF x o1 @ c x o2 @ c IADD g ENDOF
SUB OF x o1 @ c x o2 @ c ISUB g ENDOF
LT OF x o1 @ c x o2 @ c ILT g ENDOF
SET OF x o2 @ c ISTORE g x o1 @ val @ g ENDOF
IF1 OF x o1 @ c JZ g hole TO p1 x o2 @ c p1 =here fix ENDOF
IF2 OF x o1 @ c JZ g hole TO p1 x o2 @ c JMP g hole TO p2
p1 =here fix x o3 @ c p2 =here fix ENDOF
WHILE1 OF =here TO p1 x o1 @ c JZ g hole TO p2 x o2 @ c
JMP g hole p1 fix p2 =here fix ENDOF
DO1 OF =here TO p1 x o1 @ c x o2 @ c JNZ g hole p1 fix ENDOF
EMPTY OF ENDOF
SEQ OF x o1 @ c x o2 @ c ENDOF
EXPR1 OF x o1 @ c IPOP g ENDOF
PROG OF x o1 @ c HALT g ENDOF
ENDCASE ;

\ ---------------------------------------------------------------------------

\ Virtual machine.

CREATE globals #26 CELLS ALLOT

0 VALUE pc
0 VALUE sp

: *pc++ ( -- addr ) pc @+ SWAP TO pc ;
: --sp ( -- ) =CELL -TO sp ;
: pc++ ( -- ) =CELL +TO pc ;
: !sp++ ( a -- ) sp !+ TO sp ;
: sp[-2] ( -- addr ) sp 2 CELLS - ;
: sp[-1] ( -- addr ) sp CELL- ;

: .INS ( code -- )
CR pc CELL- H. 4 SPACES
CASE
IFETCH OF ." IFETCH " pc @ 'a' + EMIT ENDOF
ISTORE OF ." ISTORE " pc @ 'a' + EMIT ENDOF
IPUSH OF ." IPUSH " pc @ DEC. ENDOF
IPOP OF ." IPOP " ENDOF
IADD OF ." IADD " ENDOF
ISUB OF ." ISUB " ENDOF
ILT OF ." ILT " ENDOF
JMP OF ." JMP " pc @+ + CELL- H. ENDOF
JZ OF ." JZ " pc @+ + CELL- H. ENDOF
JNZ OF ." JNZ " pc @+ + CELL- H. ENDOF
HALT OF ." HALT " ENDOF
." ERROR! " DUP H.
ENDCASE ;

: RUNS ( -- )
object TO pc
BEGIN CASE *pc++ DEBUG? IF DUP .INS ENDIF
IFETCH OF *pc++ globals []CELL @ !sp++ ENDOF
ISTORE OF sp[-1] @ *pc++ globals []CELL ! ENDOF
IPUSH OF *pc++ !sp++ ENDOF
IPOP OF --sp ENDOF
IADD OF sp[-1] @ sp[-2] +! --sp ENDOF
ISUB OF sp[-1] @ sp[-2] -! --sp ENDOF
ILT OF sp[-2] @ sp[-1] @ < 1 AND sp[-2] ! --sp ENDOF
JMP OF pc @ +TO pc ENDOF
JZ OF --sp sp @ 0= IF pc @ ELSE =CELL ENDIF +TO pc ENDOF
JNZ OF --sp sp @ IF pc @ ELSE =CELL ENDIF +TO pc ENDOF
HALT OF EXIT ENDOF
ENDCASE
AGAIN ;

\ ---------------------------------------------------------------------------

\ Main program.

: INITIALIZE ( -- )
BL TO ch
globals #26 CELLS ERASE
object TO =here
#1000 CELLS ALLOCATE ?ALLOCATE TO sp ;

: VPRINT ( -- )
#26 0 DO globals I CELL[] @
?DUP IF CR 'a' I + EMIT ." = " DEC.
ENDIF
LOOP ;

: FREE-nodes ( node -- ) RECURSIVE
DUP o1 @ ?DUP IF FREE-nodes ENDIF
DUP o2 @ ?DUP IF FREE-nodes ENDIF
FREE DROP ;

: EXITIALIZE ( -- )
sp FREE DROP CLEAR sp ;

: MAIN ( -- )
INITIALIZE
PROGRAM C root FREE-nodes
RUNS VPRINT
EXITIALIZE ;

\ EOF