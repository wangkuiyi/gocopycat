# `gocopycat`


## Motivations

In the project [Go+](https://github.com/goplus/gop), we tailor the source of Go
parser to Go+ parser.  In this work, we need to change a small fraction of the
syntax, particularly source code in package `go/token`, `go/ast`, `go/parser`,
and etc.

The current solution is to copy-n-paste the source code from the Go codebase.  I
am considering a new approach that includes two steps:

1. To reuse code without copy-n-paste, we can use the new syntax from Go 1.9

   ```go
   type Token = token.Token
   ```

   to reuse type declarations, and adding wrapper functions like the following
   to reuse functions.

   ```go
   func ExposedFunction(param ParamType) ReturnType {
       return ast.ExposedFunction(param)
   }
   ```

1. To change some types and functions, we replace the above inheritances.

## The Solution

To generate the boilerplate code in step 1., I wrote this tool `gocopycat`.  It
takes a directory in the `$GOPATH` directory hierarchy, which usually defines one
or more Go packages, and an optional package name, from the command-line
options.  Then, it calls Go parser to retrieve global declarations of types and
functions, and generate the above boilerplate code.  In particular, for each
input file in the source directory, it produces a counterpart file in the
specified output directory.

## Examples

I run the following commands to copycat Go's `go/ast` package.

1. Get the `gocopycat` tool and make it runnable.

   ```bash
   go get -u github.com/wangkuiyi/gocopycat
   export PATH=$GOPATH/bin:$PATH
   ```

1. Retrieve the source code, making sure it is in the `$GOPATH` directory
   hierarchy.

   ```bash
   export GOPATH=$HOME/go  # or any other directory you like.
   go get -u github.com/golang/go/src/pkg/go/ast
   ```

1. Run the tool to copycat source code in the directory
   `$GOPATH/src/github.com/golang/go/src/pkg/go/ast` into
   `$GOPATH/src/github.com/wangkuiyi/ast`.

   ```bash
   gocopycat -from=$GOPATH/src/github.com/golang/go/src/pkg/go/ast \
     -pkg=ast \
     -to=$GOPATH/src/github.com/wangkuiyi/ast
   ```

   Please be aware of the `-pkg=ast` option -- without it, gocopycat converts
   all source files; with it, gocopycat converts only source files that
   implement the package `ast` and ignores those of `ast_test`.
