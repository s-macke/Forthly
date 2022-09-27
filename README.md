# Forthly

# **[Click to run][project demo]**

This is just another Forth implementation.

- Written in Go and compiled to WebAssembly.
- [Bootstrapped](src/forth/bootstrap.fth) after the initialization of some [primitive](src/forth/primitives.go) words
- Rudimentary type checking to differ between functions, pointer to words and ints.

## Tutorial
If you are unfamiliar with Forth, try the following tutorial.
It is one of the best.

https://www.forth.com/starting-forth/

## The Heap

The heap is the most important data object in Forth. It contains data and code. 
Here it is just implemented as an array of `interface{}` or any.

```Go
heap        []any
```

Just call `dump` to see the contents of the heap. 

## How to run

To compile a standalone executable just run

```bash
cd src 
go build 
```

otherwise compile it by running `./build.sh`

## Other implementations or documentations

https://github.com/howerj/libforth

https://github.com/bshepherdson/fcc

https://gist.github.com/lbruder/10007431

https://github.com/larsbrinkhoff/forth-compiler

https://github.com/vbocan/delta-forth

https://github.com/Reschivon/movForth

https://github.com/kragen/stoneknifeforth

https://github.com/zevv/zForth

https://github.com/nornagon/jonesforth

https://github.com/sayon/forthress

https://github.com/philburk/pforth

https://github.com/larsbrinkhoff/xForth

https://gist.github.com/mbillingr/c2bdca4f618974e7e8d1449aba792b41

https://github.com/benhoyt/third

https://github.com/pzembrod/cc64

https://github.com/jkotlinski/durexforth

https://github.com/jkotlinski/acmeforth

https://github.com/skx/foth

https://github.com/unixdj/forego

https://forth.neocities.org/bootstrap/

http://lars.nocrew.org/forth2012/index.html

http://lars.nocrew.org/forth2012/core/PARSE.html

https://compilercrim.es/bootstrap/miniforth/

https://github.com/nineties/planckforth

https://bootstrapping.miraheze.org/wiki/Forth

https://github.com/tehologist/forthkit

https://github.com/hcchengithub/eforth-x86-64bits

https://github.com/mitra42/webForth

https://github.com/topics/forth?l=c&o=desc&s=updated

https://forth-standard.org/

[project demo]: https://s-macke.github.io/Forthly/
