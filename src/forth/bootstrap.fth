\ create booleans
: TRUE 1 ;
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

\ some special characters
: '\n' 10 ;
: '"' 34 ;
: CR '\n' EMIT ;
: BL   32 ;
: SPACE BL EMIT ;



\ Variables and constant )

: VARIABLE CREATE 0 , ;
\ : constant CREATE , DOES> @ ;

: RECURSE IMMEDIATE
    LATEST @  \ LATEST points to the word being compiled at the moment
    >CFA      \ get the codeword
    ,c        \ compile it
;

\ IF THEN ELSE CONTROL STRUCTURE

: IF IMMEDIATE
    ' 0BRANCH ,c	\ compile 0BRANCH
    HERE @		    \ save location of the offset on the stack
    0 ,		        \ compile a dummy offset
;

: THEN IMMEDIATE
    DUP
    HERE @ SWAP -	\ calculate the offset from the address saved on the stack
    SWAP !		    \ store the offset in the back-filled location
;


: ELSE IMMEDIATE
    ' BRANCH ,c	\ definite branch to just over the false-part
    HERE @		\ save location of the offset on the stack
    0 ,		\ compile a dummy offset
    SWAP		\ now back-fill the original (IF) offset
    DUP		\ same as for THEN word above
    HERE @ SWAP -
    SWAP !
;


\ BEGIN UNTIL CONTROL STRUCTURE

: BEGIN IMMEDIATE
    HERE @		\ save location on the stack
;

: UNTIL IMMEDIATE
    ' 0BRANCH ,c	\ compile 0BRANCH
    HERE @ -	\ calculate the offset from the address saved on the stack
    ,		\ compile the offset here
;

: AGAIN IMMEDIATE
    ' BRANCH ,c	\ compile BRANCH
    HERE @ -	\ calculate the offset back
    ,		\ compile the offset here
;

: WHILE IMMEDIATE
    ' 0BRANCH ,c	\ compile 0BRANCH
    HERE @		\ save location of the offset2 on the stack
    0 ,		\ compile a dummy offset2
;

: REPEAT IMMEDIATE
    ' BRANCH ,c	\ compile BRANCH
    SWAP		\ get the original offset (from BEGIN)
    HERE @ - ,	\ and compile it after BRANCH
    DUP
    HERE @ SWAP -	\ calculate the offset2
    SWAP !		\ and back-fill it in the original location
;

: [COMPILE] IMMEDIATE
    WORD		\ get the next word
    FIND		\ find it in the dictionary
    >CFA		\ get its codeword
    ,c		\ and compile that
;

\ : UNLESS IMMEDIATE
\     ' NOT ,c		\ compile NOT (to reverse the test)
\     [COMPILE] IF	\ continue by calling the normal IF
\ ;

\ STRINGS

\ : C,
\     HERE @ C!	\ store the character in the compiled image
\     1 HERE +!	\ increment HERE pointer by 1 byte
\ ;

\ : S" IMMEDIATE	   	    \ -- addr len
\     STATE @ IF	        \ compiling?
\         ' LITSTRING ,   \ compile LITSTRING
\         HERE @		    \ save the address of the length word on the stack
\         0 ,		        \ dummy length - we don't know what it is yet
\         BEGIN
\             KEY 		\ get next character of the string
\             DUP '"' <>
\         WHILE
\             C,		    \ copy character
\         REPEAT
\         DROP		    \ drop the double quote character at the end )
\         DUP		        \ get the saved address of the length word )
\         HERE @ SWAP -	\ calculate the length )
\         4-		        \ subtract 4 (because we measured from the start of the length word) )
\         SWAP !		    \ and back-fill the length location )
\         ALIGN		    \ round up to next multiple of 4 bytes for the remaining code )
\     ELSE		        \ immediate mode )
\         HERE @		    \ get the start address of the temporary space )
\         BEGIN
\             KEY
\             DUP '"' <>
\         WHILE
\             OVER C!		\ save next character )
\             1+		    \ increment address )
\         REPEAT
\         DROP            \ drop the final " character )
\         HERE @ -        \ calculate the length )
\         HERE @          \ push the start address )
\         SWAP            \ addr len )
\     THEN
\ ;


: ." IMMEDIATE
	STATE @ IF	\ compiling?
\ TODO
\		[COMPILE] S"	\ read the string, and compile LITSTRING, etc. )
\		' TELL ,	\ compile the final TELL )
	ELSE
		\ In immediate mode, just read characters and print them until we get  to the ending double quote.
		BEGIN
			KEY
			DUP '"' = IF
				DROP	\ drop the double quote character )
				EXIT	\ return from this function )
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


