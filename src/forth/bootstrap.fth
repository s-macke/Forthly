\ HERE and LATEST

: HERE &HERE @ ;
: LATEST &LATEST @ ;

\ more stack operations

: rot   >r swap r> swap ;           \ ( a b c -- b c a )
: -rot  swap >r swap r> ;           \ ( a b c -- c a b )
: over  >r dup r> swap ;            \ ( a b -- a b a )
: 2drop drop drop ;                 \ ( a b -- )
: 3drop 2drop drop ;                \ ( a b c -- )
: 2dup  over over ;                 \ ( a b -- a b a b )
: tuck  dup -rot ;                  \ ( a b -- b a b )
: nip   swap drop ;                 \ ( a b -- b )

\ Copy x from the return stack to the data stack.
: r@ 1 rpick ; \ this is a word, so the top of the return stack is not the intended value. Hence 1

\ read write operations

\   we don't support bytes, only ints, so just use the first 8 bit
: c! OVER 255 AND SWAP ! DROP ;
: c@ @ 255 AND ;

\ create booleans
: TRUE -1 ;
: FALSE 0 ;

: NOT FALSE = ;

\ create basic arithmetic

: NEGATE 0 SWAP - ;

: 1+ 1 + ;
: 1- 1 - ;

: >     SWAP < ;
: <=    > NOT ;
: >=    < NOT ;
: <>    = NOT ;

: 0=    0 = ;
: 0<>   0 <> ;
: 0<    0 < ;
: 0>    0 > ;
: 0<=   0 <= ;
: 0>=   0 >= ;

: +! tuck @ + swap ! ;  \ ( n a-addr -- )

: u< < ; \  TODO: This is not correct
: within over - >r - r> u< ;

\ memory sizes

\ cell size is just one entry in the heap, but which could be an integer
: cell 1 ;
: cell+ cell + ;
: cell- cell - ;
: cells cell * ;
: char+ 1 + ;
: char- 1 - ;

\ some special characters
: '\n' 10 ;
: CR '\n' EMIT ;
: BL   32 ;
: SPACE BL EMIT ;

\ comma operations

: ,
    HERE !       \ store the integer in the compiled image
    1 &HERE +!    \ increment HERE pointer by 1 byte
;

: C,
    HERE C!       \ store the character in the compiled image
    1 &HERE +!    \ increment HERE pointer by 1 byte
;

\ array operations
: ALLOT	            \ ( n -- )
    HERE SWAP       \ ( here n )
    &HERE +!        \ ( adds n to HERE, after this the old value of HERE is still on the stack )
    DROP
;

: BUFFER  \ ( n -- )
    CREATE ALLOT ;

\ Variables and constant )

: VARIABLE CREATE 0 , ;
: CONSTANT CREATE , DOES> @ ;

: RECURSE IMMEDIATE
    ?compile
    LATEST    \ LATEST points to the word being compiled at the moment
    >CFA      \ get the codeword
    ,        \ compile it
;

\ IF THEN ELSE CONTROL STRUCTURE

: [COMPILE] IMMEDIATE     \ compile but don't execute
    ?compile
    WORD        \ get the next word
    FIND        \ find it in the dictionary
    DROP        \ drop the result. We assume, that the word is definitely in the dictionary. TODO: check this
    >CFA        \ get its codeword
    ,          \ and compile that
;

: IF IMMEDIATE
    ?compile
    ' 0BRANCH ,      \ compile 0BRANCH
    HERE              \ save location of the offset on the stack
    0 ,               \ compile a dummy offset
;

: THEN IMMEDIATE
    ?compile
    DUP
    HERE SWAP -       \ calculate the offset from the address saved on the stack
    SWAP !            \ store the offset in the back-filled location
;


: ELSE IMMEDIATE
    ?compile
    ' BRANCH ,    \ definite branch to just over the false-part
    HERE           \ save location of the offset on the stack
    0 ,            \ compile a dummy offset
    SWAP           \ now back-fill the original (IF) offset
    DUP            \ same as for THEN word above
    HERE SWAP -
    SWAP !
;

: UNLESS IMMEDIATE
    ?compile
    ' NOT ,        \ compile NOT (to reverse the test)
    [COMPILE] IF    \ continue by calling the normal IF
;


: LITERAL IMMEDIATE
    ?compile
    ' LIT ,    \ compile LIT
    ,           \ compile the literal itself (from the stack)
    ;

\ Special chars

: CHAR WORD 1+ @ ;

\ compile-time version of char
: [char] immediate   \ ( compile: <spaces>ccc -- ; runtime: --- c )
    ?compile
    char
    [compile] literal
;


: '(' [CHAR] ( ;
: ')' [CHAR] ) ;
: '"' [CHAR] " ;

\ BEGIN UNTIL CONTROL STRUCTURE

: BEGIN IMMEDIATE
    ?compile
    HERE          \ save location on the stack
;

: UNTIL IMMEDIATE
    ?compile
    ' 0BRANCH ,    \ compile 0BRANCH
    HERE -          \ calculate the offset from the address saved on the stack
    ,               \ compile the offset here
;

: AGAIN IMMEDIATE
    ?compile
    ' BRANCH ,    \ compile BRANCH
    HERE -         \ calculate the offset back
    ,              \ compile the offset here
;

: WHILE IMMEDIATE
    ?compile
    ' 0BRANCH ,    \ compile 0BRANCH
    HERE            \ save location of the offset2 on the stack
    0 ,             \ compile a dummy offset2
;

: REPEAT IMMEDIATE
    ?compile
    ' BRANCH ,    \ compile BRANCH
    SWAP           \ get the original offset (from BEGIN)
    HERE - ,       \ and compile it after BRANCH
    DUP
    HERE SWAP -    \ calculate the offset2
    SWAP !         \ and back-fill it in the original location
;

\ Comment
: ( IMMEDIATE
   1           \ allowed nested parens by keeping track of depth
   BEGIN
       KEY                \ read next character
       DUP '(' = IF    \ open paren?
           DROP        \ drop the open paren
           1+        \ depth increases
       ELSE
           ')' = IF    \ close paren?
               1-        \ depth decreases
           THEN
       THEN
   DUP 0= UNTIL        \ continue until we reach matching close paren, depth 0
   DROP        \ drop the depth counter
;


\ Do Loops

: DO IMMEDIATE
    ?compile
    ' >r , \ save start
    ' >r , \ save limit
    HERE
;

: LOOP IMMEDIATE
    ?compile
    ' r> ,      \ get the limit
    ' r> ,      \ get the start
    ' 1+ ,      \ increment the start
    ' 2dup ,    \ duplicate the start and limit
    ' >r ,      \ save the new start
    ' >r ,      \ save the limit
    ' = ,       \ compare the start and limit
    ' 0branch , \ branch if start != limit
    HERE - ,     \ compile the offset
    ' rdrop ,   \ drop the limit
    ' rdrop ,   \ drop the start
;

: UNLOOP IMMEDIATE
    ?compile
    ' rdrop ,   \ drop the limit
    ' rdrop ,   \ drop the start
    ;

: i 2 rpick ;
: j 4 rpick ;
: k 6 rpick ;

\ === Integer Arithmetic (that require control flow words) ===
\ ( a b -- c )
: max 2dup > if drop else nip then ;
: min 2dup < if drop else nip then ;

: abs dup 0< if negate then ;

\ STRINGS

\ comma operator for strings, stores the string as c-addr
: ,"
         HERE                \ get the address where the length will be stored
         HERE DUP 1+         \ get the address of the first character to be stored
         0 ,                 \ just store a dummy value here, we will fill later
         BEGIN
             KEY
             DUP '"' <>
         WHILE
             C,             \ save next character
         REPEAT
         DROP               \ drop the final " character
         HERE SWAP -        \ calculate the length
         SWAP !             \ store the length
;


\ Print string
\ : type                      \ ( c-addr -- )
\    begin
\        dup c@ dup
\    while                   \ while c<>\0
\        emit 1+
\    repeat
\    2drop
\ ;

\ Print string generated by s"
: TYPE
0 do
    dup i +
    c@ emit
loop
drop
;

: S" IMMEDIATE               \ ( -- addr len )
     STATE @ IF              \ compiling?
        ' LITSTRING ,       \ compile LITSTRING
        HERE                 \ save the address of the length word on the stack
        0 ,                  \ dummy length - we don't know what it is yet
        BEGIN
            KEY              \ get next character of the string
            DUP '"' <>
        WHILE
            C,               \ copy character
        REPEAT
        DROP                 \ drop the double quote character at the end
        DUP                  \ get the saved address of the length word
        HERE SWAP -          \ calculate the length
        1 -                  \ subtract 1 (because we measured from the start of the length word)
        SWAP !               \ and back-fill the length location
     ELSE                    \ immediate mode
         HERE                \ get the start address of the temporary space
         BEGIN
             KEY
             DUP '"' <>
         WHILE
             OVER C!        \ save next character
             1+             \ increment address
         REPEAT
         DROP               \ drop the final " character
         HERE -             \ calculate the length
         HERE               \ push the start address
         SWAP               \ addr len
     THEN
;

\ ( n -- n n | n )
\ duplicate if n<>0
: ?dup dup if dup then ;


: ." IMMEDIATE
    STATE @ IF         \ compiling?
      [COMPILE] S"     \ read the string, and compile LITSTRING, etc. )
      ' TYPE ,        \ compile the final TELL )
    ELSE
        \ In immediate mode, just read characters and print them until we get  to the ending double quote.
        BEGIN
            KEY
            DUP '"' = IF
                DROP    \ drop the double quote character )
                EXIT    \ return from this function )
            THEN
            EMIT
        AGAIN
    THEN
;

: SPACES
    BEGIN
        DUP 0>
    WHILE
        SPACE
        1-
    REPEAT
    DROP
;

\ https://forth-standard.org/standard/string/DivSTRING
: /STRING
    DUP >R - SWAP R> + SWAP ;

\ Switch case

: CASE IMMEDIATE
    ?compile
    0        \ push 0 to mark the bottom of the stack
;

: OF IMMEDIATE
    ?compile
    ' OVER ,      \ compile OVER
    ' = ,           \ compile =
    [COMPILE] IF   \ compile IF
    ' DROP ,      \ compile DROP )
;

: ENDOF IMMEDIATE
    ?compile
    [COMPILE] ELSE    \ ENDOF is the same as ELSE
;

: ENDCASE IMMEDIATE
    ?compile
    ' DROP ,    \ compile DROP
    \ keep compiling THEN until we get to our zero marker
    BEGIN
        ?DUP
    WHILE
        [COMPILE] THEN
    REPEAT
;

\ ( "name" -- xt )
\ compile time tick
: ['] IMMEDIATE
    ?compile
    '                   \ read name and get xt
    [compile] literal   \ call literal
;

\ TODO really immediate?
: ABORT" IMMEDIATE
  ." Abort: "
  [COMPILE] s"
  type
  cr
  halt
;

\ TODO is this correct?
: ABORT
  cr ." Abort" cr
  quit
;


\ FILL and ERASE

: FILL ( c-addr n char -- )
    >R
    BEGIN
        DUP 0<> \ check if n is not zero
    WHILE
      OVER
      R@ SWAP C!
      1 /STRING
    REPEAT
    R> 2DROP DROP ;

: ERASE ( addr u -- )
  0 FILL ;
